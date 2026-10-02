package runtime

import "github.com/php-any/origami/data"

type errorHandlerState struct{ handlers []data.Value }

func (s *errorHandlerState) current() data.Value {
	if n := len(s.handlers); n > 0 {
		return s.handlers[n-1]
	}
	return nil
}
func (s *errorHandlerState) set(handler data.Value) data.Value {
	old := s.current()
	s.handlers = append(s.handlers, handler)
	return old
}
func (s *errorHandlerState) restore() bool {
	if n := len(s.handlers); n > 0 {
		s.handlers = s.handlers[:n-1]
	}
	return true
}
func (vm *VM) snapshotErrorHandlers() *errorHandlerState {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return &errorHandlerState{handlers: append([]data.Value(nil), vm.errorHandlers...)}
}
func (vm *VM) requestErrorHandlers() *errorHandlerState {
	if state := currentRequestCallState(); state != nil {
		if state.errorHandlers == nil {
			state.errorHandlers = vm.snapshotErrorHandlers()
		}
		return state.errorHandlers
	}
	return nil
}
func (vm *VM) SetErrorHandler(handler data.Value) data.Value {
	if state := vm.requestErrorHandlers(); state != nil {
		return state.set(handler)
	}
	vm.mu.Lock()
	defer vm.mu.Unlock()
	state := errorHandlerState{handlers: vm.errorHandlers}
	old := state.set(handler)
	vm.errorHandlers = state.handlers
	return old
}
func (vm *VM) GetErrorHandler() data.Value {
	if state := vm.requestErrorHandlers(); state != nil {
		return state.current()
	}
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return (&errorHandlerState{handlers: vm.errorHandlers}).current()
}
func (vm *VM) RestoreErrorHandler() bool {
	if state := vm.requestErrorHandlers(); state != nil {
		return state.restore()
	}
	vm.mu.Lock()
	defer vm.mu.Unlock()
	state := errorHandlerState{handlers: vm.errorHandlers}
	state.restore()
	vm.errorHandlers = state.handlers
	return true
}
func (vm *RequestVM) requestErrorHandlers() *errorHandlerState {
	if state := vm.requestCall(); state != &vm.call {
		if state.errorHandlers == nil {
			state.errorHandlers = vm.Base.snapshotErrorHandlers()
		}
		return state.errorHandlers
	}
	if vm.errorHandlers == nil {
		vm.errorHandlers = vm.Base.snapshotErrorHandlers()
	}
	return vm.errorHandlers
}
func (vm *RequestVM) SetErrorHandler(handler data.Value) data.Value {
	return vm.requestErrorHandlers().set(handler)
}
func (vm *RequestVM) GetErrorHandler() data.Value { return vm.requestErrorHandlers().current() }
func (vm *RequestVM) RestoreErrorHandler() bool   { return vm.requestErrorHandlers().restore() }
