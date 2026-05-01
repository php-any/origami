package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// CallObjectDynamicMethod 表示 $obj->{expr}(...) 或 $obj->$name(...) 这种动态方法名调用。
// 运行时先对 MethodExpr 求值得到方法名字符串，再按普通对象方法调用逻辑分派。
type CallObjectDynamicMethod struct {
	*Node      `pp:"-"`
	Object     data.GetValue
	MethodExpr data.GetValue
	Args       []data.GetValue
}

// NewCallObjectDynamicMethod 创建动态方法调用节点
func NewCallObjectDynamicMethod(from *TokenFrom, object data.GetValue, methodExpr data.GetValue, args []data.GetValue) *CallObjectDynamicMethod {
	return &CallObjectDynamicMethod{
		Node:       NewNode(from),
		Object:     object,
		MethodExpr: methodExpr,
		Args:       args,
	}
}

// GetValue 运行时先求值方法名，再委托给 CallObjectMethod 执行调用
func (pe *CallObjectDynamicMethod) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	o, ctl := pe.Object.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}

	raw, acl := pe.MethodExpr.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	if raw == nil {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("动态方法名表达式结果为 null"))
	}
	methodName := raw.(data.Value).AsString()
	if methodName == "" {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("动态方法名不能为空"))
	}

	// 借用 CallObjectMethod 的参数绑定和魔法方法逻辑
	proxy := &CallObjectMethod{
		Node:   pe.Node,
		Object: pe.Object,
		Method: methodName,
		Args:   pe.Args,
	}

	switch class := o.(type) {
	case *data.ThisValue:
		method, has := class.GetMethod(methodName)
		if has {
			fnCtx, acl := proxy.callMethodParams(class, ctx, method)
			if acl != nil {
				return nil, acl
			}
			return method.Call(fnCtx)
		}
		if magic, hasCall := class.GetMethod("__call"); hasCall {
			return proxy.invokeMagicCall(class, ctx, magic, methodName, pe.Args)
		}
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("this 对象不存在对应函数: %s", methodName))
	case *data.ClassValue:
		method, has := class.GetMethod(methodName)
		if has {
			if method.GetModifier() == data.ModifierPrivate {
				if !isCallerInClassHierarchy(ctx, class.Class) {
					return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("不能调用 private 方法: %s", methodName))
				}
			} else if method.GetModifier() == data.ModifierProtected {
				if !isCallerInClassHierarchy(ctx, class.Class) {
					return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("对象方法 %s 非公开", methodName))
				}
			}
			fnCtx, acl := proxy.callMethodParams(class, ctx, method)
			if acl != nil {
				return nil, acl
			}
			return method.Call(fnCtx)
		}
		if magic, hasCall := class.GetMethod("__call"); hasCall {
			return proxy.invokeMagicCall(class, ctx, magic, methodName, pe.Args)
		}
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("类(%s)不存在对应函数(%s)", class.Class.GetName(), methodName))
	default:
		if gm, ok := o.(data.GetMethod); ok {
			method, has := gm.GetMethod(methodName)
			if has {
				if method.GetModifier() != data.ModifierPublic {
					return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("对象方法 %s 非公开", methodName))
				}
				fnCtx, acl := proxy.callMethodParams(ctx, ctx, method)
				if acl != nil {
					return nil, acl
				}
				return method.Call(fnCtx)
			}
			if magic, hasCall := gm.GetMethod("__call"); hasCall {
				if objCtx, ok := o.(data.Context); ok {
					return proxy.invokeMagicCall(objCtx, ctx, magic, methodName, pe.Args)
				}
			}
		}
	}
	return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("当前值不支持调用函数 %s", methodName))
}
