package attribute

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReturnTypeWillChangeClass PHP 8.1+ 原生 ReturnTypeWillChange 类
// 用于标记返回类型将在未来版本中改变
type ReturnTypeWillChangeClass struct {
	node.Node
	construct data.Method
}

func NewReturnTypeWillChangeClass() *ReturnTypeWillChangeClass {
	return &ReturnTypeWillChangeClass{
		construct: &ReturnTypeWillChangeConstructMethod{},
	}
}

func (c *ReturnTypeWillChangeClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *ReturnTypeWillChangeClass) GetName() string {
	return "ReturnTypeWillChange"
}

func (c *ReturnTypeWillChangeClass) GetExtend() *string {
	return nil
}

func (c *ReturnTypeWillChangeClass) GetImplements() []string {
	return []string{}
}

func (c *ReturnTypeWillChangeClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (c *ReturnTypeWillChangeClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (c *ReturnTypeWillChangeClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return c.construct, true
	}
	return nil, false
}

func (c *ReturnTypeWillChangeClass) GetMethods() []data.Method {
	return []data.Method{
		c.construct,
	}
}

func (c *ReturnTypeWillChangeClass) GetConstruct() data.Method {
	return c.construct
}

// ReturnTypeWillChangeConstructMethod ReturnTypeWillChange 构造函数
// ReturnTypeWillChange::__construct()
type ReturnTypeWillChangeConstructMethod struct{}

func (m *ReturnTypeWillChangeConstructMethod) GetName() string {
	return "__construct"
}

func (m *ReturnTypeWillChangeConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ReturnTypeWillChangeConstructMethod) GetIsStatic() bool {
	return false
}

func (m *ReturnTypeWillChangeConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *ReturnTypeWillChangeConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *ReturnTypeWillChangeConstructMethod) GetReturnType() data.Types {
	return nil
}

func (m *ReturnTypeWillChangeConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// ReturnTypeWillChange 构造函数不需要特殊处理，只是标记返回类型将改变
	return data.NewNullValue(), nil
}
