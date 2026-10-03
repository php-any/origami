package node

import (
	"fmt"
	"github.com/php-any/origami/data"
)

func staticReferenceSlot(ctx data.Context, class data.ClassStmt, name string, from data.From) (*data.ZVal, data.Control) {
	if class != nil && ctx.GetVM() != nil {
		if current, found := ctx.GetVM().GetClass(class.GetName()); found {
			class = current
		}
	}
	if class == nil {
		return nil, data.NewErrorThrow(from, fmt.Errorf("static property has no class"))
	}
	for !hasStaticPropertySlot(class, name) && class.GetExtend() != nil {
		parent, ctl := ctx.GetVM().GetOrLoadClass(*class.GetExtend())
		if ctl != nil {
			return nil, ctl
		}
		class = parent
	}
	var declared data.Property
	var core *ClassStatement
	switch c := class.(type) {
	case *ClassStatement:
		core = c
	case *AbstractClassStatement:
		core = c.ClassStatement
	case *ClassGeneric:
		core = c.ClassStatement
	}
	if core == nil {
		return nil, data.NewErrorThrow(from, fmt.Errorf("native static property %s::$%s does not expose a reference slot", class.GetName(), name))
	}
	declared = core.StaticProperties[name]
	if p, ok := declared.(*ClassProperty); ok && p.IsConstant {
		return nil, data.NewErrorThrow(from, fmt.Errorf("Cannot take reference to class constant"))
	}
	if declared == nil {
		return nil, data.NewErrorThrow(from, fmt.Errorf("Access to undeclared static property %s::$%s", class.GetName(), name))
	}
	if declared.GetModifier() != data.ModifierPublic && !memberAccessible(ctx, declared.GetModifier(), class) {
		return nil, data.NewErrorThrow(from, fmt.Errorf("Cannot access non-public static property %s::$%s", class.GetName(), name))
	}
	value, _ := core.getStaticProperty(ctx, nil, name)
	if slot, scoped := data.StaticReferenceSlot(class.GetName(), name, value); scoped {
		slot.AddDeclaredConstraint(class.GetName()+"::$"+name, data.ResolvePropertyType(data.DeclaredTypeRef(declared.GetType()), class))
		return slot, nil
	}
	if previous, ok := core.StaticProperty.Load(name); ok {
		if slot, ok := previous.(*data.ZValValue); ok {
			return slot.ZVal, nil
		}
	}
	slot := data.NewZVal(value)
	slot.AddDeclaredConstraint(class.GetName()+"::$"+name, data.ResolvePropertyType(data.DeclaredTypeRef(declared.GetType()), class))
	core.StaticProperty.Store(name, data.NewZValValue(slot))
	return slot, nil
}

func (pe *CallStaticProperty) GetZVal(ctx data.Context) (*data.ZVal, data.Control) {
	class, ok := pe.Stmt.(data.ClassStmt)
	if !ok {
		value, ctl := pe.Stmt.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		switch value := value.(type) {
		case *data.ClassValue:
			class = value.Class
		case *data.ThisValue:
			class = value.Class
		case *data.StringValue:
			var ctl data.Control
			class, ctl = ctx.GetVM().GetOrLoadClass(value.Value)
			if ctl != nil {
				return nil, ctl
			}
		}
	}
	return staticReferenceSlot(ctx, class, pe.Property, pe.GetFrom())
}
func (pe *CallStaticPropertyLater) GetZVal(ctx data.Context) (*data.ZVal, data.Control) {
	access, ctl := pe.resolveAccess(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return access.GetZVal(ctx)
}
func (pe *CallStaticKeywordProperty) GetZVal(ctx data.Context) (*data.ZVal, data.Control) {
	class, ctl, _ := resolveLateStaticClass(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return staticReferenceSlot(ctx, class, pe.Property, pe.GetFrom())
}
func (pe *CallSelfProperty) GetZVal(ctx data.Context) (*data.ZVal, data.Control) {
	var class data.ClassStmt
	switch ctx := ctx.(type) {
	case *data.ClassMethodContext:
		class = ctx.SelfClass
		if class == nil {
			class = ctx.Class
		}
	case *data.ClassValue:
		class = ctx.Class
	}
	return staticReferenceSlot(ctx, class, pe.Property, pe.GetFrom())
}
