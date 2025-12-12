package core

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// ClosureClass 提供 PHP Closure 相关静态方法
type ClosureClass struct {
	node.Node
	bind data.Method
}

func (c *ClosureClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if c.bind == nil {
		c.bind = &ClosureBindMethod{}
	}
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *ClosureClass) GetName() string { return "Closure" }

func (c *ClosureClass) GetExtend() *string { return nil }

func (c *ClosureClass) GetImplements() []string { return nil }

func (c *ClosureClass) GetProperty(name string) (data.Property, bool) { return nil, false }

func (c *ClosureClass) GetPropertyList() []data.Property { return nil }

func (c *ClosureClass) GetMethod(name string) (data.Method, bool) { return nil, false }

func (c *ClosureClass) GetMethods() []data.Method { return nil }

// 静态方法
func (c *ClosureClass) GetStaticMethod(name string) (data.Method, bool) {
	if c.bind == nil {
		c.bind = &ClosureBindMethod{}
	}
	switch name {
	case "bind":
		return c.bind, true
	}
	return nil, false
}

// GetConstruct 无构造函数
func (c *ClosureClass) GetConstruct() data.Method { return nil }

// ClosureBindMethod 实现 Closure::bind
// 当前简化实现：返回原始闭包，不改变绑定对象/作用域
type ClosureBindMethod struct{}

func (m *ClosureBindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	closureVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数: closure"))
	}

	closureThis, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数: newThis"))
	}
	if closureThis == nil {
		return nil, utils.NewThrow(errors.New("newThis TODO"))
	}

	// 仅接受可调用类型
	switch closureVal.(type) {
	case *data.FuncValue:
		return closureVal, nil
	default:
		return nil, utils.NewThrow(errors.New("Closure::bind 需要传入闭包/可调用类型"))
	}
}

func (m *ClosureBindMethod) GetName() string { return "bind" }

func (m *ClosureBindMethod) GetModifier() data.Modifier { return data.ModifierPublic }

func (m *ClosureBindMethod) GetIsStatic() bool { return true }

func (m *ClosureBindMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "closure", 0, nil, nil),
		node.NewParameter(nil, "newThis", 1, nil, nil),  // 占位，当前未使用
		node.NewParameter(nil, "newScope", 2, nil, nil), // 占位，当前未使用
	}
}

func (m *ClosureBindMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "closure", 0, data.Mixed{}),
		node.NewVariable(nil, "newThis", 1, data.Mixed{}),
		node.NewVariable(nil, "newScope", 2, data.Mixed{}),
	}
}

func (m *ClosureBindMethod) GetReturnType() data.Types { return nil }
