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

func dis(reader *wasm_reader.WasmReader, context string) ([]Instr, Bytecode, error) {
	var err error
	var bytecode byte
	var out []Instr
	var lastOp Bytecode

	for err == nil {
		var in Instr
		bytecode, err = reader.ReadByte()
		if err != nil {
			continue
		}
		op := lookupOp(Bytecode(bytecode))
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

	if context == "if" && lastOp != endOp && lastOp != elseOp {
		return nil, lastOp, errors.New("dis: if must be terminated by else of end instruction")
	}

	if lastOp != endOp && err == nil && context == "" {
		return nil, lastOp, errors.New("dis: bytecode stream must be terminated by end instruction")
	}

	return out, lastOp, nil
}

func disIf(reader *wasm_reader.WasmReader) (*IfI, error) {
	in := new(IfI)
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
	case elseOp:
		// check if else is inside if
		if context != "if" {
			return nil, errors.New("else must be inside if")
		}
		// don't do anything with it, just return both nils
		return nil, nil
	case endOp, i32EqOp, returnOp:
		return newNoArgI(op), nil
	}
}

// Compile is a reverse of Disassemble
func Compile(is []Instr) ([]byte, error) {
	binaryFormat := binary.LittleEndian
	// let the implementation allocate bytes on demand
	byteStream := bytes.NewBuffer(nil)
	for _, i := range is {
		switch code := i.Op().Code; code {
		case localGetOp, globalGetOp, callOp, i32ConstOp:
			ins := i.(*singleArgI)
			// does allocation happen if I use b[:]?
			byteStream.WriteByte(byte(code))
			var b [4]byte
			binaryFormat.PutUint32(b[:], ins.arg0.(uint32))
			byteStream.Write(b[:])
		case i32AddOp, i32SubOp, i32MulOp, i32DivUOp, i32DivSOp, i32EqOp:
			byteStream.WriteByte(byte(code))
		case i32LoadOp, i32StoreOp:
			ins := i.(*doubleArgI)
			byteStream.WriteByte(byte(code))
			var b [8]byte
			binaryFormat.PutUint32(b[0:4], ins.arg0.(uint32))
			binaryFormat.PutUint32(b[4:8], ins.arg1.(uint32))
			byteStream.Write(b[:])
		case endOp:
			byteStream.WriteByte(byte(endOp))
		}
	}
	return byteStream.Bytes(), nil
}
