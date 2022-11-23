package exec

import (
	"errors"
	"github.com/threadedstream/wasmexperiments/internal/pkg/reporter"
)

const (
	// make an interpreter option?
	maxDepth         = 15
	maxStackFrameNum = 256
)

type Function struct {
	numLocals int
	numParams int
	code      []byte
	returns   bool
	name      string
}

func (fn *Function) call(vm *VM, index int64, args ...uint64) (any, error) {
	if len(args) != fn.numParams {
		return nil, errors.New("number of arguments do not match")
	}

	stack := make([]uint64, 0, maxDepth)

	disasmedCode, err := Disassemble(fn.code)
	if err != nil {
		return nil, err
	}

	Dump(disasmedCode)
	//compiledCode, _ := Compile(disasmedCode)

	vm.ctx = &context{
		stack:   stack,
		raw:     nil,
		ins:     disasmedCode,
		pc:      0,
		curFunc: index,
	}

	for _, arg := range args {
		vm.pushUint64(arg)
	}

	var locals []uint64
	for i := fn.numParams; i > 0; i-- {
		locals = append(locals, vm.popUint64())
	}

	vm.ctx.locals = locals

	vm.frames = append(vm.frames, vm.ctx)

	ret := fn.execCode(vm)

	vm.frames = vm.frames[:len(vm.frames)-1]
	if fn.returns {
		return ret, nil
	}

	return nil, nil
}

func (fn *Function) execCode(vm *VM) any {
	// check if function returns something
	val := vm.execCode()
	if fn.returns {
		if val != nil {
			return val
		} else {
			reporter.ReportError("expected to have return value on stack")
		}
	}
	return nil
}
