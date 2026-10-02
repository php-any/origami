package httpkernel

import (
	"strings"

	"github.com/php-any/origami/data"
)

// Livewire's flush-state callbacks retain mechanism instances as closure $this.
// Resetting their global objects is unsafe in a concurrent worker. Scope both
// the registered mechanisms and their EventBus callbacks before flush-state.
func scopeLivewire(ctx data.Context, app data.Value) {
	container, ok := app.(*data.ClassValue)
	if !ok {
		return
	}
	raw, ctl := container.GetProperty("instances")
	if ctl != nil {
		return
	}
	instances, ok := raw.(*data.ArrayValue)
	if !ok {
		return
	}
	objects := make(map[*data.ObjectValue]*data.ClassValue)
	for arraySlots74, arrayPosition74 := instances.View(), 0; arrayPosition74 < arraySlots74.Len(); arrayPosition74++ {
		slot := arraySlots74.At(arrayPosition74)
		if slot == nil {
			continue
		}
		source, ok := slot.Value.(*data.ClassValue)
		if !ok || !strings.HasPrefix(source.GetName(), "Livewire\\") {
			continue
		}
		scoped := objects[source.ObjectValue]
		if scoped == nil {
			scoped = source.CloneRequestScoped(ctx)
			objects[source.ObjectValue] = scoped
			rebindContainer(ctx, scoped, app)
		}
		slot.Value = scoped
	}
	if len(objects) == 0 {
		return
	}
	// Mechanisms can be retained only by event closures (their registration
	// instance is not necessarily the object used by boot()). Include these
	// handles before rebinding callbacks, and publish them in this container.
	var discover func(data.Value)
	discover = func(value data.Value) {
		switch v := value.(type) {
		case *data.ArrayValue:
			for arraySlots75, arrayPosition75 := v.View(), 0; arrayPosition75 < arraySlots75.Len(); arrayPosition75++ {
				slot := arraySlots75.At(arrayPosition75)
				if slot != nil {
					discover(slot.Value)
				}
			}
		case *data.FuncValue:
			if binder, ok := v.Value.(data.RequestClosureBinder); ok {
				for _, source := range binder.RequestScopeObjects() {
					if source == nil || !strings.HasPrefix(source.GetName(), "Livewire\\") || objects[source.ObjectValue] != nil {
						continue
					}
					scoped := source.CloneRequestScoped(ctx)
					objects[source.ObjectValue] = scoped
					rebindContainer(ctx, scoped, app)
					setRequestInstance(app, source.GetName(), scoped)
				}
			}
		}
	}
	for source, scoped := range objects {
		if scoped.GetName() == "Livewire\\EventBus" {
			for _, name := range []string{"listeners", "listenersAfter", "listenersBefore"} {
				value, ctl := source.GetProperty(name)
				if ctl == nil {
					discover(value)
				}
			}
		}
	}
	captures := data.NewRequestCaptureScope(ctx, objects)
	var bind func(data.Value) data.Value
	bind = func(value data.Value) data.Value {
		switch v := value.(type) {
		case *data.ClassValue:
			if scoped := objects[v.ObjectValue]; scoped != nil {
				return scoped
			}
		case *data.ArrayValue:
			// Clone each level once. Deep-cloning here and then recursively binding
			// nested arrays would copy every inner listener list twice.
			clone := data.NewArrayValueFromSlotsWithProvenance(make([]*data.ZVal, v.Len()), v.IndirectOverloadClass)
			for arraySlots76, index := v.View(), 0; index < arraySlots76.Len(); index++ {
				slot := arraySlots76.At(index)
				if slot != nil {
					clone.ReplaceSlot(index, data.CopyZValKeepName(slot, bind(slot.Value)))
				}
			}
			return clone
		case *data.FuncValue:
			if _, ok := v.Value.(data.RequestClosureBinder); ok {
				return captures.Bind(v)
			}
		}
		return value
	}
	// Only EventBus stores framework callbacks. Other mechanism properties
	// remain lazy overlays, so unused state has no copying cost.
	for source, scoped := range objects {
		if scoped.GetName() != "Livewire\\EventBus" {
			continue
		}
		for _, name := range []string{"listeners", "listenersAfter", "listenersBefore"} {
			value, ctl := source.GetProperty(name)
			if ctl == nil && value != nil {
				_ = scoped.SetProperty(name, bind(value))
			}
		}
	}
}
