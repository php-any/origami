package exception

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewThrowableInterface 定义 PHP 顶层接口 Throwable。
// 这里主要用于类型提示（catch(\Throwable) / type-hint）和 instanceof 判断，
// 方法签名先保持最小集合，与当前 Exception 实现对应。
func NewThrowableInterface() data.InterfaceStmt {
	methods := []data.Method{
		// getMessage(): string
		node.NewInterfaceMethod(nil, "getMessage", "public", nil, data.NewBaseType("string")),
		// getTraceAsString(): string
		node.NewInterfaceMethod(nil, "getTraceAsString", "public", nil, data.NewBaseType("string")),
	}
	return node.NewInterfaceStatement(nil, "Throwable", nil, methods)
}
