package node

import "github.com/php-any/origami/data"

type NamedArgument struct {
	*Node `pp:"-"`
	Name  string
	Value data.GetValue
}

func NewNamedArgument(token *TokenFrom, name string, value data.GetValue) *NamedArgument {
	return &NamedArgument{
		Node:  NewNode(token),
		Name:  name,
		Value: value,
	}
}

// GetValue 获取变量表达式的值
func (v *NamedArgument) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return v.Value.GetValue(ctx)
}
