package wasm

import (
	"errors"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wasm_reader"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wbinary"
)

type SectionID uint

const (
	CustomSectionID SectionID = iota
	TypeSectionID
	ImportSectionID
	FunctionSectionID
	TableSectionID
	LinearMemorySectionID
	GlobalSectionID
	ExportSectionID
	StartSectionID
	ElementSectionID
	CodeSectionID
	DataSectionID
)

type Serializer interface {
	Serialize() error
	Deserialize(wr *wasm_reader.WasmReader) error
}

type Section interface {
	IsSection() bool
}

type FunctionSig struct {
	params  []ValueType
	results [1]ValueType // has at most one element (before rolling support for multiple return types)
}

func (fs FunctionSig) Serialize() error { return nil }

func (fs *FunctionSig) Deserialize(reader *wasm_reader.WasmReader) error {
	// force value type to be a func
	valType, err := wbinary.ReadU8(reader)
	if err != nil {
		return err
	}
	if valType != ValueTypeFunc {
		return errors.New("readTypeSection: value type must be a function")
	}

	// start filling out function signatures slice
	paramsLen, err := wbinary.ReadVarUint32(reader)
	if err != nil {
		return err
	}

	fs.params = make([]ValueType, paramsLen, paramsLen)
	for i := 0; i < int(paramsLen); i++ {
		valTyp, err := wbinary.ReadVarUint32(reader)
		if err != nil {
			return err
		}
		fs.params[i] = ValueType(valTyp)
	}
	resultsLen, err := wbinary.ReadVarUint32(reader)
	if err != nil {
		return err
	}
	if resultsLen > 1 {
		return errors.New("readTypeSection: length of results array can't exceed the value of 1 (yet)")
	}
	for i := 0; i < int(resultsLen); i++ {
		valTyp, err := wbinary.ReadVarUint32(reader)
		if err != nil {
			return err
		}
		fs.results[i] = ValueType(valTyp)
	}
	return nil
}

type TypesSection struct {
	sigs []*FunctionSig
}

func (ts TypesSection) IsSection() bool   { return true }
func (ts *TypesSection) Serialize() error { return nil }

func (ts *TypesSection) Deserialize(reader *wasm_reader.WasmReader) error {
	// read arr length
	arrLen, err := wbinary.ReadVarUint32(reader)
	if err != nil {
		return err
	}
	ts.sigs = make([]*FunctionSig, arrLen, arrLen)

	for i, sig := range ts.sigs {
		sig = new(FunctionSig)
		if err := sig.Deserialize(reader); err != nil {
			return err
		}
		ts.sigs[i] = sig
	}

	return nil
}

type ExternalKind int

const (
	FunctionKind ExternalKind = iota
	TableKind
	MemoryKind
	GlobalKind
)

func (kind ExternalKind) String() string {
	switch kind {
	default:
		return "???"
	case FunctionKind:
		return "FunctionKind"
	case TableKind:
		return "TableKind"
	case MemoryKind:
		return "MemoryKind"
	case GlobalKind:
		return "GlobalKind"
	}
}

func (kind ExternalKind) Serialize() error { return nil }

func (kind *ExternalKind) Deserialize(reader *wasm_reader.WasmReader) error {
	bs, err := reader.ReadBytes(1)
	if err != nil {
		return err
	}
	*kind = ExternalKind(bs[0])
	return nil
}

type ImportDesc interface {
	Kind() ExternalKind
	IsImportDesc() bool
}

type FunctionKindDesc struct {
	SigIndex uint32
}

func NewFunctionKindDesc(sigIndex uint32) *FunctionKindDesc {
	fk := new(FunctionKindDesc)
	fk.SigIndex = sigIndex
	return fk
}

func (_ FunctionKindDesc) Kind() ExternalKind { return FunctionKind }
func (_ FunctionKindDesc) IsImportDesc() bool { return true }
func (_ FunctionKindDesc) Serialize() error   { return nil }

func (f *FunctionKindDesc) Deserialize(reader *wasm_reader.WasmReader) error {
	index, err := wbinary.ReadVarUint32(reader)
	if err != nil {
		return err
	}
	f.SigIndex = index
	return nil
}

type ElementType int

