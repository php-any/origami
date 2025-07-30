package node

import (
	"errors"
	"fmt"
	"github.com/php-any/origami/data"
)

// CallStaticMethod 表示对象属性访问表达式
type CallStaticMethod struct {
	*Node  `pp:"-"`
	Class  string // 类名称 Class::fn() or Class::test::one
	Method string // 函数名
}

func NewCallStaticMethod(from *TokenFrom, path string, method string) *CallStaticMethod {
	return &CallStaticMethod{
		Node:   NewNode(from),
		Class:  path,
		Method: method,
	}
}

// GetValue 获取对象属性访问表达式的值
func (pe *CallStaticMethod) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	classStmt, ok := ctx.GetVM().GetClass(pe.Class)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("类(%s)不存在。", pe.Class)))
	}

	method, has := classStmt.GetMethod(pe.Method)
	if has {
		if method.GetIsStatic() {
			return data.NewFuncValue(method), nil
		}
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("类(%s)的函数(%s)不是静态的。", pe.Class, pe.Method)))
	}

	return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("类(%s)没有函数(%s)。", pe.Class, pe.Method)))
}
