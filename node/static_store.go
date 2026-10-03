package node

import (
	"fmt"
	"github.com/php-any/origami/data"
)

func coreClass(class data.ClassStmt) *ClassStatement {
	switch c := class.(type) {
	case *ClassStatement:
		return c
	case *AbstractClassStatement:
		return c.ClassStatement
	case *ClassGeneric:
		return c.ClassStatement
	}
	return nil
}

func storeClassStatic(ctx data.Context, class data.ClassStmt, name string, value data.Value) (bool, data.Control) {
	if class != nil && ctx.GetVM() != nil {
		if current, found := ctx.GetVM().GetClass(class.GetName()); found {
			class = current
		}
	}
	if class == nil || coreClass(class) == nil {
		return false, nil
	}
	for !hasStaticPropertySlot(class, name) && class.GetExtend() != nil {
		parent, ctl := ctx.GetVM().GetOrLoadClass(*class.GetExtend())
		if ctl != nil {
			return true, ctl
		}
		class = parent
	}
	core := coreClass(class)
	if core == nil {
		return false, nil
	}
	property := core.StaticProperties[name]
	if property == nil {
		return true, data.NewErrorThrowByName(nil, fmt.Errorf("Access to undeclared static property %s::$%s", class.GetName(), name), "Error")
	}
	if p, ok := property.(*ClassProperty); ok && p.IsConstant {
		return true, data.NewErrorThrowByName(nil, fmt.Errorf("Cannot modify class constant %s::%s", class.GetName(), name), "Error")
	}
	if property.GetModifier() != data.ModifierPublic && !memberAccessible(ctx, property.GetModifier(), class) {
		return true, data.NewErrorThrowByName(nil, fmt.Errorf("Cannot access non-public static property %s::$%s", class.GetName(), name), "Error")
	}
	ref := data.ResolvePropertyType(data.DeclaredTypeRef(property.GetType()), class)
	constraint := class.GetName() + "::$" + name
	old := data.LoadRequestStaticReference(class.GetName(), name)
	if old == nil {
		if raw, ok := core.StaticProperty.Load(name); ok {
			if wrapper, ok := raw.(*data.ZValValue); ok {
				old = wrapper.ZVal
			}
		}
	}
	if source, ctl := data.ReferenceSlot(value); source != nil || ctl != nil {
		if ctl != nil {
			return true, ctl
		}
		if ctl := source.BindDeclaredConstraint(constraint, ref, ctx); ctl != nil {
			return true, ctl
		}
		if old != nil && old != source {
			old.RemoveDeclaredConstraint(constraint)
			old.ReleaseRefSlot()
		}
		source.AddRefSlot()
		if !data.BindRequestStaticReference(class.GetName(), name, source) {
			core.StaticProperty.Store(name, data.NewZValValue(source))
		}
		return true, nil
	}
	prepared, accepted, ctl := data.PrepareDeclaredValueInContext(ref, value, ctx)
	if ctl != nil {
		return true, ctl
	}
	if !accepted {
		return true, data.NewTypeError(nil, fmt.Errorf("Cannot assign value to property %s of type %s", constraint, ref.String()))
	}
	if old != nil && old.Guard() != nil {
		prepared, ctl = old.PrepareWrite(prepared, ctx)
		if ctl != nil {
			return true, ctl
		}
	}
	if data.StoreRequestStatic(class.GetName(), name, prepared) {
		return true, nil
	}
	if old != nil {
		data.CowAssign(old, prepared)
	} else {
		core.StaticProperty.Store(name, prepared)
	}
	return true, nil
}
