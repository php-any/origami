package node

import (
	"strconv"

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

func (n *Kv) GetIndex() int {
	return -1
}

func (n *Kv) GetName() string {
	return "kv TODO"
}

func (n *Kv) GetType() data.Types {
	return nil
}

func (n *Kv) SetValue(ctx data.Context, value data.Value) data.Control {
	//TODO implement me
	panic("implement me")
}

func NewKv(token *TokenFrom, v []KvPair) data.GetValue {
	return &Kv{
		Node: NewNode(token),
		V:    v,
	}
}

// GetValue 获取键值对数组字面量的值
// 当所有键为连续整数 0,1,2,...,n-1 时返回 ArrayValue（PHP 索引数组语义），
// 否则返回 ObjectValue（PHP 关联数组语义）。
func (n *Kv) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 先求值所有键和值
	type evaluatedPair struct {
		key   data.Value
		value data.Value
	}
	pairs := make([]evaluatedPair, 0, len(n.V))
	for _, pair := range n.V {
		kv, acl := pair.Key.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		vv, acl := pair.Value.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		pairs = append(pairs, evaluatedPair{key: kv.(data.Value), value: vv.(data.Value)})
	}

	// 检查是否所有键为连续整数 0..n-1
	allSequential := true
	for i, p := range pairs {
		switch k := p.key.(type) {
		case *data.IntValue:
			if k.Value != i {
				allSequential = false
			}
		case *data.StringValue:
			// 尝试将字符串键解析为整数（PHP 会自动将 "0" 转为 0）
			if iv, err := strconv.Atoi(k.Value); err != nil || iv != i {
				allSequential = false
			}
		default:
			allSequential = false
		}
		if !allSequential {
			break
		}
	}

	if allSequential {
		// 返回 ArrayValue
		values := make([]data.Value, len(pairs))
		for i, p := range pairs {
			values[i] = p.value
		}
		return data.NewArrayValue(values), nil
	}

	// 返回 ObjectValue（关联数组）
	obj := data.NewObjectValue()
	for _, p := range pairs {
		acl := obj.SetProperty(p.key.AsString(), p.value)
		if acl != nil {
			return nil, acl
		}
	}
	return obj, nil
}
