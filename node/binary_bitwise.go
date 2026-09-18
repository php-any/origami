package node

import (
	"github.com/php-any/origami/data"
)

// toIntOrZero 将值转为整数，null/非数字字符串返回 0。
// 两端都是字符串时不走这里：PHP 的 & | ^ 对字符串按字节运算。
func toIntOrZero(v data.GetValue) int {
	if v == nil {
		return 0
	}
	if _, isNull := v.(*data.NullValue); isNull {
		return 0
	}
	if _, ok := v.(*data.StringValue); ok {
		return 0
	}
	if iv, ok := v.(data.AsInt); ok {
		n, err := iv.AsInt()
		if err == nil {
			return n
		}
	}
	return 0
}

func phpBitwiseString(v data.GetValue) (string, bool) {
	s, ok := v.(*data.StringValue)
	if !ok || s == nil {
		return "", false
	}
	return s.Value, true
}

// phpBitwise 对齐 PHP：两端都是字符串则按字节 & | ^，结果仍是字符串；
// 否则转整数。mb_encode_numericentity 依赖 $s[$i] & "\xF0"。
func phpBitwise(lv, rv data.GetValue, op byte) data.Value {
	ls, lok := phpBitwiseString(lv)
	rs, rok := phpBitwiseString(rv)
	if lok && rok {
		n := len(ls)
		if len(rs) < n {
			n = len(rs)
		}
		out := make([]byte, n)
		switch op {
		case '&':
			for i := 0; i < n; i++ {
				out[i] = ls[i] & rs[i]
			}
		case '|':
			for i := 0; i < n; i++ {
				out[i] = ls[i] | rs[i]
			}
		case '^':
			for i := 0; i < n; i++ {
				out[i] = ls[i] ^ rs[i]
			}
		}
		return data.NewStringValue(string(out))
	}
	li, ri := toIntOrZero(lv), toIntOrZero(rv)
	switch op {
	case '|':
		return data.NewIntValue(li | ri)
	case '^':
		return data.NewIntValue(li ^ ri)
	default:
		return data.NewIntValue(li & ri)
	}
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
	return phpBitwise(lv, rv, '&'), nil
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
	return phpBitwise(lv, rv, '^'), nil
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
	return phpBitwise(lv, rv, '|'), nil
}
