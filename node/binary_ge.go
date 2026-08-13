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

	switch l := lv.(type) {
	case *data.IntValue:
		if ri, ok := rv.(data.AsInt); ok {
			li, err := l.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			riVal, err := ri.AsInt()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			return data.NewBoolValue(li >= riVal), nil
		}
	case *data.FloatValue:
		if rf, ok := rv.(data.AsFloat); ok {
			lf, err := l.AsFloat()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			rfVal, err := rf.AsFloat()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			return data.NewBoolValue(lf >= rfVal), nil
		}
	case *data.StringValue:
		if rs, ok := rv.(data.AsString); ok {
			ls := l.AsString()
			rsVal := rs.AsString()
			return data.NewBoolValue(ls >= rsVal), nil
		}
	}

	return data.NewBoolValue(false), nil
}
