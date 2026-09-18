package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayCombineFunction 实现 array_combine 函数
// array_combine(array $keys, array $values): array|false
//
// 必须返回 ArrayValue（含字符串键用 Name 保存），不能返回 ObjectValue：
// Laravel Arr::map 依赖 array_combine；ObjectValue 会导致 gettype=object、
// Collection::items 形态异常，进而在仪表盘 widgets 路径上失控。
type ArrayCombineFunction struct{}

func NewArrayCombineFunction() data.FuncStmt {
	return &ArrayCombineFunction{}
}

func (f *ArrayCombineFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	keysVal, _ := ctx.GetIndexValue(0)
	valuesVal, _ := ctx.GetIndexValue(1)

	if keysVal == nil || valuesVal == nil {
		return data.NewBoolValue(false), nil
	}

	keys := valueListForCombine(keysVal)
	values := valueListForCombine(valuesVal)
	if keys == nil || values == nil {
		return data.NewBoolValue(false), nil
	}

	if len(keys) != len(values) {
		return data.NewBoolValue(false), nil
	}

	result := data.NewArrayValue(nil).(*data.ArrayValue)
	for i, key := range keys {
		setCombineEntry(result, key, values[i])
	}
	return result, nil
}

func valueListForCombine(v data.Value) []data.Value {
	switch t := v.(type) {
	case *data.ArrayValue:
		return t.ToValueList()
	case *data.ObjectValue:
		out := make([]data.Value, 0)
		t.RangeProperties(func(_ string, val data.Value) bool {
			out = append(out, val)
			return true
		})
		return out
	default:
		return nil
	}
}

func setCombineEntry(av *data.ArrayValue, key, val data.Value) {
	if key == nil {
		key = data.NewStringValue("")
	}
	if val == nil {
		val = data.NewNullValue()
	}

	if iv, ok := key.(*data.IntValue); ok {
		setCombineIntKey(av, iv.Value, val)
		return
	}
	if sv, ok := key.(*data.StringValue); ok {
		if n, ok := data.ParseIntArrayKeyName(sv.Value); ok {
			setCombineIntKey(av, n, val)
			return
		}
		setCombineStringKey(av, sv.Value, val)
		return
	}
	if ai, ok := key.(data.AsInt); ok {
		if i, err := ai.AsInt(); err == nil {
			// 非纯数字字符串不能走 AsInt 成功路径（与 array 字面量一致）
			if _, isStr := key.(*data.StringValue); !isStr {
				setCombineIntKey(av, i, val)
				return
			}
		}
	}
	setCombineStringKey(av, key.AsString(), val)
}

func setCombineIntKey(av *data.ArrayValue, i int, val data.Value) {
	if i < 0 {
		setCombineStringKey(av, data.IntArrayKeyName(i), val)
		return
	}
	for len(av.List) <= i {
		av.List = append(av.List, data.NewZVal(data.NewNullValue()))
	}
	av.List[i] = data.NewZVal(val)
}

func setCombineStringKey(av *data.ArrayValue, keyStr string, val data.Value) {
	for _, z := range av.List {
		if z != nil && z.Name == keyStr {
			z.Value = val
			return
		}
	}
	av.List = append(av.List, data.NewNamedZVal(keyStr, val))
}

func (f *ArrayCombineFunction) GetName() string {
	return "array_combine"
}

func (f *ArrayCombineFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "keys", 0, nil, nil),
		node.NewParameter(nil, "values", 1, nil, nil),
	}
}

func (f *ArrayCombineFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "keys", 0, data.Mixed{}),
		node.NewVariable(nil, "values", 1, data.Mixed{}),
	}
}
