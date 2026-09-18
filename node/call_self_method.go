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
		// self:: 是词法绑定：应解析到代码定义所在的类（trait/父类宿主），而非运行时调用类
		if classCtx.SelfClass != nil {
			currentClass = classCtx.SelfClass
		} else {
			currentClass = classCtx.Class
		}
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
		vm := callVM(ctx)
		extend := currentClass.GetExtend()
		for vm != nil && extend != nil {
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
		if fn, ok := tryNewInstanceMagicCallViaStaticFunc(ctx, pe.Method, currentClass); ok {
			return data.NewFuncValue(fn), nil
		}
		// 回退到 __callStatic（含继承链）
		checkClass := currentClass
		for checkClass != nil {
			if getter, ok := checkClass.(data.GetStaticMethod); ok {
				if magic, hasMagic := getter.GetStaticMethod("__callStatic"); hasMagic {
					return data.NewFuncValue(&callStaticFunc{
						class:          currentClass,
						method:         magic,
						originalMethod: pe.Method,
					}), nil
				}
			}
			if checkClass.GetExtend() == nil {
				break
			}
			vm := callVM(ctx)
			if vm == nil {
				break
			}
			parent, acl := vm.GetOrLoadClass(*checkClass.GetExtend())
			if acl != nil || parent == nil {
				break
			}
			checkClass = parent
		}
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("当前类 %s 没有静态方法 %s", currentClass.GetName(), pe.Method))
	}

	// 返回带类信息的静态方法包装器：确保调用时上下文携带 SelfClass/StaticClass，
	// 使方法体内的 self::/parent:: 能按词法（代码定义所在类）正确解析，而不是运行时调用类。
	return data.NewFuncValue(&staticMethodFunc{
		class:     currentClass,
		callClass: currentClass,
		method:    method,
	}), nil
}
