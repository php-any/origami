package reflection

import (
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
	case "getValue":
		return &ReflectionPropertyGetValueMethod{}, true
	case "setValue":
		return &ReflectionPropertySetValueMethod{}, true
	}
	return nil, false
}
func (c *ReflectionPropertyClass) GetMethods() []data.Method {
	return []data.Method{
		&ReflectionPropertyConstructMethod{},
		&ReflectionPropertySetAccessibleMethod{},
		&ReflectionPropertyGetValueMethod{},
		&ReflectionPropertySetValueMethod{},
	}
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
	return nil, nil
}

// setAccessible PHP 8.1+ 此方法已废弃（始终为 no-op），但旧代码仍可能调用
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

// getValue 获取属性值
type ReflectionPropertyGetValueMethod struct{}

func (m *ReflectionPropertyGetValueMethod) GetName() string            { return "getValue" }
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
	return data.NewNullValue(), nil
}

// setValue 设置属性值
type ReflectionPropertySetValueMethod struct{}

func (m *ReflectionPropertySetValueMethod) GetName() string            { return "setValue" }
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
	return nil, nil
}
