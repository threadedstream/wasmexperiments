package exec

import (
	"encoding/binary"
)

func (vm *VM) i32Const() {
	val := binary.LittleEndian.Uint32(vm.ctx.compiledCode[vm.ctx.pc:])
	vm.pushInt32(int32(val))
	vm.ctx.pc += 4
}

func (vm *VM) f32Const() {
	vm.pushFloat32(vm.fetchFloat32())
}

func (vm *VM) i64Const() {
	vm.pushUint64(vm.fetchUint64())
}

func (vm *VM) f64Const() {
	vm.pushFloat64(vm.fetchFloat64())
}
