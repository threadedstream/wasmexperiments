package wasm

import (
	"bytes"
	"errors"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wasm_reader"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wbinary"
	"io"
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
	Deserialize(wr *wasm_reader.WasmReader) error
	Serialize(wr *wasm_reader.WasmReader) error
}

type Section interface {
	IsSection() bool
}

type FunctionSig struct {
	params  []ValueType
	results [1]ValueType // has at most one element (before rolling support for multiple return types)
}

type TypesSection struct {
	sigs   []*FunctionSig
	reader *wasm_reader.WasmReader
}

func (ts TypesSection) IsSection() bool { return true }

func (ts *TypesSection) Serialize() error { return nil }

func (ts *TypesSection) Deserialize() error {
	dataLen, err := wbinary.ReadVarUint32(ts.reader)
	if err != nil {
		return err
	}

	sectionData := new(bytes.Buffer)
	sectionData.Grow(int(dataLen))
	sectionReader := io.LimitReader(io.TeeReader(ts.reader.Peek().(io.Reader), sectionData), int64(dataLen))
	ts.reader.Push(sectionReader)
	_, err = ts.reader.ReadBytes(int(dataLen))
	if err != nil {
		return err
	}
	ts.reader.Pop()
	ts.reader.Push(sectionData)

	// read arr length
	arrLen, err := wbinary.ReadVarUint32(ts.reader)
	if err != nil {
		return err
	}
	ts.sigs = make([]*FunctionSig, arrLen, arrLen)

	for i, sig := range ts.sigs {
		sig = new(FunctionSig)
		// force value type to be a func
		valType, err := wbinary.ReadU8(ts.reader)
		if err != nil {
			return err
		}
		if valType != ValueTypeFunc {
			return errors.New("readTypeSection: value type must be a function")
		}

		// start filling out function signatures slice
		paramsLen, err := wbinary.ReadVarUint32(ts.reader)
		if err != nil {
			return err
		}

		sig.params = make([]ValueType, paramsLen, paramsLen)
		for i := 0; i < int(paramsLen); i++ {
			valTyp, err := wbinary.ReadVarUint32(ts.reader)
			if err != nil {
				return err
			}
			sig.params[i] = ValueType(valTyp)
		}
		resultsLen, err := wbinary.ReadVarUint32(ts.reader)
		if err != nil {
			return err
		}
		if resultsLen > 1 {
			return errors.New("readTypeSection: length of results array can't exceed the value of 1 (yet)")
		}
		for i := 0; i < int(resultsLen); i++ {
			valTyp, err := wbinary.ReadVarUint32(ts.reader)
			if err != nil {
				return err
			}
			sig.results[i] = ValueType(valTyp)
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

func (f FunctionKindDesc) Kind() ExternalKind { return FunctionKind }

func (f FunctionKindDesc) IsImport() bool { return true }

type ElementType int

const FuncRefElementType ElementType = 0x70

type Table struct {
	ElemType ElementType
	Limits   ResizableLimits
}

type TableKindDesc struct {
	Table Table
}

type ResizableLimits struct {
	Flags   uint32
	Minimum uint32
}

func NewTableKindDesc(table Table) *TableKindDesc {
	desc := new(TableKindDesc)
	desc.Table = table
	return desc
}

func (t TableKindDesc) Kind() ExternalKind { return TableKind }

func (t TableKindDesc) IsImportDesc() bool { return true }

type MemoryKindDesc struct {
	Limits ResizableLimits
}

func NewMemoryKindDesc(limits ResizableLimits) *MemoryKindDesc {
	desc := new(MemoryKindDesc)
	desc.Limits = limits
	return desc
}

func (m MemoryKindDesc) Kind() ExternalKind { return MemoryKind }

func (m MemoryKindDesc) IsImportDesc() bool { return true }

type GlobalKindDesc struct {
	Type    ValueType
	Mutable bool
}

func NewGlobalKindDesc(typ ValueType, mutable bool) *GlobalKindDesc {
	desc := new(GlobalKindDesc)
	desc.Type = typ
	desc.Mutable = mutable
	return desc
}

func (g GlobalKindDesc) Kind() ExternalKind {
	return GlobalKind
}

func (g GlobalKindDesc) IsImportDesc() bool {
	return true
}

type ImportSection struct {
	ModuleName   string
	ExportName   string
	Descriptions []ImportDesc
}
