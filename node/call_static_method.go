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
