package exec

import (
	"bytes"
	"errors"
	"fmt"
	wr "github.com/threadedstream/wasmexperiments/internal/pkg/wasm_reader"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wbinary"
	"github.com/threadedstream/wasmexperiments/internal/pkg/werrors"
	"io"
	"os"
)

const (
	magicCookie = 0x6d736100
	version     = 0x1
)

type TableEntry struct {
	Index       uint32
	Initialized bool
}

type Module struct {
	TypesSection    *TypesSection
	ImportSection   *ImportSection
	FunctionSection *FunctionSection
	TableSection    *TableSection
	MemorySection   *MemorySection
	GlobalSection   *GlobalSection
	ExportSection   *ExportSection
	StartSection    *StartSection
	ElementSection  *ElementSection
	CodeSection     *CodeSection
	DataSection     *DataSection
	CustomSections  CustomSections
	wr              *wr.WasmReader

	FunctionIndexSpace []*Function
	GlobalIndexSpace   []*GlobalDecl

	TableIndexSpace        [][]*TableEntry
	LinearMemoryIndexSpace [][]byte
}

func NewModule(path string) (*Module, error) {
	module := new(Module)
	bs, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(bs)
	module.wr = wr.NewWasmReader(r)

	module.LinearMemoryIndexSpace = make([][]byte, 1)
	fmt.Printf("-+-NewModule-+- Addr: %p, Len: %d\n", module.LinearMemoryIndexSpace, len(module.LinearMemoryIndexSpace))
	return module, nil
}

func (m *Module) Read() error {
	// validate magic_cookie
	cookie, err := wbinary.ReadU32(m.wr)
	if err != nil {
		return err
	}
	if cookie != magicCookie {
		return errors.New("cookies do not match")
	}
	ver, err := wbinary.ReadU32(m.wr)
	if err != nil {
		return err
	}
	if ver != version {
		return errors.New("versions do not match")
	}

	if err = m.readSections(); err != nil {
		return err
	}

	if m.TableSection != nil {
		m.TableIndexSpace = make([][]*TableEntry, len(m.TableSection.Entries))
	}

	if m.FunctionIndexSpace == nil {
		m.initializeFunctionIndexSpace()
	}

	return nil
}

func (m *Module) initializeFunctionIndexSpace() {
	requiredLen := 0
	if m.FunctionSection != nil {
		requiredLen += len(m.FunctionSection.Indices)
	}

	if m.ImportSection != nil {
		requiredLen += len(m.ImportSection.Entries)
	}

	m.FunctionIndexSpace = make([]*Function, requiredLen)
	if m.ImportSection != nil {
		for i := 0; i < len(m.ImportSection.Entries); i++ {
			m.FunctionIndexSpace[i] = &Function{
				code: m.CodeSection.Entries[i].Code,
				name: m.ImportSection.Entries[i].ExportName,
			}
		}
	}

	if m.FunctionSection != nil {
		for _, idx := range m.FunctionSection.Indices {
			m.FunctionIndexSpace[idx] = &Function{
				code: m.CodeSection.Entries[idx].Code,
				name: "",
			}
		}
	}
}

func (m *Module) readSections() error {
	// types section
	sectionHandlers := map[SectionID]func() error{
		TypeSectionID:         m.readTypeSection,
		ImportSectionID:       m.readImportSection,
		FunctionSectionID:     m.readFunctionSection,
		TableSectionID:        m.readTableSection,
		LinearMemorySectionID: m.readMemorySection,
		GlobalSectionID:       m.readGlobalSection,
		ExportSectionID:       m.readExportSection,
		StartSectionID:        m.readStartSection,
		ElementSectionID:      m.readElementSection,
		CodeSectionID:         m.readCodeSection,
		DataSectionID:         m.readDataSection,
	}

	var err error
	var sectionID byte
	for err == nil {
		sectionID, err = m.wr.ReadByte()
		if err != nil {
			continue
		}
		if handler, ok := sectionHandlers[SectionID(sectionID)]; ok {
			if err = m.pushRelevantReader(); err != nil {
				return err
			}
			if err = handler(); err != nil {
				return err
			}
			m.wr.Pop()
			continue
		}
		err = werrors.ErrInvalidSectionID
	}

	if err == nil || err == io.EOF {
		return nil
	}

	return err
}

func (m *Module) pushRelevantReader() error {
	dataLen, err := wbinary.ReadVarUint32(m.wr)
	if err != nil {
		return err
	}
	sectionData := new(bytes.Buffer)
	sectionData.Grow(int(dataLen))
	sectionReader := io.LimitReader(io.TeeReader(m.wr.Peek().(io.Reader), sectionData), int64(dataLen))
	m.wr.Push(sectionReader)
	_, err = m.wr.ReadBytes(int(dataLen))
	if err != nil {
		return err
	}
	m.wr.Pop()
	m.wr.Push(sectionData)
	return nil
}

func (m *Module) validateSectionID(expected SectionID) error {
	id, err := m.wr.ReadByte()
	if err != nil {
		return err
	}

	if SectionID(id) != expected {
		return werrors.ErrInvalidSectionID
	}

	return nil
}

func (m *Module) readTypeSection() error {
	ts := new(TypesSection)
	if err := ts.Deserialize(m.wr); err != nil {
		return err
	}
	m.TypesSection = ts
	return nil
}

func (m *Module) readImportSection() error {
	is := new(ImportSection)
	if err := is.Deserialize(m.wr); err != nil {
		return err
	}
	m.ImportSection = is
	return nil
}

func (m *Module) readFunctionSection() error {
	fs := new(FunctionSection)
	if err := fs.Deserialize(m.wr); err != nil {
		return err
	}
	m.FunctionSection = fs
	return nil
}

func (m *Module) readTableSection() error {
	ts := new(TableSection)
	if err := ts.Deserialize(m.wr); err != nil {
		return err
	}
	m.TableSection = ts
	return nil
}

func (m *Module) readMemorySection() error {
	ms := new(MemorySection)
	if err := ms.Deserialize(m.wr); err != nil {
		return err
	}
	m.MemorySection = ms
	return nil
}

func (m *Module) readGlobalSection() error {
	gs := new(GlobalSection)
	if err := gs.Deserialize(m.wr); err != nil {
		return err
	}
	m.GlobalSection = gs
	return nil
}

func (m *Module) readExportSection() error {
	es := new(ExportSection)
	if err := es.Deserialize(m.wr); err != nil {
		return err
	}
	m.ExportSection = es
	return nil
}

func (m *Module) readStartSection() error {
	ss := new(StartSection)
	if err := ss.Deserialize(m.wr); err != nil {
		return err
	}
	m.StartSection = ss
	return nil
}

func (m *Module) readElementSection() error {
	es := new(ElementSection)
	if err := es.Deserialize(m.wr); err != nil {
		return err
	}
	m.ElementSection = es
	return nil
}

func (m *Module) readCodeSection() error {
	cs := new(CodeSection)
	if err := cs.Deserialize(m.wr); err != nil {
		return err
	}
	m.CodeSection = cs
	return nil
}

func (m *Module) readDataSection() error {
	ds := new(DataSection)
	if err := ds.Deserialize(m.wr); err != nil {
		return err
	}
	m.DataSection = ds
	return nil
}
