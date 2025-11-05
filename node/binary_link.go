package node

import (
	"github.com/php-any/origami/data"
)

// BinaryLink 专用于字符串插值阶段的连接运算
// 行为：将左右两侧转为字符串后拼接
type BinaryLink struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryLink(from data.From, left, right data.GetValue) *BinaryLink {
	return &BinaryLink{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

func (b *BinaryLink) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}

	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	// 统一转字符串拼接，尽量与 PHP 的字符串连接行为一致
	leftStr := lv.(data.Value).AsString()
	rightStr := rv.(data.Value).AsString()

	return data.NewStringValue(leftStr + rightStr), nil
}
