package exec

import (
	"encoding/binary"
	"github.com/threadedstream/wasmexperiments/internal/types"
	"math"
)

var (
	callOp   = newOp("call", 0x10, types.ValueTypeSingleI32, types.ValueTypeVoid)
	ifOp     = newOp("if", 0x04, types.ValueTypeSingleI32, types.ValueTypeVoid)
	elseOp   = newOp("else", 0x05, types.ValueTypeSingleI32, types.ValueTypeVoid)
	endOp    = newOp("end", 0x0B, types.ValueTypeVoid, types.ValueTypeVoid)
	returnOp = newOp("return", 0x0F, types.ValueTypeVoid, types.ValueTypeVoid)
	blockOp  = newOp("block", 0x02, types.ValueTypeSingleI32, types.ValueTypeVoid)
	brOp     = newOp("br", 0x0C, types.ValueTypeSingleI32, types.ValueTypeVoid)
	brIfOp   = newOp("br_if", 0x0D, types.ValueTypeSingleI32, types.ValueTypeVoid)
	loopOp   = newOp("loop", 0x03, types.ValueTypeSingleI32, types.ValueTypeVoid)
)

func (vm *VM) execCall() {
	index := binary.LittleEndian.Uint32(vm.ctx.compiledCode[vm.ctx.pc:])
	vm.ctx.pc += 4
	fn := vm.module.GetFunction(int(index))
	var args []uint64
	for i := 0; i < fn.numParams; i++ {
		args = append(args, vm.popUint64())
	}
	val, err := fn.call(vm, int64(index), ExecutionModeRawBytecode, args...)
	if err != nil {
		panic(err)
	}
	vm.ctx = vm.frames[len(vm.frames)-1]
	switch fn.ty {
	case types.ValueTypeEmpty:
		return
	case types.ValueTypeI32:
		vm.pushInt32(int32(val.(uint64)))
	case types.ValueTypeI64:
		vm.pushInt64(int64(val.(uint64)))
	case types.ValueTypeF32:
		// not sure if that's going to work with floats, though
		vm.pushFloat32(math.Float32frombits(uint32(val.(uint64))))
	case types.ValueTypeF64:
		vm.pushFloat64(math.Float64frombits(val.(uint64)))
	}
}

func (vm *VM) execBlock() {
	stack := make([]uint64, 0, maxDepth)
	vm.ctx.stack = append(vm.ctx.stack, stack)
	// ignore type for now
	vm.ctx.pc++
}

func (vm *VM) execBr() {
	label := binary.LittleEndian.Uint32(vm.ctx.compiledCode[vm.ctx.pc:])
	brInfo := vm.branchingInfo[int(vm.ctx.pc-1)]
	vm.ctx.stack = vm.ctx.stack[:len(vm.ctx.stack)-int(label)-1]
	vm.ctx.pc = int64(brInfo.jumpToPc)
}

func (vm *VM) execBrIf() {
	val := vm.popUint32()
	if val > 0 {
		// brIf
		label := binary.LittleEndian.Uint32(vm.ctx.compiledCode[vm.ctx.pc:])
		brInfo := vm.branchingInfo[int(vm.ctx.pc-1)]
		vm.ctx.stack = vm.ctx.stack[:len(vm.ctx.stack)-int(label)-1]
		vm.ctx.pc = int64(brInfo.jumpToPc)
	} else {
		vm.ctx.pc += 4
	}
}

func (vm *VM) execLoop() {
	stack := make([]uint64, 0, maxDepth)
	vm.ctx.stack = append(vm.ctx.stack, stack)
	// ignore type for now
	vm.ctx.pc++
}

func (vm *VM) execIf() {
	value := vm.popUint64()
	vm.pushStack()
	info := vm.ifBranchingInfo[int(vm.ctx.pc-1)]
	vm.ctx.pc++
	if value > 0 {
		return
	} else {
		if info.jumpToElsePc != nil {
			vm.ctx.pc = int64(*info.jumpToElsePc + 1)
		} else {
			vm.ctx.pc = int64(info.jumpToEndPc)
		}
	}
}

func (vm *VM) execElse() {
	targetJump := vm.elseBranchingInfo[int(vm.ctx.pc-1)]
	vm.ctx.pc = int64(targetJump)
}

func (vm *VM) ret() {
	vm.returned = true
}

func (vm *VM) execEnd() {
	if len(vm.ctx.stack[len(vm.ctx.stack)-1]) > 0 {
		val := vm.popUint64()
		vm.popStack()
		vm.pushUint64(val)
	} else {
		vm.popStack()
	}
	//m := invertMap(vm.blockStartEnd)
	//endPc := vm.ctx.pc
	//if blockPc, ok := m[int(endPc)]; ok {
	//	blockPc++
	//	blockType := types.ValueType(vm.ctx.compiledCode[blockPc])
	//	switch blockType {
	//	case types.ValueTypeEmpty:
	//		return
	//	case types.ValueTypeI32:
	//		val := vm.popInt32()
	//		vm.pushInt32(val)
	//	case types.ValueTypeI64:
	//		val := vm.popInt64()
	//		vm.pushInt64(val)
	//	case types.ValueTypeF32:
	//		val := vm.popInt32()
	//		// not sure if that's going to work with floats, though
	//		vm.pushFloat32(math.Float32frombits(uint32(val)))
	//	case types.ValueTypeF64:
	//		val := vm.popInt64()
	//		vm.pushFloat64(math.Float64frombits(uint64(val)))
	//	}
	//}
	//return
}

func invertMap[K comparable, V comparable](m map[K]V) map[V]K {
	out := make(map[V]K)
	for k, v := range m {
		out[v] = k
	}
	return out
}
