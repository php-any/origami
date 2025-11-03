package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

// CallParentMethod 表示父类方法调用表达式
type CallParentMethod struct {
	*Node  `pp:"-"`
	Method string // 方法名
}

func NewCallParentMethod(from data.From, method string) *CallParentMethod {
	return &CallParentMethod{
		Node:   NewNode(from),
		Method: method,
	}
}

// GetValue 获取父类方法调用表达式的值
func (pe *CallParentMethod) GetValue(ctx data.Context) (data.GetValue, data.Control) {
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
	parentClass, acl := vm.GetOrLoadClass(parentClassName)
	if acl != nil {
		return nil, acl
	}

	// 获取父类方法
	method, has := parentClass.GetMethod(pe.Method)
	if !has {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("父类 %s 没有方法 %s", parentClassName, pe.Method)))
	}

	// 检查方法访问权限
	if method.GetModifier() == data.ModifierPrivate {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("父类方法 %s 是私有的，无法访问", pe.Method)))
	}

	return data.NewFuncValue(method), nil
}
