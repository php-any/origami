package reflection

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

type ReflectionPropertyClass struct {
	node.Node
}

func (c *ReflectionPropertyClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *ReflectionPropertyClass) GetName() string                               { return "ReflectionProperty" }
func (c *ReflectionPropertyClass) GetExtend() *string                            { return nil }
func (c *ReflectionPropertyClass) GetImplements() []string                       { return nil }
func (c *ReflectionPropertyClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *ReflectionPropertyClass) GetPropertyList() []data.Property              { return nil }
func (c *ReflectionPropertyClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case token.ConstructName:
		return &ReflectionPropertyConstructMethod{}, true
	case "setAccessible":
		return &ReflectionPropertySetAccessibleMethod{}, true
	case "getName":
		return &ReflectionPropertyGetNameMethod{}, true
	case "getDeclaringClass":
		return &ReflectionPropertyGetDeclaringClassMethod{}, true
	case "getValue":
		return &ReflectionPropertyGetValueMethod{}, true
	case "setValue":
		return &ReflectionPropertySetValueMethod{}, true
	case "setRawValue":
		// PHP 8.4：绕过 set hook；当前无 hook 运行时，行为同 setValue
		return &ReflectionPropertySetValueMethod{raw: true}, true
	case "getRawValue":
		return &ReflectionPropertyGetValueMethod{raw: true}, true
	}
	return nil, false
}
func (c *ReflectionPropertyClass) GetMethods() []data.Method {
	return []data.Method{
		&ReflectionPropertyConstructMethod{},
		&ReflectionPropertySetAccessibleMethod{},
		&ReflectionPropertyGetNameMethod{},
		&ReflectionPropertyGetDeclaringClassMethod{},
		&ReflectionPropertyGetValueMethod{},
		&ReflectionPropertySetValueMethod{},
		&ReflectionPropertySetValueMethod{raw: true},
	}
}

func newReflectionProperty(ctx data.Context, className, propertyName string) *data.ClassValue {
	propertyClass := &ReflectionPropertyClass{}
	propertyValue := data.NewClassValue(propertyClass, ctx.CreateBaseContext())
	propertyValue.ObjectValue.SetProperty("_className", data.NewStringValue(className))
	propertyValue.ObjectValue.SetProperty("_propertyName", data.NewStringValue(propertyName))
	propertyValue.ObjectValue.SetProperty("name", data.NewStringValue(propertyName))
	propertyValue.ObjectValue.SetProperty("class", data.NewStringValue(className))
	return propertyValue
}
func (c *ReflectionPropertyClass) GetConstruct() data.Method {
	return &ReflectionPropertyConstructMethod{}
}

type ReflectionPropertyConstructMethod struct{}

func (m *ReflectionPropertyConstructMethod) GetName() string            { return "__construct" }
func (m *ReflectionPropertyConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ReflectionPropertyConstructMethod) GetIsStatic() bool          { return false }
func (m *ReflectionPropertyConstructMethod) GetReturnType() data.Types  { return nil }
func (m *ReflectionPropertyConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "class", 0, nil, nil),
		node.NewParameter(nil, "property", 1, nil, nil),
	}
}
func (m *ReflectionPropertyConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "class", 0, data.Mixed{}),
		node.NewVariable(nil, "property", 1, data.Mixed{}),
	}
}
func (m *ReflectionPropertyConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	classVal, _ := ctx.GetIndexValue(0)
	propVal, _ := ctx.GetIndexValue(1)
	if classVal == nil || propVal == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("ReflectionProperty::__construct() expects class and property"))
	}

	className := classVal.AsString()
	if cv, ok := classVal.(*data.ClassValue); ok && cv.Class != nil {
		className = cv.Class.GetName()
	}
	propName := propVal.AsString()

	if cmc, ok := ctx.(*data.ClassMethodContext); ok && cmc.ObjectValue != nil {
		cmc.ObjectValue.SetProperty("_className", data.NewStringValue(className))
		cmc.ObjectValue.SetProperty("_propertyName", data.NewStringValue(propName))
		cmc.ObjectValue.SetProperty("name", data.NewStringValue(propName))
		cmc.ObjectValue.SetProperty("class", data.NewStringValue(className))
	}
	return nil, nil
}

type ReflectionPropertySetAccessibleMethod struct{}

func (m *ReflectionPropertySetAccessibleMethod) GetName() string { return "setAccessible" }
func (m *ReflectionPropertySetAccessibleMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *ReflectionPropertySetAccessibleMethod) GetIsStatic() bool         { return false }
func (m *ReflectionPropertySetAccessibleMethod) GetReturnType() data.Types { return nil }
func (m *ReflectionPropertySetAccessibleMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "accessible", 0, nil, nil),
	}
}
func (m *ReflectionPropertySetAccessibleMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "accessible", 0, data.Mixed{}),
	}
}
func (m *ReflectionPropertySetAccessibleMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}

type ReflectionPropertyGetNameMethod struct{}

