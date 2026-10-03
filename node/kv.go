package node

import (
	"errors"

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
	// PHP 7.1+ 键名解构：['action' => $action, 'uri' => $uri] = $route
	for _, pair := range n.V {
		if pair.Key == nil || pair.Value == nil {
			continue
		}
		kv, acl := pair.Key.GetValue(ctx)
		if acl != nil {
			return acl
		}
		if kv == nil {
			return data.NewErrorThrow(n.from, errors.New("数组解构键求值结果为 null"))
		}
		keyVal, ok := kv.(data.Value)
		if !ok {
			return data.NewErrorThrow(n.from, errors.New("数组解构键类型无效"))
		}
		elem, acl := kvLookup(ctx, value, keyVal.AsString())
		if acl != nil {
			return acl
		}
		if elem == nil {
			elem = data.NewNullValue()
		}

		switch target := pair.Value.(type) {
		case data.Variable:
			if ctl := target.SetValue(ctx, elem); ctl != nil {
				return ctl
			}
		case interface {
			SetValue(data.Context, data.Value) data.Control
		}:
			if ctl := target.SetValue(ctx, elem); ctl != nil {
				return ctl
			}
		default:
			return data.NewErrorThrow(n.from, errors.New("数组解构左侧目标不可赋值"))
		}
	}
	return nil
}

func kvLookup(ctx data.Context, value data.Value, key string) (data.Value, data.Control) {
	switch v := value.(type) {

	case *data.ArrayValue:
		if slot, ok := v.LookupZValByStringKey(key); ok {
			return slot.ReadValue(), nil
		}
	case *data.ClassValue:
		if method, exists := v.GetMethod("offsetGet"); exists {
			fnCtx := implicitMethodFrame(ctx, v, method)
			if vars := method.GetVariables(); len(vars) > 0 {
				_ = fnCtx.SetVariableValue(vars[0], data.NewStringValue(key))
			}
			ret, ctl := method.Call(fnCtx)
			if ctl != nil {
				return nil, ctl
			}
			if rv, ok := ret.(data.Value); ok {
				return rv, nil
			}
			return data.NewNullValue(), nil
		}
		if v.ObjectValue != nil && v.ObjectValue.HasProperty(key) {
			return v.ObjectValue.GetProperty(key)
		}
	}
	return data.NewNullValue(), nil
}

func NewKv(token *TokenFrom, v []KvPair) data.GetValue {
	return &Kv{
		Node: NewNode(token),
		V:    v,
	}
}

// GetValue 获取数字字面量的值
func (n *Kv) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	array := data.NewArrayValueFromSlots(make([]*data.ZVal, 0, len(n.V)))

	for _, pair := range n.V {
		// Key == nil：...$array 展开，保留键名合并进结果
		if pair.Key == nil {
			if acl := mergeSpreadIntoArray(ctx, array, pair.Value, n.from); acl != nil {
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
		accepted, ctl := array.AssignKey(ctx, keyVal, valVal)
		if ctl != nil {
			return nil, ctl
		}
		if !accepted {
			return nil, data.NewErrorThrow(n.from, errors.New("Illegal offset type"))
		}
	}
	return array, nil
}

func mergeSpreadIntoArray(ctx data.Context, array *data.ArrayValue, spread data.GetValue, from data.From) data.Control {
	spreadValue, acl := spread.GetValue(ctx)
	if acl != nil {
		return acl
	}
	switch sv := spreadValue.(type) {
	case *data.ArrayValue:
		for arraySlots24, i := sv.View(), 0; i < arraySlots24.Len(); i++ {
			z := arraySlots24.At(i)
			if z == nil {
				continue
			}
			if key, ok := z.PHPArrayKey(i).(*data.StringValue); ok {
				array.SetStringKey(key.Value, z.ReadValue())
			} else if !array.AppendValue(z.ReadValue()) {
				return data.NewErrorThrow(from, errors.New("Cannot add element to the array as the next element is already occupied"))
			}
		}
		return nil

	default:
		return data.NewErrorThrow(from, errors.New("展开运算符只能用于数组"))
	}
}
