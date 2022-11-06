package exec

import "errors"

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

func (fn *Function) call(vm *VM, index int64, args ...uint64) (any, error) {
	if len(args) != fn.numParams {
		return nil, errors.New("number of arguments do not match")
	}

	stack := make([]uint64, 0, maxDepth)
	var locals []uint64

	disasmedCode, err := Disassemble(fn.code)
	if err != nil {
		return nil, err
	}

	Dump(disasmedCode)
	compiledCode, _ := Compile(disasmedCode)
	prevCtx := vm.ctx
	vm.ctx = context{
		stack:   stack,
		code:    compiledCode,
		pc:      0,
		curFunc: index,
	}

	for _, arg := range args {
		vm.pushUint64(arg)
	}

	for i := fn.numParams; i > 0; i-- {
		locals = append(locals, vm.popUint64())
	}

	vm.ctx.locals = locals

	ret := fn.execCode(vm)
	vm.ctx = prevCtx
	if fn.returns {
		return ret, nil
	}

	return nil, nil
}

func (fn *Function) execCode(vm *VM) any {
	code := vm.ctx.code
	endOff := len(code)
	for int(vm.ctx.pc) < endOff {
		switch Bytecode(code[vm.ctx.pc]) {
		case i32AddOp:
			vm.i32Add()
			vm.ctx.pc++
		case callOp:
			vm.call()
		case localGetOp:
			vm.ctx.pc++
			vm.getLocal()
		case endOp:
			break
		}
	}
	if len(vm.ctx.stack) != 0 {
		return vm.popUint32()
	}
	return nil
}
