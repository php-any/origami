package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewIteratorInterface 返回 PHP Iterator 接口定义
// Iterator extends Traversable，提供 current/key/next/rewind/valid 方法
func NewIteratorInterface() data.InterfaceStmt {
	methods := []data.Method{
		node.NewInterfaceMethod(nil, "current", "public", []data.GetValue{}, data.NewBaseType("mixed")),
		node.NewInterfaceMethod(nil, "key", "public", []data.GetValue{}, data.NewBaseType("mixed")),
		node.NewInterfaceMethod(nil, "next", "public", []data.GetValue{}, data.NewBaseType("void")),
		node.NewInterfaceMethod(nil, "rewind", "public", []data.GetValue{}, data.NewBaseType("void")),
		node.NewInterfaceMethod(nil, "valid", "public", []data.GetValue{}, data.NewBaseType("bool")),
	}
	return node.NewInterfaceStatement(nil, "Iterator", []string{"Traversable"}, methods)
}
