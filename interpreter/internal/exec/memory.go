package exec

var (
	localGetOp = newVarargOp("local.get", 0x20)
)

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
