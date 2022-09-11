package exec

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/threadedstream/wasmexperiments/internal/pkg/reporter"
	"github.com/threadedstream/wasmexperiments/internal/wasm"
)

const (
	wasmPageSize = 65536
)

type context struct {
	stack   []uint64
	locals  []uint64
	code    []byte
	pc      int64
	curFunc int64
}

type VM struct {
	ctx       context
	module    *wasm.Module
	globals   []uint64
	memory    []byte
	funcs     []Function
	funcTable [256]func()

	abort bool
}

func NewVM(m *wasm.Module) (*VM, error) {
	vm := new(VM)

	if m.MemorySection != nil && len(m.MemorySection.Entries) != 0 {
		if len(m.MemorySection.Entries) > 1 {
			return nil, errors.New("newVM: expected to have exactly one instance of memory")
		}
		vm.memory = make([]byte, m.MemorySection.Entries[0].Limits.Minimum*wasmPageSize)
		fmt.Printf("##NewVM## Addr: %p, Len: %d\n", m.LinearMemoryIndexSpace, len(m.LinearMemoryIndexSpace))
		copy(vm.memory, m.LinearMemoryIndexSpace[0])
	}
	vm.globals = make([]uint64, len(m.GlobalIndexSpace))

	vm.module = m

	// initialize function index space
	fnIndexSpace := make([]*Function, 0, len(vm.module.FunctionSection.Indices))
	vm.module.FunctionIndexSpace = fnIndexSpace

	if m.StartSection != nil {
		_, err := vm.ExecFunc(int64(m.StartSection.Index))
		if err != nil {
			return nil, err
		}
	}

	return vm, nil
}

func (vm *VM) pushUint64(n uint64) {
	if len(vm.ctx.stack) >= cap(vm.ctx.stack) {
		reporter.ReportError("stack exceeding max depth: len=%d,cap=%d", len(vm.ctx.stack), cap(vm.ctx.stack))
	}
	vm.ctx.stack = append(vm.ctx.stack, n)
}

func (vm *VM) pushInt64(n int64) {
	vm.pushUint64(uint64(n))
}

func (vm *VM) pushUint32(n uint32) {
	vm.pushUint64(uint64(n))
}

func (vm *VM) pushInt32(n int32) {
	vm.pushUint64(uint64(n))
}

// pushZero is a pseudo-instruction, it has a practical utility in cmp instructions
func (vm *VM) pushZero() {
	vm.pushUint64(0)
}

// the same as pushZero
func (vm *VM) pushOne() {
	vm.pushUint64(1)
}

func (vm *VM) popUint64() uint64 {
	if len(vm.ctx.stack) == 0 {
		reporter.ReportError("popUint64: stack's empty")
	}
	idx := len(vm.ctx.stack) - 1
	return vm.ctx.stack[idx]
}

func (vm *VM) popInt64() int64 {
	return int64(vm.popUint64())
}

func (vm *VM) popUint32() uint32 {
	return uint32(vm.popUint64())
}

func (vm *VM) popInt32() int32 {
	return int32(vm.popUint64())
}

func (vm *VM) fetchUint64() uint64 {
	val := binary.LittleEndian.Uint64(vm.ctx.code[vm.ctx.pc:])
	vm.ctx.pc += 8
	return val
}

func (vm *VM) fetchInt64() int64 {
	return int64(vm.fetchUint64())
}

func (vm *VM) fetchUint32() uint32 {
	val := binary.LittleEndian.Uint32(vm.ctx.code[vm.ctx.pc:])
	vm.ctx.pc += 4
	return val
}

func (vm *VM) fetchInt32() int32 {
	return int32(vm.fetchUint32())
}

func (vm *VM) getLocal() {
	index := vm.fetchUint32()
	vm.pushUint64(vm.ctx.locals[index])
}

func (vm *VM) setLocal() {
	index := vm.fetchUint32()
	value := vm.popUint64()
	vm.ctx.locals[index] = value
}

func (vm *VM) getGlobal() {
	index := vm.fetchUint32()
	vm.pushUint64(vm.globals[index])
}

func (vm *VM) setGlobal() {
	index := vm.fetchUint32()
	value := vm.popUint64()
	vm.globals[index] = value
}

func (vm *VM) PrintInstructionStream() (string, error) {
	for _, _ = range vm.module.CodeSection.Entries {
		// TODO
	}
	return "", nil
}

func (vm *VM) ExecFunc(index int64, args ...uint64) (ret any, err error) {
	// some validation of input parameters
	if int(index) > len(vm.funcs) {
		return nil, fmt.Errorf("attempting to call a function with an index %d with length of funcs being %d", index, len(vm.funcs))
	}

	// validate number of arguments
	fn := vm.module.GetFunction(int(index))

	// assuming it's already compiled, it's true though, we don't parse any frontend

	fn.call(vm, index)

	return nil, nil
}

func (vm *VM) execCode() any {
	for int(vm.ctx.pc) < len(vm.ctx.code) && !vm.abort {
		currPc := vm.ctx.pc
		op := vm.ctx.code[currPc]
		vm.ctx.pc++
		if inst, ok := instructionMap[op]; ok {
			println("known instruction: ", inst)
		}
	}

	return nil
}
