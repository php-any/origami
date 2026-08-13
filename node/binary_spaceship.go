package node

import (
	"github.com/php-any/origami/data"
)

// BinarySpaceship 表示 spaceship 运算符 (<=>)
type BinarySpaceship struct {
	*Node
	Left  data.GetValue
	Right data.GetValue
}

// NewBinarySpaceship 创建一个新的 spaceship 运算符节点
func NewBinarySpaceship(from data.From, left, right data.GetValue) *BinarySpaceship {
	return &BinarySpaceship{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

// GetValue 获取 spaceship 运算符的值
// 返回 -1（左<右）、0（左==右）或 1（左>右）
func (b *BinarySpaceship) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	leftVal, c := b.Left.GetValue(ctx)
	if c != nil {
		return nil, c
	}

	rightVal, c := b.Right.GetValue(ctx)
	if c != nil {
		return nil, c
	}

	// 调试：直接检查值的类型
	if leftInt, ok := leftVal.(*data.IntValue); ok {
		if rightInt, ok := rightVal.(*data.IntValue); ok {
			n1 := leftInt.Value
			n2 := rightInt.Value
			if n1 < n2 {
				return data.NewIntValue(-1), nil
			} else if n1 > n2 {
				return data.NewIntValue(1), nil
			}
			return data.NewIntValue(0), nil
		}
	}

	// 尝试转换为可比较的值
	leftComparable, leftOk := leftVal.(data.Value)
	rightComparable, rightOk := rightVal.(data.Value)

	if !leftOk || !rightOk {
		// 如果无法转换，返回 0
		return data.NewIntValue(0), nil
	}

	// 使用 Compare 方法进行比较
	result := data.Compare(leftComparable, rightComparable)

	return data.NewIntValue(result), nil
}
