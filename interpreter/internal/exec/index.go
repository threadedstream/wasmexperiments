package exec

import (
	"fmt"
	"github.com/threadedstream/wasmexperiments/internal/pkg/reporter"
	"reflect"
)

// Mostly copied from wagon

type OutsizeError struct {
	ImmType string
	Size    uint64
	Max     uint64
}

func (e OutsizeError) Error() string {
	return fmt.Sprintf("validate: %s size overflow (%v), max (%v)", e.ImmType, e.Size, e.Max)
}

type InvalidTableIndexError uint32

func (e InvalidTableIndexError) Error() string {
	return fmt.Sprintf("wasm: Invalid table to table index space: %d", uint32(e))
}

type UninitializedTableEntryError uint32

func (e UninitializedTableEntryError) Error() string {
	return fmt.Sprintf("wasm: Uninitialized table entry at index: %d", uint32(e))
}

type InvalidValueTypeInitExprError struct {
	Wanted reflect.Kind
	Got    reflect.Kind
}

func (e InvalidValueTypeInitExprError) Error() string {
	return fmt.Sprintf("wasm: Wanted initializer expression to return %v value, got %v", e.Wanted, e.Got)
}

type InvalidLinearMemoryIndexError uint32

func (e InvalidLinearMemoryIndexError) Error() string {
	return fmt.Sprintf("wasm: Invalid linear memory index: %d", uint32(e))
}

func (m *Module) populateLinearMemory() error {
	if m.DataSection == nil || len(m.DataSection.Entries) == 0 {
		return nil
	}

	for _, entry := range m.DataSection.Entries {
		if entry.Index != 0 {
			return InvalidLinearMemoryIndexError(entry.Index)
		}

		val, err := m.execInitExpr(entry.Offset)
		if err != nil {
			return err
		}
		off, ok := val.(int32)
		if !ok {
			return InvalidValueTypeInitExprError{reflect.Int32, reflect.TypeOf(val).Kind()}
		}
		offset := uint32(off)

		memory := m.LinearMemoryIndexSpace[entry.Index]
		if uint64(offset)+uint64(len(entry.Data)) > uint64(len(memory)) {
			data := make([]byte, uint64(offset)+uint64(len(entry.Data)))
			copy(data, memory)
			copy(data[offset:], entry.Data)
			m.LinearMemoryIndexSpace[int(entry.Index)] = data
		} else {
			copy(memory[offset:], entry.Data)
		}
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
				code:      m.CodeSection.Entries[i].Code,
				name:      m.ImportSection.Entries[i].ExportName,
				numParams: len(m.TypesSection.sigs[i].Params),
				returns:   m.TypesSection.sigs[i].Results[0] != 0,
			}
		}
	}

	if m.FunctionSection != nil {
		for _, idx := range m.FunctionSection.Indices {
			m.FunctionIndexSpace[idx] = &Function{
				code:      m.CodeSection.Entries[idx].Code,
				name:      "",
				numParams: len(m.TypesSection.sigs[idx].Params),
				returns:   m.TypesSection.sigs[idx].Results[0] != 0,
			}
		}
	}
}

func (m *Module) GetFunction(i int) *Function {
	if i >= len(m.FunctionIndexSpace) || i < 0 {
		reporter.ReportError("module.GetFunction: attempting to use index out of bounds %d", i)
	}

	return m.FunctionIndexSpace[i]
}
