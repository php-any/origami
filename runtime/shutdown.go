package runtime

import "github.com/php-any/origami/data"

// AddShutdownCallback 注册一个 shutdown 回调。
func (vm *VM) AddShutdownCallback(cb data.Value) {
	vm.shutdownCallbacks = append(vm.shutdownCallbacks, cb)
}

// RunShutdownCallbacks 依次执行所有已注册的 shutdown 回调（仅执行一次）。
func (vm *VM) RunShutdownCallbacks() {
	vm.shutdownRunOnce.Do(func() {
		for _, cb := range vm.shutdownCallbacks {
			callShutdownCallback(vm, cb)
		}
		runHeaderCallbacks(vm)
	})
}

func callShutdownCallback(vm *VM, cb data.Value) {
	switch c := cb.(type) {
	case *data.FuncValue:
		vars := c.Value.GetVariables()
		ctx := vm.CreateContext(vars)
		if _, acl := c.Call(ctx); acl != nil {
			vm.acl(acl)
		}
	case *data.BoundFuncValue:
		vars := c.FuncValue.Value.GetVariables()
		ctx := vm.CreateContext(vars)
		if _, acl := c.Call(ctx); acl != nil {
			vm.acl(acl)
		}
	}
}
