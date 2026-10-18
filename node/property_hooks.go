package node

import (
	"fmt"
	"github.com/php-any/origami/data"
)

type propertyWriteFlags uint8

const (
	writeReadonly propertyWriteFlags = 1 << iota
	writeRestricted
	writeHooked
)

// Compile once when a declaration is bound to its owning class (including
// trait composition). Never retain request objects or caller scope in this plan.
func (p *ClassProperty) compileWritePolicy() {
	p.writeFlags = 0
	if p.IsReadonly {
		p.writeFlags |= writeReadonly
	}
	p.setModifier = data.NewModifier(p.SetVisibility)
	if p.SetVisibility != "" && p.setModifier != data.ModifierPublic {
		p.writeFlags |= writeRestricted
	}
	if p.SetHook != nil || p.Virtual {
		p.writeFlags |= writeHooked
	}
}

func propertyWriteAccess(ctx data.Context, p *ClassProperty, from data.From) data.Control {
	if p.writeFlags&writeRestricted != 0 {
		owner, _ := ctx.GetVM().GetClass(p.DeclaringClass)
		if !memberAccessible(ctx, p.setModifier, owner) {
			return data.NewErrorThrowByName(from, fmt.Errorf("Cannot modify %s(set) property %s::$%s", p.SetVisibility, p.DeclaringClass, p.Name), "Error")
		}
	}
	return nil
}

// Only declarations with a compiled special-write policy enter this path.
func prepareSpecialPropertyWrite(ctx data.Context, object *data.ClassValue, p *ClassProperty, value data.Value, from data.From) (bool, *readonlyCloneState, data.Control) {
	if ctl := propertyWriteAccess(ctx, p, from); ctl != nil {
		return true, nil, ctl
	}
	if p.writeFlags&writeHooked != 0 && !insidePropertyHook(ctx, object, p) {
		if p.SetHook != nil {
			_, ctl := invokePropertyHook(ctx, object, p, p.SetHook, value)
			return true, nil, ctl
		}
		return true, nil, data.NewErrorThrowByName(from, fmt.Errorf("Cannot write to get-only property %s", p.Name), "Error")
	}
	var cloneState *readonlyCloneState
	if p.writeFlags&writeReadonly != 0 {
		if object.ObjectValue.HasProperty(p.PropertyStorageName()) {
			cloneState = readonlyClonePermission(ctx, object.ObjectValue, p.PropertyStorageName())
			if cloneState == nil {
				return true, nil, data.NewErrorThrowByName(from, fmt.Errorf("Cannot modify readonly property %s::$%s", p.DeclaringClass, p.Name), "Error")
			}
		}
		owner, _ := ctx.GetVM().GetClass(p.DeclaringClass)
		if !memberAccessible(ctx, data.ModifierProtected, owner) {
			return true, nil, data.NewErrorThrowByName(from, fmt.Errorf("Cannot initialize readonly property %s::$%s from global scope", p.DeclaringClass, p.Name), "Error")
		}
	}
	return false, cloneState, nil
}

func insidePropertyHook(ctx data.Context, object *data.ClassValue, property *ClassProperty) bool {
	frame, ok := ctx.(*data.ClassMethodContext)
	return ok && frame.HookProperty == property.Name && frame.ClassValue.InstanceIdentity() == object.InstanceIdentity()
}

func invokePropertyHook(ctx data.Context, object *data.ClassValue, property *ClassProperty, hook data.Method, value data.Value) (data.GetValue, data.Control) {
	owner, _ := ctx.GetVM().GetClass(property.DeclaringClass)
	frame := data.WrapMethodFrame(ctx.CreateContext(hook.GetVariables()), object, owner, object.Class)
	frame.HookProperty = property.Name
	defer tryReleaseCallContext(hook, frame)
	if value != nil {
		parameters := hook.GetParams()
		if len(parameters) != 1 {
			return nil, data.NewErrorThrowByName(property.GetFrom(), fmt.Errorf("set hook requires one parameter"), "Error")
		}
		variable := hook.GetVariables()[0]
		prepared, ok, ctl := data.PreparePropertyValue(object, data.PropertyStorageName(property), variable.GetType(), value, ctx)
		if ctl != nil {
			return nil, ctl
		}
		if !ok {
			return nil, data.NewTypeError(property.GetFrom(), fmt.Errorf("Invalid property hook value"))
		}
		if ctl := frame.SetVariableValue(variable, prepared); ctl != nil {
			return nil, ctl
		}
	}
	return hook.Call(frame)
}
