package exec

import (
	"encoding/binary"
	"errors"
	"github.com/threadedstream/wasmexperiments/internal/types"
	"math"
)

var (
	localGetOp  = newVarargOp("local.get", 0x20)
	localSetOp  = newVarargOp("local.set", 0x21)
	globalGetOp = newVarargOp("global.get", 0x23)
	globalSetOp = newVarargOp("global.set", 0x24)
	i32LoadOp   = newOp("i32.load", 0x28, types.ValueTypeDoubleI32, types.ValueTypeSingleI32)
	f32LoadOp   = newOp("f32.load", 0x2a, types.ValueTypeDoubleI32, types.ValueTypeSingleF32)
	i32StoreOp  = newOp("i32.store", 0x36, types.ValueTypeDoubleI32, types.ValueTypeEmpty)
	f32StoreOp  = newOp("f32.store", 0x38, types.ValueTypeDoubleI32, types.ValueTypeEmpty)
)

var (
	ErrOutOfMemory = errors.New("exec: out of memory")
)

func (vm *VM) getLocal() {
	vm.ctx.pc++
	index := vm.fetchUint32()
	vm.pushUint64(vm.ctx.locals[index])
}

func (vm *VM) setLocal() {
	index := vm.fetchUint32()
	value := vm.popUint64()
	vm.ctx.locals[index] = value
}

func (vm *VM) getGlobal() {
	index := vm.fetchUint32()
	vm.pushUint64(vm.globals[index])
}

func (vm *VM) setGlobal() {
	index := vm.fetchUint32()
	value := vm.popUint64()
	vm.globals[index] = value
}

func (vm *VM) i32Load() {
	if !vm.inBounds(3) {
		panic(ErrOutOfMemory)
	}
	vm.pushUint32(binary.LittleEndian.Uint32(vm.currMem()))
}

func (vm *VM) i32Store() {
	if !vm.inBounds(3) {
		panic(ErrOutOfMemory)
	}
	v := vm.popUint32()
	binary.LittleEndian.PutUint32(vm.currMem(), v)
}

func (vm *VM) f32Load() {
	if !vm.inBounds(3) {
		panic(ErrOutOfMemory)
	}
	vm.pushFloat32(math.Float32frombits(binary.LittleEndian.Uint32(vm.currMem())))
}

func (vm *VM) f32Store() {
	if !vm.inBounds(3) {
		panic(ErrOutOfMemory)
	}
	v := math.Float32bits(vm.popFloat32())
	binary.LittleEndian.PutUint32(vm.currMem(), v)
}

func (vm *VM) inBounds(offset int) bool {
	addr := uint64(binary.LittleEndian.Uint32(vm.ctx.code[vm.ctx.pc:])) + uint64(vm.peekUint32())
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
