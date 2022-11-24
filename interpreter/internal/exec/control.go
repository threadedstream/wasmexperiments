package exec

import (
	"github.com/threadedstream/wasmexperiments/internal/pkg/reporter"
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

func (vm *VM) execBlock() {
	var retVal any
	in := vm.currIns().(*BlockI)
	newCtx := &context{
		parent:  vm.ctx,
		stack:   make([]uint64, 0, maxDepth),
		locals:  vm.ctx.locals,
		raw:     vm.ctx.raw,
		ins:     in.body,
		curFunc: vm.ctx.curFunc,
		isBlock: true,
	}
	vm.ctx.pc++
	vm.ctx = newCtx
	vm.ctxchain = append(vm.ctxchain, newCtx)
	retVal = vm.execCode()
	switch bt := in.blockType.(type) {
	case types.EmptyBlockType:
		break
	case types.ResultBlockType:
		v := retVal.(int64)
		switch bt.Ty {
		case types.ValueTypeI32:
			vm.pushInt32(int32(v))
		case types.ValueTypeI64:
			vm.pushInt64(v)
		case types.ValueTypeF32:
			// not sure if that's going to work with floats, though
			vm.pushFloat32(math.Float32frombits(uint32(v)))
		case types.ValueTypeF64:
			vm.pushFloat64(math.Float64frombits(uint64(v)))
		}
	case types.OtherBlockType:
		reporter.ReportError("unknown result type with value %d", bt.X)
	}

	if newCtx.breakExecuted || vm.returned {
		return
	}

}

func (vm *VM) execBr() {
	in := vm.currIns().(*BrI)
	if in.context == "block" {
		vm.ctx.breakExecuted = true
		var temp uint64
		if len(vm.ctx.stack) > 0 {
			temp = vm.popUint64()
		}
		off := in.arg0.(uint32)
		neededCtx := vm.ctx
		idx := uint32(0)
		for idx < off {
			neededCtx = neededCtx.parent
			idx++
		}
		vm.ctx = neededCtx.parent
		vm.pushUint64(temp)
	} else if in.context == "loop" {
		off := in.arg0.(uint32)
		neededCtx := vm.ctx
		idx := uint32(0)
		for idx < off {
			neededCtx = neededCtx.parent
			idx++
		}

		vm.ctx.breakExecuted = true
		vm.ctx.pendingContext = neededCtx
		if neededCtx.isBlock {
			neededCtx.breakExecuted = true
		}
		neededCtx.pc = 0
	}
	vm.ctx.pc++
}

func (vm *VM) execBrBlock() {

}

func (vm *VM) execBrLoop() {
	in := vm.currIns().(*BrI)
	// specify location to jump to and start executing this thing over again
	id := in.arg0.(uint32)
	idx := len(vm.ctxchain) - int(id) - 1
	neededCtx := vm.ctxchain[idx]
	if neededCtx.isBlock {
		neededCtx.pc = int64(len(neededCtx.ins))
		vm.ctxchain = vm.ctxchain[:idx]
		vm.ctx = vm.ctxchain[len(vm.ctxchain)-1]
	} else {
		vm.ctxchain = vm.ctxchain[:idx+1]
		vm.ctx = vm.ctxchain[len(vm.ctxchain)-1]
		vm.ctx.pc = 0
	}
}

func (vm *VM) execBrIfBlock() {

}

func (vm *VM) execBrIfLoop() {
	var _ any
	in := vm.currIns().(*BrIfI)
	val := vm.popUint32()
	if val > 0 {
		// specify location to jump to and start executing this thing over again
		id := in.arg0.(uint32)
		idx := len(vm.ctxchain) - int(id) - 1
		neededCtx := vm.ctxchain[idx]
		if neededCtx.isBlock {
			neededCtx.pc = int64(len(neededCtx.ins))
			vm.ctxchain = vm.ctxchain[:idx]
			vm.ctx = vm.ctxchain[len(vm.ctxchain)-1]
		} else {
			// neededCtx
			vm.ctxchain = vm.ctxchain[:idx+1]
			vm.ctx = vm.ctxchain[len(vm.ctxchain)-1]
			vm.ctx.pc = 0
		}
		// block->loop
	} else {
		// just increment pc
		vm.ctx.pc++
	}
}

func (vm *VM) execBrIf() {
	var _ any
	in := vm.currIns().(*BrIfI)
	// decide if we enter "if" body
	val := vm.popUint32()

	if val > 0 {
		if in.context == "block" {
			vm.ctx.breakExecuted = true
			var temp uint64
			if len(vm.ctx.stack) > 0 {
				temp = vm.popUint64()
			}
			off := in.arg0.(uint32)
			neededCtx := vm.ctx
			idx := uint32(0)
			for idx < off {
				neededCtx = neededCtx.parent
				idx++
			}
			// TODO(threadedstream): make use of pendingCtx
			vm.ctx = neededCtx.parent
			vm.pushUint64(temp)
		} else if in.context == "loop" {
			off := in.arg0.(uint32)
			neededCtx := vm.ctx
			idx := uint32(0)
			for idx < off {
				neededCtx = neededCtx.parent
				idx++
			}
			vm.ctx.breakExecuted = true
			vm.ctx.pendingContext = neededCtx
			if neededCtx.isBlock {
				neededCtx.breakExecuted = true
			}
		}
	}
	vm.ctx.pc++
}

func (vm *VM) execLoop() {
	var retVal any
	in := vm.currIns().(*LoopI)
	newCtx := &context{
		parent:  vm.ctx,
		stack:   make([]uint64, 0, maxDepth),
		locals:  vm.ctx.locals,
		raw:     vm.ctx.raw,
		ins:     in.body,
		curFunc: vm.ctx.curFunc,
	}
	vm.ctx.pc++
	vm.ctx = newCtx
	vm.ctxchain = append(vm.ctxchain, newCtx)
	retVal = vm.execCode()
	switch bt := in.blockType.(type) {
	case types.EmptyBlockType:
		break
	case types.ResultBlockType:
		v := retVal.(int64)
		switch bt.Ty {
		case types.ValueTypeI32:
			vm.pushInt32(int32(v))
		case types.ValueTypeI64:
			vm.pushInt64(v)
		case types.ValueTypeF32:
			// not sure if that's going to work with floats, though
			vm.pushFloat32(math.Float32frombits(uint32(v)))
		case types.ValueTypeF64:
			vm.pushFloat64(math.Float64frombits(uint64(v)))
		}
	case types.OtherBlockType:
		reporter.ReportError("unknown result type with value %d", bt.X)
	}

	if newCtx.breakExecuted || vm.returned {
		return
	}
}

func (vm *VM) execIf() {
	var retVal any
	in := vm.currIns().(*IfI)
	// decide if we enter "if" body
	val := vm.popUint32()
	if val > 0 {
		newCtx := &context{
			stack:   make([]uint64, 0, maxDepth),
			locals:  vm.ctx.locals,
			raw:     vm.ctx.raw,
			ins:     in.body,
			curFunc: vm.ctx.curFunc,
		}
		vm.frames = append(vm.frames, newCtx)
		vm.ctx = newCtx
		retVal = vm.execCode()
		vm.frames = vm.frames[:len(vm.frames)-1]
		// blindly accept that vm.frames always has at least one frame at this point
		vm.ctx = vm.frames[len(vm.frames)-1]
	} else {
		if in.elseBody != nil {
			newCtx := &context{
				stack:   make([]uint64, 0, maxDepth),
				locals:  vm.ctx.locals,
				raw:     vm.ctx.raw,
				ins:     in.elseBody,
				curFunc: vm.ctx.curFunc,
			}
			vm.frames = append(vm.frames, newCtx)
			vm.ctx = newCtx
			retVal = vm.execCode()
			vm.frames = vm.frames[:len(vm.frames)-1]
			// blindly accept that vm.frames always has at least one frame at this point
			vm.ctx = vm.frames[len(vm.frames)-1]
		}
	}

	switch bt := in.blockType.(type) {
	case types.EmptyBlockType:
		return
	case types.ResultBlockType:
		v := retVal.(int64)
		switch bt.Ty {
		case types.ValueTypeI32:
			vm.pushInt32(int32(v))
		case types.ValueTypeI64:
			vm.pushInt64(v)
		case types.ValueTypeF32:
			// not sure if that's going to work with floats, though
			vm.pushFloat32(math.Float32frombits(uint32(v)))
		case types.ValueTypeF64:
			vm.pushFloat64(math.Float64frombits(uint64(v)))
		}
	case types.OtherBlockType:
		reporter.ReportError("unknown result type with value %d", bt.X)
	}
	vm.ctx.pc++
}

func (vm *VM) ret() {
	vm.returned = true
}
