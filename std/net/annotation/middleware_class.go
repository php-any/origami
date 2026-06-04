package annotation

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// MiddlewareClass @Middleware 注解类 - 特性注解
// 用于标记控制器或方法使用的中间件类
//
// 用法:
//
//	#[Middleware(AuthInterceptor::class)]
//	class UserController { }
type MiddlewareClass struct {
	node.Node
	source    *Middleware
	construct data.Method
}

func (m *MiddlewareClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := newMiddleware()
	return data.NewClassValue(&MiddlewareClass{
		source:    source,
		construct: &MiddlewareConstructMethod{source},
	}, ctx.CreateBaseContext()), nil
}

func (m *MiddlewareClass) GetName() string { return "Net\\Annotation\\Middleware" }

func (m *MiddlewareClass) GetExtend() *string {
	return nil
}

func (m *MiddlewareClass) GetImplements() []string {
	return []string{node.TypeFeature}
}

// GetProperty 返回动态属性，从 source 中读取 className
func (m *MiddlewareClass) GetProperty(name string) (data.Property, bool) {
	if name == "className" && m.source != nil {
		return &middlewareClassNameProperty{source: m.source}, true
	}
	return nil, false
}

func (m *MiddlewareClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (m *MiddlewareClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return m.construct, true
	}
	return nil, false
}

func (m *MiddlewareClass) GetMethods() []data.Method {
	return []data.Method{m.construct}
}

func (m *MiddlewareClass) GetConstruct() data.Method {
	return m.construct
}

// ClassName 获取中间件类名
func (m *MiddlewareClass) ClassName() string {
	if m.source != nil {
		return m.source.className
	}
	return ""
}

// middlewareClassNameProperty 动态属性，从 source 中读取 className
type middlewareClassNameProperty struct {
	source *Middleware
}

func (p *middlewareClassNameProperty) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(p.source.className), nil
}

func (p *middlewareClassNameProperty) GetName() string                { return "className" }
func (p *middlewareClassNameProperty) GetModifier() data.Modifier     { return data.ModifierPublic }
func (p *middlewareClassNameProperty) GetIsStatic() bool              { return false }
func (p *middlewareClassNameProperty) GetDefaultValue() data.GetValue { return data.NewStringValue("") }
func (p *middlewareClassNameProperty) GetType() data.Types            { return data.NewBaseType("string") }
func (p *middlewareClassNameProperty) SetType(data.Types)             {}
func (p *middlewareClassNameProperty) GetZVal(ctx data.GetPropertyZVal) (*data.ZVal, data.Control) {
	return data.NewZVal(data.NewStringValue(p.source.className)), nil
}

// Middleware 中间件实例
type Middleware struct {
	className string
	target    interface{}
}

func newMiddleware() *Middleware {
	return &Middleware{}
}

// MiddlewareConstructMethod 构造函数
type MiddlewareConstructMethod struct {
	middleware *Middleware
}

func (m *MiddlewareConstructMethod) GetName() string {
	return "__construct"
}

func (m *MiddlewareConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *MiddlewareConstructMethod) GetIsStatic() bool {
	return false
}

func (m *MiddlewareConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "className", 0, nil, nil),
		node.NewParameter(nil, node.TargetName, 1, nil, nil),
	}
}

func (m *MiddlewareConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "className", 0, nil),
		node.NewVariable(nil, "target", 1, nil),
	}
}

func (m *MiddlewareConstructMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *MiddlewareConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return data.NewStringValue("Middleware annotation constructed"), nil
	}

	if sv, ok := a0.(*data.StringValue); ok {
		m.middleware.className = sv.AsString()
	}

	// 读取 target（被注解的类）
	if tv, ok := ctx.GetIndexValue(1); ok {
		if anyT, ok := tv.(*data.AnyValue); ok {
			if cls, ok := anyT.Value.(*node.ClassStatement); ok {
				// 注册到控制器中间件
				fmt.Printf("[DEBUG] MiddlewareConstructMethod: controller=%s, className=%s\n", cls.GetName(), m.middleware.className)
				AddControllerMiddleware(cls.GetName(), m.middleware.className)
			}
		}
	}

	return data.NewStringValue("Middleware annotation constructed"), nil
}
