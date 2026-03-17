package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewOuterIteratorInterface 返回 PHP OuterIterator 接口定义
// OuterIterator extends Iterator，额外提供 getInnerIterator() 方法
func NewOuterIteratorInterface() data.InterfaceStmt {
	methods := []data.Method{
		node.NewInterfaceMethod(nil, "getInnerIterator", "public", []data.GetValue{}, data.NewBaseType("Iterator")),
	}
	return node.NewInterfaceStatement(nil, "OuterIterator", []string{"Iterator"}, methods)
}
