package exec

import "github.com/threadedstream/wasmexperiments/internal/wasm"

const (
	wasmPageSize = 65536
)

type context struct {
	stack   []uint64
	locals  []uint64
	code    []byte
	pc      int64
	curFunc int64
}

type VM struct {
	ctx       context
	module    *wasm.Module
	memory    []byte
	funcs     []function
	funcTable [256]func()
}

func (vm *VM) PrintInstructionStream() (string, error) {
	for _, _ = range vm.module.CodeSection.Entries {
		// TODO
	}
	return "", nil
}
