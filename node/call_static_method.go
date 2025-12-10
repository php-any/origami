package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

// CallStaticMethod 表示对象属性访问表达式
type CallStaticMethod struct {
	*Node  `pp:"-"`
	stmt   data.GetValue // 类名称 Class::fn() or Class::test::one
	Method string        // 函数名
}

func NewCallStaticMethod(from *TokenFrom, path data.GetValue, method string) *CallStaticMethod {
	return &CallStaticMethod{
		Node:   NewNode(from),
		stmt:   path,
		Method: method,
	}
}

// GetValue 获取对象属性访问表达式的值
func (pe *CallStaticMethod) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	switch expr := pe.stmt.(type) {
	case data.GetStaticMethod:
		method, has := expr.GetStaticMethod(pe.Method)
		if has {
			return data.NewFuncValue(method), nil
		}

		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("无法调用函数(%s)。", pe.Method)))
	default:
		c, acl := pe.stmt.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		switch expr := c.(type) {
		case data.GetStaticMethod:
			method, has := expr.GetStaticMethod(pe.Method)
			if has {
				return data.NewFuncValue(method), nil
			}
		case data.GetMethod:
			method, has := expr.GetMethod(pe.Method)
			if has {
				return data.NewFuncValue(method), nil
			}
		}
	}

	name := ""
	if getName, ok := pe.stmt.(data.ClassStmt); ok {
		name = getName.GetName()
	}

	return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("(%v)无法调用函数(%s)。", name, pe.Method)))
}

// CallStaticMethodLater 延迟的静态方法调用（类未加载时）
type CallStaticMethodLater struct {
	*Node
	className string // 类名（字符串形式）
	method    string // 方法名
	namespace string // 命名空间
}

// NewCallStaticMethodLater 创建延迟的静态方法调用
func NewCallStaticMethodLater(from *TokenFrom, className, method, namespace string) *CallStaticMethodLater {
	return &CallStaticMethodLater{
		Node:      NewNode(from),
		className: className,
		method:    method,
		namespace: namespace,
	}
}

// GetValue 获取延迟静态方法调用的值
func (pe *CallStaticMethodLater) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 尝试加载类
	stmt, acl := ctx.GetVM().GetOrLoadClass(pe.className)
	if acl != nil {
		return nil, acl
	}
	if stmt == nil {
		// 如果还是找不到，尝试使用命名空间
		fullClassName := pe.className
		if pe.namespace != "" {
			fullClassName = pe.namespace + "\\" + pe.className
		}
		stmt, acl = ctx.GetVM().GetOrLoadClass(fullClassName)
		if acl != nil {
			return nil, acl
		}
		if stmt == nil {
			return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("无法调用静态方法(%s::%s), 未找到类", pe.className, pe.method))
		}
	}

	// 创建实际的静态方法调用
	tokenFrom, ok := pe.GetFrom().(*TokenFrom)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("无法获取TokenFrom信息"))
	}
	callStaticMethod := NewCallStaticMethod(tokenFrom, stmt, pe.method)
	return callStaticMethod.GetValue(ctx)
}
