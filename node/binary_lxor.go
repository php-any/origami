package node

import (
	"github.com/php-any/origami/data"
)

type BinaryLxor struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryLxor(from data.From, left, right data.GetValue) *BinaryLxor {
	return &BinaryLxor{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

func (b *BinaryLxor) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}
	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	lb, err := lv.(data.AsBool).AsBool()
	if err != nil {
		return nil, data.NewErrorThrow(b.from, err)
	}
	rb, err := rv.(data.AsBool).AsBool()
	if err != nil {
		return nil, data.NewErrorThrow(b.from, err)
	}

	return data.NewBoolValue(lb != rb), nil
}
