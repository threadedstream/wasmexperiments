package exec

import (
	"errors"
	"fmt"
)

const (
	wasmPageSize = 65536
)

type context struct {
	stack   []uint64
	locals  []uint64
	raw     []byte
	ins     []Instr
	pc      int64
	curFunc int64
}

type VM struct {
	ctx       context
	module    *Module
	globals   []uint64
	memory    []byte
	funcs     []Function
	funcTable map[Bytecode]func()
	// for quick querying
	funcMap map[string]uint32
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

	return vm, nil
}

func (vm *VM) initFuncTable() {
	if vm.funcTable == nil {
		vm.funcTable = map[Bytecode]func(){
			ifOp:        vm.execIf,
			i32EqOp:     vm.i32Eq,
			i32AddOp:    vm.i32Add,
			i32SubOp:    vm.i32Sub,
			i32MulOp:    vm.i32Mul,
			i32DivSOp:   vm.i32DivS,
			i32DivUOp:   vm.i32DivU,
			i32RemSOp:   vm.i32RemS,
			i32RemUOp:   vm.i32RemU,
			f32AddOp:    vm.f32Add,
			f32SubOp:    vm.f32Sub,
			f32MulOp:    vm.f32Mul,
			callOp:      vm.call,
			localGetOp:  vm.getLocal,
			localSetOp:  vm.setLocal,
			globalGetOp: vm.getGlobal,
			globalSetOp: vm.setGlobal,
			i32ConstOp:  vm.i32Const,
			i32LoadOp:   vm.i32Load,
			f32LoadOp:   vm.f32Load,
			i32StoreOp:  vm.i32Store,
			f32StoreOp:  vm.f32Store,
		}
	}
}

func (vm *VM) ExecFunc(index int64, args ...uint64) (ret any, err error) {
	//// some validation of input parameters
	//if int(index) > len(vm.funcs) {
	//	return nil, fmt.Errorf("attempting to call a function with an index %d with length of funcs being %d", index, len(vm.funcs))
	//}

	// TODO(threadedstream): resolves to nil should the function of the following form be called
	// (func $fac (export "fac") (param $x i32) (result i32)
	fn := vm.module.GetFunction(int(index))

	return fn.call(vm, index, args...)
}

func (vm *VM) execCode() any {
	startIndex := vm.module.StartSection.Index
	fn := vm.module.GetFunction(int(startIndex))

	fn.call(vm, int64(startIndex))
	return nil
}
