package node

import (
	"github.com/php-any/origami/data"
)

// KvPair 表示一个键值对
type KvPair struct {
	Key   data.GetValue
	Value data.GetValue
}

type Kv struct {
	*Node `pp:"-"`
	V     []KvPair // 使用切片保证顺序
}

func NewKv(token *TokenFrom, v []KvPair) data.GetValue {
	return &Kv{
		Node: NewNode(token),
		V:    v,
	}
}

// GetValue 获取数字字面量的值
func (n *Kv) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	obj := data.NewObjectValue()

	for _, pair := range n.V {
		kv, acl := pair.Key.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		vv, acl := pair.Value.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}

		obj.SetProperty(kv.(data.Value).AsString(), vv.(data.Value))
	}
	return obj, nil
}
