package main

import (
	"bytes"
	"flag"
	"github.com/threadedstream/wasmexperiments/internal/pkg/reporter"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wasm_reader"
	"github.com/threadedstream/wasmexperiments/internal/wasm"
	"os"
)

func main() {
	var path = flag.String("path", "", "path to a wasm binary")
	flag.Parse()
	bs, err := os.ReadFile(*path)
	if err != nil {
		reporter.ReportError(err.Error())
	}
	r := bytes.NewReader(bs)
	wr := wasm_reader.NewWasmReader(r)
	module := wasm.NewModule(wr)
	if err := module.Read(); err != nil {
		reporter.ReportError(err.Error())
	}
}
