package container

import "github.com/php-any/origami/data"

func init() {
	data.RegisterNativeRequestPolicy[*ctnState](data.NativeClone, func(scope *data.RequestObjectScope, source *ctnState) *ctnState {
		if source == nil {
			return nil
		}
		clone := newCtnState()
		scope.RememberNativeState(source, clone)
		source.mu.Lock()
		defer source.mu.Unlock()
		for name, binding := range source.bindings {
			binding.concrete = scope.Bind(binding.concrete)
			clone.bindings[name] = binding
		}
		for name, value := range source.instances {
			clone.instances[name] = scope.Bind(value)
		}
		for name, alias := range source.aliases {
			clone.aliases[name] = alias
		}
		for name, resolved := range source.resolved {
			clone.resolved[name] = resolved
		}
		clone.buildStack = append([]string(nil), source.buildStack...)
		return clone
	})
}
