package node

import (
	"strconv"

	"github.com/php-any/origami/data"
)

type BinaryDot struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryDot(from data.From, left, right data.GetValue) *BinaryDot {
	return &BinaryDot{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

func unwrapValue(v data.GetValue) data.GetValue {
	if v == nil {
		return v
	}
	switch t := v.(type) {
	case *data.ZValValue:
		if t != nil && t.ZVal != nil && t.ZVal.Value != nil {
			if val, ok := t.ZVal.Value.(data.GetValue); ok {
				return unwrapValue(val)
			}
		}
		return v
	case *data.ReferenceValue:
		actualVal, _ := t.Val.GetValue(t.Ctx)
		return unwrapValue(actualVal)
	case *data.IndexReferenceValue:
		actualVal, _ := t.Expr.GetValue(t.Ctx)
		return unwrapValue(actualVal)
	default:
		return v
	}
}

func concatOperandString(ctx data.Context, v data.GetValue) (string, data.Control) {
	switch t := v.(type) {
	case *data.StringValue:
		return t.AsString(), nil
	case *data.IntValue:
		return strconv.Itoa(t.Value), nil
	case *data.FloatValue:
		return t.AsString(), nil
	case *data.BoolValue:
		if t.Value {
			return "1", nil
		}
		return "", nil
	case *data.NullValue:
		return "", nil
	case *data.ArrayValue:
		return t.AsString(), nil
	default:
		return ValueToDisplayString(ctx, v)
	}
}

func (b *BinaryDot) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if b == nil || b.Left == nil {
		return data.NewNullValue(), nil
	}
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}
	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	leftStr, acl := concatOperandString(ctx, unwrapValue(lv))
	if acl != nil {
		return nil, acl
	}
	rightStr, acl := concatOperandString(ctx, unwrapValue(rv))
	if acl != nil {
		return nil, acl
	}
	return data.NewStringValue(leftStr + rightStr), nil
}
