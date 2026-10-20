package runtime

import (
	"context"

	"github.com/php-any/origami/data"
)

type autoloadPhase uint8

const (
	autoloadLoading autoloadPhase = iota + 1
	autoloadDefined
	autoloadFailed
)

type autoloadFlight struct {
	phase autoloadPhase
	owner uint64
	done  chan struct{}
}
type autoloadFlights struct {
	symbols map[data.SymbolID]*autoloadFlight
}

// 同一个请求只有一个执行线程。同符号递归直接返回；不会有另一个线程在等这次加载。
func (s *autoloadFlights) begin(ctx context.Context, id data.SymbolID) (bool, func(bool)) {
	owner := goid()
	if s.symbols == nil {
		s.symbols = make(map[data.SymbolID]*autoloadFlight)
	}
	if pending := s.symbols[id]; pending != nil && pending.phase == autoloadLoading {
		if pending.owner == owner {
			return false, nil
		}
		select {
		case <-pending.done:
			return false, nil
		case <-ctx.Done():
			panic(data.ErrRequestCanceled)
		}
	}
	flight := &autoloadFlight{phase: autoloadLoading, owner: owner, done: make(chan struct{})}
	s.symbols[id] = flight
	return true, func(defined bool) {
		if flight.phase != autoloadLoading {
			return
		}
		if defined {
			flight.phase = autoloadDefined
		} else {
			flight.phase = autoloadFailed
		}
		close(flight.done)
	}
}
func (vm *VM) BeginAutoload(ctx data.Context, id data.SymbolID) (bool, func(bool)) {
	return vm.autoloadFlights.begin(ctx.GoContext(), id)
}
func (vm *RequestVM) BeginAutoload(ctx data.Context, id data.SymbolID) (bool, func(bool)) {
	return vm.autoloadFlights.begin(ctx.GoContext(), id)
}
