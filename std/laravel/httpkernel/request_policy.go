package httpkernel

import "github.com/php-any/origami/data"

func init() {
	data.RegisterNativeRequestPolicy[*kernelState](data.NativeClone, func(scope *data.RequestObjectScope, source *kernelState) *kernelState {
		if source == nil {
			return nil
		}
		clone := *source
		scope.RememberNativeState(source, &clone)
		clone.app, clone.router = scope.Bind(source.app), scope.Bind(source.router)
		clone.bootstrappers, clone.middleware, clone.middlewarePriority = append([]string(nil), source.bootstrappers...), append([]string(nil), source.middleware...), append([]string(nil), source.middlewarePriority...)
		clone.middlewareGroups, clone.middlewareAliases = cloneGroupsMap(source.middlewareGroups), cloneStringMap(source.middlewareAliases)
		if source.tel != nil {
			clone.tel = &telescopeCache{disabled: source.tel.disabled, resolved: source.tel.resolved, patterns: append([]string(nil), source.tel.patterns...), only: append([]string(nil), source.tel.only...)}
			if source.tel.resolved {
				clone.tel.obj = scope.Bind(source.tel.obj)
				clone.tel.once.Do(func() {})
			}
		}
		return &clone
	})
}
