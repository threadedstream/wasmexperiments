package exec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wasm_reader"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wbinary"
	"io"
)

var (
	errInvalidOp = errors.New("invalid op")
)

// Disassemble transforms code into sequence of instructions
func Disassemble(code []byte) ([]Instr, error) {
	reader := wasm_reader.NewWasmReader(bytes.NewReader(code))
	out, _, err := dis(reader, "")
	return out, err
}

func dis(reader *wasm_reader.WasmReader, context string) ([]Instr, Opcode, error) {
	var err error
	var bytecode byte
	var out []Instr
	var lastOp Opcode

	for err == nil {
		var in Instr
		bytecode, err = reader.ReadByte()
		if err != nil {
			continue
		}
		op := lookupOp(Opcode(bytecode))
		if !op.IsValid() {
			err = errInvalidOp
			continue
		}
		in, err = decodeIns(op, reader, context)
		if err != nil {
			continue
		}

		if op.Code == elseOp || op.Code == endOp {
			lastOp = op.Code
			break
		}
		out = append(out, in)
		lastOp = op.Code
	}
	if err != nil && err != io.EOF {
		return nil, lastOp, err
	}

	if (context == "if" || context == "block" || context == "loop") && lastOp != endOp && lastOp != elseOp {
		return nil, lastOp, errors.New("dis: if must be terminated by else of end instruction")
	}

	if lastOp != endOp && err == nil && context == "" {
		return nil, lastOp, errors.New("dis: bytecode stream must be terminated by end instruction")
	}

	return out, lastOp, nil
}

func disBlock(reader *wasm_reader.WasmReader) (*BlockI, error) {
	in := new(BlockI)
	in.commonI = commonI{op: lookupOp(blockOp)}
	err := in.resolveBlockType(reader)
	if err != nil {
		return nil, err
	}

	blockBody, _, err := dis(reader, "block")
	if err != nil {
		return nil, err
	}
	in.body = blockBody
	return in, nil
}

func disLoop(reader *wasm_reader.WasmReader) (*LoopI, error) {
	in := new(LoopI)
	in.commonI = commonI{op: lookupOp(loopOp)}
	err := in.resolveBlockType(reader)
	if err != nil {
		return nil, err
	}
	loopBody, _, err := dis(reader, "loop")
	if err != nil {
		return nil, err
	}
	in.body = loopBody
	return in, nil
}

func disIf(reader *wasm_reader.WasmReader) (*IfI, error) {
	in := new(IfI)
	in.commonI = commonI{op: lookupOp(ifOp)}
	err := in.resolveBlockType(reader)
	if err != nil {
		return nil, err
	}
	ifBody, lastOp, err := dis(reader, "if")
	if err != nil {
		return nil, err
	}

	in.body = ifBody
	if lastOp == elseOp {
		elseBody, _, err := dis(reader, "else")
		if err != nil {
			return nil, err
		}
		in.elseBody = elseBody
	}

	return in, nil
}

func decodeIns(op Op, reader *wasm_reader.WasmReader, context string) (Instr, error) {
	switch op.Code {
	default:
		return nil, fmt.Errorf("decodeIns: unknown instruction with opcode %v", op.Code)
	case i32AddOp, i32SubOp, i32MulOp, i32DivUOp, i32DivSOp:
		// we ain't got any operand stack yet
		return newDoubleArgI(op, nil, nil), nil
	case globalGetOp, localGetOp, localSetOp, globalSetOp, callOp:
		index, e := wbinary.ReadVarUint32(reader)
		if e != nil {
			return nil, e
		}
		return newSingleArgI(op, index), nil
	case i32LoadOp, i32StoreOp:
		var align, off uint32
		align, err := wbinary.ReadVarUint32(reader)
		if err != nil {
			return nil, err
		}
		off, err = wbinary.ReadVarUint32(reader)
		if err != nil {
			return nil, err
		}
		return newDoubleArgI(op, align, off), nil
	case i32ConstOp:
		imm, e := wbinary.ReadVarUint32(reader)
		if e != nil {
			return nil, e
		}
		return newSingleArgI(op, imm), nil
	case ifOp:
		ifI, err := disIf(reader)
		if err != nil {
			return nil, err
		}
		return ifI, nil
	case blockOp:
		blockI, err := disBlock(reader)
		if err != nil {
			return nil, err
		}
		return blockI, nil
	case loopOp:
		loopI, err := disLoop(reader)
		if err != nil {
			return nil, err
		}
		return loopI, nil
	case brOp:
		imm, e := wbinary.ReadVarUint32(reader)
		if e != nil {
			return nil, e
		}
		brI := newSingleArgI(op, imm).(*BrI)
		brI.context = context
		return brI, nil
	case brIfOp:
		imm, e := wbinary.ReadVarUint32(reader)
		if e != nil {
			return nil, e
		}
		brIfI := newSingleArgI(op, imm).(*BrIfI)
		brIfI.context = context
		return brIfI, nil
	case elseOp:
		// check if else is inside if
		if context != "if" {
			return nil, errors.New("else must be inside if")
		}
		// don't do anything with it, just return both nils
		return nil, nil
	case endOp, i32EqOp, returnOp, i32LtSOp:
		return newNoArgI(op), nil
	}
}

func Compile(code []byte) ([]byte, error) {
	binaryFormat := binary.LittleEndian
	// let the implementation allocate bytes on demand
	writer := bytes.NewBuffer(nil)
	reader := wasm_reader.NewWasmReader(bytes.NewBuffer(code))
	var err error
	var opcode byte
	for err == nil {
		opcode, err = reader.ReadByte()
		switch Opcode(opcode) {
		case i32LoadOp:
			writer.WriteByte(opcode)
			var align, off uint32
			align, err = wbinary.ReadVarUint32(reader)
			if err != nil {
				return nil, err
			}
			off, err = wbinary.ReadVarUint32(reader)
			if err != nil {
				return nil, err
			}
			var b [8]byte
			binaryFormat.PutUint32(b[0:4], align)
			binaryFormat.PutUint32(b[4:8], off)
			writer.Write(b[:])
		case localGetOp, globalGetOp, localSetOp, callOp, i32ConstOp, brIfOp, brOp:
			writer.WriteByte(opcode)
			var imm uint32
			imm, err = wbinary.ReadVarUint32(reader)
			if err != nil {
				continue
			}
			var b [4]byte
			binaryFormat.PutUint32(b[:], imm)
			writer.Write(b[:])
		case i32AddOp, i32SubOp, i32MulOp, i32DivUOp, i32DivSOp, i32EqOp, i32LtSOp, returnOp, endOp:
			writer.WriteByte(opcode)
		case blockOp, loopOp:
			writer.WriteByte(opcode)
			var imm uint8
			imm, err = wbinary.ReadVarUint7(reader)
			if err != nil {
				continue
			}
			writer.WriteByte(imm)
		}

	}
	//	case i32LoadOp, i32StoreOp:
	//		ins := i.(*doubleArgI)
	//		byteStream.WriteByte(byte(code))
	//		var b [8]byte
	//		binaryFormat.PutUint32(b[0:4], ins.arg0.(uint32))
	//		binaryFormat.PutUint32(b[4:8], ins.arg1.(uint32))
	//		byteStream.Write(b[:])
	//	case endOp:
	//		byteStream.WriteByte(byte(endOp))
	//	}
	//}
	return writer.Bytes(), nil
}
