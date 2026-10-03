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
	return callValue(pe.GetReferenceValue(ctx))
}

func (pe *CallObjectDynamicMethod) GetReferenceValue(ctx data.Context) (data.GetValue, data.Control) {
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

	proxy := &CallObjectMethod{Node: pe.Node, Object: &literalGetValue{v: o}, Method: methodName, Args: pe.Args}
	return proxy.GetReferenceValue(ctx)
}
