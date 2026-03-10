package system

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// newDateTimeInterface 定义 PHP 顶层接口 DateTimeInterface。
// 目前仅包含与 DateTime / DateTimeImmutable 最小公共子集一致的方法集合，用于类型提示和 instanceof 判断。
func newDateTimeInterface() data.InterfaceStmt {
	methods := []data.Method{
		// format(string $format): string
		node.NewInterfaceMethod(
			nil,
			"format",
			"public",
			[]data.GetValue{
				node.NewVariable(nil, "format", 0, data.NewBaseType("string")),
			},
			data.NewBaseType("string"),
		),
		// getTimestamp(): int
		node.NewInterfaceMethod(
			nil,
			"getTimestamp",
			"public",
			nil,
			data.NewBaseType("int"),
		),
	}

	return node.NewInterfaceStatement(nil, "DateTimeInterface", nil, methods)
}
