package node

import (
	"fmt"
	"github.com/php-any/origami/data"
)

func checkStaticPropertyAccess(ctx data.Context, class data.ClassStmt, name string, from data.From) data.Control {
	for class != nil {
		if core := coreClass(class); core != nil {
			if property := core.StaticProperties[name]; property != nil {
				if !memberAccessible(ctx, property.GetModifier(), class) {
					return data.NewErrorThrowByName(from, fmt.Errorf("Cannot access non-public static property %s::$%s", class.GetName(), name), "Error")
				}
				if declaration, ok := property.(*ClassProperty); ok && declaration.IsConstant {
					if _, cached := core.StaticProperty.Load(name); !cached {
						_, ctl := core.initStaticProperty(property, ctx.GetVM())
						return ctl
					}
				}
				return nil
			}
		}
		if class.GetExtend() == nil {
			break
		}
		class, _ = ctx.GetVM().GetClass(*class.GetExtend())
	}
	return nil
}

// LookupStaticProperty 在类、父类及 implements 的接口继承链上查找类常量/静态属性。
func LookupStaticProperty(vm data.VM, class data.ClassStmt, name string) (data.Value, bool) {
	if class == nil {
		return nil, false
	}
	if v, ok := staticPropertyOn(vm, class, name); ok {
		return v, true
	}
	for _, ifaceName := range class.GetImplements() {
		if v, ok := lookupStaticPropertyOnInterface(vm, ifaceName, name); ok {
			return v, true
		}
	}
	last := class
	for last.GetExtend() != nil {
		parentName := last.GetExtend()
		if parentName == nil || *parentName == "" {
			break
		}
		parent, acl := vm.GetOrLoadClass(*parentName)
		if acl != nil || parent == nil {
			break
		}
		if v, ok := staticPropertyOn(vm, parent, name); ok {
			return v, true
		}
		for _, ifaceName := range parent.GetImplements() {
			if v, ok := lookupStaticPropertyOnInterface(vm, ifaceName, name); ok {
				return v, true
			}
		}
		last = parent
	}
	return nil, false
}

func staticPropertyOn(vm data.VM, class data.ClassStmt, name string) (data.Value, bool) {
	switch class := class.(type) {
	case *ClassStatement:
		return class.getStaticProperty(nil, vm, name)
	case *AbstractClassStatement:
		return class.getStaticProperty(nil, vm, name)
	case *ClassGeneric:
		return class.getStaticProperty(nil, vm, name)
	}
	if gsp, ok := class.(data.GetStaticProperty); ok {
		return gsp.GetStaticProperty(name)
	}
	return nil, false
}

func staticPropertyValue(ctx data.Context, getter data.GetStaticProperty, name string) (data.Value, bool) {
	switch class := getter.(type) {
	case *ClassStatement:
		return class.getStaticProperty(ctx, nil, name)
	case *AbstractClassStatement:
		return class.getStaticProperty(ctx, nil, name)
	case *ClassGeneric:
		return class.getStaticProperty(ctx, nil, name)
	default:
		return getter.GetStaticProperty(name)
	}
}

func lookupStaticPropertyOnInterface(vm data.VM, ifaceName, name string) (data.Value, bool) {
	iface, acl := vm.LoadPkg(ifaceName)
	if acl != nil || iface == nil {
		return nil, false
	}
	if gsp, ok := iface.(data.GetStaticProperty); ok {
		if v, ok := gsp.GetStaticProperty(name); ok {
			return v, true
		}
	}
	if is, ok := iface.(data.InterfaceStmt); ok {
		for _, parent := range is.GetExtends() {
			if v, ok := lookupStaticPropertyOnInterface(vm, parent, name); ok {
				return v, true
			}
		}
	}
	return nil, false
}
