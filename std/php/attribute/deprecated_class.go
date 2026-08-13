package attribute

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// DeprecatedClass PHP 8.3+ 原生 Deprecated 类
// 用于标记已弃用的代码元素
type DeprecatedClass struct {
	node.Node
	construct data.Method
}

func NewDeprecatedClass() *DeprecatedClass {
	return &DeprecatedClass{
		construct: &DeprecatedConstructMethod{},
	}
}

func (c *DeprecatedClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *DeprecatedClass) GetName() string {
	return "Deprecated"
}

func (c *DeprecatedClass) GetExtend() *string {
	return nil
}

func (c *DeprecatedClass) GetImplements() []string {
	return []string{}
}

func (c *DeprecatedClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (c *DeprecatedClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (c *DeprecatedClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return c.construct, true
	}
	return nil, false
}

func (c *DeprecatedClass) GetMethods() []data.Method {
	return []data.Method{
		c.construct,
	}
}

func (c *DeprecatedClass) GetConstruct() data.Method {
	return c.construct
}

// DeprecatedConstructMethod Deprecated 构造函数
// Deprecated::__construct(?string $reason = null, ?string $replacement = null)
type DeprecatedConstructMethod struct{}

func (m *DeprecatedConstructMethod) GetName() string {
	return "__construct"
}

func (m *DeprecatedConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *DeprecatedConstructMethod) GetIsStatic() bool {
	return false
}

func (m *DeprecatedConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "reason", 0, data.NewNullValue(), nil),
		node.NewParameter(nil, "replacement", 1, data.NewNullValue(), nil),
	}
}

func (m *DeprecatedConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "reason", 0, nil),
		node.NewVariable(nil, "replacement", 1, nil),
	}
}

func (m *DeprecatedConstructMethod) GetReturnType() data.Types {
	return nil
}

func (m *DeprecatedConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// Deprecated 构造函数不需要特殊处理，只是标记代码已弃用
	return data.NewNullValue(), nil
}
