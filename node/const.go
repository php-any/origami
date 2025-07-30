package node

import "github.com/php-any/origami/data"

func (u *ConstStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// use语句本身不返回值
	return nil, nil
}

// ConstStatement 表示常量声明语句
type ConstStatement struct {
	*Node       `pp:"-"`
	Name        string
	Initializer data.GetValue
}

// NewConstStatement 创建一个新的常量声明语句
func NewConstStatement(token *TokenFrom, name string, initializer data.GetValue) *ConstStatement {
	return &ConstStatement{
		Node:        NewNode(token),
		Name:        name,
		Initializer: initializer,
	}
}
