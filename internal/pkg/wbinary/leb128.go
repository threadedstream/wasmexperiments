package wbinary

import (
	"errors"
	"github.com/threadedstream/wasmexperiments/internal/pkg/reporter"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wasm_reader"
)

func readVarUint(wr *wasm_reader.WasmReader, n int) (uint64, error) {
	if n > 64 {
		reporter.ReportError("readVarUint: n can't be larger than 64")
	}
	var p byte
	var shift uint
	var res uint64
	for {
		p, _ = wr.ReadByte()
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

func ReadVarUint32(wr *wasm_reader.WasmReader) (uint32, error) {
	val, err := readVarUint(wr, 32)
	return uint32(val), err
}

func ReadVarUint64(reader *wasm_reader.WasmReader) (uint64, error) {
	return readVarUint(reader, 64)
}
