package exec

import "github.com/threadedstream/wasmexperiments/internal/types"

var (
	callOp = newOp("call", 0x10, types.ValueTypeSingleI32, types.ValueTypeEmpty)
	ifOp   = newOp("if", 0x04, types.ValueTypeSingleI32, types.ValueTypeEmpty)
	elseOp = newOp("else", 0x05, types.ValueTypeSingleI32, types.ValueTypeEmpty)
	endOp  = newOp("end", 0x0B, types.ValueTypeEmpty, types.ValueTypeEmpty)
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
