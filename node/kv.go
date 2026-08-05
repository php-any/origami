package node

import (
	"errors"
	"strconv"

	"github.com/php-any/origami/data"
)

// KvPair 表示一个键值对。Key == nil 表示 ...$spread 展开项，Value 为 ArraySpread。
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

// GetValue 获取数字字面量的值
func (n *Kv) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	obj := data.NewObjectValue()

	for _, pair := range n.V {
		// Key == nil：...$array 展开，保留键名合并进结果
		if pair.Key == nil {
			if acl := mergeSpreadIntoObject(ctx, obj, pair.Value, n.from); acl != nil {
				return nil, acl
			}
			continue
		}

		kv, acl := pair.Key.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		vv, acl := pair.Value.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		if kv == nil {
			return nil, data.NewErrorThrow(n.from, errors.New("数组键求值结果为 null"))
		}
		if vv == nil {
			vv = data.NewNullValue()
		}

		keyVal, ok := kv.(data.Value)
		if !ok {
			return nil, data.NewErrorThrow(n.from, errors.New("数组键类型无效"))
		}
		valVal, ok := vv.(data.Value)
		if !ok {
			return nil, data.NewErrorThrow(n.from, errors.New("数组值类型无效"))
		}
		acl = obj.SetProperty(keyVal.AsString(), valVal)
		if acl != nil {
			return nil, acl
		}
	}
	return obj, nil
}

func mergeSpreadIntoObject(ctx data.Context, obj *data.ObjectValue, spread data.GetValue, from data.From) data.Control {
	spreadValue, acl := spread.GetValue(ctx)
	if acl != nil {
		return acl
	}
	switch sv := spreadValue.(type) {
	case *data.ArrayValue:
		for i, z := range sv.List {
			if z == nil {
				continue
			}
			key := z.Name
			if key == "" {
				key = strconv.Itoa(i)
			}
			if acl := obj.SetProperty(key, z.Value); acl != nil {
				return acl
			}
		}
		return nil
	case *data.ObjectValue:
		var mergeAcl data.Control
		sv.RangeProperties(func(key string, value data.Value) bool {
			if acl := obj.SetProperty(key, value); acl != nil {
				mergeAcl = acl
				return false
			}
			return true
		})
		return mergeAcl
	default:
		return data.NewErrorThrow(from, errors.New("展开运算符只能用于数组"))
	}
}
