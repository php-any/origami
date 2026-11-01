package runtime

import (
	"context"
	"sync"

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
	mu      sync.Mutex
	symbols map[data.SymbolID]*autoloadFlight
}

// 同 goroutine 的同符号递归直接返回；其他 goroutine 释放元数据锁后等待加载完成。
func (s *autoloadFlights) begin(ctx context.Context, id data.SymbolID) (bool, func(bool)) {
	owner := goid()
	s.mu.Lock()
	if s.symbols == nil {
		s.symbols = make(map[data.SymbolID]*autoloadFlight)
	}
	if pending := s.symbols[id]; pending != nil && pending.phase == autoloadLoading {
		if pending.owner == owner {
			s.mu.Unlock()
			return false, nil
		}
		s.mu.Unlock()
		select {
		case <-pending.done:
			return false, nil
		case <-ctx.Done():
			panic(data.ErrRequestCanceled)
		}
	}
	flight := &autoloadFlight{phase: autoloadLoading, owner: owner, done: make(chan struct{})}
	s.symbols[id] = flight
	s.mu.Unlock()
	return true, func(defined bool) {
		s.mu.Lock()
		defer s.mu.Unlock()
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
