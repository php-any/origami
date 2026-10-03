package exception

import "github.com/php-any/origami/data"

func init() {
	data.RegisterNativeRequestPolicy[*Exception](data.NativeClone, func(scope *data.RequestObjectScope, source *Exception) *Exception {
		if source == nil {
			return nil
		}
		copy := *source
		copy.trace = append([]traceFrame(nil), source.trace...)
		return &copy
	})
}
