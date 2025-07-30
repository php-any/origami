package node

import (
	"github.com/php-any/origami/data"
)

type Kv struct {
	*Node `pp:"-"`
	V     map[Statement]Statement
}

func NewKv(token *TokenFrom, v map[Statement]Statement) data.GetValue {
	return &Kv{
		Node: NewNode(token),
		V:    v,
	}
}

// GetValue 获取数字字面量的值
func (n *Kv) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	obj := data.NewObjectValue()

	for k, v := range n.V {
		kv, acl := k.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		vv, acl := v.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}

		obj.SetProperty(kv.(data.Value).AsString(), vv.(data.Value))
	}
	return obj, nil
}
