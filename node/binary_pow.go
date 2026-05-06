package node

import (
	"math"

	"github.com/php-any/origami/data"
)

type BinaryPow struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryPow(from data.From, left, right data.GetValue) *BinaryPow {
	return &BinaryPow{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

func (b *BinaryPow) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}

	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	lf, err := lv.(data.AsFloat).AsFloat()
	if err != nil {
		return nil, data.NewErrorThrow(b.from, err)
	}
	rf, err := rv.(data.AsFloat).AsFloat()
	if err != nil {
		return nil, data.NewErrorThrow(b.from, err)
	}

	result := math.Pow(lf, rf)

	// 如果两个操作数都是整数且结果是精确整数，返回整数
	if _, lok := lv.(*data.IntValue); lok {
		if _, rok := rv.(*data.IntValue); rok {
			if result == math.Trunc(result) && result <= math.MaxInt && result >= math.MinInt {
				return data.NewIntValue(int(result)), nil
			}
		}
	}

	return data.NewFloatValue(result), nil
}
