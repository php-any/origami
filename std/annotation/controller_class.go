package annotation

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ControllerClass Controller注解类 - 特性注解
type ControllerClass struct {
	node.Node
	process   data.Method
	register  data.Method
	construct data.Method
}

func (c *ControllerClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := newController()

	return data.NewClassValue(&ControllerClass{
		process:   &ControllerProcessMethod{source},
		register:  &ControllerRegisterMethod{source},
		construct: &ControllerConstructMethod{source},
	}, ctx.CreateBaseContext()), nil
}

func (c *ControllerClass) GetName() string {
	return "Annotation\\Controller"
}

func (c *ControllerClass) GetExtend() *string {
	return nil
}

func (c *ControllerClass) GetImplements() []string {
	return []string{node.TypeFeature} // 特性注解
}

func (c *ControllerClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (c *ControllerClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (c *ControllerClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "process":
		return c.process, true
	case "register":
		return c.register, true
	case "__construct":
		return c.construct, true
	}
	return nil, false
}

func (c *ControllerClass) GetMethods() []data.Method {
	return []data.Method{
		c.process,
		c.register,
		c.construct,
	}
}

func (c *ControllerClass) GetConstruct() data.Method {
	return c.construct
}

// Controller 控制器实例
type Controller struct {
	name string
}

func newController() *Controller {
	return &Controller{}
}

// ControllerConstructMethod 构造函数 - 特性注解只接收注解参数
type ControllerConstructMethod struct {
	controller *Controller
}

func (m *ControllerConstructMethod) GetName() string {
	return "__construct"
}

func (m *ControllerConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ControllerConstructMethod) GetIsStatic() bool {
	return false
}

func (m *ControllerConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, data.NewNullValue(), data.NewBaseType("string")),
	}
}

func (m *ControllerConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, nil),
	}
}

func (m *ControllerConstructMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *ControllerConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 特性注解构造函数：只接收注解声明的参数
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, name: 0"))
	}

	name := ""
	if v, ok := a0.(*data.StringValue); ok {
		name = v.AsString()
	}

	m.controller.name = name
	return data.NewStringValue("Controller annotation constructed with name: " + name), nil
}

// ControllerProcessMethod 处理注解的方法
type ControllerProcessMethod struct {
	controller *Controller
}

func (m *ControllerProcessMethod) GetName() string {
	return "process"
}

func (m *ControllerProcessMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ControllerProcessMethod) GetIsStatic() bool {
	return false
}

func (m *ControllerProcessMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *ControllerProcessMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *ControllerProcessMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *ControllerProcessMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 特性注解处理逻辑
	return data.NewStringValue("Controller processed with name: " + m.controller.name), nil
}

// ControllerRegisterMethod 注册控制器的方法
type ControllerRegisterMethod struct {
	controller *Controller
}

func (m *ControllerRegisterMethod) GetName() string {
	return "register"
}

func (m *ControllerRegisterMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ControllerRegisterMethod) GetIsStatic() bool {
	return false
}

func (m *ControllerRegisterMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *ControllerRegisterMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *ControllerRegisterMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

func (m *ControllerRegisterMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 特性注解注册逻辑
	return data.NewStringValue("Controller registered with name: " + m.controller.name), nil
}
