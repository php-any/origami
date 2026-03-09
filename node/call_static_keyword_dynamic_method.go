package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

// CallStaticKeywordDynamicMethod 表示 static::{expr}() 这种“动态静态方法名”的调用。
// 运行时语义：先对 expr 求值取字符串方法名，再按 static::method() 的规则分派。
type CallStaticKeywordDynamicMethod struct {
	*Node          `pp:"-"`
	NameExpression data.GetValue // 运行时求值得到方法名
}

func NewDynamicCallStaticKeywordMethod(from data.From, nameExpr data.GetValue) *CallStaticKeywordDynamicMethod {
	return &CallStaticKeywordDynamicMethod{
		Node:           NewNode(from),
		NameExpression: nameExpr,
	}
}

func (pe *CallStaticKeywordDynamicMethod) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 与 static::method 一样，必须在类方法上下文中使用
	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("static:: 只能在类方法中使用"))
	}

	// 先计算方法名表达式
	raw, acl := pe.NameExpression.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	if raw == nil {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("static::{...} 方法名表达式结果为 null"))
	}
	methodName := raw.(data.Value).AsString()
	if methodName == "" {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("static::{...} 方法名不能为空"))
	}

	// 后续逻辑与 CallStaticKeywordMethod 相同：在当前类及其继承链中查找静态方法
	currentClass := classCtx.Class

	getter, ok := currentClass.(data.GetStaticMethod)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("当前类 %s 不支持静态方法访问", currentClass.GetName()))
	}

	method, has := getter.GetStaticMethod(methodName)
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
				method, has = getter.GetStaticMethod(methodName)
				if has {
					return data.NewFuncValue(method), nil
				}
				extend = ext.GetExtend()
			}
		}

		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("当前类 %s 没有静态方法 %s", currentClass.GetName(), methodName))
	}

	return data.NewFuncValue(method), nil
}
