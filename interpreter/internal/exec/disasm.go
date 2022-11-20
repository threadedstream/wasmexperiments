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
		switch i.Op().Code {
		case localGetOp:
			ins := i.(*I32LocalGetI)
			// does allocation happen if I use b[:]?
			byteStream.WriteByte(byte(localGetOp))
			var b [4]byte
			binaryFormat.PutUint32(b[:], ins.arg0)
			byteStream.Write(b[:])
		case i32AddOp:
			byteStream.WriteByte(byte(i32AddOp))
		case callOp:
			byteStream.WriteByte(byte(callOp))
			var b [4]byte
			binaryFormat.PutUint32(b[:], uint32(i.Args[0].(int)))
			byteStream.Write(b[:])
		case i32LoadOp:
			byteStream.WriteByte(byte(i32LoadOp))
			var b [8]byte
			binaryFormat.PutUint32(b[0:4], uint32(i.Args[0].(int)))
			binaryFormat.PutUint32(b[4:8], uint32(i.Args[1].(int)))
			byteStream.Write(b[:])
		case i32StoreOp:
			byteStream.WriteByte(byte(i32StoreOp))
			var b [8]byte
			binaryFormat.PutUint32(b[0:4], uint32(i.Args[0].(int)))
			binaryFormat.PutUint32(b[4:8], uint32(i.Args[1].(int)))
			byteStream.Write(b[:])
		case i32ConstOp:
			byteStream.WriteByte(byte(i32ConstOp))
			var b [4]byte
			binaryFormat.PutUint32(b[:], uint32(i.Args[0].(int)))
			byteStream.Write(b[:])
		case endOp:
			byteStream.WriteByte(byte(endOp))
		}
	}
	return byteStream.Bytes(), nil
}
