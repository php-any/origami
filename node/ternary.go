package node

import (
	"github.com/php-any/origami/data"
)

// TernaryExpression 表示三目运算符表达式
type TernaryExpression struct {
	*Node      `pp:"-"`
	Condition  data.GetValue // 条件表达式
	TrueValue  data.GetValue // 真值表达式
	FalseValue data.GetValue // 假值表达式
}

// NewTernaryExpression 创建一个新的三目运算符表达式
func NewTernaryExpression(from *TokenFrom, condition, trueValue, falseValue data.GetValue) *TernaryExpression {
	return &TernaryExpression{
		Node:       NewNode(from),
		Condition:  condition,
		TrueValue:  trueValue,
		FalseValue: falseValue,
	}
}

// GetValue 获取三目运算符表达式的值
func (t *TernaryExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 计算条件表达式的值
	conditionValue, ctl := t.Condition.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}

	// 将条件值转换为布尔值
	var condition bool
	switch v := conditionValue.(type) {
	case *data.BoolValue:
		condition = v.Value
	case *data.IntValue:
		condition = v.Value != 0
	case *data.FloatValue:
		condition = v.Value != 0
	case *data.StringValue:
		condition = len(v.Value) > 0
	case *data.NullValue:
		condition = false
	case *data.ArrayValue:
		condition = len(v.Value) > 0
	default:
		// 对于其他类型，尝试转换为布尔值
		if boolValue, ok := v.(data.AsBool); ok {
			if val, err := boolValue.AsBool(); err == nil {
				condition = val
			} else {
				condition = false
			}
		} else {
			condition = false
		}
	}

	// 根据条件选择真值或假值
	if condition {
		return t.TrueValue.GetValue(ctx)
	} else {
		return t.FalseValue.GetValue(ctx)
	}
}

// AsString 返回三目运算符表达式的字符串表示
func (t *TernaryExpression) AsString() string {
	return "ternary_expression"
}
