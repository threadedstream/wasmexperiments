package exec

import (
	"bytes"
	"errors"

	"github.com/threadedstream/wasmexperiments/internal/pkg/wasm_reader"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wbinary"
	"github.com/threadedstream/wasmexperiments/internal/types"
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

type Validatable interface {
	Validate() error
}

type FunctionSig struct {
	Params  []types.ValueType
	Results [1]types.ValueType // has at most one element (before rolling support for multiple return types)
}

func (fs FunctionSig) Serialize() error { return nil }

func (fs *FunctionSig) Deserialize(reader *wasm_reader.WasmReader) error {
	// force value type to be a func
	valType, err := wbinary.ReadU8(reader)
	if err != nil {
		return err
	}
	if valType != types.ValueTypeFunc {
		return errors.New("readTypeSection: value type must be a function")
	}

	// start filling out function signatures slice
	paramsLen, err := wbinary.ReadVarUint32(reader)
	if err != nil {
		return err
	}

	fs.Params = make([]types.ValueType, paramsLen, paramsLen)
	for i := 0; i < int(paramsLen); i++ {
		valTyp, err := wbinary.ReadVarUint32(reader)
		if err != nil {
			return err
		}
		fs.Params[i] = types.ValueType(valTyp)
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
		fs.Results[i] = types.ValueType(valTyp)
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
	Maximum uint32
}

func (_ ResizableLimits) Serialize() error { panic("unimplemented!") }

func (rl *ResizableLimits) Deserialize(reader *wasm_reader.WasmReader) error {
	var err error
	if rl.Flags, err = wbinary.ReadVarUint32(reader); err != nil {
		return err
	}
	if rl.Minimum, err = wbinary.ReadVarUint32(reader); err != nil {
		return err
	}

	// see if 0x1 bit is set
	if rl.Flags&0x1 != 0 {
		if rl.Maximum, err = wbinary.ReadVarUint32(reader); err != nil {
			return err
		}
	}
	return nil
}

type MemoryKindDesc struct {
	Limits ResizableLimits
}

func (_ MemoryKindDesc) Kind() ExternalKind { return MemoryKind }
func (_ MemoryKindDesc) IsImportDesc() bool { return true }
func (_ MemoryKindDesc) Serialize() error   { return nil }

func (m *MemoryKindDesc) Deserialize(reader *wasm_reader.WasmReader) error {
	return m.Limits.Deserialize(reader)
}

type GlobalKindDesc struct {
	Type    types.ValueType
	Mutable bool
}

func (_ GlobalKindDesc) Kind() ExternalKind { return GlobalKind }
func (_ GlobalKindDesc) IsImportDesc() bool { return true }
func (_ GlobalKindDesc) Serialize() error   { return nil }

func (g *GlobalKindDesc) Deserialize(reader *wasm_reader.WasmReader) error {
	if err := g.Type.Deserialize(reader); err != nil {
		return err
	}

	mut, err := reader.ReadByte()
	if err != nil {
		return err
	}

	if mut != 0x0 && mut != 0x1 {
		return errors.New("section: expected Mutable to be 0x0 or 0x1")
	}

	g.Mutable = mut == 0x1
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
	var err error
	if err = g.Description.Deserialize(reader); err != nil {
		return err
	}

	if g.Init, err = readInitExpr(reader); err != nil {
		return err
	}

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

type ExportEntry struct {
	Name  string
	Kind  ExternalKind
	Index uint32
}

func (e ExportEntry) Serialize() error { return nil }

func (e *ExportEntry) Deserialize(reader *wasm_reader.WasmReader) error {
	var err error
	if e.Name, err = wbinary.ReadUTF8StringUint(reader); err != nil {
		return err
	}

	if err = e.Kind.Deserialize(reader); err != nil {
		return err
	}

	if e.Index, err = wbinary.ReadVarUint32(reader); err != nil {
		return err
	}
	return nil
}

type ExportSection struct {
	Entries map[string]*ExportEntry
}

var ErrDuplicateExport = errors.New("section: duplicate exports not allowed")

func (_ ExportSection) IsSection() bool  { return true }
func (_ ExportSection) Serialize() error { return nil }
func (_ ExportSection) Validate() error  { return nil }

func (e *ExportSection) Deserialize(reader *wasm_reader.WasmReader) error {
	count, err := wbinary.ReadVarUint32(reader)
	if err != nil {
		return err
	}
	e.Entries = make(map[string]*ExportEntry, count)

	for i := uint32(0); i < count; i++ {
		entry := new(ExportEntry)
		if err = entry.Deserialize(reader); err != nil {
			return err
		}
		if _, exists := e.Entries[entry.Name]; exists {
			return ErrDuplicateExport
		}
		e.Entries[entry.Name] = entry
	}
	return nil
}

type StartSection struct {
	Index uint32
}

func (_ StartSection) IsSection() bool  { return true }
func (_ StartSection) Serialize() error { return nil }

func (s *StartSection) Deserialize(reader *wasm_reader.WasmReader) error {
	var err error
	if s.Index, err = wbinary.ReadVarUint32(reader); err != nil {
		return err
	}
	return nil
}

type TableInitializer struct {
	Index  uint32
	Offset []byte   // must return i32
	Elems  []uint32 // in case if table's element_type is funcref
}

func (_ TableInitializer) Serializer() error { return nil }

func (t *TableInitializer) Deserialize(reader *wasm_reader.WasmReader) error {
	var err error
	if t.Index, err = wbinary.ReadVarUint32(reader); err != nil {
		return err
	}
	if t.Offset, err = readInitExpr(reader); err != nil {
		return err
	}

	elemsNum, err := wbinary.ReadVarUint32(reader)
	if err != nil {
		return err
	}
	t.Elems = make([]uint32, elemsNum, elemsNum)
	for i := uint32(0); i < elemsNum; i++ {
		elem, err := wbinary.ReadVarUint32(reader)
		if err != nil {
			return err
		}
		t.Elems[i] = elem
	}
	return nil
}

type ElementSection struct {
	Entries []*TableInitializer
}

func (_ ElementSection) IsSection() bool  { return true }
func (_ ElementSection) Serialize() error { return nil }

// TODO(threadedstream): has not yet been tested
func (e *ElementSection) Deserialize(reader *wasm_reader.WasmReader) error {
	count, err := wbinary.ReadVarUint32(reader)
	if err != nil {
		return err
	}
	e.Entries = make([]*TableInitializer, count, count)

	for i := uint32(0); i < count; i++ {
		tableInit := new(TableInitializer)
		if err = tableInit.Deserialize(reader); err != nil {
			return err
		}
		e.Entries[i] = tableInit
	}
	return nil
}

type LocalEntry struct {
	Count uint32
	Type  types.ValueType
}

func (_ LocalEntry) Serialize() error { return nil }

func (l *LocalEntry) Deserialize(reader *wasm_reader.WasmReader) error {
	var err error
	if l.Count, err = wbinary.ReadVarUint32(reader); err != nil {
		return err
	}
	if err = l.Type.Deserialize(reader); err != nil {
		return err
	}
	return nil
}

type FunctionBody struct {
	Size   uint32
	Locals []*LocalEntry
	Code   []byte
}

var ErrFunctionNoEnd = errors.New("section: missing 'end' instruction at the end of function body")

func (fb FunctionBody) Serialize() error { return nil }

func (fb *FunctionBody) Deserialize(reader *wasm_reader.WasmReader) error {
	var err error
	if fb.Size, err = wbinary.ReadVarUint32(reader); err != nil {
		return err
	}

	var body []byte
	if body, err = reader.ReadBytes(int(fb.Size)); err != nil {
		return err
	}

	bodyReader := bytes.NewBuffer(body)
	reader.Push(bodyReader)
	defer reader.Pop()

	// read number of locals
	var localCount uint32
	if localCount, err = wbinary.ReadVarUint32(reader); err != nil {
		return err
	}
	fb.Locals = make([]*LocalEntry, 0, localCount)
	for i := uint32(0); i < localCount; i++ {
		local := new(LocalEntry)
		if err = local.Deserialize(reader); err != nil {
			return err
		}
		fb.Locals = append(fb.Locals, local)
	}

	code := bodyReader.Bytes()
	if code[len(code)-1] != end {
		return ErrFunctionNoEnd
	}

	fb.Code = code[:len(code)-1]

	return nil
}

type CodeSection struct {
	Entries []*FunctionBody
}

func (_ CodeSection) IsSection() bool  { return true }
func (_ CodeSection) Serialize() error { return nil }

func (c *CodeSection) Deserialize(reader *wasm_reader.WasmReader) error {
	count, err := wbinary.ReadVarUint32(reader)
	if err != nil {
		return nil
	}
	c.Entries = make([]*FunctionBody, 0, count)

	for i := uint32(0); i < count; i++ {
		functionBody := new(FunctionBody)
		if err = functionBody.Deserialize(reader); err != nil {
			return err
		}
		c.Entries = append(c.Entries, functionBody)
	}
	return nil
}

type DataInitializer struct {
	Index  uint32
	Offset []byte
	Data   []byte
}

func (_ DataInitializer) Serialize() error { return nil }

func (d *DataInitializer) Deserialize(reader *wasm_reader.WasmReader) error {
	var err error
	if d.Index, err = wbinary.ReadVarUint32(reader); err != nil {
		return err
	}
	if d.Offset, err = readInitExpr(reader); err != nil {
		return err
	}
	if d.Data, err = wbinary.ReadByteArray(reader); err != nil {
		return err
	}
	return nil
}

type DataSection struct {
	Entries []*DataInitializer
}

func (_ DataSection) IsSection() bool  { return true }
func (_ DataSection) Serialize() error { return nil }

func (d *DataSection) Deserialize(reader *wasm_reader.WasmReader) error {
	count, err := wbinary.ReadVarUint32(reader)
	if err != nil {
		return err
	}
	d.Entries = make([]*DataInitializer, 0, count)

	for i := uint32(0); i < count; i++ {
		dataInit := new(DataInitializer)
		if err = dataInit.Deserialize(reader); err != nil {
			return err
		}
		d.Entries = append(d.Entries, dataInit)
	}
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
