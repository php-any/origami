package annotation

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewTableClass 创建 Table 注解类
func NewTableClass() *TableClass {
	return &TableClass{}
}

// TableClass Table注解类 - 特性注解
type TableClass struct {
	node.Node
	construct data.Method
	name      string // 直接存储注解数据
}

func (t *TableClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 创建新的注解实例，直接使用字段存储数据
	instance := &TableClass{}

	// 创建构造函数，传入注解实例引用
	instance.construct = &TableConstructMethod{tableClass: instance}

	return data.NewClassValue(instance, ctx.CreateBaseContext()), nil
}

func (t *TableClass) GetName() string {
	return "database\\annotation\\Table"
}

func (t *TableClass) GetExtend() *string {
	return nil
}

func (t *TableClass) GetImplements() []string {
	return []string{node.TypeFeature} // 特性注解
}

func (t *TableClass) GetProperty(name string) (data.Property, bool) {
	// 返回注解实例的属性
	switch name {
	case "name":
		return node.NewProperty(nil, "name", "public", false, data.NewStringValue(t.name)), true
	}
	return nil, false
}

func (t *TableClass) GetProperties() map[string]data.Property {
	properties := make(map[string]data.Property)

	properties["name"] = node.NewProperty(nil, "name", "public", false, data.NewStringValue(t.name))

	return properties
}

func (t *TableClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return t.construct, true
	}
	return nil, false
}

func (t *TableClass) GetMethods() []data.Method {
	return []data.Method{
		t.construct,
	}
}

func (t *TableClass) GetConstruct() data.Method {
	return t.construct
}

// TableConstructMethod 构造函数 - 特性注解只接收注解参数
type TableConstructMethod struct {
	tableClass *TableClass // 引用注解实例
}

func (m *TableConstructMethod) GetName() string {
	return "__construct"
}

func (m *TableConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *TableConstructMethod) GetIsStatic() bool {
	return false
}

func (m *TableConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, data.NewNullValue(), data.NewBaseType("string")),
	}
}

func (m *TableConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, nil),
	}
}

func (m *TableConstructMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *TableConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 特性注解构造函数：只接收注解声明的参数
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, name: 0"))
	}

	name := ""
	if v, ok := a0.(*data.StringValue); ok {
		name = v.AsString()
	}

	// 将数据直接存储到注解实例的字段上
	if m.tableClass != nil {
		m.tableClass.name = name
	}

	// 构造函数只负责初始化，返回 nil 表示成功
	return nil, nil
}
