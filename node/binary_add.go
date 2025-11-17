package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

type BinaryAdd struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryAdd(from data.From, left, right data.GetValue) *BinaryAdd {
	return &BinaryAdd{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

func (b *BinaryAdd) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}

	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	switch l := lv.(type) {
	case *data.StringValue:
		lStr := l.AsString()
		rStr := rv.(data.Value).AsString()

		return data.NewStringValue(lStr + rStr), nil
	case *data.IntValue:
		switch r := rv.(type) {
		case data.AsInt:
			li, err := l.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			ri, err := r.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}

			return data.NewIntValue(li + ri), nil
		case data.AsString:
			return data.NewStringValue(l.AsString() + r.AsString()), nil
		}
	case *data.FloatValue:
		switch r := rv.(type) {
		case data.AsInt:
			lf, err := l.AsFloat()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			ri, err := r.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}

			return data.NewFloatValue(lf + float64(ri)), nil
		case *data.StringValue:
			lf, err := l.AsFloat()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			rf := r.AsString()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			return data.NewStringValue(fmt.Sprintf("%f", lf) + rf), nil
		case data.AsFloat:
			lf, err := l.AsFloat()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			rf, err := r.AsFloat()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			return data.NewFloatValue(lf + rf), nil
		}

	case *data.NullValue:
		if riv, ok := rv.(data.AsInt); ok {
			ri, err := riv.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			return data.NewIntValue(0 + ri), nil
		}

		return nil, data.NewErrorThrow(b.from, errors.New("左边 null 值无法运行+符号"))

	case *data.ArrayValue:
		lStr := l.AsString()
		rStr := rv.(data.Value).AsString()

		return data.NewStringValue(lStr + rStr), nil
	case *data.ObjectValue:
		lStr := l.AsString()
		rStr := rv.(data.Value).AsString()

		return data.NewStringValue(lStr + rStr), nil
	case *data.AnyValue:
		lStr := l.AsString()
		rStr := rv.(data.Value).AsString()

		return data.NewStringValue(lStr + rStr), nil
	}

	return nil, data.NewErrorThrow(b.from, fmt.Errorf("TODO 有未支持的类型加法 %v", lv))
}
