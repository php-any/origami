package directory

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewSeekableIteratorInterface 返回 PHP SeekableIterator 接口定义
// SeekableIterator extends Iterator，额外提供 seek(int $position) 方法
func NewSeekableIteratorInterface() data.InterfaceStmt {
	methods := []data.Method{
		node.NewInterfaceMethod(nil, "seek", "public", []data.GetValue{
			node.NewParameter(nil, "offset", 0, nil, data.NewBaseType("int")),
		}, data.NewBaseType("void")),
	}
	return node.NewInterfaceStatement(nil, "SeekableIterator", []string{"Iterator"}, methods)
}
