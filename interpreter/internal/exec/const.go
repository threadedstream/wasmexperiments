package exec

func (vm *VM) i32Const() {
	in := vm.currIns().(*I32ConstI)
	vm.pushUint32(in.arg0.(uint32))
	vm.ctx.pc++
}

func (vm *VM) f32Const() {
	vm.pushFloat32(vm.fetchFloat32())
}

func (vm *VM) i64Const() {
	vm.pushUint64(vm.fetchUint64())
}

func (vm *VM) f64Const() {
	vm.pushFloat64(vm.fetchFloat64())
}
