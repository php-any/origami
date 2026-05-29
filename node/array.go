package node

import (
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
	arr := []data.Value{}
	for _, statement := range n.V {
		v, acl := statement.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}

		// 检查是否是展开运算符
		if spread, ok := statement.(*ArraySpread); ok {
			// 展开数组元素
			spreadValue, acl := spread.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}
			if arrayValue, ok := spreadValue.(*data.ArrayValue); ok {
				// 将展开的数组元素添加到结果数组中
				arr = append(arr, arrayValue.ToValueList()...)
			} else {
				return nil, data.NewErrorThrow(n.from, data.NewError(n.from, "展开运算符只能用于数组", nil))
			}
		} else {
			// 普通元素
			arr = append(arr, v.(data.Value))
		}
	}
	av := data.NewArrayValue(arr).(*data.ArrayValue)
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
