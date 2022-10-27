package exec

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/threadedstream/wasmexperiments/internal/pkg/wasm_reader"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wbinary"
)

const (
	i32Const  byte = 0x41
	i64Const  byte = 0x42
	f32Const  byte = 0x43
	f64Const  byte = 0x44
	globalGet byte = 0x23
	end       byte = 0x0b
)

var ErrEmptyInitExpr = errors.New("wasm: initializer expression produces no value")

type InvalidInitExprOpError byte

func (e InvalidInitExprOpError) Error() string {
	return fmt.Sprintf("wasm: invalid opcode in initializer expression: %#x", byte(e))
}

type InvalidGlobalIndexError uint32

func (e InvalidGlobalIndexError) Error() string {
	return fmt.Sprintf("wasm: invalid index to global index space: %#x", uint32(e))
}

func readInitExpr(reader *wasm_reader.WasmReader) ([]byte, error) {
	var b [1]byte
	buf := new(bytes.Buffer)
	r := io.TeeReader(reader.Peek().(io.Reader), buf)
	reader.Push(r)
	defer reader.Pop()

outer:
	for {
		_, err := io.ReadFull(r, b[:])
		if err != nil {
			return nil, err
		}
		switch b[0] {
		default:
			return nil, InvalidInitExprOpError(b[0])
		case i32Const:
			if _, err = wbinary.ReadVarInt32(reader); err != nil {
				return nil, err
			}
		case i64Const:
			if _, err = wbinary.ReadVarInt64(reader); err != nil {
				return nil, err
			}
		case f32Const:
			if _, err = wbinary.ReadU32(reader); err != nil {
				return nil, err
			}
		case f64Const:
			if _, err = wbinary.ReadU64(reader); err != nil {
				return nil, err
			}
		case globalGet:
			if _, err = wbinary.ReadVarUint32(reader); err != nil {
				return nil, err
			}
		case end:
			break outer
		}
	}

	if buf.Len() == 0 {
		return nil, ErrEmptyInitExpr
	}

	return buf.Bytes(), nil
}
