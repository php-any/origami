package runtime

import "github.com/php-any/origami/data"

// AddShutdownCallback 注册一个 shutdown 回调。
func (vm *VM) AddShutdownCallback(cb data.Value) {
	if st := currentRequestCallState(); st != nil {
		st.shutdown.add(cb)
		return
	}
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.shutdownCallbacks = append(vm.shutdownCallbacks, cb)
}

// RunShutdownCallbacks 依次执行所有已注册的 shutdown 回调（仅执行一次）。
func (vm *VM) RunShutdownCallbacks() {
	if st := currentRequestCallState(); st != nil {
		st.shutdown.run(vm)
		return
	}
	vm.shutdownRunOnce.Do(func() {
		defer func() {
			vm.mu.Lock()
			vm.shutdownCallbacks = nil
			vm.mu.Unlock()
		}()
		for i := 0; ; i++ {
			vm.mu.Lock()
			if i >= len(vm.shutdownCallbacks) {
				vm.mu.Unlock()
				break
			}
			cb := vm.shutdownCallbacks[i]
			vm.mu.Unlock()
			if ctl := callShutdownCallback(vm, cb); ctl != nil {
				if exit, ok := ctl.(data.ExitControl); ok && exit.IsExit() {
					return
				}
				vm.acl(ctl)
				return
			}
		}
		runHeaderCallbacks(vm)
	})
}

func callShutdownCallback(vm data.VM, cb data.Value) data.Control {
	switch c := cb.(type) {
	case *data.FuncValue:
		vars := c.Value.GetVariables()
		ctx := vm.CreateContext(vars)
		_, acl := c.Call(ctx)
		return acl
	case *data.BoundFuncValue:
		vars := c.FuncValue.Value.GetVariables()
		ctx := vm.CreateContext(vars)
		_, acl := c.Call(ctx)
		return acl
	}
	return nil
}

// A request owns its queue even when a startup closure still calls the base VM.
// Set done before invoking callbacks so reentrant RunShutdownCallbacks is safe;
// callbacks registered during shutdown are consumed in registration order.
type shutdownQueue struct {
	callbacks []data.Value
	done      bool
}

func (q *shutdownQueue) add(cb data.Value) { q.callbacks = append(q.callbacks, cb) }

func (q *shutdownQueue) run(vm data.VM) {
	if q.done {
		return
	}
	q.done = true
	defer func() { q.callbacks = nil }()
	for i := 0; i < len(q.callbacks); i++ {
		if ctl := callShutdownCallback(vm, q.callbacks[i]); ctl != nil {
			if exit, ok := ctl.(data.ExitControl); ok && exit.IsExit() {
				return
			}
			vm.ThrowControl(ctl)
			return
		}
	}
	runHeaderCallbacks(vm)
}
