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
	// TODO(threadedstream): implement call logic
}
