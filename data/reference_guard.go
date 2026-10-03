package data

import "fmt"

type propertyConstraint struct {
	owner   *ObjectValue
	object  *ClassValue
	name    string
	typeRef TypeRef
}

// Constraints belong to the shared reference cell, not the variable name.
// The slice is replaced on binding/unbinding so request overlays can share
// immutable constraint metadata until its owner identities are rebound.
type ReferenceGuard struct{ properties []propertyConstraint }

func (z *ZVal) AddDeclaredConstraint(name string, ref TypeRef) {
	if ref == TypeInvalid {
		return
	}
	var constraints []propertyConstraint
	if z.Guard() != nil {
		for _, previous := range z.Guard().properties {
			if previous.owner == nil && previous.name == name {
				return
			}
		}
		constraints = append(constraints, z.Guard().properties...)
	}
	z.setGuard(&ReferenceGuard{properties: append(constraints, propertyConstraint{name: name, typeRef: ref})})
}

func (z *ZVal) BindDeclaredConstraint(name string, ref TypeRef, ctx Context) Control {
	copy := *z
	if z.reference != nil {
		cell := *z.reference
		copy.reference = &cell
	}
	copy.AddDeclaredConstraint(name, ref)
	prepared, ctl := copy.PrepareWrite(z.ReadValue(), ctx)
	if ctl != nil {
		return ctl
	}
	z.setGuard(copy.Guard())
	CowAssign(z, prepared)
	return nil
}

func (z *ZVal) RemoveDeclaredConstraint(name string) { z.removePropertyConstraint(nil, name) }

func ResolvePropertyType(ref TypeRef, owner ClassStmt) TypeRef {
	return resolvePropertyType(ref, owner)
}

func (z *ZVal) PrepareWrite(value Value, ctx Context) (Value, Control) {
	if z.Guard() == nil {
		return value, nil
	}
	var result Value
	for _, constraint := range z.Guard().properties {
		prepared, ok, ctl := PrepareDeclaredValueInContext(constraint.typeRef, value, ctx)
		if ctl != nil {
			return nil, ctl
		}
		if !ok {
			return nil, NewTypeError(nil, fmt.Errorf("Cannot assign %s to reference held by property %s of type %s", value.AsString(), constraint.name, constraint.typeRef.String()))
		}
		if result != nil && !samePreparedReferenceValue(result, prepared) {
			return nil, NewTypeError(nil, fmt.Errorf("Cannot assign value to reference held by properties with incompatible types"))
		}
		result = prepared
	}
	if result == nil {
		result = value
	}
	return result, nil
}

func samePreparedReferenceValue(left, right Value) bool {
	if left == right {
		return true
	}
	switch l := left.(type) {
	case *IntValue:
		r, ok := right.(*IntValue)
		return ok && l.Value == r.Value
	case *FloatValue:
		r, ok := right.(*FloatValue)
		return ok && l.Value == r.Value
	case *StringValue:
		r, ok := right.(*StringValue)
		return ok && l.Value == r.Value
	case *BoolValue:
		r, ok := right.(*BoolValue)
		return ok && l.Value == r.Value
	case *NullValue:
		_, ok := right.(*NullValue)
		return ok
	}
	return false
}

func (z *ZVal) addPropertyConstraint(object *ClassValue, name string, ref TypeRef) {
	if ref == TypeInvalid {
		return
	}
	if z.Guard() != nil {
		for _, previous := range z.Guard().properties {
			if previous.owner == object.ObjectValue && previous.name == name {
				return
			}
		}
	}
	constraint := propertyConstraint{owner: object.ObjectValue, object: object, name: name, typeRef: ref}
	var properties []propertyConstraint
	if z.Guard() != nil {
		properties = append(properties, z.Guard().properties...)
	}
	z.setGuard(&ReferenceGuard{properties: append(properties, constraint)})
}

func (z *ZVal) removePropertyConstraint(owner *ObjectValue, name string) {
	if z == nil || z.Guard() == nil {
		return
	}
	var properties []propertyConstraint
	for _, constraint := range z.Guard().properties {
		if constraint.owner != owner || constraint.name != name {
			properties = append(properties, constraint)
		}
	}
	if len(properties) == 0 {
		z.setGuard(nil)
	} else {
		z.setGuard(&ReferenceGuard{properties: properties})
	}
}

