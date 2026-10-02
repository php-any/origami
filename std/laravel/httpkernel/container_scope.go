package httpkernel

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PHP services and retained callback objects receive request handles before
// execution. Native instance state still requires its explicit service policy.
func scopeContainerCallbacks(ctx data.Context, worker, request *data.ClassValue, objects map[*data.ObjectValue]*data.ClassValue) {
	if worker == nil || request == nil {
		return
	}
	objects[worker.ObjectValue] = request
	objects[request.ObjectValue] = request
	// Specialized services have already been scoped (including native state).
	// Reuse their final identities when rebinding captured factory receivers.
	if raw, ctl := worker.GetProperty("instances"); ctl == nil {
		if original, ok := raw.(*data.ArrayValue); ok {
			if raw, ctl := request.GetProperty("instances"); ctl == nil {
				if scoped, ok := raw.(*data.ArrayValue); ok {
					for slots, i := original.View(), 0; i < slots.Len(); i++ {
						slot := slots.At(i)
						from, ok := slot.Value.(*data.ClassValue)
						if !ok {
							continue
						}
						if next, ok := scoped.LookupZValByStringKey(slot.Name); ok {
							if to, ok := next.Value.(*data.ClassValue); ok && to.ObjectValue != from.ObjectValue {
								objects[from.ObjectValue] = to
								objects[to.ObjectValue] = to
							}
						}
					}
				}
			}
		}
	}
	var scopeObject func(*data.ClassValue) *data.ClassValue
	scopeObject = func(source *data.ClassValue) *data.ClassValue {
		if source == nil || source.ObjectValue == nil {
			return source
		}
		if target := objects[source.ObjectValue]; target != nil {
			return target
		}
		containerBound := false
		for _, name := range []string{"app", "container"} {
			if !source.ObjectValue.HasProperty(name) {
				continue
			}
			value, ctl := source.ObjectValue.GetProperty(name)
			if ctl != nil {
				continue
			}
			switch owner := value.(type) {
			case *data.ClassValue:
				containerBound = containerBound || owner.ObjectValue == worker.ObjectValue
			case *data.ThisValue:
				containerBound = containerBound || owner.ObjectValue == worker.ObjectValue
			}
		}
		_, phpObject := source.Class.(*node.ClassStatement)
		if !containerBound && !phpObject {
			return source
		}
		target := source.CloneRequestScoped(ctx)
		objects[source.ObjectValue] = target
		objects[target.ObjectValue] = target
		rebindContainer(ctx, target, request)
		return target
	}
	// Clone retained PHP instances once, preserving all aliases.
	if value, ctl := request.GetProperty("instances"); ctl == nil {
		if instances, ok := value.(*data.ArrayValue); ok {
			for slots, i := instances.View(), 0; i < slots.Len(); i++ {
				slot := slots.At(i)
				switch source := slot.Value.(type) {
				case *data.ClassValue:
					slot.Value = scopeObject(source)
				case *data.ThisValue:
					if target := scopeObject(source.ClassValue); target != source.ClassValue {
						slot.Value = data.NewThisValue(target)
					}
				}
			}
		}
	}
	captures := data.NewRequestCaptureScope(ctx, objects)
	seen := make(map[data.Value]data.Value)
	var bind func(data.Value) data.Value
	bind = func(value data.Value) data.Value {
		if value == nil {
			return nil
		}
		if previous, ok := seen[value]; ok {
			return previous
		}
		original := value
		seen[value] = value
		switch source := value.(type) {
		case *data.ClassValue:
			value = scopeObject(source)
		case *data.ThisValue:
			if target := scopeObject(source.ClassValue); target != source.ClassValue {
				value = data.NewThisValue(target)
			}
		case *data.FuncValue:
			if closure, ok := source.Value.(data.RequestClosureBinder); ok {
				for _, owner := range closure.RequestScopeObjects() {
					scopeObject(owner)
				}
				value = captures.Bind(source)
			}
		case *data.BoundFuncValue:
			scopeObject(source.BoundObject)
			if closure, ok := source.Value.(data.RequestClosureBinder); ok {
				for _, owner := range closure.RequestScopeObjects() {
					scopeObject(owner)
				}
				value = captures.Bind(source)
			}
		case *data.ArrayValue:
			var clone *data.ArrayValue
			for slots, i := source.View(), 0; i < slots.Len(); i++ {
				slot := slots.At(i)
				if slot == nil {
					continue
				}
				child := bind(slot.Value)
				if child != slot.Value {
					if clone == nil {
						clone = data.CloneArrayValue(source)
					}
					clone.ReplaceSlot(i, data.CopyZValKeepName(slot, child))
				}
			}
			if clone != nil {
				value = clone
			}
		}
		seen[original] = value
		return value
	}
	// PHP clone preserves a closure's receiver. Captured callbacks on cloned
	// managers must therefore be rebound explicitly to the same scoped identity.
	visited := make(map[*data.ObjectValue]bool)
	for _, target := range objects {
		if target == nil || target.ObjectValue == request.ObjectValue || visited[target.ObjectValue] {
			continue
		}
		visited[target.ObjectValue] = true
		target.ObjectValue.RangeProperties(func(name string, value data.Value) bool {
			if bound := bind(value); bound != value {
				_ = target.SetProperty(name, bound)
			}
			return true
		})
	}
	for _, name := range []string{"bindings", "methodBindings", "extenders", "reboundCallbacks", "contextual", "globalBeforeResolvingCallbacks", "beforeResolvingCallbacks", "globalResolvingCallbacks", "resolvingCallbacks", "globalAfterResolvingCallbacks", "afterResolvingCallbacks", "afterResolvingAttributeCallbacks"} {
		value, ctl := request.GetProperty(name)
		if ctl != nil || value == nil {
			continue
		}
		if bound := bind(value); bound != value {
			_ = request.SetProperty(name, bound)
		}
	}
}
