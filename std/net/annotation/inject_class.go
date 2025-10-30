package annotation

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// InjectClass Inject注解类 - 宏注解
type InjectClass struct {
	node.Node
	process   data.Method
	inject    data.Method
	construct data.Method
}

func (i *InjectClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := newInject()

	return data.NewClassValue(&InjectClass{
		process:   &InjectProcessMethod{source},
		inject:    &InjectInjectMethod{source},
		construct: &InjectConstructMethod{source},
	}, ctx.CreateBaseContext()), nil
}

func (i *InjectClass) GetName() string {
	return "Annotation\\Inject"
}

func (i *InjectClass) GetExtend() *string {
	return nil
}

func (i *InjectClass) GetImplements() []string {
	return []string{node.TypeMacro} // 宏注解
}

func (i *InjectClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (i *InjectClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (i *InjectClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "process":
		return i.process, true
	case "inject":
		return i.inject, true
	case "__construct":
		return i.construct, true
	}
	return nil, false
}

func (i *InjectClass) GetMethods() []data.Method {
	return []data.Method{
		i.process,
		i.inject,
		i.construct,
	}
}

func (i *InjectClass) GetConstruct() data.Method {
	return i.construct
}

// Inject 注入实例
type Inject struct {
	service string
	target  any // 被注解的节点
}

func newInject() *Inject {
	return &Inject{}
}

// InjectConstructMethod 构造函数 - 宏注解能接收被注解的节点
type InjectConstructMethod struct {
	inject *Inject
}

func (m *InjectConstructMethod) GetName() string {
	return "__construct"
}

func (m *InjectConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *InjectConstructMethod) GetIsStatic() bool {
	return false
}

func (m *InjectConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "service", 0, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, node.TargetName, 1, nil, nil), // 被注解的节点
	}
}

func (m *InjectConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "service", 0, nil),
		node.NewVariable(nil, "target", 1, nil),
	}
}

func (m *InjectConstructMethod) GetReturnType() data.Types {
	return nil
}

func (m *InjectConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 宏注解构造函数：接收注解参数和被注解的节点
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, service: 0"))
	}

	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, target: 1"))
	}

	service := ""
	if v, ok := a0.(*data.StringValue); ok {
		service = v.AsString()
	}

	// 获取被注解的节点
	target := a1

	m.inject.service = service
	m.inject.target = target.(*data.AnyValue).Value

	// 检查属性是否有默认值, 没有没有就可以修改节点
	// TODO
	switch t := m.inject.target.(type) {
	case *node.ClassProperty:
		if t.DefaultValue == nil {
			t.DefaultValue = data.NewStringValue("TODO inject 注解未实现具体功能")
		}
	}

	return nil, nil
}

// InjectProcessMethod 处理注解的方法
type InjectProcessMethod struct {
	inject *Inject
}

func (m *InjectProcessMethod) GetName() string {
	return "process"
}

func (m *InjectProcessMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *InjectProcessMethod) GetIsStatic() bool {
	return false
}

func (m *InjectProcessMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *InjectProcessMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *InjectProcessMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *InjectProcessMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 宏注解处理逻辑
	return data.NewStringValue("Inject processed with service: " + m.inject.service), nil
}

// InjectInjectMethod 执行注入的方法
type InjectInjectMethod struct {
	inject *Inject
}

func (m *InjectInjectMethod) GetName() string {
	return "inject"
}

func (m *InjectInjectMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *InjectInjectMethod) GetIsStatic() bool {
	return false
}

func (m *InjectInjectMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *InjectInjectMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *InjectInjectMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *InjectInjectMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 宏注解依赖注入逻辑
	return data.NewStringValue("Dependency injected for service: " + m.inject.service), nil
}
