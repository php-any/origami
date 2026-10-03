package node

import (
	"github.com/php-any/origami/data"
)

// BinaryEqStrict 表示严格相等表达式
type BinaryEqStrict struct {
	*Node
	Left  data.GetValue
	Right data.GetValue
}

// NewBinaryEqStrict 创建一个新的严格相等表达式
func NewBinaryEqStrict(from data.From, left data.GetValue, right data.GetValue) *BinaryEqStrict {
	return &BinaryEqStrict{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

// GetValue 获取严格相等表达式的值
func (b *BinaryEqStrict) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 计算左操作数
	leftValue, c := b.Left.GetValue(ctx)
	if c != nil {
		return nil, c
	}

	// 计算右操作数
	rightValue, c := b.Right.GetValue(ctx)
	if c != nil {
		return nil, c
	}
	if leftValue == nil {
		leftValue = data.NewNullValue()
	}
	if rightValue == nil {
		rightValue = data.NewNullValue()
	}

	// 严格相等比较：类型和值都必须相等
	result := isStrictEqual(leftValue, rightValue)
	return data.NewBoolValue(result), nil
}

func isPHPNull(v data.GetValue) bool {
	if v == nil {
		return true
	}
	_, ok := v.(*data.NullValue)
	return ok
}

func isTypedNilValue(v data.GetValue) bool {
	if v == nil {
		return true
	}
	switch t := v.(type) {
	case *data.IntValue:
		return t == nil
	case *data.FloatValue:
		return t == nil
	case *data.BoolValue:
		return t == nil
	case *data.StringValue:
		return t == nil
	case *data.NullValue:
		return t == nil
	case *data.ArrayValue:
		return t == nil

	case *data.ClassValue:
		return t == nil
	case *data.ThisValue:
		return t == nil
	case *data.FuncValue:
		return t == nil
	case *data.BoundFuncValue:
		return t == nil
	default:
		return false
	}
}

func typedNilEquals(other data.GetValue) bool {
	return isPHPNull(other) || isTypedNilValue(other)
}

// isStrictEqual 进行严格相等比较
// 这是一个独立的工具函数，供 BinaryEqStrict 和 BinaryNeStrict 复用
func isStrictEqual(value1, value2 data.GetValue) bool {
	if value1 == nil {
		return typedNilEquals(value2)
	}
	switch v1 := value1.(type) {
	case *data.IntValue:
		if v1 == nil {
			return typedNilEquals(value2)
		}
		if v2, ok2 := value2.(*data.IntValue); ok2 {
			return v2 != nil && v1.Value == v2.Value
		}
		return false
	case *data.FloatValue:
		if v1 == nil {
			return typedNilEquals(value2)
		}
		if v2, ok2 := value2.(*data.FloatValue); ok2 {
			return v2 != nil && v1.Value == v2.Value
		}
		return false
	case *data.BoolValue:
		if v1 == nil {
			return typedNilEquals(value2)
		}
		if v2, ok2 := value2.(*data.BoolValue); ok2 {
			return v2 != nil && v1.Value == v2.Value
		}
		return false
	case *data.StringValue:
		if v1 == nil {
			return typedNilEquals(value2)
		}
		if v2, ok2 := value2.(*data.StringValue); ok2 {
			return v2 != nil && v1.Value == v2.Value
		}
		return false
	case *data.NullValue:
		return isPHPNull(value2) || isTypedNilValue(value2)
	case *data.ArrayValue:
		if v1 == nil {
			return typedNilEquals(value2)
		}
		if v2, ok2 := value2.(*data.ArrayValue); ok2 {
			if v2 == nil {
				return false
			}
			// 数组比较：长度和每个元素都相等
			if v1.Len() != v2.Len() {
				return false
			}
			for arraySlots16, i := v1.View(), 0; i < arraySlots16.Len(); i++ {
				zval1 := arraySlots16.At(i)
				zval2 := v2.At(i)
				if zval1 == nil || zval2 == nil {
					if zval1 == nil && zval2 == nil {
						continue
					}
					return false
				}
				// 递归比较数组元素
				if !zval1.SameArrayKey(i, zval2, i) {
					return false
				}
				if !isStrictEqual(zval1.ReadValue(), zval2.ReadValue()) {
					return false
				}
			}
			return true
		}
		// 空 ArrayValue 与空 ObjectValue（关联数组）在 PHP 中均为 []

		return false

	case *data.ClassValue:
		if v1 == nil {
			return typedNilEquals(value2)
		}
		// PHP === 对对象比较同一性（同一实例）。
		// 方法调用会为 $this 再包一层 ClassValue，但共享同一 ObjectValue。
		switch v2 := value2.(type) {
		case *data.ClassValue:
			return sameClassInstance(v1, v2)
		case *data.ThisValue:
			return sameClassInstance(v1, v2.ClassValue)
		}
		return false
	case *data.ThisValue:
		if v1 == nil {
			return typedNilEquals(value2)
		}
		switch v2 := value2.(type) {
		case *data.ThisValue:
			return sameClassInstance(v1.ClassValue, v2.ClassValue)
		case *data.ClassValue:
			return sameClassInstance(v1.ClassValue, v2)
		}
		return false
	case *data.FuncValue:
		if v1 == nil {
			return typedNilEquals(value2)
		}
		// 闭包 === ：同一实例；禁止与 null 经 AsString("") 误判相等（Livewire EventBus finish）
		if v2, ok := value2.(*data.FuncValue); ok {
			return sameFuncValue(v1, v2)
		}
		if v2, ok := value2.(*data.BoundFuncValue); ok {
			return sameFuncValue(v1, &v2.FuncValue)
		}
		return false
	case *data.BoundFuncValue:
		if v1 == nil {
			return typedNilEquals(value2)
		}
		if v2, ok := value2.(*data.BoundFuncValue); ok {
			return sameFuncValue(&v1.FuncValue, &v2.FuncValue) && v1.BoundObject == v2.BoundObject
		}
		if v2, ok := value2.(*data.FuncValue); ok {
			return sameFuncValue(&v1.FuncValue, v2)
		}
		return false
	default:
		// 保留历史：部分值经 AsString 比较；但 null 绝不能与非 null 因 "" 相等
		if _, ok := value2.(*data.NullValue); ok {
			return false
		}
		if _, ok := value1.(*data.NullValue); ok {
			return false
		}
		if strValue1, ok := value1.(data.AsString); ok {
			if strValue2, ok := value2.(data.AsString); ok {
				return strValue1.AsString() == strValue2.AsString()
			}
		}
		return false
	}
}

// sameFuncValue 判断两个 FuncValue 是否为同一闭包实例。
func sameFuncValue(a, b *data.FuncValue) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a == b {
		return true
	}
	return a.Value != nil && a.Value == b.Value
}

// sameClassInstance 判断两个 ClassValue 是否指向同一 PHP 对象实例。
func sameClassInstance(a, b *data.ClassValue) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a == b {
		return true
	}
	return a.ObjectValue != nil && a.ObjectValue == b.ObjectValue
}
