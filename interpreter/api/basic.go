package api

import (
	"github.com/threadedstream/wasmexperiments/internal/exec"
)

type WasmApi struct {
	vm *exec.VM
}

func NewWasmApi(path string) (*WasmApi, error) {
	api := new(WasmApi)
	mod, err := exec.NewModule(path)
	if err != nil {
		return nil, err
	}
	api.vm, err = exec.NewVM(mod)
	if err != nil {
		return nil, err
	}
	return api, nil
}

func (api *WasmApi) Call(name string, args ...uint64) (any, error) {
	// resolve function name
	index, err := api.vm.QueryFunction(name)
	if err != nil {
		return nil, err
	}
	return api.vm.ExecFunc(int64(index), args...)
}
