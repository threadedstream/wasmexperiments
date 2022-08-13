package wasm

import (
	"github.com/threadedstream/wasmexperiments/internal/pkg/wasm_reader"
)

type SectionID uint

const (
	CustomSection SectionID = iota
	TypeSection
	ImportSection
	FunctionSection
	TableSection
	LinearMemorySection
	GlobalSection
	ExportSection
	StartSection
	ElementSection
	CodeSection
	DataSection
)

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
