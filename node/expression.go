package node

import "github.com/php-any/origami/data"

// UnaryExpression 表示一元表达式
type UnaryExpression struct {
	*Node    `pp:"-"`
	Operator string
	Right    data.GetValue
}

// NewUnaryExpression 创建一个新的一元表达式
func NewUnaryExpression(token *TokenFrom, operator string, right data.GetValue) data.GetValue {
	return &UnaryExpression{
		Node:     NewNode(token),
		Operator: operator,
		Right:    right,
	}
}

// GetValue 获取一元表达式的值
func (u *UnaryExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	right, control := u.Right.GetValue(ctx)
	if control != nil {
		return nil, control
	}

	switch u.Operator {
	case "-":
		if b, ok := right.(data.AsFloat); ok {
			bv, err := b.AsFloat()
			if err != nil {
				return nil, data.NewErrorThrow(u.from, err)
			}
			return data.NewFloatValue(-bv), nil
		}
	case "!":
		if b, ok := right.(data.AsBool); ok {
			bv, err := b.AsBool()
			if err != nil {
				return nil, data.NewErrorThrow(u.from, err)
			}
			return data.NewBoolValue(!bv), nil
		}
	}
	return nil, nil
}
