package wasm_reader

import (
	"github.com/threadedstream/wasmexperiments/internal/pkg/reporter"
	"github.com/threadedstream/wasmexperiments/internal/pkg/utils"
	"io"
)

type WasmReader struct {
	readers []io.Reader
}

func NewWasmReader(r io.Reader) *WasmReader {
	wr := &WasmReader{
		readers: make([]io.Reader, 1),
	}
	wr.readers[0] = r
	return wr
}

func (wr *WasmReader) Push(val interface{}) {
	if r, ok := val.(io.Reader); ok {
		wr.readers = append(wr.readers, r)
		return
	}
	reporter.ReportError("expected the val to be of type io.Reader")
}

func (wr *WasmReader) Peek() interface{} {
	utils.Assert(!wr.Empty(), "peeking from an empty queue")
	l := len(wr.readers)
	return wr.readers[l-1]
}

func (wr *WasmReader) Pop() {
	utils.Assert(!wr.Empty(), "popping from an empty queue")
	l := len(wr.readers)
	wr.readers = wr.readers[:l-1]
}

func (wr *WasmReader) Empty() bool {
	return len(wr.readers) == 0
}

func (wr *WasmReader) ReadByte() (byte, error) {
	bs, err := wr.ReadBytes(1)
	if err != nil {
		return 0, err
	}
	return bs[0], nil
}

func (wr *WasmReader) ReadBytes(n int) ([]byte, error) {
	r := wr.Peek().(io.Reader)
	bs := make([]byte, n, n)
	if _, err := r.Read(bs); err != nil {
		return nil, err
	}
	return bs, nil
}
