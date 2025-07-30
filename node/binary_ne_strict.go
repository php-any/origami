package node

import (
	"github.com/php-any/origami/data"
)

// BinaryNeStrict 表示严格不等表达式
type BinaryNeStrict struct {
	*Node
	left  data.GetValue
	right data.GetValue
}

// NewBinaryNeStrict 创建一个新的严格不等表达式
func NewBinaryNeStrict(from data.From, left data.GetValue, right data.GetValue) *BinaryNeStrict {
	return &BinaryNeStrict{
		Node:  NewNode(from),
		left:  left,
		right: right,
	}
}

// GetValue 获取严格不等表达式的值
func (b *BinaryNeStrict) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 计算左操作数
	leftValue, c := b.left.GetValue(ctx)
	if c != nil {
		return nil, c
	}

	// 计算右操作数
	rightValue, c := b.right.GetValue(ctx)
	if c != nil {
		return nil, c
	}

	// 严格不等比较：类型和值都必须不相等
	result := !isStrictEqual(leftValue, rightValue)
	return data.NewBoolValue(result), nil
}
