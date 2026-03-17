package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewRecursiveIteratorInterface 返回 PHP RecursiveIterator 接口定义
// RecursiveIterator extends Iterator，额外提供 hasChildren()/getChildren() 方法
func NewRecursiveIteratorInterface() data.InterfaceStmt {
	methods := []data.Method{
		node.NewInterfaceMethod(nil, "hasChildren", "public", []data.GetValue{}, data.NewBaseType("bool")),
		node.NewInterfaceMethod(nil, "getChildren", "public", []data.GetValue{}, data.NewBaseType("RecursiveIterator")),
	}
	return node.NewInterfaceStatement(nil, "RecursiveIterator", []string{"Iterator"}, methods)
}
