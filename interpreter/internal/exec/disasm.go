package exec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wasm_reader"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wbinary"
	"io"
)

var (
	errInvalidOp = errors.New("invalid op")
)

func Disassemble(code []byte) ([]Instr, error) {
	var out []Instr
	var err error
	var bytecode byte
	var lastOp Bytecode
	reader := wasm_reader.NewWasmReader(bytes.NewReader(code))

	for err == nil {
		bytecode, err = reader.ReadByte()
		if err != nil {
			continue
		}
		op := lookupOp(Bytecode(bytecode))
		if !op.IsValid() {
			return nil, errInvalidOp
		}
		switch op.Code {
		case i32AddOp, i32SubOp, i32MulOp, i32DivUOp, i32DivSOp, i32EqOp:
			// we ain't got any operand stack yet
			ins := newDoubleArgI(op, nil, nil)
			out = append(out, ins)
		case globalGetOp, localGetOp, localSetOp, globalSetOp, callOp:
			index, e := wbinary.ReadVarUint32(reader)
			if e != nil {
				return nil, e
			}
			ins := newSingleArgI(op, index)
			out = append(out, ins)
		case i32LoadOp, i32StoreOp:
			alignment, e := wbinary.ReadVarUint32(reader)
			if e != nil {
				return nil, err
			}
			off, e := wbinary.ReadVarUint32(reader)
			if e != nil {
				return nil, err
			}
			ins := newDoubleArgI(op, alignment, off)
			out = append(out, ins)
		case i32ConstOp:
			imm, e := wbinary.ReadVarUint32(reader)
			if e != nil {
				return nil, e
			}
			ins := newSingleArgI(op, imm)
			out = append(out, ins)
		case endOp:
			// do nothing
		}
		lastOp = op.Code
	}

	if lastOp != endOp {
		return nil, errors.New("disasm: bytecode stream must be terminated by end instruction")
	}

	if err != nil {
		if err == io.EOF {
			return out, nil
		} else {
			return nil, err
		}
	}
	return out, nil
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
