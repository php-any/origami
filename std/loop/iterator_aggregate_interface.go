package loop

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IteratorAggregate 语言级接口：与 SPL 中 IteratorAggregate 对应
// 在 PHP 中它继承 Traversable，并要求实现 getIterator(): Traversable。
// 这里至少建出 getIterator() 方法的接口 AST，避免注册完全空的接口。
func newIteratorAggregateInterface() data.InterfaceStmt {
	methods := []data.Method{
		// 先使用“未指定”返回类型（等价于 mixed），后续如果引入 Traversable 接口再精细化类型。
		node.NewInterfaceMethod(nil, "getIterator", "public", nil, data.NewBaseType("")),
	}
	return node.NewInterfaceStatement(nil, "IteratorAggregate", nil, methods)
}
