package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

// CallStaticKeywordMethod 表示 static::method() （late static binding 风格）的静态方法调用表达式
// 注意：当前实现语义上仍等同于 self::method()，但通过单独节点类型与 self:: 区分，便于后续增强
type CallStaticKeywordMethod struct {
	*Node  `pp:"-"`
	Method string // 方法名
}

func NewCallStaticKeywordMethod(from data.From, method string) *CallStaticKeywordMethod {
	return &CallStaticKeywordMethod{
		Node:   NewNode(from),
		Method: method,
	}
}

// GetValue 获取 static::method() 调用的值
func (pe *CallStaticKeywordMethod) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 与 self:: 一样，必须在类方法上下文中使用
	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("static:: 只能在类方法中使用"))
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
		extend := currentClass.GetExtend()
		for extend != nil {
			vm := ctx.GetVM()
			ext, acl := vm.GetOrLoadClass(*extend)
			if acl != nil {
				return nil, acl
			}
			extend = nil
			getter, ok = ext.(data.GetStaticMethod)
			if ok {
				method, has = getter.GetStaticMethod(pe.Method)
				if has {
					return data.NewFuncValue(method), nil
				}
				extend = ext.GetExtend()
			}
		}

		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("当前类 %s 没有静态方法 %s", currentClass.GetName(), pe.Method))
	}

	return data.NewFuncValue(method), nil
}
