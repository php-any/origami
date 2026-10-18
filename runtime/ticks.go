package runtime

import "github.com/php-any/origami/data"

func (vm *VM) PHPTickState() *data.TickState {
	if call := currentRequestCallState(); call != nil {
		return &call.ticks
	}
	return &vm.ticks
}

func (vm *RequestVM) PHPTickState() *data.TickState { return &vm.requestCall().ticks }
