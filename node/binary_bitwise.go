package node

import (
	"github.com/php-any/origami/data"
)

// toIntOrZero 将值转为整数，null/非数字字符串返回 0
func toIntOrZero(v data.GetValue) int {
	if v == nil {
		return 0
	}
	if _, isNull := v.(*data.NullValue); isNull {
		return 0
	}
	if _, ok := v.(*data.StringValue); ok {
		return 0 // 字符串按位操作返回 0（PHP 兼容）
	}
	if iv, ok := v.(data.AsInt); ok {
		n, err := iv.AsInt()
		if err == nil {
			return n
		}
	}
	return 0
}

type BinaryBitAnd struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryBitAnd(from data.From, left, right data.GetValue) *BinaryBitAnd {
	return &BinaryBitAnd{Node: NewNode(from), Left: left, Right: right}
}

func (b *BinaryBitAnd) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}
	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}
	return data.NewIntValue(toIntOrZero(lv) & toIntOrZero(rv)), nil
}

type BinaryBitXor struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryBitXor(from data.From, left, right data.GetValue) *BinaryBitXor {
	return &BinaryBitXor{Node: NewNode(from), Left: left, Right: right}
}

func (b *BinaryBitXor) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}
	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}
	return data.NewIntValue(toIntOrZero(lv) ^ toIntOrZero(rv)), nil
}

type BinaryBitOr struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryBitOr(from data.From, left, right data.GetValue) *BinaryBitOr {
	return &BinaryBitOr{Node: NewNode(from), Left: left, Right: right}
}

func (b *BinaryBitOr) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}
	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}
	return data.NewIntValue(toIntOrZero(lv) | toIntOrZero(rv)), nil
}
