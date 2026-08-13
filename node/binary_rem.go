package node

import (
	"errors"
	"github.com/php-any/origami/data"
)

type BinaryRem struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryRem(from data.From, left, right data.GetValue) *BinaryRem {
	return &BinaryRem{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

func (b *BinaryRem) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}

	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	switch lv.(type) {
	case *data.IntValue:
		li, err := lv.(data.AsInt).AsInt()
		if err != nil {
			return nil, data.NewErrorThrow(b.from, err)
		}
		ri, err := rv.(data.AsInt).AsInt()
		if err != nil {
			return nil, data.NewErrorThrow(b.from, err)
		}

		if ri == 0 {
			return nil, data.NewErrorThrow(b.from, errors.New("除零错误"))
		}

		return data.NewIntValue(li % ri), nil
	case *data.FloatValue:
		lf, err := lv.(data.AsFloat).AsFloat()
		if err != nil {
			return nil, data.NewErrorThrow(b.from, err)
		}
		rf, err := rv.(data.AsFloat).AsFloat()
		if err != nil {
			return nil, data.NewErrorThrow(b.from, err)
		}

		if rf == 0 {
			return nil, data.NewErrorThrow(b.from, errors.New("除零错误"))
		}

		return data.NewFloatValue(float64(int64(lf) % int64(rf))), nil
	}

	return nil, data.NewErrorThrow(b.from, errors.New("TODO 有未支持的类型取模"))
}
