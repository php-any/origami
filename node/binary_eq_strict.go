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

	// 严格相等比较：类型和值都必须相等
	result := isStrictEqual(leftValue, rightValue)
	return data.NewBoolValue(result), nil
}

// isStrictEqual 进行严格相等比较
// 这是一个独立的工具函数，供 BinaryEqStrict 和 BinaryNeStrict 复用
func isStrictEqual(value1, value2 data.GetValue) bool {
	switch v1 := value1.(type) {
	case *data.IntValue:
		if v2, ok2 := value2.(*data.IntValue); ok2 {
			return v1.Value == v2.Value
		}
		return false
	case *data.FloatValue:
		if v2, ok2 := value2.(*data.FloatValue); ok2 {
			return v1.Value == v2.Value
		}
		return false
	case *data.BoolValue:
		if v2, ok2 := value2.(*data.BoolValue); ok2 {
			return v1.Value == v2.Value
		}
		return false
	case *data.StringValue:
		if v2, ok2 := value2.(*data.StringValue); ok2 {
			return v1.Value == v2.Value
		}
		return false
	case *data.NullValue:
		if _, ok2 := value2.(*data.NullValue); ok2 {
			return true
		}
		return false
	case *data.ArrayValue:
		if v2, ok2 := value2.(*data.ArrayValue); ok2 {
			// 数组比较：长度和每个元素都相等
			if len(v1.List) != len(v2.List) {
				return false
			}
			for i, zval1 := range v1.List {
				zval2 := v2.List[i]
				// 递归比较数组元素
				if !isStrictEqual(zval1.Value, zval2.Value) {
					return false
				}
			}
			return true
		}
		// 空 ArrayValue 与空 ObjectValue（关联数组）在 PHP 中均为 []
		if v2, ok2 := value2.(*data.ObjectValue); ok2 {
			return len(v1.List) == 0 && len(v2.GetProperties()) == 0
		}
		return false
	case *data.ObjectValue:
		if v2, ok2 := value2.(*data.ArrayValue); ok2 {
			return len(v1.GetProperties()) == 0 && len(v2.List) == 0
		}
		if v2, ok2 := value2.(*data.ObjectValue); ok2 {
			// 对象比较：属性数量和每个属性都相等
			props1 := v1.GetProperties()
			props2 := v2.GetProperties()
			if len(props1) != len(props2) {
				return false
			}
			for key, val1 := range props1 {
				val2, exists := props2[key]
				if !exists {
					return false
				}
				// 递归比较对象属性
				if !isStrictEqual(val1, val2) {
					return false
				}
			}
			return true
		}
		return false
	case *data.ClassValue:
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
		switch v2 := value2.(type) {
		case *data.ThisValue:
			return sameClassInstance(v1.ClassValue, v2.ClassValue)
		case *data.ClassValue:
			return sameClassInstance(v1.ClassValue, v2)
		}
		return false
	default:
		// 对于其他类型，尝试字符串比较
		if strValue1, ok := value1.(data.AsString); ok {
			if strValue2, ok := value2.(data.AsString); ok {
				return strValue1.AsString() == strValue2.AsString()
			}
		}
		return false
	}
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
