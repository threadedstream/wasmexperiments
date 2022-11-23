package exec

func (vm *VM) i32Const() {
	in := vm.currIns().(*I32ConstI)
	var val uint32
	if v, ok := in.arg0.(uint64); ok {
		val = uint32(v)
	}
	val = in.arg0.(uint32)
	vm.pushUint32(val)
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
