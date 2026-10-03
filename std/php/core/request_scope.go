package core

import (
	"fmt"
	"weak"

	"github.com/php-any/origami/data"
)

func init() {
	data.RegisterNativeRequestPolicy[*WeakMapClass](data.NativeClone, func(scope *data.RequestObjectScope, source *WeakMapClass) *WeakMapClass {
		if source == nil {
			return nil
		}
		clone := NewWeakMapClass()
		scope.RememberNativeState(source, clone)
		source.mu.RLock()
		defer source.mu.RUnlock()
		for name, value := range source.store {
			if key, exists := source.keys[name]; exists {
				object := key.object.Value()
				if object == nil {
					continue
				}
				target := scope.Object(&data.ClassValue{ObjectValue: object, Class: key.class})
				name = fmt.Sprintf("obj:%p", target.ObjectValue)
				clone.keys[name] = weakMapObjectKey{object: weak.Make(target.ObjectValue), class: key.class}
			}
			clone.store[name] = scope.BindSlot(value)
		}
		return clone
	})
}
