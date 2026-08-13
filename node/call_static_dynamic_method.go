package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// CallStaticDynamicMethod 表示 Class::$method() / $class::$method() 动态静态方法名调用。
// 运行时先对 MethodExpr 求值得到方法名，再按 CallStaticMethod 分派（含 __callStatic）。
type CallStaticDynamicMethod struct {
	*Node      `pp:"-"`
	stmt       data.GetValue // 类引用：ClassStmt / 类名字符串 / $class 变量等
	MethodExpr data.GetValue
}

func NewCallStaticDynamicMethod(from data.From, class data.GetValue, methodExpr data.GetValue) *CallStaticDynamicMethod {
	return &CallStaticDynamicMethod{
		Node:       NewNode(from),
		stmt:       class,
		MethodExpr: methodExpr,
	}
}

func (pe *CallStaticDynamicMethod) GetStmt() data.GetValue {
	return pe.stmt
}

func (pe *CallStaticDynamicMethod) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	raw, acl := pe.MethodExpr.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	if raw == nil {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("动态静态方法名表达式结果为 null"))
	}
	methodName := ""
	if v, ok := raw.(data.Value); ok {
		methodName = v.AsString()
	}
	if methodName == "" {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("动态静态方法名不能为空"))
	}

	var from *TokenFrom
	if tf, ok := pe.GetFrom().(*TokenFrom); ok {
		from = tf
	}
	proxy := NewCallStaticMethod(from, pe.stmt, methodName)
	return proxy.GetValue(ctx)
}
