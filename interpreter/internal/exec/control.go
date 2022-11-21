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
)

func (vm *VM) execIf() {
	var retVal any
	in := vm.currIns().(*IfI)
	// decide if we enter "if" body
	val := vm.popUint32()
	if val > 0 {
		newCtx := &context{
			stack:   make([]uint64, maxDepth),
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
	// don't do anything yet
	vm.ctx.pc++
}
