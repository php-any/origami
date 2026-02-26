package exception

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewJsonSerializableInterface 定义 PHP 顶层接口 JsonSerializable。
// 该接口只有一个方法：jsonSerialize(): mixed
// json_encode 在遇到实现了 JsonSerializable 的对象时，应优先调用 jsonSerialize() 的返回值进行编码。
func NewJsonSerializableInterface() data.InterfaceStmt {
	methods := []data.Method{
		node.NewInterfaceMethod(nil, "jsonSerialize", "public", nil, data.NewBaseType("mixed")),
	}
	return node.NewInterfaceStatement(nil, "JsonSerializable", nil, methods)
}
