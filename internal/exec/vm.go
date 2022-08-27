package exec

import (
	"encoding/binary"
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
	memory    []byte
	funcs     []function
	funcTable [256]func()
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

func (vm *VM) PrintInstructionStream() (string, error) {
	for _, _ = range vm.module.CodeSection.Entries {
		// TODO
	}
	return "", nil
}

func (vm *VM) ExecFunc(index int64, args uint64) (ret any, err error) {
	// some validation of input parameters
	if int(index) > len(vm.funcs) {
		return nil, fmt.Errorf("attempting to call a function with an index %d with length of funcs being %d", index, len(vm.funcs))
	}

	return nil, nil
}

func (vm *VM) execCode() any {
	return nil
}
