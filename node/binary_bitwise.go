package node

import (
	"errors"

	"github.com/php-any/origami/data"
)

// Bitwise AND: $a & $b
type BinaryBitAnd struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryBitAnd(from data.From, left, right data.GetValue) *BinaryBitAnd {
	return &BinaryBitAnd{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
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

	li, ok := lv.(data.AsInt)
	if !ok {
		return nil, data.NewErrorThrow(b.from, errors.New("按位与左操作数必须是整数"))
	}
	ri, ok := rv.(data.AsInt)
	if !ok {
		return nil, data.NewErrorThrow(b.from, errors.New("按位与右操作数必须是整数"))
	}

	ln, err := li.AsInt()
	if err != nil {
		return nil, data.NewErrorThrow(b.from, err)
	}
	rn, err := ri.AsInt()
	if err != nil {
		return nil, data.NewErrorThrow(b.from, err)
	}

	return data.NewIntValue(int(ln & rn)), nil
}

// Bitwise XOR: $a ^ $b
type BinaryBitXor struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryBitXor(from data.From, left, right data.GetValue) *BinaryBitXor {
	return &BinaryBitXor{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
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

	li, ok := lv.(data.AsInt)
	if !ok {
		return nil, data.NewErrorThrow(b.from, errors.New("按位异或左操作数必须是整数"))
	}
	ri, ok := rv.(data.AsInt)
	if !ok {
		return nil, data.NewErrorThrow(b.from, errors.New("按位异或右操作数必须是整数"))
	}

	ln, err := li.AsInt()
	if err != nil {
		return nil, data.NewErrorThrow(b.from, err)
	}
	rn, err := ri.AsInt()
	if err != nil {
		return nil, data.NewErrorThrow(b.from, err)
	}

	return data.NewIntValue(int(ln ^ rn)), nil
}

// Bitwise OR: $a | $b
type BinaryBitOr struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryBitOr(from data.From, left, right data.GetValue) *BinaryBitOr {
	return &BinaryBitOr{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
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

	li, ok := lv.(data.AsInt)
	if !ok {
		return nil, data.NewErrorThrow(b.from, errors.New("按位或左操作数必须是整数"))
	}
	ri, ok := rv.(data.AsInt)
	if !ok {
		return nil, data.NewErrorThrow(b.from, errors.New("按位或右操作数必须是整数"))
	}

	ln, err := li.AsInt()
	if err != nil {
		return nil, data.NewErrorThrow(b.from, err)
	}
	rn, err := ri.AsInt()
	if err != nil {
		return nil, data.NewErrorThrow(b.from, err)
	}

	return data.NewIntValue(int(ln | rn)), nil
}
