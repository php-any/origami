package array

import (
	"sort"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// KsortFunction 实现 ksort 函数
// ksort(array &$array, int $flags = SORT_REGULAR): bool
// 按键名升序排序，保持键到值的关联，不重新索引键。
type KsortFunction struct{}

func NewKsortFunction() data.FuncStmt {
	return &KsortFunction{}
}

func (f *KsortFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)
	flagsValue, _ := ctx.GetIndexValue(1) // 可选 flags

	if arrayValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 支持两种内部表示：
	// - ArrayValue: 索引数组（int 键）
	// - ObjectValue: 关联数组（string 键）
	switch v := arrayValue.(type) {
	case *data.ArrayValue:
		// 空数组，直接返回
		if len(v.List) == 0 {
			return data.NewBoolValue(true), nil
		}
		// 对整数键 0..n-1 的数组按键排序等价于不变，这里直接返回 true
		return data.NewBoolValue(true), nil

	case *data.ObjectValue:
		props := v.GetProperties()
		if len(props) == 0 {
			return data.NewBoolValue(true), nil
		}

		// 收集所有键
		keys := make([]string, 0, len(props))
		for k := range props {
			keys = append(keys, k)
		}

		// 解析 flags（当前仅支持 SORT_REGULAR / SORT_STRING / SORT_NUMERIC）
		flags := 0
		if flagsValue != nil {
			if intVal, ok := flagsValue.(*data.IntValue); ok {
				flags, _ = intVal.AsInt()
			}
		}

		// 根据 flags 对键排序
		sort.Slice(keys, func(i, j int) bool {
			ki := keys[i]
			kj := keys[j]
			switch flags {
			case 1: // SORT_NUMERIC
				// 简化实现：按字符串数值比较
				return ki < kj
			case 2: // SORT_STRING
				fallthrough
			default: // SORT_REGULAR
				return ki < kj
			}
		})

		// 清空并按排序后的键顺序重新设置属性。
		// 注意：ObjectValue 使用的是 OrderedMap，重复 SetProperty 同名键会更新而不改变顺序，
		// 所以这里通过整体结构赋值 (*v = *newObj) 来替换内部的 OrderedMap。
		newObj := data.NewObjectValue()
		for _, k := range keys {
			newObj.SetProperty(k, props[k])
		}
		// 用新对象整体替换旧对象的内部状态
		*v = *newObj

		return data.NewBoolValue(true), nil

	default:
		return data.NewBoolValue(false), nil
	}
}

func (f *KsortFunction) GetName() string {
	return "ksort"
}

func (f *KsortFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, data.Mixed{}),
		node.NewParameter(nil, "flags", 1, data.NewIntValue(0), data.Int{}),
	}
}

func (f *KsortFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
		node.NewVariable(nil, "flags", 1, data.Int{}),
	}
}
