package wbinary

import (
	"errors"
	"github.com/threadedstream/wasmexperiments/internal/pkg/reporter"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wasm_reader"
)

func readVarInt(wr *wasm_reader.WasmReader, n int) (int64, error) {
	if n > 64 {
		panic(errors.New("leb128: n must <= 64"))
	}
	var (
		p     byte
		res   int64
		shift uint
		err   error
	)
	for {
		p, err = wr.ReadByte()
		if err != nil {
			return 0, err
		}
		b := int64(p)
		switch {
		case b < 1<<6 && uint64(b) < uint64(1<<(n-1)):
			res += (1 << shift) * b
			return res, nil
		case b >= 1<<6 && b < 1<<7 && uint64(b)+1<<(n-1) >= 1<<7:
			res += (1 << shift) * (b - 1<<7)
			return res, nil
		case b >= 1<<7 && n > 7:
			res += (1 << shift) * (b - 1<<7)
			shift += 7
			n -= 7
		default:
			return 0, errors.New("leb128: invalid int")
		}
	}
}

func readVarUint(wr *wasm_reader.WasmReader, n int) (uint64, error) {
	if n > 64 {
		reporter.ReportError("readVarUint: n can't be larger than 64")
	}
	var (
		p     byte
		shift uint
		res   uint64
		err   error
	)
	for {
		p, err = wr.ReadByte()
		if err != nil {
			return 0, err
		}
		b := uint64(p)
		switch {
		default:
			return 0, errors.New("readVarUint: invalid uint")
		case b < 1<<7 && b <= 1<<n-1:
			res += (1 << shift) * b
			return res, nil
		case b >= 1<<7 && n > 7:
			res += (1 << shift) * (b - 1<<7)
			shift += 7
			n -= 7
		}
	}
}

func ReadVarUint7(wr *wasm_reader.WasmReader) (uint8, error) {
	val, err := readVarUint(wr, 7)
	return uint8(val), err
}

func ReadVarUint32(wr *wasm_reader.WasmReader) (uint32, error) {
	val, err := readVarUint(wr, 32)
	return uint32(val), err
}

func ReadVarUint64(reader *wasm_reader.WasmReader) (uint64, error) {
	return readVarUint(reader, 64)
}

func ReadVarInt32(wr *wasm_reader.WasmReader) (int32, error) {
	val, err := readVarInt(wr, 32)
	return int32(val), err
}

func ReadVarInt64(wr *wasm_reader.WasmReader) (int64, error) {
	return readVarInt(wr, 64)
}