// Resolve property-relative names once on a cold reference binding or when
// a property declaration actually needs lexical class scope.
func PropertyType(object *ClassValue, name string, ref TypeRef) TypeRef {
	if !ref.RequiresClassScope() {
		return ref
	}
	owner := object.Class
	for owner != nil {
		if _, found := owner.GetProperty(name); found {
			break
		}
		if owner.GetExtend() == nil {
			break
		}
		parent, found := object.GetVM().GetClass(*owner.GetExtend())
		if !found {
			break
		}
		owner = parent
	}
	return resolvePropertyType(ref, owner)
}

func resolvePropertyType(ref TypeRef, owner ClassStmt) TypeRef {
	if owner == nil {
		return ref
	}
	if ref == TypeSelf {
		return DeclaredTypeRef(NewDeclaredType(owner.GetName()))
	}
	if ref == TypeParent && owner.GetExtend() != nil {
		return DeclaredTypeRef(NewDeclaredType(*owner.GetExtend()))
	}
	if ref.Kind() == TypeKindUnion || ref.Kind() == TypeKindIntersection {
		parts := ref.Members()
		types := make([]Types, parts.Len())
		for i := 0; i < parts.Len(); i++ {
			types[i] = resolvePropertyType(parts.At(i), owner)
		}
		if ref.Kind() == TypeKindUnion {
			return DeclaredTypeRef(NewDeclaredUnionType(types))
		}
		return DeclaredTypeRef(NewDeclaredIntersectionType(types))
	}
	return ref
}

func ReferenceSlot(value Value) (*ZVal, Control) {
	switch ref := value.(type) {
	case *ArraySlotRef:
		return ref.Slot, nil
	case *ZValValue:
		return ref.ZVal, nil
	case *ReferenceValue:
		if property, ok := ref.Val.(interface {
			GetZVal(Context) (*ZVal, Control)
		}); ok {
			return property.GetZVal(ref.Ctx)
		}
		return ref.Ctx.GetIndexZVal(ref.Val.GetIndex()), nil
	}
	return nil, nil
}

func PreparePropertyValue(object *ClassValue, name string, ty Types, value Value, ctx Context) (Value, bool, Control) {
	switch value.(type) {
	case *ReferenceValue, *ArraySlotRef, *ZValValue:
		// Binding validates all property constraints before mutating the source.
		return value, true, nil
	}
	ref := PropertyType(object, name, DeclaredTypeRef(ty))
	prepared, accepted, ctl := PrepareDeclaredValueInContext(ref, value, ctx)
	if ctl != nil || !accepted {
		return prepared, accepted, ctl
	}
	if slot, ok := object.property.GetZVal(name); ok && slot.Guard() != nil {
		prepared, ctl = slot.PrepareWrite(prepared, ctx)
		if ctl != nil {
			return nil, false, ctl
		}
	}
	return prepared, true, nil
}

func (c *ClassValue) BindPropertyReference(ctx Context, name string, slot *ZVal) Control {
	if slot == nil {
		return NewErrorThrow(nil, fmt.Errorf("Cannot bind missing reference"))
	}
	if property, found := c.GetPropertyStmt(name); found {
		ref := PropertyType(c, name, DeclaredTypeRef(property.GetType()))
		copy := *slot
		if slot.reference != nil {
			cell := *slot.reference
			copy.reference = &cell
		}
		copy.addPropertyConstraint(c, name, ref)
		value, ctl := copy.PrepareWrite(slot.ReadValue(), ctx)
		if ctl != nil {
			return ctl
		}
		slot.setGuard(copy.Guard())
		CowAssign(slot, value)
	}
	if old, ok := c.property.GetZVal(name); ok && old != slot {
		old.removePropertyConstraint(c.ObjectValue, name)
		old.ReleaseRefSlot()
	}
	c.property.(interface{ BindZVal(string, *ZVal) }).BindZVal(name, slot)
	slot.AddRefSlot()
	return nil
}
