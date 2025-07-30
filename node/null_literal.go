package node

import "github.com/php-any/origami/data"

// NullLiteral 表示null字面量
type NullLiteral struct {
	*Node `pp:"-"`
}

// NewNullLiteral 创建一个新的null字面量
func NewNullLiteral(token *TokenFrom) data.GetValue {
	return &NullLiteral{
		Node: NewNode(token),
	}
}

// GetValue 获取null字面量的值
func (n *NullLiteral) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}
