package node

import (
	"errors"
	"github.com/php-any/origami/data"
)

type BinaryMul struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryMul(from data.From, left, right data.GetValue) *BinaryMul {
	return &BinaryMul{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

func (b *BinaryMul) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}

	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	switch lvv := lv.(type) {
	case *data.IntValue:
		li, err := lvv.AsInt()
		if err != nil {
			return nil, data.NewErrorThrow(b.from, err)
		}
		ri, err := rv.(data.AsInt).AsInt()
		if err != nil {
			return nil, data.NewErrorThrow(b.from, err)
		}

		return data.NewIntValue(li * ri), nil
	case *data.FloatValue:
		lf, err := lvv.AsFloat()
		if err != nil {
			return nil, data.NewErrorThrow(b.from, err)
		}
		rf, err := rv.(data.AsFloat).AsFloat()
		if err != nil {
			return nil, data.NewErrorThrow(b.from, err)
		}

		return data.NewFloatValue(lf * rf), nil
	}

	return nil, data.NewErrorThrow(b.from, errors.New("TODO 有未支持的类型乘法"))
}