const FuncRefElementType ElementType = 0x70

type Table struct {
	ElemType ElementType
	Limits   ResizableLimits
}

func (_ Table) Serialize() error { return nil }

func (_ *Table) Deserialize(reader *wasm_reader.WasmReader) error {
	return nil
}

type TableKindDesc struct {
	Table Table
}

func NewTableKindDesc(table Table) *TableKindDesc {
	desc := new(TableKindDesc)
	desc.Table = table
	return desc
}

func (_ TableKindDesc) Kind() ExternalKind { return TableKind }
func (_ TableKindDesc) IsImportDesc() bool { return true }
func (_ TableKindDesc) Serialize() error   { return nil }

func (tk *TableKindDesc) Deserialize(reader *wasm_reader.WasmReader) error {
	if err := tk.Table.Deserialize(reader); err != nil {
		return err
	}
	return nil
}

type ResizableLimits struct {
	Flags   uint32
	Minimum uint32
}

type MemoryKindDesc struct {
	Limits ResizableLimits
}

func (_ MemoryKindDesc) Kind() ExternalKind { return MemoryKind }
func (_ MemoryKindDesc) IsImportDesc() bool { return true }
func (_ MemoryKindDesc) Serialize() error   { return nil }

func (m *MemoryKindDesc) Deserialize(reader *wasm_reader.WasmReader) error {
	var err error
	if m.Limits.Flags, err = wbinary.ReadVarUint32(reader); err != nil {
		return err
	}
	if m.Limits.Minimum, err = wbinary.ReadVarUint32(reader); err != nil {
		return err
	}
	return nil
}

type GlobalKindDesc struct {
	Type    ValueType
	Mutable bool
}

func (_ GlobalKindDesc) Kind() ExternalKind { return GlobalKind }
func (_ GlobalKindDesc) IsImportDesc() bool { return true }
func (_ GlobalKindDesc) Serialize() error   { return nil }

func (g *GlobalKindDesc) Deserialize(reader *wasm_reader.WasmReader) error {
	return nil
}

type ImportEntry struct {
	ModuleName  string
	ExportName  string
	Description ImportDesc
}

func (ie *ImportEntry) Serialize() error { return nil }

func (ie *ImportEntry) Deserialize(reader *wasm_reader.WasmReader) error {
	var err error
	if ie.ModuleName, err = wbinary.ReadUTF8StringUint(reader); err != nil {
		return err
	}
	if ie.ExportName, err = wbinary.ReadUTF8StringUint(reader); err != nil {
		return err
	}
	var kind ExternalKind
	if err = kind.Deserialize(reader); err != nil {
		return err
	}
	switch kind {
	case FunctionKind:
		fk := new(FunctionKindDesc)
		if e := fk.Deserialize(reader); e != nil {
			return e
		}
		ie.Description = fk
	case TableKind:
		tk := new(TableKindDesc)
		if e := tk.Deserialize(reader); e != nil {
			return e
		}
		ie.Description = tk
	case MemoryKind:
		mk := new(MemoryKindDesc)
		if e := mk.Deserialize(reader); e != nil {
			return e
		}
		ie.Description = mk
	case GlobalKind:
		gk := new(GlobalKindDesc)
		if e := gk.Deserialize(reader); e != nil {
			return e
		}
		ie.Description = gk
	}
	return nil
}

type ImportSection struct {
	Entries []*ImportEntry
}

func (i ImportSection) IsSection() bool { return true }

func (i *ImportSection) Deserialize(reader *wasm_reader.WasmReader) error {
	count, err := wbinary.ReadVarUint32(reader)
	if err != nil {
		return err
	}
	i.Entries = make([]*ImportEntry, 0, count)
	for n := uint32(0); n < count; n++ {
		entry := new(ImportEntry)
		if err = entry.Deserialize(reader); err != nil {
			return err
		}
		i.Entries = append(i.Entries, entry)
	}
	return nil
}

type TableSection struct {
	Entries []*Table
}

func (_ TableSection) IsSection() bool  { return true }
func (_ TableSection) Serialize() error { return nil }