func (m *ReflectionPropertyGetNameMethod) GetName() string            { return "getName" }
func (m *ReflectionPropertyGetNameMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ReflectionPropertyGetNameMethod) GetIsStatic() bool          { return false }
func (m *ReflectionPropertyGetNameMethod) GetReturnType() data.Types  { return data.String{} }
func (m *ReflectionPropertyGetNameMethod) GetParams() []data.GetValue { return nil }
func (m *ReflectionPropertyGetNameMethod) GetVariables() []data.Variable {
	return nil
}

type ReflectionPropertyGetDeclaringClassMethod struct{}

func (m *ReflectionPropertyGetDeclaringClassMethod) GetName() string { return "getDeclaringClass" }
func (m *ReflectionPropertyGetDeclaringClassMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *ReflectionPropertyGetDeclaringClassMethod) GetIsStatic() bool { return false }
func (m *ReflectionPropertyGetDeclaringClassMethod) GetReturnType() data.Types {
	return data.Mixed{}
}
func (m *ReflectionPropertyGetDeclaringClassMethod) GetParams() []data.GetValue {
	return nil
}
func (m *ReflectionPropertyGetDeclaringClassMethod) GetVariables() []data.Variable {
	return nil
}
func (m *ReflectionPropertyGetDeclaringClassMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok && cmc.ObjectValue != nil {
		if className, ctl := cmc.ObjectValue.GetProperty("_className"); ctl == nil && className != nil {
			reflectionClass := &ReflectionClassClass{}
			value := data.NewClassValue(reflectionClass, ctx.CreateBaseContext())
			value.ObjectValue.SetProperty("_className", data.NewStringValue(className.AsString()))
			return value, nil
		}
	}
	return data.NewNullValue(), nil
}
func (m *ReflectionPropertyGetNameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok && cmc.ObjectValue != nil {
		if v, ctl := cmc.ObjectValue.GetProperty("_propertyName"); ctl == nil && v != nil {
			return v, nil
		}
	}
	return data.NewStringValue(""), nil
}

type ReflectionPropertyGetValueMethod struct {
	raw bool
}

func (m *ReflectionPropertyGetValueMethod) GetName() string {
	if m.raw {
		return "getRawValue"
	}
	return "getValue"
}
func (m *ReflectionPropertyGetValueMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ReflectionPropertyGetValueMethod) GetIsStatic() bool          { return false }
func (m *ReflectionPropertyGetValueMethod) GetReturnType() data.Types  { return nil }
func (m *ReflectionPropertyGetValueMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object", 0, data.NewNullValue(), nil),
	}
}
func (m *ReflectionPropertyGetValueMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object", 0, data.Mixed{}),
	}
}
func (m *ReflectionPropertyGetValueMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	propName := reflectionPropertyName(ctx)
	objVal, _ := ctx.GetIndexValue(0)
	if propName == "" || objVal == nil {
		return data.NewNullValue(), nil
	}
	switch o := objVal.(type) {
	case *data.ClassValue:
		if o.ObjectValue != nil {
			if v, _ := o.ObjectValue.GetProperty(propName); v != nil {
				return v, nil
			}
		}
		if z, ctl := o.GetPropertyZVal(propName); ctl == nil && z != nil && z.Value != nil {
			return z.Value, nil
		}
	case data.GetProperty:
		if v, ctl := o.GetProperty(propName); ctl == nil && v != nil {
			return v, nil
		}
	}
	return data.NewNullValue(), nil
}

type ReflectionPropertySetValueMethod struct {
	raw bool
}

func (m *ReflectionPropertySetValueMethod) GetName() string {
	if m.raw {
		return "setRawValue"
	}
	return "setValue"
}
func (m *ReflectionPropertySetValueMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ReflectionPropertySetValueMethod) GetIsStatic() bool          { return false }
func (m *ReflectionPropertySetValueMethod) GetReturnType() data.Types  { return nil }
func (m *ReflectionPropertySetValueMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object", 0, nil, nil),
		node.NewParameter(nil, "value", 1, nil, nil),
	}
}
func (m *ReflectionPropertySetValueMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object", 0, data.Mixed{}),
		node.NewVariable(nil, "value", 1, data.Mixed{}),
	}
}
func (m *ReflectionPropertySetValueMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	propName := reflectionPropertyName(ctx)
	objVal, _ := ctx.GetIndexValue(0)
	value, _ := ctx.GetIndexValue(1)
	if propName == "" || objVal == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("ReflectionProperty::setValue() missing object or property"))
	}
	if value == nil {
		value = data.NewNullValue()
	}
	switch o := objVal.(type) {
	case *data.ClassValue:
		return nil, o.SetProperty(propName, value)
	case data.SetProperty:
		return nil, o.SetProperty(propName, value)
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("ReflectionProperty::setValue() expects object"))
}

func reflectionPropertyName(ctx data.Context) string {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok && cmc.ObjectValue != nil {
		if v, ctl := cmc.ObjectValue.GetProperty("_propertyName"); ctl == nil && v != nil {
			return v.AsString()
		}
	}
	return ""
}
