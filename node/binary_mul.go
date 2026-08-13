package node

import (
	"errors"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
)

type BinaryMul struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryMul(from data.From, left, right data.GetValue) *BinaryMul {
	return &BinaryMul{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

func (b *BinaryMul) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}

	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	leftNumber, leftIsInt, leftOK := numericForMul(lv)
	rightNumber, rightIsInt, rightOK := numericForMul(rv)
	if leftOK && rightOK {
		if leftIsInt && rightIsInt {
			return data.NewIntValue(int(leftNumber) * int(rightNumber)), nil
		}
		return data.NewFloatValue(leftNumber * rightNumber), nil
	}

	return nil, data.NewErrorThrow(b.from, errors.New("TODO 有未支持的类型乘法"))
}

// numericForMul 实现 PHP 算术运算中的常用数字转换，包括 numeric-string。
func numericForMul(value data.GetValue) (number float64, isInt bool, ok bool) {
	switch v := value.(type) {
	case *data.IntValue:
		n, err := v.AsInt()
		return float64(n), true, err == nil
	case *data.FloatValue:
		n, err := v.AsFloat()
		return n, false, err == nil
	case *data.BoolValue:
		if v.Value {
			return 1, true, true
		}
		return 0, true, true
	case *data.NullValue:
		return 0, true, true
	case *data.StringValue:
		s := strings.TrimSpace(v.Value)
		if s == "" {
			return 0, false, false
		}
		if !strings.ContainsAny(s, ".eE") {
			if n, err := strconv.ParseInt(s, 10, 64); err == nil {
				return float64(n), true, true
			}
		}
		n, err := strconv.ParseFloat(s, 64)
		return n, false, err == nil
	default:
		return 0, false, false
	}
}
