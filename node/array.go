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
			if ok {
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
			// 支持 Generator 展开到数组字面量
			vals, spreadCtl := spreadToValues(ctx, spreadValue)
			if spreadCtl != nil {
				return nil, spreadCtl
			}
			for _, val := range vals {
				setArrayLiteralEntry(av, data.NewIntValue(nextIndex), val)
				nextIndex++
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
	// PHP：纯数字字符串键当作 int；非数字字符串必须走字符串键。
	// 旧逻辑对实现了 AsInt 的 StringValue 一律走 int 分支，且忽略 AsInt 错误，
	// 导致 "lazy" 等键被写成 List[0] 且 Name 为空（[...$assoc] 丢键）。
	if sv, ok := key.(*data.StringValue); ok {
		if n, isIntKey := data.ParseIntArrayKeyName(sv.Value); isIntKey {
			setArrayLiteralIntKey(av, n, val)
			return
		}
		setArrayLiteralStringKey(av, sv.Value, val)
		return
	}
	if iv, ok := key.(*data.IntValue); ok {
		setArrayLiteralIntKey(av, iv.Value, val)
		return
	}
	if ai, ok := key.(data.AsInt); ok {
		if i, err := ai.AsInt(); err == nil {
			setArrayLiteralIntKey(av, i, val)
			return
		}
	}
	setArrayLiteralStringKey(av, key.AsString(), val)
}

func setArrayLiteralIntKey(av *data.ArrayValue, i int, val data.Value) {
	if i < 0 {
		setArrayLiteralStringKey(av, data.IntArrayKeyName(i), val)
		return
	}
	for len(av.List) <= i {
		av.List = append(av.List, data.NewZVal(data.NewNullValue()))
	}
	av.List[i] = data.NewZVal(val)
}

func setArrayLiteralStringKey(av *data.ArrayValue, keyStr string, val data.Value) {
	for _, z := range av.List {
		if z != nil && z.Name == keyStr {
			z.Value = val
			return
		}
	}
	av.List = append(av.List, data.NewNamedZVal(keyStr, val))
}
