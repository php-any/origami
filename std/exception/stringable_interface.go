package exception

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewStringableInterface 定义 PHP 顶层接口 Stringable。
// 主要用于类型提示（Stringable $x）以及与 __toString 相关的约束。
func NewStringableInterface() data.InterfaceStmt {
	methods := []data.Method{
		// __toString(): string
		node.NewInterfaceMethod(nil, "__toString", "public", nil, data.NewBaseType("string")),
	}
	return node.NewInterfaceStatement(nil, "Stringable", nil, methods)
}
