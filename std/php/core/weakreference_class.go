package core

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"strings"
	"weak"
)

type WeakReferenceClass struct {
	node.Node
	target weak.Pointer[data.ClassValue]
}

func (*WeakReferenceClass) GetName() string                          { return "WeakReference" }
func (*WeakReferenceClass) GetExtend() *string                       { return nil }
func (*WeakReferenceClass) GetImplements() []string                  { return nil }
func (*WeakReferenceClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (*WeakReferenceClass) GetPropertyList() []data.Property         { return nil }
func (*WeakReferenceClass) GetConstruct() data.Method {
	return &weakReferenceMethod{name: "__construct"}
}
func (c *WeakReferenceClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(&WeakReferenceClass{}, ctx.CreateBaseContext()), nil
}
func (c *WeakReferenceClass) GetMethods() []data.Method {
	return []data.Method{&weakReferenceMethod{name: "create"}, &weakReferenceMethod{name: "get"}}
}
func (c *WeakReferenceClass) GetMethod(name string) (data.Method, bool) {
	switch strings.ToLower(name) {
	case "create":
		return &weakReferenceMethod{name: "create"}, true
	case "get":
		return &weakReferenceMethod{name: "get"}, true
	}
	return nil, false
}
func (c *WeakReferenceClass) GetStaticMethod(name string) (data.Method, bool) {
	return c.GetMethod(name)
}

type weakReferenceMethod struct{ name string }

func (m *weakReferenceMethod) GetName() string          { return m.name }
func (*weakReferenceMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *weakReferenceMethod) GetIsStatic() bool        { return m.name == "create" }
func (*weakReferenceMethod) GetReturnType() data.Types  { return nil }
func (m *weakReferenceMethod) GetParams() []data.GetValue {
	if m.name == "create" {
		return []data.GetValue{node.NewParameter(nil, "object", 0, nil, nil)}
	}
	return nil
}
func (m *weakReferenceMethod) GetVariables() []data.Variable {
	if m.name == "create" {
		return []data.Variable{node.NewVariable(nil, "object", 0, nil)}
	}
	return nil
}
func (m *weakReferenceMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	switch m.name {
	case "__construct":
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Direct instantiation of WeakReference is not allowed, use WeakReference::create instead"), "Error")
	case "create":
		v, _ := ctx.GetIndexValue(0)
		var object *data.ClassValue
		switch v := v.(type) {
		case *data.ClassValue:
			object = v
		case *data.ThisValue:
			object = v.ClassValue
		}
		if object == nil {
			return nil, data.NewTypeError(nil, fmt.Errorf("WeakReference::create(): argument must be an object"))
		}
		class := &WeakReferenceClass{target: weak.Make(object.InstanceIdentity())}
		return data.NewClassValue(class, ctx.CreateBaseContext()), nil
	case "get":
		if receiver, ok := ctx.(*data.ClassMethodContext); ok {
			if class, ok := receiver.Class.(*WeakReferenceClass); ok {
				if target := class.target.Value(); target != nil {
					return target, nil
				}
			}
		}
		return data.NewNullValue(), nil
	}
	return nil, nil
}
