package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewArrayAccessInterface 定义 PHP 内置接口 ArrayAccess。
// 该接口允许对象像数组一样被访问，包含以下方法：
// - offsetExists(mixed $offset): bool
// - offsetGet(mixed $offset): mixed
// - offsetSet(mixed $offset, mixed $value): void
// - offsetUnset(mixed $offset): void
func NewArrayAccessInterface() data.InterfaceStmt {
	methods := []data.Method{
		// offsetExists(mixed $offset): bool
		node.NewInterfaceMethod(
			nil,
			"offsetExists",
			"public",
			[]data.GetValue{
				node.NewParameter(nil, "offset", 0, nil, data.NewBaseType("mixed")),
			},
			data.NewBaseType("bool"),
		),
		// offsetGet(mixed $offset): mixed
		node.NewInterfaceMethod(
			nil,
			"offsetGet",
			"public",
			[]data.GetValue{
				node.NewParameter(nil, "offset", 0, nil, data.NewBaseType("mixed")),
			},
			data.NewBaseType("mixed"),
		),
		// offsetSet(mixed $offset, mixed $value): void
		node.NewInterfaceMethod(
			nil,
			"offsetSet",
			"public",
			[]data.GetValue{
				node.NewParameter(nil, "offset", 0, nil, data.NewBaseType("mixed")),
				node.NewParameter(nil, "value", 1, nil, data.NewBaseType("mixed")),
			},
			data.NewBaseType("void"),
		),
		// offsetUnset(mixed $offset): void
		node.NewInterfaceMethod(
			nil,
			"offsetUnset",
			"public",
			[]data.GetValue{
				node.NewParameter(nil, "offset", 0, nil, data.NewBaseType("mixed")),
			},
			data.NewBaseType("void"),
		),
	}
	return node.NewInterfaceStatement(nil, "ArrayAccess", nil, methods)
}
