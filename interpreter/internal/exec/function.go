package exec

import (
	"encoding/binary"
	"errors"
	"github.com/threadedstream/wasmexperiments/internal/pkg/reporter"
	"github.com/threadedstream/wasmexperiments/internal/types"
)

const (
	// make an interpreter option?
	maxDepth = 32
)

type ExecutionMode int

const (
	ExecutionModeInstructionSequence ExecutionMode = iota
	ExecutionModeRawBytecode
)

type branchInfo struct {
	nestingLevel int
	jumpToPc     int
}

type ifInfo struct {
	nestingLevel int
	jumpToElsePc *int
	jumpToEndPc  int
}

type Function struct {
	name              string
	numLocals         int
	numParams         int
	code              []byte
	returns           bool
	ty                types.ValueType
	blockStartEndInfo map[int]int
	branchingInfo     map[int]branchInfo
	ifBranchingInfo   map[int]ifInfo
	elseBranchingInfo map[int]int
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
	vm.blockStartEnd = fn.blockStartEndInfo
	vm.branchingInfo = fn.branchingInfo
	vm.ifBranchingInfo = fn.ifBranchingInfo
	vm.elseBranchingInfo = fn.elseBranchingInfo

	vm.ctx.stack = append(vm.ctx.stack, make([]uint64, 0, maxDepth))

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
		case localGetOp, globalGetOp, localSetOp, callOp, i32ConstOp, brIfOp, brOp:
			derefPc += 5
		case i32AddOp, i32SubOp, i32MulOp, i32DivUOp, i32DivSOp, i32EqOp, i32LtSOp, elseOp:
			derefPc++
		case blockOp, loopOp:
			blockRecords[blockIdx] = derefPc
			blockIdx++
			derefPc += 2
		case ifOp:
			derefPc += 2
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

	l := len(blockRecords)
	blockStartEnd := map[int]int{}
	for k, v := range blockRecords {
		blockStartEnd[v] = endRecords[l-k-1]
	}
	fn.blockStartEndInfo = blockStartEnd
}

func (fn *Function) gatherBranchingInfo(code []byte) {
	type blockinfo struct {
		name         string
		pc           int
		nestingLevel int
	}
	branchingInfo := make(map[int]branchInfo)
	ifBranchingInfo := make(map[int]ifInfo)
	elseBranchingInfo := make(map[int]int)
	blockchain := make([]blockinfo, 0)
	ifBlockchain := make([]blockinfo, 0)
	elseBlockchain := make([]blockinfo, 0)
	var nestingLevel int
	do := func(opcode Opcode, pc *int) error {
		if pc == nil {
			panic("gatherBranchingInfo.do: nil pc")
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
			block := blockchain[len(blockchain)-int(idx)-1]
			if block.name == "loop" {
				branchingInfo[brPc] = branchInfo{jumpToPc: block.pc, nestingLevel: block.nestingLevel}
			} else if block.name == "block" {
				endPc := fn.blockStartEndInfo[block.pc]
				branchingInfo[brPc] = branchInfo{jumpToPc: endPc + 1, nestingLevel: block.nestingLevel}
				blockchain = blockchain[:len(blockchain)-1]
			}
			derefPc += 4
		case blockOp:
			blockchain = append(blockchain, blockinfo{"block", derefPc, nestingLevel})
			nestingLevel++
			derefPc += 2
		case loopOp:
			blockchain = append(blockchain, blockinfo{"loop", derefPc, nestingLevel})
			nestingLevel++
			derefPc += 2
		case ifOp:
			ifBlockchain = append(ifBlockchain, blockinfo{"if", derefPc, nestingLevel})
			ifBranchingInfo[derefPc] = ifInfo{}
			derefPc += 2
			// todo
		case elseOp:
			elseBlockchain = append(elseBlockchain, blockinfo{"else", derefPc, nestingLevel})
			block := ifBlockchain[len(ifBlockchain)-1]
			elsePc := derefPc
			ifBranchingInfo[block.pc] = ifInfo{jumpToElsePc: &elsePc}
			derefPc++
		case endOp:
			if len(ifBlockchain) > 0 {
				block := ifBlockchain[len(ifBlockchain)-1]
				info := ifBranchingInfo[block.pc]
				info.jumpToEndPc = derefPc
				ifBranchingInfo[block.pc] = info
				if len(elseBlockchain) > 0 {
					elseBlock := elseBlockchain[len(elseBlockchain)-1]
					if info.jumpToElsePc != nil && *info.jumpToElsePc == elseBlock.pc {
						elseBranchingInfo[elseBlock.pc] = derefPc
					}
				}
			}
			derefPc++
		case returnOp:
			derefPc++
		}
		*pc = derefPc
		return nil
	}

	_ = Visit(code, do)
	fn.ifBranchingInfo = ifBranchingInfo
	fn.elseBranchingInfo = elseBranchingInfo
	fn.branchingInfo = branchingInfo
}

func (fn *Function) execRawBytecode(vm *VM) any {
	for int(vm.ctx.pc) < len(vm.ctx.compiledCode) {
		if vm.returned {
			vm.returned = false
			break
		}
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

func (vm *VM) execRange(start, end int) any {
	for vm.ctx.pc = int64(start); vm.ctx.pc < int64(end); {
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
