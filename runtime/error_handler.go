package runtime

import "github.com/php-any/origami/data"

type errorHandlerState struct {
	handlers []data.Value
	masks    []int
}

func (s *errorHandlerState) current() data.Value {
	if n := len(s.handlers); n > 0 {
		return s.handlers[n-1]
	}
	return nil
}
func (s *errorHandlerState) set(handler data.Value) data.Value {
	return s.setMask(handler, 32767)
}
func (s *errorHandlerState) setMask(handler data.Value, mask int) data.Value {
	old := s.current()
	s.handlers = append(s.handlers, handler)
	s.masks = append(s.masks, mask)
	return old
}
func (s *errorHandlerState) restore() bool {
	if n := len(s.handlers); n > 0 {
		s.handlers = s.handlers[:n-1]
		s.masks = s.masks[:n-1]
	}
	return true
}
func (vm *VM) snapshotErrorHandlers() *errorHandlerState {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return &errorHandlerState{handlers: append([]data.Value(nil), vm.errorHandlers...), masks: append([]int(nil), vm.errorHandlerMasks...)}
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
	return vm.SetErrorHandlerMask(handler, 32767)
}
func (vm *VM) SetErrorHandlerMask(handler data.Value, mask int) data.Value {
	if state := vm.requestErrorHandlers(); state != nil {
		return state.setMask(handler, mask)
	}
	vm.mu.Lock()
	defer vm.mu.Unlock()
	state := errorHandlerState{handlers: vm.errorHandlers, masks: vm.errorHandlerMasks}
	old := state.setMask(handler, mask)
	vm.errorHandlers = state.handlers
	vm.errorHandlerMasks = state.masks
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
	state := errorHandlerState{handlers: vm.errorHandlers, masks: vm.errorHandlerMasks}
	state.restore()
	vm.errorHandlers = state.handlers
	vm.errorHandlerMasks = state.masks
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
func (vm *RequestVM) SetErrorHandlerMask(handler data.Value, mask int) data.Value {
	return vm.requestErrorHandlers().setMask(handler, mask)
}
func (s *errorHandlerState) forLevel(level int) data.Value {
	if n := len(s.handlers); n != 0 && (len(s.masks) < n || s.masks[n-1]&level != 0) {
		return s.current()
	}
	return nil
}
func (vm *VM) GetErrorHandlerFor(level int) data.Value {
	if state := vm.requestErrorHandlers(); state != nil {
		return state.forLevel(level)
	}
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return (&errorHandlerState{handlers: vm.errorHandlers, masks: vm.errorHandlerMasks}).forLevel(level)
}
func (vm *RequestVM) GetErrorHandlerFor(level int) data.Value {
	return vm.requestErrorHandlers().forLevel(level)
}
func (vm *RequestVM) GetErrorHandler() data.Value { return vm.requestErrorHandlers().current() }
func (vm *RequestVM) RestoreErrorHandler() bool   { return vm.requestErrorHandlers().restore() }
