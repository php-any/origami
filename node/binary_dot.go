package node

import (
	"fmt"

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

func (b *BinaryDot) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 获取左操作数的值
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}

	// 获取右操作数的值
	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	// 将两个操作数都转换为字符串
	leftStr := ""
	rightStr := ""

	// 处理左操作数
	switch l := lv.(type) {
	case *data.StringValue:
		leftStr = l.AsString()
	case *data.IntValue:
		leftStr = fmt.Sprintf("%d", l.Value)
	case *data.FloatValue:
		leftStr = fmt.Sprintf("%f", l.Value)
	case *data.BoolValue:
		if l.Value {
			leftStr = "1"
		} else {
			leftStr = ""
		}
	case *data.NullValue:
		leftStr = ""
	case *data.ArrayValue:
		leftStr = l.AsString()
	default:
		// 对于其他类型，尝试使用 AsString 方法
		if strValue, ok := l.(data.AsString); ok {
			leftStr = strValue.AsString()
		} else {
			leftStr = fmt.Sprintf("%v", l)
		}
	}

	// 处理右操作数
	switch r := rv.(type) {
	case *data.StringValue:
		rightStr = r.AsString()
	case *data.IntValue:
		rightStr = fmt.Sprintf("%d", r.Value)
	case *data.FloatValue:
		rightStr = fmt.Sprintf("%f", r.Value)
	case *data.BoolValue:
		if r.Value {
			rightStr = "1"
		} else {
			rightStr = ""
		}
	case *data.NullValue:
		rightStr = ""
	case *data.ArrayValue:
		rightStr = r.AsString()
	default:
		// 对于其他类型，尝试使用 AsString 方法
		if strValue, ok := r.(data.AsString); ok {
			rightStr = strValue.AsString()
		} else {
			rightStr = fmt.Sprintf("%v", r)
		}
	}

	// 连接字符串并返回
	return data.NewStringValue(leftStr + rightStr), nil
}
