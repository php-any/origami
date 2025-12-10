package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

// CallSelfProperty 表示当前类静态属性访问表达式
type CallSelfProperty struct {
	*Node    `pp:"-"`
	Property string // 属性名
}

func NewCallSelfProperty(from data.From, property string) *CallSelfProperty {
	return &CallSelfProperty{
		Node:     NewNode(from),
		Property: property,
	}
}

// GetValue 获取当前类静态属性访问表达式的值
func (pe *CallSelfProperty) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 检查是否在类方法上下文中
	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("self:: 只能在类方法中使用"))
	}

	// 获取当前类
	currentClass := classCtx.Class

	// 检查类是否实现了 GetStaticProperty 接口
	getter, ok := currentClass.(data.GetStaticProperty)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("当前类 %s 不支持静态属性访问", currentClass.GetName()))
	}

	// 获取当前类的静态属性
	property, has := getter.GetStaticProperty(pe.Property)
	if !has {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("当前类 %s 没有静态属性 %s", currentClass.GetName(), pe.Property))
	}

	return property, nil
}

// SetProperty 设置当前类静态属性的值
func (pe *CallSelfProperty) SetProperty(ctx data.Context, name string, value data.Value) data.Control {
	// 检查是否在类方法上下文中
	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return data.NewErrorThrow(pe.GetFrom(), errors.New("self:: 只能在类方法中使用"))
	}

	// 获取当前类
	currentClass := classCtx.Class

	// 检查类是否实现了 SetProperty 接口
	if setter, ok := currentClass.(data.SetProperty); ok {
		return setter.SetProperty(pe.Property, value)
	}

	return data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("当前类 %s 无法设置静态属性 %s", currentClass.GetName(), pe.Property))
}
