package annotation

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewGeneratedValueClass 创建 GeneratedValue 注解类
func NewGeneratedValueClass() *GeneratedValueClass {
	return &GeneratedValueClass{}
}

// GeneratedValueClass GeneratedValue注解类 - 特性注解
type GeneratedValueClass struct {
	node.Node
	construct data.Method
	strategy  string // 直接存储注解数据
}

func (g *GeneratedValueClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 创建新的注解实例，直接使用字段存储数据
	instance := &GeneratedValueClass{}

	// 创建构造函数，传入注解实例引用
	instance.construct = &GeneratedValueConstructMethod{generatedValueClass: instance}

	return data.NewClassValue(instance, ctx.CreateBaseContext()), nil
}

func (g *GeneratedValueClass) GetName() string {
	return "Database\\Annotation\\GeneratedValue"
}

func (g *GeneratedValueClass) GetExtend() *string {
	return nil
}

func (g *GeneratedValueClass) GetImplements() []string {
	return []string{node.TypeFeature} // 特性注解
}

func (g *GeneratedValueClass) GetProperty(name string) (data.Property, bool) {
	// 返回注解实例的属性
	switch name {
	case "strategy":
		return node.NewProperty(nil, "strategy", "public", false, data.NewStringValue(g.strategy)), true
	}
	return nil, false
}

func (g *GeneratedValueClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (g *GeneratedValueClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return g.construct, true
	}
	return nil, false
}

func (g *GeneratedValueClass) GetMethods() []data.Method {
	return []data.Method{
		g.construct,
	}
}

func (g *GeneratedValueClass) GetConstruct() data.Method {
	return g.construct
}

// GeneratedValueConstructMethod 构造函数 - 特性注解只接收注解参数
type GeneratedValueConstructMethod struct {
	generatedValueClass *GeneratedValueClass // 引用注解实例
}

func (m *GeneratedValueConstructMethod) GetName() string {
	return "__construct"
}

func (m *GeneratedValueConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *GeneratedValueConstructMethod) GetIsStatic() bool {
	return false
}

func (m *GeneratedValueConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "strategy", 0, data.NewStringValue("AUTO"), data.NewBaseType("string")),
	}
}

func (m *GeneratedValueConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "strategy", 0, nil),
	}
}

func (m *GeneratedValueConstructMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *GeneratedValueConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 特性注解构造函数：只接收注解声明的参数
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, strategy: 0"))
	}

	strategy := "AUTO"
	if v, ok := a0.(*data.StringValue); ok {
		strategy = v.AsString()
	}

	// 将数据直接存储到注解实例的字段上
	if m.generatedValueClass != nil {
		m.generatedValueClass.strategy = strategy
	}

	// 构造函数只负责初始化，返回 nil 表示成功
	return nil, nil
}
