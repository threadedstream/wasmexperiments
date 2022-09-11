package exec

const (
	// make an interpreter option?
	maxDepth = 15
)

type Function struct {
	numLocals int
	numParams int
	code      []byte
	returns   bool
	name      string
}

func (fn *Function) call(vm *VM, index int64) {
	stack := make([]uint64, 0, maxDepth)
	locals := make([]uint64, fn.numLocals)

	for i := fn.numParams - 1; i > 0; i-- {
		locals[i] = vm.popUint64()
	}

	prevCtx := vm.ctx
	vm.ctx = context{
		stack:   stack,
		locals:  locals,
		code:    fn.code,
		pc:      0,
		curFunc: index,
	}

	ret := vm.execCode()
	vm.ctx = prevCtx
	if fn.returns {
		vm.pushUint64(ret.(uint64))
	}
}
