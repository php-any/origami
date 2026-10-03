package cookie

import "github.com/php-any/origami/data"

func init() {
	data.RegisterNativeRequestPolicy[*jarState](data.NativeClone, func(scope *data.RequestObjectScope, source *jarState) *jarState {
		if source == nil {
			return nil
		}
		clone := *source
		scope.RememberNativeState(source, &clone)
		clone.domain, clone.secure = scope.Bind(source.domain), scope.Bind(source.secure)
		clone.queued = make(map[string]map[string]*data.ClassValue, len(source.queued))
		for name, paths := range source.queued {
			cookies := make(map[string]*data.ClassValue, len(paths))
			for path, cookie := range paths {
				cookies[path] = scope.Object(cookie)
			}
			clone.queued[name] = cookies
		}
		return &clone
	})
}
