package exec

import (
	"encoding/binary"
	"errors"
	"github.com/threadedstream/wasmexperiments/internal/types"
)

var (
	localGetOp  = newVarargOp("local.get", 0x20)
	localSetOp  = newVarargOp("local.set", 0x21)
	globalGetOp = newVarargOp("global.get", 0x23)
	globalSetOp = newVarargOp("global.set", 0x24)
	i32LoadOp   = newOp("i32.load", 0x28, types.ValueTypeDoubleI32, types.ValueTypeSingleI32)
	f32LoadOp   = newOp("f32.load", 0x2a, types.ValueTypeDoubleI32, types.ValueTypeSingleF32)
	i32StoreOp  = newOp("i32.store", 0x36, types.ValueTypeDoubleI32, types.ValueTypeVoid)
	f32StoreOp  = newOp("f32.store", 0x38, types.ValueTypeDoubleI32, types.ValueTypeVoid)
)

var (
	ErrOutOfMemory = errors.New("exec: out of memory")
)

func (vm *VM) getLocal() {
	index := binary.LittleEndian.Uint32(vm.ctx.compiledCode[vm.ctx.pc:])
	vm.pushUint64(vm.ctx.locals[index])
	vm.ctx.pc += 4
}

func (vm *VM) setLocal() {
	index := binary.LittleEndian.Uint32(vm.ctx.compiledCode[vm.ctx.pc:])
	value := vm.popUint64()
	vm.ctx.locals[index] = value
	vm.ctx.pc += 4
}

func (vm *VM) getGlobal() {
	index := binary.LittleEndian.Uint32(vm.ctx.compiledCode[vm.ctx.pc:])
	vm.pushUint64(vm.globals[index])
	vm.ctx.pc += 4
}

func (vm *VM) setGlobal() {
	index := binary.LittleEndian.Uint32(vm.ctx.compiledCode[vm.ctx.pc:])
	value := vm.popUint64()
	vm.globals[index] = value
	vm.ctx.pc += 4
}

func (vm *VM) teeLocal() {
	panic("unreachable")
}

func (vm *VM) i32Load() {
	base := int(vm.peekUint32())
	align := int(binary.LittleEndian.Uint32(vm.ctx.compiledCode[vm.ctx.pc:]))
	if align != 2 {
		panic("should have alignment of 2")
	}
	vm.ctx.pc += 4
	off := int(binary.LittleEndian.Uint32(vm.ctx.compiledCode[vm.ctx.pc:]))
	vm.ctx.pc += 4
	if !vm.inBounds(base, off, 3) {
		panic(ErrOutOfMemory)
	}
	valueAt := vm.memory[(base + off):]
	vm.pushUint32(binary.LittleEndian.Uint32(valueAt))
}

func (vm *VM) i32Store() {
	val := vm.popUint32()
	base := int(vm.peekUint32())
	ioff := int(binary.LittleEndian.Uint32(vm.ctx.compiledCode[vm.ctx.pc:]))
	if !vm.inBounds(base, ioff, 3) {
		panic(ErrOutOfMemory)
	}
	effAddr := vm.memory[(base + ioff):]
	binary.LittleEndian.PutUint32(effAddr, val)
	vm.ctx.pc += 4
}

func (vm *VM) f32Load() {
	panic("unreachable")
	//if !vm.inBounds(3) {
	//	panic(ErrOutOfMemory)
	//}
	//vm.pushFloat32(math.Float32frombits(binary.LittleEndian.Uint32(vm.currMem())))
}

func (vm *VM) f32Store() {
	panic("unreachable")
	//if !vm.inBounds(3) {
	//	panic(ErrOutOfMemory)
	//}
	//v := math.Float32bits(vm.popFloat32())
	//binary.LittleEndian.PutUint32(vm.currMem(), v)
}

func (vm *VM) inBounds(base, ioffset, offset int) bool {
	addr := uint64(base) + uint64(ioffset)
	return addr+uint64(offset) < uint64(len(vm.memory))
}

func (vm *VM) fetchEffectiveAddr() int {
	baseAddr := vm.fetchUint32()
	offset := vm.popInt32()
	return int(baseAddr + uint32(offset))
}

func (vm *VM) currMem() []byte {
	return vm.memory[vm.fetchEffectiveAddr():]
}
