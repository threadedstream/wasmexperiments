package wbinary

import (
	"encoding/binary"
	"errors"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wasm_reader"
	"unicode/utf8"
)

func ReadU64(wr *wasm_reader.WasmReader) (uint64, error) {
	val, err := wr.ReadBytes(8)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(val), nil
}

func ReadU32(wr *wasm_reader.WasmReader) (uint32, error) {
	val, err := wr.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(val), nil
}

func ReadU8(wr *wasm_reader.WasmReader) (uint8, error) {
	val, err := wr.ReadByte()
	if err != nil {
		return 0, err
	}
	return val, nil
}

func readUTF8String(reader *wasm_reader.WasmReader, n uint32) (string, error) {
	bs, err := reader.ReadBytes(int(n))
	if err != nil {
		return "", err
	}
	if !utf8.Valid(bs) {
		return "", errors.New("invalid utf8 string")
	}
	return string(bs), nil
}

func ReadUTF8StringUint(reader *wasm_reader.WasmReader) (string, error) {
	n, err := ReadVarUint32(reader)
	if err != nil {
		return "", err
	}
	return readUTF8String(reader, n)
}
