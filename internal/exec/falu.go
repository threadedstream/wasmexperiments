package exec

import "math"

func (vm *VM) f32Add() {
	rhs := float32(vm.popUint32())
	lhs := float32(vm.popUint32())
	vm.pushUint32(uint32(lhs + rhs))
}

func (vm *VM) f32Sub() {
	rhs := float32(vm.popUint32())
	lhs := float32(vm.popUint32())
	vm.pushUint32(uint32(lhs - rhs))
}

func (vm *VM) f32Mul() {
	rhs := float32(vm.popUint32())
	lhs := float32(vm.popUint32())
	vm.pushUint32(uint32(lhs * rhs))
}

func (vm *VM) f32Div() {
	rhs := float32(vm.popUint32())
	lhs := float32(vm.popUint32())
	vm.pushUint32(uint32(lhs / rhs))
}

func (vm *VM) f32Sqrt() {
	target := float32(vm.popUint32())
	vm.pushUint32(uint32(math.Sqrt(float64(target))))
}

func (vm *VM) f32Min() {
	rhs := float32(vm.popUint32())
	lhs := float32(vm.popUint32())
	vm.pushUint32(uint32(math.Min(float64(lhs), float64(rhs))))
}

func (vm *VM) f32Max() {
	rhs := float32(vm.popUint32())
	lhs := float32(vm.popUint32())
	vm.pushUint32(uint32(math.Max(float64(lhs), float64(rhs))))
}

func (vm *VM) f32Ceil() {
	target := float32(vm.popUint32())
	vm.pushUint32(uint32(math.Ceil(float64(target))))
}

func (vm *VM) f32Floor() {
	target := float32(vm.popUint32())
	vm.pushUint32(uint32(math.Floor(float64(target))))
}

func (vm *VM) f32Trunc() {
	target := float32(vm.popUint32())
	vm.pushUint32(uint32(math.Trunc(float64(target))))
}

func (vm *VM) f32Nearest() {
	target := float32(vm.popUint32())
	vm.pushUint32(uint32(math.Round(float64(target))))
}

func (vm *VM) f32Abs() {
	target := float32(vm.popUint32())
	vm.pushUint32(uint32(math.Abs(float64(target))))
}

func (vm *VM) f32Neg() {
	target := float32(vm.popUint32())
	if math.Signbit(float64(target)) {
		vm.pushUint32(uint32(math.Copysign(float64(target), -1)))
	} else {
		vm.pushUint32(uint32(math.Copysign(float64(target), 1)))
	}
}

func (vm *VM) f32Copysign() {
	sign := float32(vm.popUint32())
	target := float32(vm.popUint32())
	vm.pushUint32(uint32(math.Copysign(float64(target), float64(sign))))
}

func (vm *VM) f32Eq() {
	rhs := float32(vm.popUint32())
	lhs := float32(vm.popUint32())
	if lhs == rhs {
		vm.pushOne()
	} else {
		vm.pushZero()
	}
}

func (vm *VM) f32Ne() {
	rhs := float32(vm.popUint32())
	lhs := float32(vm.popUint32())
	if lhs != rhs {
		vm.pushOne()
	} else {
		vm.pushZero()
	}
}

func (vm *VM) f32Lt() {
	rhs := float32(vm.popUint32())
	lhs := float32(vm.popUint32())
	if lhs < rhs {
		vm.pushOne()
	} else {
		vm.pushZero()
	}
}

func (vm *VM) f32Le() {
	rhs := float32(vm.popUint32())
	lhs := float32(vm.popUint32())
	if lhs <= rhs {
		vm.pushOne()
	} else {
		vm.pushZero()
	}
}

func (vm *VM) f32Gt() {
	rhs := float32(vm.popUint32())
	lhs := float32(vm.popUint32())
	if lhs > rhs {
		vm.pushOne()
	} else {
		vm.pushZero()
	}
}

func (vm *VM) f32Ge() {
	rhs := float32(vm.popUint32())
	lhs := float32(vm.popUint32())
	if lhs >= rhs {
		vm.pushOne()
	} else {
		vm.pushZero()
	}
}

func (vm *VM) i32TruncF32S() {
	target := float32(vm.popUint32())
	vm.pushInt32(int32(target))
}

func (vm *VM) i32TruncF64S() {
	target := float64(vm.popInt64())
	vm.pushUint32(uint32(float32(target)))
}

func (vm *VM) i32TruncF32U() {
	target := float32(vm.popUint64())
	vm.pushUint32(uint32(target))
}

func (vm *VM) i32TruncF64U() {
	target := float64(vm.popUint64())
	vm.pushUint32(uint32(float32(target)))
}

func (vm *VM) f32DemoteF64() {
	target := float64(vm.popUint64())
	vm.pushUint32(uint32(float32(target)))
}

func (vm *VM) f64PromoteF32() {
	target := float32(vm.popUint32())
	vm.pushUint64(uint64(target))
}

func (vm *VM) f32ConvertI32S() {
	target := float32(vm.popInt32())
	vm.pushUint32(uint32(target))
}

func (vm *VM) f32ConvertI64S() {
	target := vm.popInt64()
	vm.pushUint32(uint32(float32(target)))
}

func (vm *VM) f32ConvertI32U() {
	target := float32(vm.popUint32())
	vm.pushUint32(uint32(target))
}

func (vm *VM) f32ConvertI64U() {
	target := vm.popUint64()
	vm.pushUint32(uint32(float32(target)))
}

func (vm *VM) f32ReinterpretI32() {
	target := vm.popUint32()
	reinterpreted := math.Float32frombits(target)
	vm.pushUint32(uint32(reinterpreted))
}
