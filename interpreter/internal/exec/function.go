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

type ExecutionMode int

const (
	ExecutionModeInstructionSequence ExecutionMode = iota
	ExecutionModeRawBytecode
)

type Function struct {
	numLocals int
	numParams int
	code      []byte
	returns   bool
	name      string
}

func (fn *Function) call(vm *VM, index int64, mode ExecutionMode, args ...uint64) (any, error) {
	if len(args) != fn.numParams {
		return nil, errors.New("number of arguments do not match")
	}

	stack := make([]uint64, 0, maxDepth)

	//disasmedCode, err := Disassemble(fn.code)
	//if err != nil {
	//	return nil, err
	//}

	//Dump(disasmedCode)
	compiledCode, _ := Compile(fn.code)

	_ = fn.staticallyAnalyze(compiledCode)

	vm.ctx = &context{
		stack:        stack,
		compiledCode: compiledCode,
		ins:          nil,
		pc:           0,
		curFunc:      index,
	}

	vm.ctxchain = append(vm.ctxchain, vm.ctx)

	for _, arg := range args {
		vm.pushUint64(arg)
	}

	var locals []uint64
	for i := fn.numParams; i > 0; i-- {
		locals = append(locals, vm.popUint64())
	}

	vm.ctx.locals = locals

	vm.frames = append(vm.frames, vm.ctx)

	var ret any
	if mode == ExecutionModeRawBytecode {
		ret = fn.execRawBytecode(vm)
	} else {
		ret = fn.execInstrSeq(vm)
	}

	vm.frames = vm.frames[:len(vm.frames)-1]
	if fn.returns {
		return ret, nil
	}

	return nil, nil
}

func (fn *Function) execInstrSeq(vm *VM) any {
	// check if function returns something
	_ = vm.execCode()
	if fn.returns {
		if len(vm.ctx.stack) > 0 {
			return vm.popUint64()
		} else {
			reporter.ReportError("expected to have return value on stack")
		}
	}
	return nil
}

func (fn *Function) staticallyAnalyze(code []byte) map[int]int {
	blockRecords := map[int]int{}
	endRecords := map[int]int{}

	pc := 0
	blockIdx := 0
	endIdx := 0
	for ; pc < len(code); pc++ {
		if (Bytecode(code[pc]) == blockOp) || (Bytecode(code[pc]) == loopOp) {
			blockRecords[blockIdx] = pc
			blockIdx++
		} else if Bytecode(code[pc]) == endOp {
			endRecords[endIdx] = pc
			endIdx++
		}
	}
	if len(blockRecords) != len(endRecords) {
		panic("unmatched number of block/loop and ends")
	}
	l := len(blockRecords)
	blockStartEnd := map[int]int{}
	for k, v := range blockRecords {
		blockStartEnd[v] = endRecords[l-k]
	}
	return blockStartEnd
}

func (fn *Function) execRawBytecode(vm *VM) any {
	for int(vm.ctx.pc) < len(vm.ctx.compiledCode) {
		// skip instruction
		if handler, ok := funcTable[Bytecode(vm.ctx.compiledCode[vm.ctx.pc])]; ok {
			vm.ctx.pc++
			handler()
			continue
		}
	}
	if len(vm.ctx.stack) > 0 {
		return vm.popUint64()
	}
	return nil
}
