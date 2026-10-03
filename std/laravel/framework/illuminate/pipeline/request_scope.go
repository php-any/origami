package pipeline

import "github.com/php-any/origami/data"

func init() {
	data.RegisterNativeRequestPolicy[*pipeState](data.NativeClone, func(scope *data.RequestObjectScope, source *pipeState) *pipeState {
		if source == nil {
			return nil
		}
		clone := *source
		scope.RememberNativeState(source, &clone)
		clone.passable, clone.finallyCb, clone.withinTransaction = scope.Bind(source.passable), scope.Bind(source.finallyCb), scope.Bind(source.withinTransaction)
		clone.container = scope.Object(source.container)
		clone.pipes = make([]data.Value, len(source.pipes))
		for i, pipe := range source.pipes {
			clone.pipes[i] = scope.Bind(pipe)
		}
		return &clone
	})
}
