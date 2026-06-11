package container

import (
	"sync"

	"github.com/php-any/origami/data"
)

var (
	registeringMu     sync.Mutex
	registeringEngine *Engine
)

func setRegisteringEngine(e *Engine) func() {
	registeringMu.Lock()
	prev := registeringEngine
	registeringEngine = e
	registeringMu.Unlock()
	return func() {
		registeringMu.Lock()
		registeringEngine = prev
		registeringMu.Unlock()
	}
}

func activeEngine(ctx data.Context) *Engine {
	registeringMu.Lock()
	e := registeringEngine
	registeringMu.Unlock()
	if e != nil {
		return e
	}
	if engine, acl := engineFromCtx(ctx); acl == nil && engine != nil {
		return engine
	}
	return nil
}
