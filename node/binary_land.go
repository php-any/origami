package node

import (
	"github.com/php-any/origami/data"
)

type BinaryLand struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryLand(from data.From, left, right data.GetValue) *BinaryLand {
	return &BinaryLand{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

func (b *BinaryLand) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}

	// 短路求值：如果左操作数为假，直接返回假
	lb, err := lv.(data.AsBool).AsBool()
	if err != nil {
		return nil, data.NewErrorThrow(b.from, err)
	}
	if !lb {
		return data.NewBoolValue(false), nil
	}

	// 左操作数为真，继续求值右操作数
	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	rb, err := rv.(data.AsBool).AsBool()
	if err != nil {
		return nil, data.NewErrorThrow(b.from, err)
	}

	return data.NewBoolValue(rb), nil
}
