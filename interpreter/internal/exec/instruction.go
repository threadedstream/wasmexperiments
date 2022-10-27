package exec

import (
	"github.com/threadedstream/wasmexperiments/internal/types"
	"log"
)

type Bytecode byte

const (
	I32Add   Bytecode = 0x6A
	Call              = 0x10
	LocalGet          = 0x20
)

var (
	instrLookup = make(map[Bytecode]*Instr)
)

type Instr struct {
	Name        string
	Code        Bytecode
	InputTypes  []types.ValueType
	OutputTypes []types.ValueType
}

func NewInstr(name string, code Bytecode, inputTypes []types.ValueType, outputTypes []types.ValueType) *Instr {
	if _, ok := instrLookup[code]; ok {
		log.Panic("instruction already registered")
	}
	i := &Instr{
		Name:        name,
		Code:        code,
		InputTypes:  inputTypes,
		OutputTypes: outputTypes,
	}
	instrLookup[code] = i
	return i
}

//// branch instruction
//var instructionMap = map[byte]byte{
//	0x02: iBlock,
//	0x03: iLoop,
//	0x0c: iBr,
//	0x0d: iBrIf,
//	0x0e: iBrTable,
//	0x04: iIf,
//	0x05: iElse,
//	0x0b: iEnd,
//	0x0f: iRet,
//	0x00: iUnreachable,
//	0x01: iNop,
//	0x1a: iDrop,
//	0x41: iI32Const,
//	0x42: iI64Const,
//	0x43: iF32Const,
//	0x44: iF64Const,
//	0x20: iLocalGet,
//	0x21: iLocalSet,
//	0x22: iLocalTee,
//	0x23: iGlobalGet,
//	0x24: iGlobalSet,
//	0x1b: iSelect,
//	0x10: iCall,
//	0x11: iCallIndirect,
//	0x6a: iI32Add,
//	0x7c: iI64Add,
//	0x6b: iI32Sub,
//	0x7d: iI64Sub,
//	0x6c: iI32Mul,
//	0x7e: iI64Mul,
//	0x6d: iI32DivS,
//	0x7f: iI64SivS,
//	0x6e: iI32DivU,
//	0x80: iI64DivU,
//	0x6f: iI32RemS,
//	0x81: iI64RemS,
//	0x70: iI32RemU,
//	0x82: iI64RemU,
//	0x71: iI32And,
//	0x83: iI64And,
//	0x72: iI32Or,
//	0x84: iI64Or,
//	0x73: iI32Xor,
//	0x85: iI64Xor,
//	0x74: iI32Shl,
//	0x86: iI64Shl,
//	0x75: iI32ShrS,
//	0x87: iI64ShrS,
//	0x76: iI32ShrU,
//	0x88: iI64ShrU,
//	0x77: iI32RotL,
//	0x89: iI64RotL,
//	0x78: iI32RotR,
//	0x8a: iI64RotR,
//	0x67: iI32Clz,
//	0x79: iI64Clz,
//	0x68: iI32Ctz,
//	0x7a: iI64Ctz,
//	0x69: iI32PopCnt,
//	0x7b: iI64PopCnt,
//	0x45: iI32EqZ,
//	0x50: iI64EqZ,
//}
//
//const (
////go:generate stringer -type Instruction -linecomment instruction.go
//	iBlock        byte = 0x02
//	iLoop              = 0x03
//	iBr                = 0x0c
//	iBrIf              = 0x0d
//	iBrTable           = 0x0e
//	iIf                = 0x04
//	iElse              = 0x05
//	iEnd               = 0x0b
//	iRet               = 0x0f
//	iUnreachable       = 0x00
//	iNop               = 0x01
//	iDrop              = 0x1a
//	iI32Const          = 0x41
//	iI64Const          = 0x42
//	iF32Const          = 0x43
//	iF64Const          = 0x44
//	iLocalGet          = 0x20
//	iLocalSet          = 0x21
//	iLocalTee          = 0x22
//	iGlobalGet         = 0x23
//	iGlobalSet         = 0x24
//	iSelect            = 0x1b
//	iCall              = 0x10
//	iCallIndirect      = 0x11
//	iI32Add            = 0x6a
//	iI64Add            = 0x7c
//	iI32Sub            = 0x6b
//	iI64Sub            = 0x7d
//	iI32Mul            = 0x6c
//	iI64Mul            = 0x7e
//	iI32DivS           = 0x6d
//	iI64SivS           = 0x7f
//	iI32DivU           = 0x6e
//	iI64DivU           = 0x80
//	iI32RemS           = 0x6f
//	iI64RemS           = 0x81
//	iI32RemU           = 0x70
//	iI64RemU           = 0x82
//	iI32And            = 0x71
//	iI64And            = 0x83
//	iI32Or             = 0x72
//	iI64Or             = 0x84
//	iI32Xor            = 0x73
//	iI64Xor            = 0x85
//	iI32Shl            = 0x74
//	iI64Shl            = 0x86
//	iI32ShrS           = 0x75
//	iI64ShrS           = 0x87
//	iI32ShrU           = 0x76
//	iI64ShrU           = 0x88
//	iI32RotL           = 0x77
//	iI64RotL           = 0x89
//	iI32RotR           = 0x78
//	iI64RotR           = 0x8a
//	iI32Clz            = 0x67
//	iI64Clz            = 0x79
//	iI32Ctz            = 0x68
//	iI64Ctz            = 0x7a
//	iI32PopCnt         = 0x69
//	iI64PopCnt         = 0x7b
//	iI32EqZ            = 0x45
//	iI64EqZ            = 0x50
//)
