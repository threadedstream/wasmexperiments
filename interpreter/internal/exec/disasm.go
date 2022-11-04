package exec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	errInvalidOp = errors.New("invalid op")
)

type Instr struct {
	Op            Op
	Args          []any
	isUnreachable bool
}

func (i Instr) String() string {
	s := strings.Builder{}
	s.WriteString(i.Op.Name + " ")
	for _, arg := range i.Args {
		s.WriteString(fmt.Sprintf("%v,", arg))
	}
	return s.String()
}

func Dump(is []Instr) {
	s := strings.Builder{}
	for _, i := range is {
		s.WriteString(i.String())
		s.WriteRune('\n')
	}
	println(s.String())
}

func NewInstr(op Op, args []any, isUnreachable bool) Instr {
	return Instr{
		Op:            op,
		Args:          args,
		isUnreachable: isUnreachable,
	}
}

func Disassemble(code []byte) ([]Instr, error) {
	var out []Instr
	var err error
	var bytecode byte
	reader := bytes.NewReader(code)

	for err == nil {
		bytecode, err = reader.ReadByte()
		if err != nil {
			continue
		}
		op := lookupOp(Bytecode(bytecode))
		if !op.IsValid() {
			return nil, errInvalidOp
		}
		instr := Instr{
			Op: op,
		}
		switch instr.Op.Code {
		case i32AddOp:
			// done
			out = append(out, instr)
		case localGetOp:
			instr.Args = append(instr.Args, int(ignoreError(reader.ReadByte)))
			out = append(out, instr)
		case callOp:
			instr.Args = append(instr.Args, int(ignoreError(reader.ReadByte)))
			out = append(out, instr)
		case endOp:
			out = append(out, instr)
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
		switch i.Op.Code {
		case localGetOp:
			// does allocation happen if I use b[:]?
			byteStream.WriteByte(byte(localGetOp))
			var b [4]byte
			binaryFormat.PutUint32(b[:], uint32(i.Args[0].(int)))
			byteStream.Write(b[:])
		case i32AddOp:
			byteStream.WriteByte(byte(i32AddOp))
		case callOp:
			byteStream.WriteByte(byte(callOp))
			var b [4]byte
			binaryFormat.PutUint32(b[:], uint32(i.Args[0].(int)))
			byteStream.Write(b[:])
		case endOp:
			byteStream.WriteByte(byte(endOp))
		}
	}
	return byteStream.Bytes(), nil
}
