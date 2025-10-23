package loop

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// 定义 Iterator 接口的 AST（语言级接口）
type iteratorInterface struct {
	*node.InterfaceStatement
}

func newIteratorInterface() data.InterfaceStmt {
	// 方法签名：current(): mixed, key(): mixed, next(): void, rewind(): void, valid(): bool
	methods := []data.Method{
		node.NewInterfaceMethod(nil, "current", "public", nil, data.NewBaseType("")),
		node.NewInterfaceMethod(nil, "key", "public", nil, data.NewBaseType("")),
		node.NewInterfaceMethod(nil, "next", "public", nil, data.NewBaseType("void")),
		node.NewInterfaceMethod(nil, "rewind", "public", nil, data.NewBaseType("void")),
		node.NewInterfaceMethod(nil, "valid", "public", nil, data.NewBaseType("bool")),
	}
	return node.NewInterfaceStatement(nil, "Iterator", nil, methods)
}

func Load(vm data.VM) {
	vm.AddInterface(newIteratorInterface())
	vm.AddClass(NewListClass())
	vm.AddClass(NewHashMapClass())
}
