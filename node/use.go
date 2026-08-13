package node

import "github.com/php-any/origami/data"

// UseStatement 表示use语句节点
type UseStatement struct {
	*Node     `pp:"-"`
	Namespace string // 命名空间名称
	Alias     string // 别名（可选）
}

// NewUseStatement 创建一个新的use语句节点
func NewUseStatement(from data.From, namespace string, alias string) *UseStatement {
	return &UseStatement{
		Node:      NewNode(from),
		Namespace: namespace,
		Alias:     alias,
	}
}

// GetNamespace 返回命名空间名称
func (u *UseStatement) GetNamespace() string {
	return u.Namespace
}

// GetAlias 返回别名
func (u *UseStatement) GetAlias() string {
	return u.Alias
}

// GetValue 获取use语句节点的值
func (u *UseStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// use语句本身不会引入文件, 只是声明一个别名
	//  如果调用本包类会导致引入，因为需要尝试拼接完整地址
	//  运行时, GetClass 才需要自动引入	·
	return nil, nil
}
