package node

import (
	"errors"

	"github.com/php-any/origami/data"
)

type BinaryQuo struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryQuo(from data.From, left, right data.GetValue) *BinaryQuo {
	return &BinaryQuo{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

func (b *BinaryQuo) GetValue(ctx data.Context) (data.GetValue, data.Control) {
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

		if ai, ok := rv.(data.AsFloat); ok {
			ri, err := ai.AsFloat()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			if ri == 0 {
				return nil, data.NewErrorThrow(b.from, errors.New("除零错误"))
			}
			return data.NewFloatValue(float64(li) / ri), nil
		}

		if ai, ok := rv.(data.AsInt); ok {
			ri, err := ai.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			if ri == 0 {
				return nil, data.NewErrorThrow(b.from, errors.New("除零错误"))
			}
			return data.NewFloatValue(float64(li) / float64(ri)), nil
		}
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

		return data.NewFloatValue(lf / rf), nil
	}

	return nil, data.NewErrorThrow(b.from, errors.New("TODO 有未支持的类型除法"))
}
