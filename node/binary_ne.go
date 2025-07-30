package node

import (
	"github.com/php-any/origami/data"
)

type BinaryNe struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryNe(from data.From, left, right data.GetValue) *BinaryNe {
	return &BinaryNe{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

func (b *BinaryNe) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}

	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	// 简单的相等性比较，然后取反
	if lv == rv {
		return data.NewBoolValue(false), nil
	}

	// 类型转换比较
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
			return data.NewBoolValue(li != riVal), nil
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
			return data.NewBoolValue(lf != rfVal), nil
		}
	case *data.StringValue:
		if rs, ok := rv.(data.AsString); ok {
			ls := l.AsString()
			rsVal := rs.AsString()
			return data.NewBoolValue(ls != rsVal), nil
		}
	case *data.BoolValue:
		if rb, ok := rv.(data.AsBool); ok {
			lb, err := l.AsBool()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			rbVal, err := rb.AsBool()
			if err != nil {
				return nil, data.NewErrorThrow(b.from, err)
			}
			return data.NewBoolValue(lb != rbVal), nil
		}
	}

	return data.NewBoolValue(true), nil
}
