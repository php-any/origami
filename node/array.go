package node

import (
	"github.com/php-any/origami/data"
)

type Array struct {
	*Node `pp:"-"`
	V     []Statement
}

func NewArray(token *TokenFrom, arr []Statement) data.GetValue {
	return &Array{
		Node: NewNode(token),
		V:    arr,
	}
}

// GetValue 获取数字字面量的值
func (n *Array) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	arr := make([]data.Value, len(n.V))
	for i, statement := range n.V {
		v, acl := statement.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		arr[i] = v.(data.Value)
	}
	return data.NewArrayValue(arr), nil
}
