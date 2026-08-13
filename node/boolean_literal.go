package node

import "github.com/php-any/origami/data"

// BooleanLiteral 表示布尔字面量
type BooleanLiteral struct {
	*Node `pp:"-"`
	Value bool
}

// NewBooleanLiteral 创建一个新的布尔字面量
func NewBooleanLiteral(token *TokenFrom, value bool) data.GetValue {
	return &BooleanLiteral{
		Node:  NewNode(token),
		Value: value,
	}
}

// GetValue 获取布尔字面量的值
func (b *BooleanLiteral) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(b.Value), nil
}
