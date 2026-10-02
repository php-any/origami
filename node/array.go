package node

import "github.com/php-any/origami/data"

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

	for _, statement := range n.V {
		// 检查是否是展开运算符
		if spread, ok := statement.(*ArraySpread); ok {
			spreadValue, acl := spread.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}
			arrayValue, ok := spreadValue.(*data.ArrayValue)
			if ok {
				for arraySlots12, arrayPosition12 := arrayValue.View(), 0; arrayPosition12 < arraySlots12.Len(); arrayPosition12++ {
					z := arraySlots12.At(arrayPosition12)
					if z == nil {
						continue
					}
					if key, ok := z.PHPArrayKey(arrayPosition12).(*data.StringValue); ok {
						av.SetStringKey(key.Value, z.Value)
					} else if !av.AppendValue(z.Value) {
						return nil, arrayLiteralAppendError(n.GetFrom())
					}
				}
				continue
			}
			// 支持 Generator 展开到数组字面量
			vals, spreadCtl := spreadToValues(ctx, spreadValue)
			if spreadCtl != nil {
				return nil, spreadCtl
			}
			for _, val := range vals {
				if !av.AppendValue(val) {
					return nil, arrayLiteralAppendError(n.GetFrom())
				}
			}
			continue
		}

		v, acl := statement.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		if !av.AppendValue(v.(data.Value)) {
			return nil, arrayLiteralAppendError(n.GetFrom())
		}
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
		if !av.SetKey(kv.(data.Value), vv.(data.Value)) {
			return nil, data.NewErrorThrowByName(n.GetFrom(), data.NewError(n.GetFrom(), "Illegal offset type", nil), "TypeError")
		}
	}
	return av, nil
}

func arrayLiteralAppendError(from data.From) data.Control {
	return data.NewErrorThrow(from, data.NewError(nil, "Cannot add element to the array as the next element is already occupied", nil))
}
