package exec

import (
	"fmt"
)

func (vm *VM) QueryFunction(name string) (uint32, error) {
	if index, ok := vm.funcMap[name]; ok {
		return index, nil
	}
	return 0, fmt.Errorf("no index associated with function '%s'", name)
}

func (vm *VM) PushValuesToStack(values ...uint64) {
	for _, val := range values {
		vm.pushUint64(val)
	}
}
