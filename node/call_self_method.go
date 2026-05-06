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
	// 检查是否在类上下文中（类方法或类级初始化器）
	var currentClass data.ClassStmt
	if classCtx, ok := ctx.(*data.ClassMethodContext); ok {
		currentClass = classCtx.Class
	} else if classVal, ok := ctx.(*data.ClassValue); ok {
		currentClass = classVal.Class
	} else {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("self:: 只能在类方法中使用"))
	}

	// 检查类是否实现了 GetStaticMethod 接口
	getter, ok := currentClass.(data.GetStaticMethod)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("当前类 %s 不支持静态方法访问", currentClass.GetName()))
	}

	// 获取当前类的静态方法
	method, has := getter.GetStaticMethod(pe.Method)
	if !has {
		// 沿继承链向上查找（trait 中 self:: 应能访问使用类继承链上的方法）
		vm := ctx.GetVM()
		extend := currentClass.GetExtend()
		for extend != nil {
			parent, acl := vm.GetOrLoadClass(*extend)
			if acl != nil || parent == nil {
				break
			}
			if parentGetter, ok := parent.(data.GetStaticMethod); ok {
				if m, ok := parentGetter.GetStaticMethod(pe.Method); ok {
					method = m
					has = true
					break
				}
			}
			extend = parent.GetExtend()
		}
	}
	if !has {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("当前类 %s 没有静态方法 %s", currentClass.GetName(), pe.Method))
	}

	return data.NewFuncValue(method), nil
}
