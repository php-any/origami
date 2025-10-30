package annotation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RouteClass Route注解类
type RouteClass struct {
	node.Node
	source    *Route
	construct data.Method
}

func (r *RouteClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := newRoute()

	return data.NewClassValue(&RouteClass{
		source:    source,
		construct: &RouteConstructMethod{source},
	}, ctx.CreateBaseContext()), nil
}

func (r *RouteClass) GetName() string { return "Net\\Annotation\\Route" }

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
	case "__construct":
		return r.construct, true
	case "prefix":
		return &RoutePrefixMethod{r.source}, true
	}
	return nil, false
}

func (r *RouteClass) GetMethods() []data.Method {
	return []data.Method{
		r.construct,
		&RoutePrefixMethod{r.source},
	}
}

func (r *RouteClass) GetConstruct() data.Method {
	return r.construct
}

// Prefix 便捷访问当前实例的前缀值
func (r *RouteClass) Prefix() string {
	if r.source != nil {
		return r.source.prefix
	}
	return ""
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
		node.NewParameter(nil, node.TargetName, 1, nil, nil),
	}
}

func (m *RouteConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "prefix", 0, nil),
		node.NewVariable(nil, "target", 1, nil),
	}
}

func (m *RouteConstructMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *RouteConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	vv, ok := ctx.GetIndexValue(0)
	if ok {
		m.route.prefix = vv.AsString()
	}
	if tv, ok := ctx.GetIndexValue(1); ok {
		if anyT, ok := tv.(*data.AnyValue); ok {
			m.route.target = anyT.Value
		}
	}
	return data.NewStringValue("Route annotation constructed"), nil
}

// 暴露 prefix 值给控制器注解读取
type RoutePrefixMethod struct{ route *Route }

func (m *RoutePrefixMethod) GetName() string            { return "prefix" }
func (m *RoutePrefixMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *RoutePrefixMethod) GetIsStatic() bool          { return false }
func (m *RoutePrefixMethod) GetParams() []data.GetValue { return []data.GetValue{} }
func (m *RoutePrefixMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (m *RoutePrefixMethod) GetReturnType() data.Types { return data.NewBaseType("string") }
func (m *RoutePrefixMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(m.route.prefix), nil
}

// RouteProcessMethod 处理注解的方法
// 删除无用的 process/register 方法
