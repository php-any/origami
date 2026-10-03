package node

import "github.com/php-any/origami/data"

func assignPromotedProperty(frame, receiver data.Context, parameter *PromotedParameter, value data.Value) data.Control {
	if this, ok := value.(*data.ThisValue); ok {
		value = this.ClassValue.InstanceIdentity()
	}
	if owner, ok := frame.(*data.ClassMethodContext); ok && owner.SelfClass != nil {
		if property, found := owner.SelfClass.GetProperty(parameter.GetName()); found {
			if object, ok := receiver.(data.SetProperty); ok {
				if instance, ok := object.(*data.ClassValue); ok {
					_, ctl := (&CallObjectProperty{Node: parameter.Node, Object: instance, Property: parameter.GetName()}).AssignValue(frame, value)
					return ctl
				}
				return object.SetProperty(data.PropertyStorageName(property), value)
			}
		}
	}
	return receiver.SetVariableValue(parameter, value)
}

func propertyObject(object data.Value) *data.ClassValue {
	if this, ok := object.(*data.ThisValue); ok {
		return this.ClassValue
	}
	return object.(*data.ClassValue)
}

func assignPropertyValue(ctx data.Context, object *data.ClassValue, name string, value data.Value) (data.Value, data.Control) {
	if this, ok := value.(*data.ThisValue); ok {
		value = this.ClassValue.InstanceIdentity()
	}
	if slot, ctl := data.ReferenceSlot(value); ctl != nil {
		return nil, ctl
	} else if slot != nil {
		if ctl := object.BindPropertyReference(ctx, name, slot); ctl != nil {
			return nil, ctl
		}
		return slot.ReadValue(), nil
	}
	return value, object.SetProperty(name, value)
}
