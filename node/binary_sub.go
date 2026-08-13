package node

import (
	"errors"
	"github.com/php-any/origami/data"
)

type BinarySub struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinarySub(from data.From, left, right data.GetValue) *BinarySub {
	return &BinarySub{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

func (b *BinarySub) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}

	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	switch lv.(type) {
	case *data.StringValue:
		li, err := lv.(data.AsFloat).AsFloat()
		if err != nil {
			return nil, data.NewErrorThrow(b.from, err)
		}
		ri, err := rv.(data.AsFloat).AsFloat()
		if err != nil {
			return nil, data.NewErrorThrow(b.from, err)
		}

		return data.NewFloatValue(li - ri), nil
	case *data.IntValue:
		li, err := lv.(data.AsInt).AsInt()
		if err != nil {
			return nil, data.NewErrorThrow(b.from, err)
		}
		// 右操作数为数字字符串时走 int；含小数点等非纯 int 字符串走 float（对齐 PHP）
		if rs, ok := rv.(*data.StringValue); ok {
			if ri, err := rs.AsInt(); err == nil {
				return data.NewIntValue(li - ri), nil
			}
			rf, err := rs.AsFloat()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			return data.NewFloatValue(float64(li) - rf), nil
		}
		if riv, ok := rv.(data.AsInt); ok {
			ri, err := riv.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			return data.NewIntValue(li - ri), nil
		}
		if rfv, ok := rv.(data.AsFloat); ok {
			rf, err := rfv.AsFloat()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			return data.NewFloatValue(float64(li) - rf), nil
		}
		return nil, data.NewErrorThrow(b.from, errors.New("TODO 有未支持的类型减法"))
	case *data.FloatValue:
		li, err := lv.(data.AsFloat).AsFloat()
		if err != nil {
			return nil, data.NewErrorThrow(b.from, err)
		}
		ri, err := rv.(data.AsFloat).AsFloat()
		if err != nil {
			return nil, data.NewErrorThrow(b.from, err)
		}

		return data.NewFloatValue(li - ri), nil
	case *data.NullValue:
		li, err := lv.(data.AsFloat).AsFloat()
		if err != nil {
			return nil, data.NewErrorThrow(b.from, err)
		}
		ri, err := rv.(data.AsFloat).AsFloat()
		if err != nil {
			return nil, data.NewErrorThrow(b.from, err)
		}

		return data.NewFloatValue(li - ri), nil
	}

	return nil, data.NewErrorThrow(b.from, errors.New("TODO 有未支持的类型减法"))
}
