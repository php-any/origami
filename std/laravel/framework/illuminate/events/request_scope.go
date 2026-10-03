package events

import "github.com/php-any/origami/data"

func init() {
	data.RegisterNativeRequestPolicy[*dispatcherState](data.NativeClone, cloneDispatcherState)
}

func cloneDispatcherState(scope *data.RequestObjectScope, source *dispatcherState) *dispatcherState {
	if source == nil {
		return nil
	}
	state := newDispatcherState()
	scope.RememberNativeState(source, state)
	source.mu.Lock()
	defer source.mu.Unlock()
	cloneListeners := func(listeners map[string][]data.Value) map[string][]data.Value {
		out := make(map[string][]data.Value, len(listeners))
		for name, list := range listeners {
			values := make([]data.Value, len(list))
			for i, value := range list {
				values[i] = scope.Bind(value)
			}
			out[name] = values
		}
		return out
	}
	state.listeners, state.wildcards = cloneListeners(source.listeners), cloneListeners(source.wildcards)
	state.queueResolver, state.txResolver = scope.Bind(source.queueResolver), scope.Bind(source.txResolver)
	return state
}

// ScopeRequest isolates native state as well as PHP properties. PHP clone alone
// copies the AnyValue pointer and would keep listeners/caches shared globally.
func ScopeRequest(ctx data.Context, original *data.ClassValue) *data.ClassValue {
	if provider, ok := ctx.GetVM().(data.RequestScopeProvider); ok {
		return provider.RequestObjectScope().Object(original)
	}
	_, native := original.Class.(*DispatcherClass)
	if !native && (original.ObjectValue == nil || !original.ObjectValue.HasProperty(dispatcherStateProp)) {
		return original.CloneSandbox(ctx)
	}
	scoped := original.CloneRequestScoped(ctx)
	var scope *data.RequestObjectScope
	if provider, ok := ctx.GetVM().(data.RequestScopeProvider); ok {
		scope = provider.RequestObjectScope()
		scope.Remember(original, scoped)
	}
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
				if scope != nil {
					values[i] = scope.Bind(value)
				}
				if fn, ok := value.(*data.FuncValue); ok {
					if wrapper, ok := fn.Value.(*listenerWrapper); ok {
						clone := *wrapper
						clone.owner = scoped
						if scope != nil {
							clone.listener = scope.Bind(wrapper.listener)
						}
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
	if scope != nil {
		state.queueResolver = scope.Bind(state.queueResolver)
		state.txResolver = scope.Bind(state.txResolver)
	}
	// Pushed/deferred events and wildcard results are request caches.
	_ = scoped.SetProperty(dispatcherStateProp, data.NewAnyValue(state))
	return scoped
}
