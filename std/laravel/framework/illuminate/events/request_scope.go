package events

import "github.com/php-any/origami/data"

// ScopeRequest isolates native state as well as PHP properties. PHP clone alone
// copies the AnyValue pointer and would keep listeners/caches shared globally.
func ScopeRequest(ctx data.Context, original *data.ClassValue) *data.ClassValue {
	_, native := original.Class.(*DispatcherClass)
	if !native && (original.ObjectValue == nil || !original.ObjectValue.HasProperty(dispatcherStateProp)) {
		return original.CloneSandbox(ctx)
	}
	scoped := original.CloneRequestScoped(ctx)
	source := stateOf(original)
	source.mu.Lock()
	defer source.mu.Unlock()
	state := newDispatcherState()
	cloneListeners := func(source map[string][]data.Value) map[string][]data.Value {
		out := make(map[string][]data.Value, len(source))
		for name, list := range source {
			values := make([]data.Value, len(list))
			for i, value := range list {
				values[i] = value
				if fn, ok := value.(*data.FuncValue); ok {
					if wrapper, ok := fn.Value.(*listenerWrapper); ok {
						clone := *wrapper
						clone.owner = scoped
						values[i] = data.NewFuncValue(&clone)
					}
				}
			}
			out[name] = values
		}
		return out
	}
	state.listeners = cloneListeners(source.listeners)
	state.wildcards = cloneListeners(source.wildcards)
	state.queueResolver, state.txResolver = source.queueResolver, source.txResolver
	// Pushed/deferred events and wildcard results are request caches.
	_ = scoped.SetProperty(dispatcherStateProp, data.NewAnyValue(state))
	return scoped
}
