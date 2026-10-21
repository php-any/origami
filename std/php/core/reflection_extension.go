package core

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"strings"
)

type ReflectionExtensionClass struct{ node.Node }

func (*ReflectionExtensionClass) GetName() string                          { return "ReflectionExtension" }
func (*ReflectionExtensionClass) GetExtend() *string                       { return nil }
func (*ReflectionExtensionClass) GetImplements() []string                  { return nil }
func (*ReflectionExtensionClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (*ReflectionExtensionClass) GetPropertyList() []data.Property         { return nil }
func (*ReflectionExtensionClass) GetConstruct() data.Method {
	return &reflectionExtensionMethod{name: "__construct"}
}
func (c *ReflectionExtensionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *ReflectionExtensionClass) GetMethods() []data.Method {
	var methods []data.Method
	for _, name := range []string{"__construct", "getName", "getVersion", "info", "isPersistent", "isTemporary"} {
		methods = append(methods, &reflectionExtensionMethod{name: name})
	}
	return methods
}
func (c *ReflectionExtensionClass) GetMethod(name string) (data.Method, bool) {
	for _, m := range c.GetMethods() {
		if strings.EqualFold(name, m.GetName()) {
			return m, true
		}
	}
	return nil, false
}

type reflectionExtensionMethod struct{ name string }

func (m *reflectionExtensionMethod) GetName() string          { return m.name }
func (*reflectionExtensionMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (*reflectionExtensionMethod) GetIsStatic() bool          { return false }
func (*reflectionExtensionMethod) GetReturnType() data.Types  { return nil }
func (m *reflectionExtensionMethod) GetParams() []data.GetValue {
	if m.name == "__construct" {
		return []data.GetValue{node.NewParameter(nil, "name", 0, nil, nil)}
	}
	return nil
}
func (m *reflectionExtensionMethod) GetVariables() []data.Variable {
	if m.name == "__construct" {
		return []data.Variable{node.NewVariable(nil, "name", 0, nil)}
	}
	return nil
}
func (m *reflectionExtensionMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	receiver, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("ReflectionExtension: missing receiver"))
	}
	nameValue, _ := receiver.ObjectValue.GetProperty("name")
	name := ""
	if nameValue != nil {
		name = nameValue.AsString()
	}
	switch m.name {
	case "__construct":
		v, _ := ctx.GetIndexValue(0)
		if v != nil {
			name = strings.ToLower(v.AsString())
		}
		if !extensionAvailable(ctx.GetVM(), name) {
			return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Extension %s does not exist", name), "ReflectionException")
		}
		receiver.ObjectValue.SetProperty("name", data.NewStringValue(name))
		return nil, nil
	case "getName":
		return data.NewStringValue(name), nil
	case "getVersion":
		v, _ := ctx.GetVM().GetConstant("PHP_VERSION")
		return v, nil
	case "isPersistent":
		return data.NewBoolValue(true), nil
	case "isTemporary":
		return data.NewBoolValue(false), nil
	case "info":
		return nil, data.EmitOutput(ctx, name+"\nImplementation => Origami PHP compatibility module\n")
	}
	return nil, nil
}
