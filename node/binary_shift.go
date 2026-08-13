package node

import (
	"github.com/php-any/origami/data"
)

// BinaryShl 表示左移运算符 <<
type BinaryShl struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryShl(from data.From, left, right data.GetValue) *BinaryShl {
	return &BinaryShl{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

func (b *BinaryShl) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}

	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	// 左移运算：将左操作数转换为整数，右操作数也是整数
	li, err := lv.(data.AsInt).AsInt()
	if err != nil {
		return nil, data.NewErrorThrow(b.from, err)
	}
	ri, err := rv.(data.AsInt).AsInt()
	if err != nil {
		return nil, data.NewErrorThrow(b.from, err)
	}

	return data.NewIntValue(li << ri), nil
}

// BinaryShr 表示右移运算符 >>
type BinaryShr struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryShr(from data.From, left, right data.GetValue) *BinaryShr {
	return &BinaryShr{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

func (b *BinaryShr) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}

	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	// 右移运算：将左操作数转换为整数，右操作数也是整数
	li, err := lv.(data.AsInt).AsInt()
	if err != nil {
		return nil, data.NewErrorThrow(b.from, err)
	}
	ri, err := rv.(data.AsInt).AsInt()
	if err != nil {
		return nil, data.NewErrorThrow(b.from, err)
	}

	return data.NewIntValue(li >> ri), nil
}
