package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewTraversableInterface 返回 PHP Traversable 接口定义
// Traversable 是空接口，所有可遍历接口的基接口
func NewTraversableInterface() data.InterfaceStmt {
	return node.NewInterfaceStatement(nil, "Traversable", nil, nil)
}
