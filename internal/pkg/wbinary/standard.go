package wbinary

import (
	"encoding/binary"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wasm_reader"
)

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
