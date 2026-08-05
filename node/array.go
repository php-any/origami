package node

import (
	"strconv"

	"github.com/php-any/origami/data"
)

type Array struct {
	*Node `pp:"-"`
	V     []data.GetValue
	Keys  []KvPair // array(1, 2=>3) 中 => 之后的键值对
}

func NewArray(token *TokenFrom, arr []data.GetValue) data.GetValue {
	return &Array{
		Node: NewNode(token),
		V:    arr,
	}
}

func NewArrayWithKeys(token *TokenFrom, list []data.GetValue, keys []KvPair) data.GetValue {
	return &Array{
		Node: NewNode(token),
		V:    list,
		Keys: keys,
	}
}

// GetValue 获取数字字面量的值
func (n *Array) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	av := data.NewArrayValue(nil).(*data.ArrayValue)
	nextIndex := 0

	for _, statement := range n.V {
		// 检查是否是展开运算符
		if spread, ok := statement.(*ArraySpread); ok {
			spreadValue, acl := spread.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}
			arrayValue, ok := spreadValue.(*data.ArrayValue)
			if !ok {
				return nil, data.NewErrorThrow(n.from, data.NewError(n.from, "展开运算符只能用于数组", nil))
			}
			for _, z := range arrayValue.List {
				if z == nil {
					continue
				}
				if z.Name != "" {
					setArrayLiteralEntry(av, data.NewStringValue(z.Name), z.Value)
					if ik, err := strconv.Atoi(z.Name); err == nil && strconv.Itoa(ik) == z.Name && ik >= nextIndex {
						nextIndex = ik + 1
					}
				} else {
					setArrayLiteralEntry(av, data.NewIntValue(nextIndex), z.Value)
					nextIndex++
				}
			}
			continue
		}

		v, acl := statement.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		setArrayLiteralEntry(av, data.NewIntValue(nextIndex), v.(data.Value))
		nextIndex++
	}
	for _, pair := range n.Keys {
		kv, acl := pair.Key.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		vv, acl := pair.Value.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		setArrayLiteralEntry(av, kv.(data.Value), vv.(data.Value))
	}
	return av, nil
}

func setArrayLiteralEntry(av *data.ArrayValue, key, val data.Value) {
	if iv, ok := key.(data.AsInt); ok {
		i, _ := iv.AsInt()
		for len(av.List) <= i {
			av.List = append(av.List, data.NewZVal(data.NewNullValue()))
		}
		av.List[i] = data.NewZVal(val)
		return
	}
	keyStr := key.AsString()
	for _, z := range av.List {
		if z != nil && z.Name == keyStr {
			z.Value = val
			return
		}
	}
	av.List = append(av.List, &data.ZVal{Name: keyStr, Value: val})
}
