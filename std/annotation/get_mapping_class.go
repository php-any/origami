package annotation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// GetMappingClass GetMapping注解类
type GetMappingClass struct {
	node.Node
	process   data.Method
	mapping   data.Method
	construct data.Method
}

func (g *GetMappingClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := newGetMapping()

	return data.NewClassValue(&GetMappingClass{
		process:   &GetMappingProcessMethod{source},
		mapping:   &GetMappingMappingMethod{source},
		construct: &GetMappingConstructMethod{source},
	}, ctx.CreateBaseContext()), nil
}

func (g *GetMappingClass) GetName() string {
	return "Annotation\\GetMapping"
}

func (g *GetMappingClass) GetExtend() *string {
	return nil
}

func (g *GetMappingClass) GetImplements() []string {
	return []string{node.TypeFeature} // 特性注解
}

func (g *GetMappingClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (g *GetMappingClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (g *GetMappingClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "process":
		return g.process, true
	case "mapping":
		return g.mapping, true
	case "__construct":
		return g.construct, true
	}
	return nil, false
}

func (g *GetMappingClass) GetMethods() []data.Method {
	return []data.Method{
		g.process,
		g.mapping,
		g.construct,
	}
}

func (g *GetMappingClass) GetConstruct() data.Method {
	return g.construct
}

// GetMapping 映射实例
type GetMapping struct {
	path   string
	target interface{} // 被注解的反射实例
}

func newGetMapping() *GetMapping {
	return &GetMapping{}
}

// GetMappingConstructMethod 构造函数
type GetMappingConstructMethod struct {
	mapping *GetMapping
}

func (m *GetMappingConstructMethod) GetName() string {
	return "__construct"
}

func (m *GetMappingConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *GetMappingConstructMethod) GetIsStatic() bool {
	return false
}

func (m *GetMappingConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "path", 0, data.NewStringValue("/"), data.NewBaseType("string")),
	}
}

func (m *GetMappingConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "path", 0, nil),
	}
}

func (m *GetMappingConstructMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *GetMappingConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 构造函数逻辑：从上下文获取注解参数和目标
	// 这里应该从ctx中获取传入的参数
	// 简化实现，返回构造成功消息
	return data.NewStringValue("GetMapping annotation constructed with parameters"), nil
}

// GetMappingProcessMethod 处理注解的方法
type GetMappingProcessMethod struct {
	mapping *GetMapping
}

func (m *GetMappingProcessMethod) GetName() string {
	return "process"
}

func (m *GetMappingProcessMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *GetMappingProcessMethod) GetIsStatic() bool {
	return false
}

func (m *GetMappingProcessMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *GetMappingProcessMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *GetMappingProcessMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *GetMappingProcessMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 实现GetMapping注解处理逻辑
	// 可以访问 m.mapping.arguments 和 m.mapping.target
	return data.NewStringValue("GetMapping processed"), nil
}

// GetMappingMappingMethod 注册映射的方法
type GetMappingMappingMethod struct {
	mapping *GetMapping
}

func (m *GetMappingMappingMethod) GetName() string {
	return "mapping"
}

func (m *GetMappingMappingMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *GetMappingMappingMethod) GetIsStatic() bool {
	return false
}

func (m *GetMappingMappingMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *GetMappingMappingMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *GetMappingMappingMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *GetMappingMappingMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 实现GET路由映射逻辑
	// 可以访问 m.mapping.arguments 和 m.mapping.target
	return data.NewStringValue("GET mapping registered"), nil
}
