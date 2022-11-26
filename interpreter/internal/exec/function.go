package exec

import (
	"encoding/binary"
	"errors"
	"github.com/threadedstream/wasmexperiments/internal/pkg/reporter"
)

const (
	// make an interpreter option?
	maxDepth = 15
)

type ExecutionMode int

const (
	ExecutionModeInstructionSequence ExecutionMode = iota
	ExecutionModeRawBytecode
)

type Function struct {
	numLocals         int
	numParams         int
	code              []byte
	returns           bool
	name              string
	blockStartEndInfo map[int]int
	branchingInfo     map[int]int
}

type ifRecord struct {
	elsePc *int
	endPc  int
}

func (fn *Function) call(vm *VM, index int64, mode ExecutionMode, args ...uint64) (any, error) {
	if len(args) != fn.numParams {
		return nil, errors.New("number of arguments do not match")
	}

	stack := make([][]uint64, 0)

	//disasmedCode, err := Disassemble(fn.code)
	//if err != nil {
	//	return nil, err
	//}

	//Dump(disasmedCode)
	compiledCode, _ := Compile(fn.code)

	fn.gatherBlockInfo(compiledCode)
	fn.gatherBranchingInfo(compiledCode)

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

func (fn *Function) gatherBlockInfo(code []byte) {
	blockRecords := map[int]int{}
	endRecords := map[int]int{}
	blockIdx := 0
	endIdx := 0
	do := func(opcode Opcode, pc *int) error {
		if pc == nil {
			panic("staticallyAnalyze.do: nil pc")
		}
		derefPc := *pc
		switch opcode {
		case i32LoadOp:
			derefPc += 9
		case localGetOp, globalGetOp, localSetOp, callOp, i32ConstOp:
			derefPc += 5
		case i32AddOp, i32SubOp, i32MulOp, i32DivUOp, i32DivSOp, i32EqOp, i32LtSOp:
			derefPc++
		case blockOp, loopOp:
			blockRecords[blockIdx] = derefPc
			blockIdx++
			derefPc += 2
		case ifOp, elseOp:
			derefPc++
			// todo
		case endOp:
			endRecords[endIdx] = derefPc
			endIdx++
			derefPc++
		case returnOp:
			derefPc++
		}
		*pc = derefPc
		return nil
	}

	_ = Visit(code, do)

	if len(blockRecords) != len(endRecords) {
		panic("unmatched number of block/loop and ends")
	}
	l := len(blockRecords)
	blockStartEnd := map[int]int{}
	for k, v := range blockRecords {
		blockStartEnd[v] = endRecords[l-k]
	}
	fn.blockStartEndInfo = blockStartEnd
}

func (fn *Function) gatherBranchingInfo(code []byte) {
	type blockinfo struct {
		name string
		pc   int
	}
	branchInfo := make(map[int]int)
	blockchain := make([]blockinfo, 0)
	do := func(opcode Opcode, pc *int) error {
		if pc == nil {
			panic("staticallyAnalyze.do: nil pc")
		}
		derefPc := *pc
		switch opcode {
		case i32LoadOp:
			derefPc += 9
		case localGetOp, globalGetOp, localSetOp, callOp, i32ConstOp:
			derefPc += 5
		case i32AddOp, i32SubOp, i32MulOp, i32DivUOp, i32DivSOp, i32EqOp, i32LtSOp:
			derefPc++
		case brOp, brIfOp:
			brPc := derefPc
			derefPc++
			idx := binary.LittleEndian.Uint32(code[derefPc:])
			block := blockchain[len(blockchain)-int(idx)]
			branchInfo[brPc] = block.pc
			derefPc += 4
		case blockOp:
			blockchain = append(blockchain, blockinfo{"block", derefPc})
			derefPc += 2
		case loopOp:
			blockchain = append(blockchain, blockinfo{"loop", derefPc})
			derefPc += 2
		case ifOp, elseOp:
			derefPc++
			// todo
		case endOp:
			derefPc++
		case returnOp:
			derefPc++
		}
		*pc = derefPc
		return nil
	}

	_ = Visit(code, do)
	fn.branchingInfo = branchInfo
}

func (fn *Function) execRawBytecode(vm *VM) any {
	for int(vm.ctx.pc) < len(vm.ctx.compiledCode) {
		// skip instruction
		if handler, ok := funcTable[Opcode(vm.ctx.compiledCode[vm.ctx.pc])]; ok {
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
