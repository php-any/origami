package node

import (
	"github.com/php-any/origami/data"
)

type BinaryLor struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryLor(from data.From, left, right data.GetValue) *BinaryLor {
	return &BinaryLor{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

func (b *BinaryLor) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}

	// 短路求值：如果左操作数为真，直接返回真
	lb, err := lv.(data.AsBool).AsBool()
	if err != nil {
		return nil, data.NewErrorThrow(b.from, err)
	}
	if lb {
		return data.NewBoolValue(true), nil
	}

	// 左操作数为假，继续求值右操作数
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
