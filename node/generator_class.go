package node

import (
	"github.com/php-any/origami/data"
)

// GeneratorClass 表示生成器类
type GeneratorClass struct {
	generator data.Generator
}

// NewGeneratorClass 创建一个新的生成器类实例
func NewGeneratorClass(generator data.Generator) data.ClassStmt {
	return &GeneratorClass{
		generator: generator,
	}
}

// GetValue 返回类实例
func (g *GeneratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(g, ctx.CreateBaseContext()), nil
}

// GetFrom 返回来源信息
func (g *GeneratorClass) GetFrom() data.From {
	return nil
}

// GetName 返回类名
func (g *GeneratorClass) GetName() string {
	return "Generator"
}

// GetExtend 返回父类
func (g *GeneratorClass) GetExtend() *string {
	return nil
}

// GetImplements 返回实现的接口
func (g *GeneratorClass) GetImplements() []string {
	return nil
}

// GetProperty 获取属性
func (g *GeneratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

// GetPropertyList 获取所有属性列表
func (g *GeneratorClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

// GetMethod 获取方法
func (g *GeneratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "valid":
		return &GeneratorValidMethod{generator: g.generator}, true
	case "current":
		return &GeneratorCurrentMethod{generator: g.generator}, true
	case "key":
		return &GeneratorKeyMethod{generator: g.generator}, true
	case "next":
		return &GeneratorNextMethod{generator: g.generator}, true
	case "rewind":
		return &GeneratorRewindMethod{generator: g.generator}, true
	case "send":
		return &GeneratorSendMethod{generator: g.generator}, true
	case "throw":
		return &GeneratorThrowMethod{generator: g.generator}, true
	case "getReturn":
		return &GeneratorGetReturnMethod{generator: g.generator}, true
	}
	return nil, false
}

// GetMethods 获取所有方法
func (g *GeneratorClass) GetMethods() []data.Method {
	return []data.Method{
		&GeneratorValidMethod{generator: g.generator},
		&GeneratorCurrentMethod{generator: g.generator},
		&GeneratorKeyMethod{generator: g.generator},
		&GeneratorNextMethod{generator: g.generator},
		&GeneratorRewindMethod{generator: g.generator},
		&GeneratorSendMethod{generator: g.generator},
		&GeneratorThrowMethod{generator: g.generator},
		&GeneratorGetReturnMethod{generator: g.generator},
	}
}

// GetConstruct 获取构造函数
func (g *GeneratorClass) GetConstruct() data.Method {
	return nil
}

// GeneratorValidMethod 实现 valid() 方法
type GeneratorValidMethod struct {
	generator data.Generator
}

func (m *GeneratorValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return m.generator.Valid(ctx)
}

func (m *GeneratorValidMethod) GetName() string {
	return "valid"
}

func (m *GeneratorValidMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *GeneratorValidMethod) GetIsStatic() bool {
	return false
}

func (m *GeneratorValidMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *GeneratorValidMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *GeneratorValidMethod) GetReturnType() data.Types {
	return data.NewBaseType("bool")
}

// GeneratorCurrentMethod 实现 current() 方法
type GeneratorCurrentMethod struct {
	generator data.Generator
}

func (m *GeneratorCurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return m.generator.Current(ctx)
}

func (m *GeneratorCurrentMethod) GetName() string {
	return "current"
}

func (m *GeneratorCurrentMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *GeneratorCurrentMethod) GetIsStatic() bool {
	return false
}

func (m *GeneratorCurrentMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *GeneratorCurrentMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *GeneratorCurrentMethod) GetReturnType() data.Types {
	return data.NewBaseType("mixed")
}

// GeneratorKeyMethod 实现 key() 方法
type GeneratorKeyMethod struct {
	generator data.Generator
}

func (m *GeneratorKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return m.generator.Key(ctx)
}

func (m *GeneratorKeyMethod) GetName() string {
	return "key"
}

func (m *GeneratorKeyMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *GeneratorKeyMethod) GetIsStatic() bool {
	return false
}

func (m *GeneratorKeyMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *GeneratorKeyMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *GeneratorKeyMethod) GetReturnType() data.Types {
	return data.NewBaseType("mixed")
}

// GeneratorNextMethod 实现 next() 方法
type GeneratorNextMethod struct {
	generator data.Generator
}

func (m *GeneratorNextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	ctl := m.generator.Next(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewNullValue(), nil
}

func (m *GeneratorNextMethod) GetName() string {
	return "next"
}

func (m *GeneratorNextMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *GeneratorNextMethod) GetIsStatic() bool {
	return false
}

func (m *GeneratorNextMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *GeneratorNextMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *GeneratorNextMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}

// GeneratorRewindMethod 实现 rewind() 方法
type GeneratorRewindMethod struct {
	generator data.Generator
}

func (m *GeneratorRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return m.generator.Rewind(ctx)
}

func (m *GeneratorRewindMethod) GetName() string {
	return "rewind"
}

func (m *GeneratorRewindMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *GeneratorRewindMethod) GetIsStatic() bool {
	return false
}

func (m *GeneratorRewindMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *GeneratorRewindMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *GeneratorRewindMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}

// GeneratorSendMethod 实现 send() 方法
type GeneratorSendMethod struct {
	generator data.Generator
}

func (m *GeneratorSendMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	if value == nil {
		value = data.NewNullValue()
	}
	ctl := m.generator.Send(ctx, value)
	if ctl != nil {
		return nil, ctl
	}
	return m.generator.Current(ctx)
}

func (m *GeneratorSendMethod) GetName() string {
	return "send"
}

func (m *GeneratorSendMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *GeneratorSendMethod) GetIsStatic() bool {
	return false
}

func (m *GeneratorSendMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		NewParameter(nil, "value", 0, nil, nil),
	}
}

func (m *GeneratorSendMethod) GetVariables() []data.Variable {
	return []data.Variable{
		NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
	}
}

func (m *GeneratorSendMethod) GetReturnType() data.Types {
	return data.NewBaseType("mixed")
}

// GeneratorThrowMethod 实现 throw() 方法
type GeneratorThrowMethod struct {
	generator data.Generator
}

func (m *GeneratorThrowMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	ctl := m.generator.Throw(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewNullValue(), nil
}

func (m *GeneratorThrowMethod) GetName() string {
	return "throw"
}

func (m *GeneratorThrowMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *GeneratorThrowMethod) GetIsStatic() bool {
	return false
}

func (m *GeneratorThrowMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *GeneratorThrowMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *GeneratorThrowMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}

// GeneratorGetReturnMethod 实现 getReturn() 方法
type GeneratorGetReturnMethod struct {
	generator data.Generator
}

func (m *GeneratorGetReturnMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return m.generator.GetReturn(ctx)
}

func (m *GeneratorGetReturnMethod) GetName() string {
	return "getReturn"
}

func (m *GeneratorGetReturnMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *GeneratorGetReturnMethod) GetIsStatic() bool {
	return false
}

func (m *GeneratorGetReturnMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *GeneratorGetReturnMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *GeneratorGetReturnMethod) GetReturnType() data.Types {
	return data.NewBaseType("mixed")
}
