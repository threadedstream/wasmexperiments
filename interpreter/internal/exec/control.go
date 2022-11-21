package exec

import "github.com/threadedstream/wasmexperiments/internal/types"

var (
	callOp   = newOp("call", 0x10, types.ValueTypeSingleI32, types.ValueTypeVoid)
	ifOp     = newOp("if", 0x04, types.ValueTypeSingleI32, types.ValueTypeVoid)
	elseOp   = newOp("else", 0x05, types.ValueTypeSingleI32, types.ValueTypeVoid)
	endOp    = newOp("end", 0x0B, types.ValueTypeVoid, types.ValueTypeVoid)
	returnOp = newOp("return", 0x0F, types.ValueTypeVoid, types.ValueTypeVoid)
)

func (vm *VM) execIf() {
	/*	// decide if we enter if body
		val := vm.popUint32()
		if val > 0 {
			// enter if body
			// save previous stack
			old := vm.ctx.stack
			vm.ctx.stack = make([]uint64, maxDepth)
		} else {

		}
	*/
}
