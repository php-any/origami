package runtime

import "github.com/php-any/origami/data"

type exceptionHandlerBinding struct{ original, callable data.Value }

// ExceptionHandlerState is owned by a single PHP request. Copies preserve the
// startup registrations while keeping later set/restore and recursion local.
type ExceptionHandlerState struct {
	current  exceptionHandlerBinding
	previous []exceptionHandlerBinding
	running  bool
}

func (s *ExceptionHandlerState) Set(original, callable data.Value) data.Value {
	old := s.current.original
	s.previous = append(s.previous, s.current)
	s.current = exceptionHandlerBinding{original, callable}
	return old
}
func (s *ExceptionHandlerState) Current() data.Value { return s.current.original }
func (s *ExceptionHandlerState) Restore() bool {
	if n := len(s.previous); n != 0 {
		s.current = s.previous[n-1]
		s.previous = s.previous[:n-1]
	}
	return true
}
func (s *ExceptionHandlerState) Clone() *ExceptionHandlerState {
	return &ExceptionHandlerState{current: s.current, previous: append([]exceptionHandlerBinding(nil), s.previous...)}
}

// Handle is an uncaught-exception boundary, never an operation in PHP Call().
// A failure from the handler replaces the original control and is not rehandled.
func (s *ExceptionHandlerState) Handle(vm data.VM, control data.Control) (bool, data.Control) {
	thrown, ok := control.(*data.ThrowValue)
	if !ok || thrown == nil || thrown.PHPCompileFatal || s.running || s.current.callable == nil {
		return false, control
	}
	var function data.FuncStmt
	var invoke func(data.Context) (data.GetValue, data.Control)
	switch callback := s.current.callable.(type) {
	case *data.FuncValue:
		function, invoke = callback.Value, callback.Call
	case *data.BoundFuncValue:
		function, invoke = callback.Value, callback.Call
	default:
		return false, control
	}
	s.running = true
	defer func() { s.running = false }()
	ctx := vm.CreateContext(function.GetVariables())
	ctx.SetStrictTypes(false)
	if ctl := data.BindDeclaredArgs(ctx, function, []data.Value{thrown.PHPValue()}); ctl != nil {
		return true, ctl
	}
	_, ctl := invoke(ctx)
	return true, ctl
}

func (vm *VM) SnapshotExceptionHandlers() *ExceptionHandlerState {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.exceptionHandlers.Clone()
}

func (vm *VM) requestExceptionHandlers() *ExceptionHandlerState {
	if st := currentRequestCallState(); st != nil {
		if st.exceptionHandlers == nil {
			st.exceptionHandlers = vm.SnapshotExceptionHandlers()
		}
		return st.exceptionHandlers
	}
	return nil
}

func (vm *VM) SetExceptionHandler(handler data.Value) data.Value {
	return vm.SetExceptionHandlerBinding(handler, handler)
}
func (vm *VM) SetExceptionHandlerBinding(original, callable data.Value) data.Value {
	if state := vm.requestExceptionHandlers(); state != nil {
		return state.Set(original, callable)
	}
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.exceptionHandlers.Set(original, callable)
}
func (vm *VM) GetExceptionHandler() data.Value {
	if state := vm.requestExceptionHandlers(); state != nil {
		return state.Current()
	}
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.exceptionHandlers.Current()
}
func (vm *VM) RestoreExceptionHandler() bool {
	if state := vm.requestExceptionHandlers(); state != nil {
		return state.Restore()
	}
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.exceptionHandlers.Restore()
}
func (vm *VM) HandleUnhandledException(control data.Control) (bool, data.Control) {
	if state := vm.requestExceptionHandlers(); state != nil {
		return state.Handle(vm, control)
	}
	return vm.exceptionHandlers.Handle(vm, control)
}

func (vm *RequestVM) requestExceptionHandlers() *ExceptionHandlerState {
	if state := vm.requestCall(); state != &vm.call {
		if state.exceptionHandlers == nil {
			state.exceptionHandlers = vm.Base.SnapshotExceptionHandlers()
		}
		return state.exceptionHandlers
	}
	if vm.exceptionHandlers == nil {
		vm.exceptionHandlers = vm.Base.SnapshotExceptionHandlers()
	}
	return vm.exceptionHandlers
}
func (vm *RequestVM) SetExceptionHandler(handler data.Value) data.Value {
	return vm.SetExceptionHandlerBinding(handler, handler)
}
func (vm *RequestVM) SetExceptionHandlerBinding(original, callable data.Value) data.Value {
	return vm.requestExceptionHandlers().Set(original, callable)
}
func (vm *RequestVM) GetExceptionHandler() data.Value { return vm.requestExceptionHandlers().Current() }
func (vm *RequestVM) RestoreExceptionHandler() bool   { return vm.requestExceptionHandlers().Restore() }
func (vm *RequestVM) HandleUnhandledException(control data.Control) (bool, data.Control) {
	return vm.requestExceptionHandlers().Handle(vm, control)
}
