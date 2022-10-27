package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/threadedstream/wasmexperiments/internal/exec"
	"github.com/threadedstream/wasmexperiments/internal/pkg/reporter"
	"github.com/threadedstream/wasmexperiments/internal/pkg/wasm_reader"
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
	module := exec.NewModule(wr)
	if err := module.Read(); err != nil {
		reporter.ReportError(err.Error())
	}
	executor, err := exec.NewVM(module)
	fmt.Printf("executor's address: %p", executor)
}
