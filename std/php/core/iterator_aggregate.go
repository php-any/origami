package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewIteratorAggregateInterface 返回 PHP IteratorAggregate 接口定义
// IteratorAggregate extends Traversable，提供 getIterator() 方法
func NewIteratorAggregateInterface() data.InterfaceStmt {
	methods := []data.Method{
		node.NewInterfaceMethod(nil, "getIterator", "public", []data.GetValue{}, data.NewBaseType("Traversable")),
	}
	return node.NewInterfaceStatement(nil, "IteratorAggregate", []string{"Traversable"}, methods)
}
