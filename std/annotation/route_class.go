package annotation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RouteClass Route注解类
type RouteClass struct {
	node.Node
	process   data.Method
	register  data.Method
	construct data.Method
}

func (r *RouteClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := newRoute()

	return data.NewClassValue(&RouteClass{
		process:   &RouteProcessMethod{source},
		register:  &RouteRegisterMethod{source},
		construct: &RouteConstructMethod{source},
	}, ctx.CreateBaseContext()), nil
}

func (r *RouteClass) GetName() string {
	return "Annotation\\Route"
}

func (r *RouteClass) GetExtend() *string {
	return nil
}

func (r *RouteClass) GetImplements() []string {
	return []string{node.TypeFeature} // 特性注解
}

func (r *RouteClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (r *RouteClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (r *RouteClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "process":
		return r.process, true
	case "register":
		return r.register, true
	case "__construct":
		return r.construct, true
	}
	return nil, false
}

func (r *RouteClass) GetMethods() []data.Method {
	return []data.Method{
		r.process,
		r.register,
		r.construct,
	}
}

func (r *RouteClass) GetConstruct() data.Method {
	return r.construct
}

// Route 路由实例
type Route struct {
	prefix string
	target interface{} // 被注解的反射实例
}

func newRoute() *Route {
	return &Route{}
}

// RouteConstructMethod 构造函数
type RouteConstructMethod struct {
	route *Route
}

func (m *RouteConstructMethod) GetName() string {
	return "__construct"
}

func (m *RouteConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *RouteConstructMethod) GetIsStatic() bool {
	return false
}

func (m *RouteConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "prefix", 0, data.NewStringValue("/"), data.NewBaseType("string")),
	}
}

func (m *RouteConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "prefix", 0, nil),
	}
}

func (m *RouteConstructMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *RouteConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 构造函数逻辑：从上下文获取注解参数和目标
	// 这里应该从ctx中获取传入的参数
	// 简化实现，返回构造成功消息
	return data.NewStringValue("Route annotation constructed with parameters"), nil
}

// RouteProcessMethod 处理注解的方法
type RouteProcessMethod struct {
	route *Route
}

func (m *RouteProcessMethod) GetName() string {
	return "process"
}

func (m *RouteProcessMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *RouteProcessMethod) GetIsStatic() bool {
	return false
}

func (m *RouteProcessMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *RouteProcessMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *RouteProcessMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *RouteProcessMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 实现路由注解处理逻辑
	// 可以访问 m.route.arguments 和 m.route.target
	return data.NewStringValue("Route processed"), nil
}

// RouteRegisterMethod 注册路由的方法
type RouteRegisterMethod struct {
	route *Route
}

func (m *RouteRegisterMethod) GetName() string {
	return "register"
}

func (m *RouteRegisterMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *RouteRegisterMethod) GetIsStatic() bool {
	return false
}

func (m *RouteRegisterMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *RouteRegisterMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *RouteRegisterMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *RouteRegisterMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 实现路由注册逻辑
	// 可以访问 m.route.arguments 和 m.route.target
	return data.NewStringValue("Route registered"), nil
}
