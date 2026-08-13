package node

import "github.com/php-any/origami/data"

// GlobalStatement 表示 global 变量声明语句
// PHP 中 global $var 将函数局部变量绑定到同名全局变量（共享 ZVal）
type GlobalStatement struct {
	*Node `pp:"-"`
	// Names 是声明为全局变量的变量名列表（不含 $ 前缀）
	Names []string
	// Indexes 是各变量在当前函数作用域中的索引
	Indexes []int
}

// NewGlobalStatement 创建一个新的 global 声明语句
func NewGlobalStatement(from data.From, names []string, indexes []int) *GlobalStatement {
	return &GlobalStatement{
		Node:    NewNode(from),
		Names:   names,
		Indexes: indexes,
	}
}

// GlobalVarProvider 是 VM 需要实现的全局变量访问接口
type GlobalVarProvider interface {
	EnsureGlobalZVal(name string) *data.ZVal
}

// GetValue 执行 global 语句：将本地 slot 替换为全局 ZVal，实现共享
func (g *GlobalStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	vm := ctx.GetVM()
	provider, ok := vm.(GlobalVarProvider)
	if !ok {
		// VM 不支持全局变量管理，忽略
		return nil, nil
	}

	for i, name := range g.Names {
		if i >= len(g.Indexes) {
			break
		}
		localIndex := g.Indexes[i]
		// 获取（或创建）全局变量的 ZVal
		globalZVal := provider.EnsureGlobalZVal(name)
		// 将本地 ctx 的该 slot 替换为全局 ZVal，实现共享
		ctx.SetIndexZVal(localIndex, globalZVal)
	}

	return nil, nil
}
