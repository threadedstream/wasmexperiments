package exec

func (vm *VM) call() {
	index := vm.popInt32()

	fn := vm.module.FunctionIndexSpace[index]
	fn.call(vm, int64(index))
}
