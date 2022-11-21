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
	in := vm.currIns().(*LocalGetI)
	index := in.arg0.(uint32)
	vm.pushUint64(vm.ctx.locals[index])
	vm.ctx.pc++
}

func (vm *VM) setLocal() {
	in := vm.currIns().(*LocalSetI)
	index := in.arg0.(uint32)
	value := vm.popUint64()
	vm.ctx.locals[index] = value
	vm.ctx.pc++
}

func (vm *VM) getGlobal() {
	in := vm.currIns().(*GlobalGetI)
	index := in.arg0.(uint32)
	vm.pushUint64(vm.globals[index])
	vm.ctx.pc++
}

func (vm *VM) setGlobal() {
	in := vm.currIns().(*GlobalSetI)
	index := in.arg0.(uint32)
	value := vm.popUint64()
	vm.globals[index] = value
	vm.ctx.pc++
}

func (vm *VM) teeLocal() {
	panic("unreachable")
}

func (vm *VM) i32Load() {
	in := vm.currIns().(*I32LoadI)
	base := int(vm.peekUint32())
	ioff := int(in.arg1.(int32))
	off := 3
	if !vm.inBounds(base, ioff, off) {
		panic(ErrOutOfMemory)
	}
	valueAt := vm.memory[(base + ioff):]
	vm.pushUint32(binary.LittleEndian.Uint32(valueAt))
	vm.ctx.pc++
}

func (vm *VM) i32Store() {
	in := vm.currIns().(*I32StoreI)
	val := vm.popUint32()
	base := int(vm.peekUint32())
	ioff := int(in.arg1.(int32))
	off := 3
	if !vm.inBounds(base, ioff, off) {
		panic(ErrOutOfMemory)
	}
	effAddr := vm.memory[(base + ioff):]
	binary.LittleEndian.PutUint32(effAddr, val)
	vm.ctx.pc++
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
