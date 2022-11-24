package exec

import (
	"errors"
	"fmt"
	"github.com/threadedstream/wasmexperiments/internal/pkg/reporter"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wasm_reader"
)

const (
	wasmPageSize = 65536
)

type context struct {
	parent         *context // useful in blocks
	pendingContext *context
	stack          []uint64
	locals         []uint64
	compiledCode   []byte
	ins            []Instr
	pc             int64
	curFunc        int64
	isBlock        bool
	breakExecuted  bool
}

type VM struct {
	ctx      *context
	frames   []*context
	ctxchain []*context
	module   *Module
	globals  []uint64
	memory   []byte
	funcs    []Function
	reader   *wasm_reader.WasmReader

	// for quick querying
	funcMap      map[string]uint32
	returned     bool
	blockCounter uint32
}

func NewVM(m *Module) (*VM, error) {
	vm := new(VM)

	vm.initFuncTable()
	if m.MemorySection != nil && len(m.MemorySection.Entries) != 0 {
		if len(m.MemorySection.Entries) > 1 {
			return nil, errors.New("newVM: expected to have exactly one instance of memory")
		}
		vm.memory = make([]byte, m.MemorySection.Entries[0].Limits.Minimum*wasmPageSize)
		fmt.Printf("##NewVM## Addr: %p, Len: %d\n", m.LinearMemoryIndexSpace, len(m.LinearMemoryIndexSpace))
		copy(vm.memory, m.LinearMemoryIndexSpace[0])
	}

	if m.FunctionIndexSpace == nil {
		m.initializeFunctionIndexSpace()
	}

	vm.globals = make([]uint64, len(m.GlobalIndexSpace))
	vm.module = m

	if m.ExportSection != nil {
		vm.funcMap = make(map[string]uint32)
		for _, entry := range m.ExportSection.Entries {
			vm.funcMap[entry.Name] = entry.Index
		}
	}

	if m.StartSection != nil {
		_, err := vm.ExecFunc(int64(m.StartSection.Index))
		if err != nil {
			return nil, err
		}
	}

	vm.frames = make([]*context, 0, maxStackFrameNum)

	return vm, nil
}

func (vm *VM) ExecFunc(index int64, args ...uint64) (ret any, err error) {
	//// some validation of input parameters
	//if int(index) > len(vm.funcs) {
	//	return nil, fmt.Errorf("attempting to call a function with an index %d with length of funcs being %d", index, len(vm.funcs))
	//}

	// TODO(threadedstream): resolves to nil should the function of the following form be called
	// (func $fac (export "fac") (param $x i32) (result i32)
	fn := vm.module.GetFunction(int(index))

	return fn.call(vm, index, ExecutionModeRawBytecode, args...)
}

func (vm *VM) execCode() any {
	for int(vm.ctx.pc) < len(vm.ctx.ins) {
		if vm.returned || (vm.ctx.breakExecuted && vm.ctx.isBlock) {
			break
		}

		currCode := vm.ctx.ins[vm.ctx.pc].Op().Code
		if currCode == endOp {
			vm.ctx.pc++
			break
		}
		if handler, ok := funcTable[currCode]; ok {
			handler()
			if vm.ctx.breakExecuted {
				vm.ctx = vm.ctx.pendingContext
				if vm.ctx.isBlock {
					vm.ctx.breakExecuted = true
					vm.ctx = vm.ctx.parent
					goto end
				}
				vm.ctx.breakExecuted = false
				vm.ctx.pc = 0
			}
			continue
		}
		reporter.ReportError("execCode: unknown instruction with code %v\n", currCode)
	}

end:
	return nil
}
