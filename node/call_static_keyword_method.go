package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

// staticMethodFuncWithLateBinding 支持后期静态绑定的静态方法包装器
type staticMethodFuncWithLateBinding struct {
	callClass data.ClassStmt // 调用时的类（用于后期静态绑定）
	method    data.Method
}

func newStaticMethodFuncWithLateBinding(callClass data.ClassStmt, method data.Method) *staticMethodFuncWithLateBinding {
	return &staticMethodFuncWithLateBinding{
		callClass: callClass,
		method:    method,
	}
}

func (s *staticMethodFuncWithLateBinding) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return s, nil
}

func (s *staticMethodFuncWithLateBinding) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 创建类方法上下文，绑定调用时的类（用于 static:: 后期静态绑定）
	classValue := data.NewClassValue(s.callClass, ctx)
	fnCtx := classValue.CreateContext(s.method.GetVariables())
	// 设置后期静态绑定类
	if cmc, ok := fnCtx.(*data.ClassMethodContext); ok {
		cmc.StaticClass = s.callClass
	}
	// 复制参数
	for i := 0; i < len(s.method.GetVariables()); i++ {
		fnCtx.SetIndexZVal(i, ctx.GetIndexZVal(i))
	}
	return s.method.Call(fnCtx)
}

func (s *staticMethodFuncWithLateBinding) AsString() string {
	return fmt.Sprintf("static::%s", s.method.GetName())
}

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
	// 与 self:: 一样，必须在类上下文中使用
	var currentClass data.ClassStmt
	if classCtx, ok := ctx.(*data.ClassMethodContext); ok {
		currentClass = classCtx.Class
		if classCtx.StaticClass != nil {
			currentClass = classCtx.StaticClass
		}
	} else if classVal, ok := ctx.(*data.ClassValue); ok {
		currentClass = classVal.Class
	} else {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("static:: 只能在类方法中使用"))
	}

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
					// 返回包装器，携带调用类信息用于后期静态绑定
					return newStaticMethodFuncWithLateBinding(currentClass, method), nil
				}
				extend = ext.GetExtend()
			}
		}

		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("当前类 %s 没有静态方法 %s", currentClass.GetName(), pe.Method))
	}

	// 返回包装器，携带调用类信息用于后期静态绑定
	return newStaticMethodFuncWithLateBinding(currentClass, method), nil
}
