package node

import (
	"fmt"
	"reflect"

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

// unwrapValue 解包可能嵌套在 ZVal 或 ReferenceValue 中的值
func unwrapValue(v data.GetValue) data.GetValue {
	if v == nil {
		return v
	}

	// 使用反射检查 ZVal 类型（因为 ZVal 没有实现 GetValue 接口）
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr && !rv.IsNil() {
		// 检查是否是 ZVal 类型（通过类型名称）
		if rv.Type().String() == "*data.ZVal" {
			// 获取 Value 字段
			valueField := rv.Elem().FieldByName("Value")
			if valueField.IsValid() && !valueField.IsNil() {
				if val, ok := valueField.Interface().(data.GetValue); ok {
					return unwrapValue(val)
				}
			}
		}
	}

	// 处理 ReferenceValue - 获取引用的实际值
	if ref, ok := v.(*data.ReferenceValue); ok {
		actualVal, _ := ref.Val.GetValue(ref.Ctx)
		return unwrapValue(actualVal)
	}

	// 处理 IndexReferenceValue
	if ref, ok := v.(*data.IndexReferenceValue); ok {
		actualVal, _ := ref.Expr.GetValue(ref.Ctx)
		return unwrapValue(actualVal)
	}

	return v
}

func (b *BinaryDot) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if b == nil {
		return data.NewNullValue(), nil
	}
	// 获取左操作数的值
	if b.Left == nil {
		return data.NewNullValue(), nil
	}
	lv, lCtl := b.Left.GetValue(ctx)
	if lCtl != nil {
		return nil, lCtl
	}

	// 获取右操作数的值
	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	// 解包可能的 ZVal 或 ReferenceValue
	lv = unwrapValue(lv)
	rv = unwrapValue(rv)

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
