package attribute

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// AllowDynamicPropertiesClass PHP 8.2+ 原生 AllowDynamicProperties 类
// 用于允许类使用动态属性（PHP 8.2+ 默认不允许动态属性）
type AllowDynamicPropertiesClass struct {
	node.Node
	construct data.Method
}

func NewAllowDynamicPropertiesClass() *AllowDynamicPropertiesClass {
	return &AllowDynamicPropertiesClass{
		construct: &AllowDynamicPropertiesConstructMethod{},
	}
}

func (c *AllowDynamicPropertiesClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *AllowDynamicPropertiesClass) GetName() string {
	return "AllowDynamicProperties"
}

func (c *AllowDynamicPropertiesClass) GetExtend() *string {
	return nil
}

func (c *AllowDynamicPropertiesClass) GetImplements() []string {
	return []string{}
}

func (c *AllowDynamicPropertiesClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (c *AllowDynamicPropertiesClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (c *AllowDynamicPropertiesClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return c.construct, true
	}
	return nil, false
}

func (c *AllowDynamicPropertiesClass) GetMethods() []data.Method {
	return []data.Method{
		c.construct,
	}
}

func (c *AllowDynamicPropertiesClass) GetConstruct() data.Method {
	return c.construct
}

// AllowDynamicPropertiesConstructMethod AllowDynamicProperties 构造函数
// AllowDynamicProperties::__construct()
type AllowDynamicPropertiesConstructMethod struct{}

func (m *AllowDynamicPropertiesConstructMethod) GetName() string {
	return "__construct"
}

func (m *AllowDynamicPropertiesConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *AllowDynamicPropertiesConstructMethod) GetIsStatic() bool {
	return false
}

func (m *AllowDynamicPropertiesConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *AllowDynamicPropertiesConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *AllowDynamicPropertiesConstructMethod) GetReturnType() data.Types {
	return nil
}

func (m *AllowDynamicPropertiesConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// AllowDynamicProperties 构造函数不需要特殊处理，只是标记类允许动态属性
	return data.NewNullValue(), nil
}
