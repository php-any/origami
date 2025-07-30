package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

// CallParentProperty 表示父类属性访问表达式
type CallParentProperty struct {
	*Node    `pp:"-"`
	Property string // 属性名
}

func NewCallParentProperty(from data.From, property string) *CallParentProperty {
	return &CallParentProperty{
		Node:     NewNode(from),
		Property: property,
	}
}

// GetValue 获取父类属性访问表达式的值
func (pe *CallParentProperty) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 检查是否在类方法上下文中
	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("parent:: 只能在类方法中使用"))
	}

	// 获取父类
	if classCtx.Class.GetExtend() == nil {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("当前类没有父类"))
	}

	parentClassName := *classCtx.Class.GetExtend()
	vm := ctx.GetVM()
	parentClass, ok := vm.GetClass(parentClassName)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("找不到父类: %s", parentClassName)))
	}

	// 获取父类属性
	property, has := parentClass.GetProperty(pe.Property)
	if !has {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("父类 %s 没有属性 %s", parentClassName, pe.Property)))
	}

	// 检查属性访问权限
	if property.GetModifier() == data.ModifierPrivate {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("父类属性 %s 是私有的，无法访问", pe.Property)))
	}

	return property.GetValue(ctx)
}
