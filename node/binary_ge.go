package node

import (
	"github.com/php-any/origami/data"
)

type BinaryGe struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryGe(from data.From, left, right data.GetValue) *BinaryGe {
	return &BinaryGe{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

func (b *BinaryGe) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}

	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	if cmp, ok := phpCompareValues(ctx, lv, rv); ok {
		return data.NewBoolValue(cmp >= 0), nil
	}
	return data.NewBoolValue(false), nil
}
