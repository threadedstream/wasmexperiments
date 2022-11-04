package exec

import (
	"github.com/threadedstream/wasmexperiments/internal/pkg/reporter"
	"github.com/threadedstream/wasmexperiments/internal/types"
	"math"
	"math/bits"
)

var (
	i32AddOp = newOp("i32.add", 0x6A, []types.ValueType{types.ValueTypeI32, types.ValueTypeI32}, []types.ValueType{types.ValueTypeI32})
	callOp   = newOp("call", 0x10, []types.ValueType{types.ValueTypeI32}, []types.ValueType{})
	endOp    = newOp("end", 0x0B, types.ValueTypeEmpty, types.ValueTypeEmpty)
)

func (vm *VM) i32Clz() {
	vm.pushUint64(uint64(bits.LeadingZeros32(vm.popUint32())))
}

func (vm *VM) i32Ctz() {
	vm.pushUint64(uint64(bits.TrailingZeros32(vm.popUint32())))
}

func (vm *VM) i32Popcnt() {
	vm.pushUint64(uint64(bits.OnesCount32(vm.popUint32())))
}

func (vm *VM) i32Add() {
	vm.pushUint32(vm.popUint32() + vm.popUint32())
}

func (vm *VM) i32Mul() {
	// TODO(threadedstream): add overflow checks?
	vm.pushUint32(vm.popUint32() * vm.popUint32())
}

func (vm *VM) i32DivS() {
	rhs := vm.popInt32()
	lhs := vm.popInt32()
	if lhs == math.MinInt32 && rhs == -1 {
		reporter.ReportError("detected integer overflow")
	}
	vm.pushInt32(lhs / rhs)
}

func (vm *VM) i32Sub() {
	rhs := vm.popUint32()
	lhs := vm.popUint32()
	vm.pushUint32(lhs - rhs)
}

func (vm *VM) i32DivU() {
	rhs := vm.popUint32()
	lhs := vm.popUint32()
	vm.pushUint32(lhs / rhs)
}

func (vm *VM) i32RemS() {
	rhs := vm.popInt32()
	lhs := vm.popInt32()
	vm.pushInt32(lhs % rhs)
}

func (vm *VM) i32RemU() {
	rhs := vm.popUint32()
	lhs := vm.popUint32()
	vm.pushUint32(lhs % rhs)
}

func (vm *VM) i32And() {
	vm.pushUint32(vm.popUint32() & vm.popUint32())
}

func (vm *VM) i32Or() {
	vm.pushUint32(vm.popUint32() | vm.popUint32())
}

func (vm *VM) i32Xor() {
	vm.pushUint32(vm.popUint32() ^ vm.popUint32())
}

func (vm *VM) i32Shl() {
	shift := vm.popUint32()
	target := vm.popUint32()
	vm.pushUint32(target << (shift % 32))
}

func (vm *VM) i32Shr() {
	shift := vm.popUint32()
	target := vm.popUint32()
	vm.pushUint32(target >> (shift % 32))
}

func (vm *VM) i32ShrS() {
	shift := vm.popUint32()
	target := vm.popInt32()
	vm.pushInt32(target >> (shift % 32))
}

func (vm *VM) i32RotL() {
	factor := vm.popUint32()
	target := vm.popUint32()
	vm.pushUint32(bits.RotateLeft32(target, int(factor)))
}

func (vm *VM) i32RotR() {
	factor := vm.popUint32()
	target := vm.popUint32()
	vm.pushUint32(bits.RotateLeft32(target, int(-factor)))
}

func (vm *VM) i32Eqz() {
	target := vm.popUint32()
	if target == 0 {
		vm.pushUint32(1)
	} else {
		vm.pushUint32(0)
	}
}

func (vm *VM) i32Eq() {
	rhs := vm.popUint32()
	lhs := vm.popUint32()
	if lhs == rhs {
		vm.pushOne()
	} else {
		vm.pushZero()
	}
}

func (vm *VM) i32Ne() {
	rhs := vm.popUint32()
	lhs := vm.popUint32()
	if lhs != rhs {
		vm.pushOne()
	} else {
		vm.pushZero()
	}
}

func (vm *VM) i32LtS() {
	rhs := vm.popInt32()
	lhs := vm.popInt32()
	if lhs < rhs {
		vm.pushOne()
	} else {
		vm.pushZero()
	}
}

func (vm *VM) i32LtU() {
	rhs := vm.popUint32()
	lhs := vm.popUint32()
	if lhs < rhs {
		vm.pushOne()
	} else {
		vm.pushZero()
	}
}

func (vm *VM) i32LeS() {
	rhs := vm.popInt32()
	lhs := vm.popInt32()
	if lhs <= rhs {
		vm.pushOne()
	} else {
		vm.pushZero()
	}
}

func (vm *VM) i32LeU() {
	rhs := vm.popUint32()
	lhs := vm.popUint32()
	if lhs <= rhs {
		vm.pushOne()
	} else {
		vm.pushZero()
	}
}

func (vm *VM) i32GtS() {
	rhs := vm.popInt32()
	lhs := vm.popInt32()
	if lhs > rhs {
		vm.pushOne()
	} else {
		vm.pushZero()
	}
}

func (vm *VM) i32GtU() {
	rhs := vm.popUint32()
	lhs := vm.popUint32()
	if lhs > rhs {
		vm.pushOne()
	} else {
		vm.pushZero()
	}
}

func (vm *VM) i32GeS() {
	rhs := vm.popInt32()
	lhs := vm.popInt32()
	if lhs >= rhs {
		vm.pushOne()
	} else {
		vm.pushZero()
	}
}

func (vm *VM) i32GeU() {
	rhs := vm.popUint32()
	lhs := vm.popUint32()
	if lhs >= rhs {
		vm.pushOne()
	} else {
		vm.pushZero()
	}
}