func (ts *TableSection) Deserialize(reader *wasm_reader.WasmReader) error {
	count, err := wbinary.ReadVarUint32(reader)
	if err != nil {
		return err
	}

	ts.Entries = make([]*Table, 0, count)
	for i := uint32(0); i < count; i++ {
		entry := new(Table)
		if err = entry.Deserialize(reader); err != nil {
			return err
		}
		ts.Entries = append(ts.Entries, entry)
	}
	return nil
}

type FunctionSection struct {
	Indices []uint32
}

func (_ FunctionSection) IsSection() bool  { return true }
func (_ FunctionSection) Serialize() error { return nil }

func (f *FunctionSection) Deserialize(reader *wasm_reader.WasmReader) error {
	count, err := wbinary.ReadVarUint32(reader)
	if err != nil {
		return err
	}
	f.Indices = make([]uint32, 0, count)
	for i := uint32(0); i < count; i++ {
		index, e := wbinary.ReadVarUint32(reader)
		if e != nil {
			return e
		}
		f.Indices = append(f.Indices, index)
	}
	return nil
}

type MemorySection struct {
	Entries []*MemoryKindDesc
}

func (_ MemorySection) IsSection() bool  { return true }
func (_ MemorySection) Serialize() error { return nil }

func (m *MemorySection) Deserialize(reader *wasm_reader.WasmReader) error {
	count, err := wbinary.ReadVarUint32(reader)
	if err != nil {
		return err
	}

	m.Entries = make([]*MemoryKindDesc, 0, count)
	for i := uint32(0); i < count; i++ {
		mkd := new(MemoryKindDesc)
		if err = mkd.Deserialize(reader); err != nil {
			return err
		}
		m.Entries = append(m.Entries, mkd)
	}
	return nil
}

type GlobalDecl struct {
	// Description of a Global declaration (variable, perhaps)
	Description GlobalKindDesc
	// Init expression (instruction) to compute initial value of global decl
	Init []byte
}

func (_ GlobalDecl) Serialize() error { return nil }

func (g *GlobalDecl) Deserialize(reader *wasm_reader.WasmReader) error {
	return nil
}

type GlobalSection struct {
	Entries []*GlobalDecl
}

func (_ GlobalSection) IsSection() bool  { return true }
func (_ GlobalSection) Serialize() error { return nil }

func (g *GlobalSection) Deserialize(reader *wasm_reader.WasmReader) error {
	count, err := wbinary.ReadVarUint32(reader)
	if err != nil {
		return err
	}
	g.Entries = make([]*GlobalDecl, 0, count)

	for i := uint32(0); i < count; i++ {
		decl := new(GlobalDecl)
		if err = decl.Deserialize(reader); err != nil {
			return err
		}
		g.Entries = append(g.Entries, decl)
	}
	return nil
}

type ExportSection struct {
}

func (_ ExportSection) IsSection() bool  { return true }
func (_ ExportSection) Serialize() error { return nil }

func (e *ExportSection) Deserialize(reader *wasm_reader.WasmReader) error {
	return nil
}

type StartSection struct {
}

func (_ StartSection) IsSection() bool  { return true }
func (_ StartSection) Serialize() error { return nil }

func (s *StartSection) Deserialize(reader *wasm_reader.WasmReader) error {
	return nil
}

type ElementSection struct {
}

func (_ ElementSection) IsSection() bool  { return true }
func (_ ElementSection) Serialize() error { return nil }

func (e *ElementSection) Deserialize(reader *wasm_reader.WasmReader) error {
	return nil
}

type CodeSection struct {
}

func (_ CodeSection) IsSection() bool  { return true }
func (_ CodeSection) Serialize() error { return nil }

func (c *CodeSection) Deserialize(reader *wasm_reader.WasmReader) error {
	return nil
}

type DataSection struct {
}

func (_ DataSection) IsSection() bool  { return true }
func (_ DataSection) Serialize() error { return nil }

func (d *DataSection) Deserialize(reader *wasm_reader.WasmReader) error {
	return nil
}

type CustomSection struct {
}

type CustomSections []*CustomSection

func (_ CustomSection) IsSection() bool  { return true }
func (_ CustomSection) Serialize() error { return nil }

func (c *CustomSection) Deserialize(reader *wasm_reader.WasmReader) error {
	return nil
}
