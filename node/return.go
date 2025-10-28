package node

import "github.com/php-any/origami/data"

func (u *ReturnStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if u.Value == nil {
		return nil, data.NewReturnControl(data.NewNullValue())
	}
	v, ctl := u.Value.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return nil, data.NewReturnControl(v.(data.Value))
}

// ReturnStatement 表示return语句
type ReturnStatement struct {
	*Node `pp:"-"`
	Value data.GetValue
}

// NewReturnStatement 创建一个新的return语句
func NewReturnStatement(from *TokenFrom, value data.GetValue) *ReturnStatement {
	return &ReturnStatement{
		Node:  NewNode(from),
		Value: value,
	}
}

// ReturnsStatement 表示多值 return 语句

type ReturnsStatement struct {
	*Node  `pp:"-"`
	Values []data.GetValue
}

// NewReturnsStatement 创建一个新的多值 return 语句
func NewReturnsStatement(from *TokenFrom, values []data.GetValue) *ReturnsStatement {
	return &ReturnsStatement{
		Node:   NewNode(from),
		Values: values,
	}
}

func (u *ReturnsStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	var result []data.Value
	for _, expr := range u.Values {
		v, ctl := expr.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		if val, ok := v.(data.Value); ok {
			result = append(result, val)
		} else {
			result = append(result, data.NewNullValue())
		}
	}
	return nil, data.NewReturnControl(data.NewArrayValue(result))
}
