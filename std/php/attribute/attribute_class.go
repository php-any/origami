package attribute

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// AttributeClass PHP 8.0+ 原生 Attribute 类
// 用于标记一个类可以作为属性（注解）使用
type AttributeClass struct {
	node.Node
	construct data.Method
}

func NewAttributeClass() *AttributeClass {
	return &AttributeClass{
		construct: &AttributeConstructMethod{},
	}
}

func (c *AttributeClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *AttributeClass) GetName() string {
	return "Attribute"
}

func (c *AttributeClass) GetExtend() *string {
	return nil
}

func (c *AttributeClass) GetImplements() []string {
	return []string{}
}

func (c *AttributeClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (c *AttributeClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (c *AttributeClass) GetStaticProperty(name string) (data.Value, bool) {
	switch name {
	case "TARGET_CLASS":
		return data.NewIntValue(1), true
	case "TARGET_FUNCTION":
		return data.NewIntValue(2), true
	case "TARGET_METHOD":
		return data.NewIntValue(4), true
	case "TARGET_PROPERTY":
		return data.NewIntValue(8), true
	case "TARGET_CLASS_CONSTANT":
		return data.NewIntValue(16), true
	case "TARGET_PARAMETER":
		return data.NewIntValue(32), true
	case "TARGET_ALL":
		return data.NewIntValue(63), true
	}
	return nil, false
}

func (c *AttributeClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return c.construct, true
	}
	return nil, false
}

func (c *AttributeClass) GetMethods() []data.Method {
	return []data.Method{
		c.construct,
	}
}

func (c *AttributeClass) GetConstruct() data.Method {
	return c.construct
}

// AttributeConstructMethod Attribute 构造函数
// Attribute::__construct(int $flags = Attribute::TARGET_ALL)
type AttributeConstructMethod struct{}

func (m *AttributeConstructMethod) GetName() string {
	return "__construct"
}

func (m *AttributeConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *AttributeConstructMethod) GetIsStatic() bool {
	return false
}

func (m *AttributeConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "flags", 0, node.NewIntLiteral(nil, "1"), data.NewBaseType("int")), // TARGET_ALL = 1
	}
}

func (m *AttributeConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "flags", 0, data.NewBaseType("int")),
	}
}

func (m *AttributeConstructMethod) GetReturnType() data.Types {
	return nil
}

func (m *AttributeConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// Attribute 构造函数不需要特殊处理，只是标记类可以作为属性使用
	return data.NewNullValue(), nil
}

// Attribute 常量
const (
	TARGET_CLASS          = 1  // 可以应用于类
	TARGET_FUNCTION       = 2  // 可以应用于函数
	TARGET_METHOD         = 4  // 可以应用于方法
	TARGET_PROPERTY       = 8  // 可以应用于属性
	TARGET_CLASS_CONSTANT = 16 // 可以应用于类常量
	TARGET_PARAMETER      = 32 // 可以应用于参数
	TARGET_ALL            = 63 // 可以应用于所有目标
)
