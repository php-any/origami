package node

import (
	"github.com/php-any/origami/data"
)

type BinaryGt struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryGt(from data.From, left, right data.GetValue) *BinaryGt {
	return &BinaryGt{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

func (b *BinaryGt) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}

	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	if cmp, ok := phpCompareValues(ctx, lv, rv); ok {
		return data.NewBoolValue(cmp > 0), nil
	}
	return data.NewBoolValue(false), nil
}
