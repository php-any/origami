package annotation

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewColumnClass 创建 Column 注解类
func NewColumnClass() *ColumnClass {
	return &ColumnClass{}
}

// ColumnClass Column注解类 - 特性注解
type ColumnClass struct {
	node.Node
	construct data.Method

	name     string
	nullable bool
	length   int
}

func (c *ColumnClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 创建新的注解实例，直接使用字段存储数据
	instance := &ColumnClass{}

	// 创建构造函数，传入注解实例引用
	instance.construct = &ColumnConstructMethod{columnClass: instance}

	return data.NewClassValue(instance, ctx.CreateBaseContext()), nil
}

func (c *ColumnClass) GetName() string {
	return "database\\annotation\\Column"
}

func (c *ColumnClass) GetExtend() *string {
	return nil
}

func (c *ColumnClass) GetImplements() []string {
	return []string{node.TypeFeature} // 特性注解
}

func (c *ColumnClass) GetProperty(name string) (data.Property, bool) {
	// 返回注解实例的属性
	switch name {
	case "name":
		return node.NewProperty(nil, "name", "public", false, data.NewStringValue(c.name)), true
	case "nullable":
		return node.NewProperty(nil, "nullable", "public", false, data.NewBoolValue(c.nullable)), true
	case "length":
		return node.NewProperty(nil, "length", "public", false, data.NewIntValue(c.length)), true
	}
	return nil, false
}

func (c *ColumnClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "name", "public", false, data.NewStringValue(c.name)),
		node.NewProperty(nil, "length", "public", false, data.NewIntValue(c.length)),
		node.NewProperty(nil, "nullable", "public", false, data.NewBoolValue(c.nullable)),
	}
}

func (c *ColumnClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return c.construct, true
	}
	return nil, false
}

func (c *ColumnClass) GetMethods() []data.Method {
	return []data.Method{
		c.construct,
	}
}

func (c *ColumnClass) GetConstruct() data.Method {
	return c.construct
}

// ColumnConstructMethod 构造函数 - 特性注解只接收注解参数
type ColumnConstructMethod struct {
	columnClass *ColumnClass // 引用注解实例
}

func (m *ColumnConstructMethod) GetName() string {
	return "__construct"
}

func (m *ColumnConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ColumnConstructMethod) GetIsStatic() bool {
	return false
}

func (m *ColumnConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, data.NewNullValue(), data.NewBaseType("string")),
		node.NewParameter(nil, "nullable", 1, data.NewBoolValue(true), data.NewBaseType("bool")),
		node.NewParameter(nil, "length", 2, data.NewIntValue(255), data.NewBaseType("int")),
	}
}

func (m *ColumnConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, nil),
		node.NewVariable(nil, "nullable", 1, nil),
		node.NewVariable(nil, "length", 2, nil),
	}
}

func (m *ColumnConstructMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *ColumnConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 特性注解构造函数：只接收注解声明的参数
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, name: 0"))
	}

	name := ""
	if v, ok := a0.(*data.StringValue); ok {
		name = v.AsString()
	}

	nullable := true
	if a1, ok := ctx.GetIndexValue(1); ok {
		if v, ok := a1.(*data.BoolValue); ok {
			nullable, _ = v.AsBool()
		}
	}

	length := 255
	if a2, ok := ctx.GetIndexValue(2); ok {
		if v, ok := a2.(*data.IntValue); ok {
			length, _ = v.AsInt()
		}
	}

	// 将数据直接存储到注解实例的字段上
	if m.columnClass != nil {
		m.columnClass.name = name
		m.columnClass.nullable = nullable
		m.columnClass.length = length
	}

	// 构造函数只负责初始化，返回 nil 表示成功
	return nil, nil
}
