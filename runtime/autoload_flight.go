package runtime

import (
	"context"
	"github.com/php-any/origami/data"
	"sync"
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

// Recursive probes of the same symbol decline immediately. Other callers
// wait for a declaration attempt without holding registry or callback locks.
func (s *autoloadFlights) begin(ctx context.Context, id data.SymbolID) (bool, func(bool)) {
	owner := goid()
	for {
		s.mu.Lock()
		if s.symbols == nil {
			s.symbols = make(map[data.SymbolID]*autoloadFlight)
		}
		if pending := s.symbols[id]; pending != nil && pending.phase == autoloadLoading {
			if pending.owner == owner {
				s.mu.Unlock()
				return false, nil
			}
			done := pending.done
			s.mu.Unlock()
			select {
			case <-done:
				return false, nil
			case <-ctx.Done():
				panic(data.ErrRequestCanceled)
			}
		}
		flight := &autoloadFlight{phase: autoloadLoading, owner: owner, done: make(chan struct{})}
		s.symbols[id] = flight
		s.mu.Unlock()
		return true, finishAutoload(s, flight)
	}
}

func finishAutoload(s *autoloadFlights, flight *autoloadFlight) func(bool) {
	var once sync.Once
	return func(defined bool) {
		once.Do(func() {
			s.mu.Lock()
			defer s.mu.Unlock()
			if defined {
				flight.phase = autoloadDefined
			} else {
				flight.phase = autoloadFailed
			}
			close(flight.done)
		})
	}
}
func (vm *VM) BeginAutoload(ctx data.Context, id data.SymbolID) (bool, func(bool)) {
	return vm.autoloadFlights.begin(ctx.GoContext(), id)
}
func (vm *RequestVM) BeginAutoload(ctx data.Context, id data.SymbolID) (bool, func(bool)) {
	return vm.autoloadFlights.begin(ctx.GoContext(), id)
}
