package exec

type function interface {
	call(vm *VM, index int64)
}

type compiledFunction struct {
	code           []byte
	maxDepth       int
	totalLocalVars int
	args           int
	returns        bool
}

func (cf *compiledFunction) call(vm *VM, index int64) {
	stack := make([]uint64, 0, cf.maxDepth+1)
	locals := make([]uint64, cf.totalLocalVars)

	for i := cf.args - 1; i > 0; i-- {
		locals[i] = vm.popUint64()
	}

	prevCtx := vm.ctx
	vm.ctx = context{
		stack:   stack,
		locals:  locals,
		code:    cf.code,
		pc:      0,
		curFunc: index,
	}

	ret := vm.execCode()
	vm.ctx = prevCtx
	if cf.returns {
		vm.pushUint64(ret.(uint64))
	}
}
