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
	i32LoadOp   = newOp("i32.load", 0x28, []types.ValueType{types.ValueTypeI32, types.ValueTypeI32}, []types.ValueType{types.ValueTypeI32})
	f32LoadOp   = newOp("f32.load", 0x2a, []types.ValueType{types.ValueTypeI32, types.ValueTypeI32}, []types.ValueType{types.ValueTypeF32})
	//i32StoreOp  = newOp("i32.store", 0x36)
)

var (
	ErrOutOfMemory = errors.New("exec: out of memory")
)

func (vm *VM) getLocal() {
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

func (vm *VM) f32Load() {
	if !vm.inBounds(3) {
		panic(ErrOutOfMemory)
	}
	vm.pushFloat32(math.Float32frombits(binary.LittleEndian.Uint32(vm.currMem())))
}

func (vm *VM) f32Store() {
	v := math.Float32bits(vm.popFloat32())
	if !vm.inBounds(3) {
		panic(ErrOutOfMemory)
	}
	binary.LittleEndian.PutUint32(vm.currMem(), v)
}

func (vm *VM) inBounds(offset int) bool {
	addr := uint64(binary.LittleEndian.Uint32(vm.ctx.code[vm.ctx.pc:])) + uint64(uint32(vm.ctx.stack[len(vm.ctx.stack)-1]))
	return addr+uint64(offset) < uint64(len(vm.memory))
}

func (vm *VM) fetchBaseAddr() int {
	return int(vm.fetchUint32() + uint32(vm.popInt32()))
}

func (vm *VM) currMem() []byte {
	return vm.memory[vm.fetchBaseAddr():]
}
