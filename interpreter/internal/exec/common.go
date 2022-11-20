package exec

import (
	"encoding/binary"
	"github.com/threadedstream/wasmexperiments/internal/pkg/reporter"
	"math"
)

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

// pushZero is a pseudo-instruction, it has a practical utility in cmp instruction
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
	val := vm.ctx.stack[idx]
	vm.ctx.stack = vm.ctx.stack[:idx]
	return val
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

// the same as popUint64, but doesn't pop value off the stack
func (vm *VM) peekUint64() uint64 {
	if len(vm.ctx.stack) == 0 {
		reporter.ReportError("popUint64: stack's empty")
	}
	idx := len(vm.ctx.stack) - 1
	val := vm.ctx.stack[idx]
	return val
}

func (vm *VM) peekInt64() int64 {
	return int64(vm.peekUint64())
}

func (vm *VM) peekUint32() uint32 {
	return uint32(vm.peekUint64())
}

func (vm *VM) peekInt32() int32 {
	return int32(vm.peekUint64())
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

func (vm *VM) fetchFloat32() float32 {
	val := binary.LittleEndian.Uint32(vm.ctx.code[vm.ctx.pc:])
	vm.ctx.pc += 4
	return math.Float32frombits(val)
}

func (vm *VM) fetchFloat64() float64 {
	val := binary.LittleEndian.Uint64(vm.ctx.code[vm.ctx.pc:])
	vm.ctx.pc += 8
	return math.Float64frombits(val)
}

func (vm *VM) pushFloat32(v float32) {
	vm.pushUint32(uint32(v))
}

func (vm *VM) pushFloat64(v float64) {
	vm.pushUint64(uint64(v))
}

func (vm *VM) popFloat32() float32 {
	return math.Float32frombits(vm.popUint32())
}
