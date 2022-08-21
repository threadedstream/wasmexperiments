package exec

// branch instructions
var instructionMap = map[byte]Instruction{
	0x02: iBlock,
	0x03: iLoop,
	0x0c: iBr,
	0x0d: iBrIf,
	0x0e: iBrTable,
	0x04: iIf,
	0x05: iElse,
	0x0b: iEnd,
	0x0f: iRet,
	0x00: iUnreachable,
	0x01: iNop,
	0x1a: iDrop,
	0x41: iI32Const,
	0x42: iI64Const,
	0x43: iF32Const,
	0x44: iF64Const,
	0x20: iLocalGet,
	0x21: iLocalSet,
	0x22: iLocalTee,
	0x23: iGlobalGet,
	0x24: iGlobalSet,
	0x1b: iSelect,
	0x10: iCall,
	0x11: iCallIndirect,
	0x6a: iI32Add,
	0x7c: iI64Add,
	0x6b: iI32Sub,
	0x7d: iI64Sub,
	0x6c: iI32Mul,
	0x7e: iI64Mul,
	0x6d: iI32DivS,
	0x7f: iI64SivS,
	0x6e: iI32DivU,
	0x80: iI64DivU,
	0x6f: iI32RemS,
	0x81: iI64RemS,
	0x70: iI32RemU,
	0x82: iI64RemU,
	0x71: iI32And,
	0x83: iI64And,
	0x72: iI32Or,
	0x84: iI64Or,
	0x73: iI32Xor,
	0x85: iI64Xor,
	0x74: iI32Shl,
	0x86: iI64Shl,
	0x75: iI32ShrS,
	0x87: iI64ShrS,
	0x76: iI32ShrU,
	0x88: iI64ShrU,
	0x77: iI32RotL,
	0x89: iI64RotL,
	0x78: iI32RotR,
	0x8a: iI64RotR,
	0x67: iI32Clz,
	0x79: iI64Clz,
	0x68: iI32Ctz,
	0x7a: iI64Ctz,
	0x69: iI32PopCnt,
	0x7b: iI64PopCnt,
	0x45: iI32EqZ,
	0x50: iI64EqZ,
}

type Instruction uint8

//go:generate stringer -type Instruction -linecomment instruction.go
const (
	iBlock Instruction = iota
	iLoop
	iBr
	iBrIf
	iBrTable
	iIf
	iElse
	iEnd
	iRet
	iUnreachable
	iNop
	iDrop
	iI32Const
	iI64Const
	iF32Const
	iF64Const
	iLocalGet
	iLocalSet
	iLocalTee
	iGlobalGet
	iGlobalSet
	iSelect
	iCall
	iCallIndirect
	iI32Add
	iI64Add
	iI32Sub
	iI64Sub
	iI32Mul
	iI64Mul
	iI32DivS
	iI64SivS
	iI32DivU
	iI64DivU
	iI32RemS
	iI64RemS
	iI32RemU
	iI64RemU
	iI32And
	iI64And
	iI32Or
	iI64Or
	iI32Xor
	iI64Xor
	iI32Shl
	iI64Shl
	iI32ShrS
	iI64ShrS
	iI32ShrU
	iI64ShrU
	iI32RotL
	iI64RotL
	iI32RotR
	iI64RotR
	iI32Clz
	iI64Clz
	iI32Ctz
	iI64Ctz
	iI32PopCnt
	iI64PopCnt
	iI32EqZ
	iI64EqZ
)
