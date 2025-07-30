package node

import (
	"github.com/php-any/origami/data"
)

// NullCoalesceExpression 表示空合并运算符表达式
type NullCoalesceExpression struct {
	*Node `pp:"-"`
	Left  data.GetValue // 左操作数
	Right data.GetValue // 右操作数
}

// NewNullCoalesceExpression 创建一个新的空合并运算符表达式
func NewNullCoalesceExpression(from *TokenFrom, left, right data.GetValue) *NullCoalesceExpression {
	return &NullCoalesceExpression{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

// GetValue 获取空合并运算符表达式的值
func (n *NullCoalesceExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 计算左操作数的值
	leftValue, ctl := n.Left.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}

	// 检查左操作数是否为 null
	if leftValue == nil {
		// 如果左操作数为 nil，返回右操作数
		return n.Right.GetValue(ctx)
	}

	// 检查是否为 NullValue 类型
	switch leftValue.(type) {
	case *data.NullValue:
		// 如果左操作数为 null，返回右操作数
		return n.Right.GetValue(ctx)
	default:
		// 如果左操作数不为 null，返回左操作数
		return leftValue, nil
	}
}

// AsString 返回空合并运算符表达式的字符串表示
func (n *NullCoalesceExpression) AsString() string {
	return "null_coalesce_expression"
}
