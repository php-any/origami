package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

// CallSelfMethod 表示当前类静态方法调用表达式
type CallSelfMethod struct {
	*Node  `pp:"-"`
	Method string // 方法名
}

func NewCallSelfMethod(from data.From, method string) *CallSelfMethod {
	return &CallSelfMethod{
		Node:   NewNode(from),
		Method: method,
	}
}

// GetValue 获取当前类静态方法调用表达式的值
func (pe *CallSelfMethod) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 检查是否在类方法上下文中
	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("self:: 只能在类方法中使用"))
	}

	// 获取当前类
	currentClass := classCtx.Class

	// 检查类是否实现了 GetStaticMethod 接口
	getter, ok := currentClass.(data.GetStaticMethod)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("当前类 %s 不支持静态方法访问", currentClass.GetName()))
	}

	// 获取当前类的静态方法
	method, has := getter.GetStaticMethod(pe.Method)
	if !has {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("当前类 %s 没有静态方法 %s", currentClass.GetName(), pe.Method))
	}

	return data.NewFuncValue(method), nil
}
