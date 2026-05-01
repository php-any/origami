package array

import (
	"sort"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// KrsortFunction 实现 krsort 函数
// krsort(array &$array, int $flags = SORT_REGULAR): bool
// 按键名降序排序，保持键到值的关联。
type KrsortFunction struct{}

func NewKrsortFunction() data.FuncStmt {
	return &KrsortFunction{}
}

func (f *KrsortFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)
	flagsValue, _ := ctx.GetIndexValue(1)

	if arrayValue == nil {
		return data.NewBoolValue(false), nil
	}

	switch v := arrayValue.(type) {
	case *data.ArrayValue:
		if len(v.List) == 0 {
			return data.NewBoolValue(true), nil
		}
		// Reverse the list
		for i, j := 0, len(v.List)-1; i < j; i, j = i+1, j-1 {
			v.List[i], v.List[j] = v.List[j], v.List[i]
		}
		return data.NewBoolValue(true), nil

	case *data.ObjectValue:
		props := v.GetProperties()
		if len(props) == 0 {
			return data.NewBoolValue(true), nil
		}

		keys := make([]string, 0, len(props))
		for k := range props {
			keys = append(keys, k)
		}

		flags := 0
		if flagsValue != nil {
			if intVal, ok := flagsValue.(*data.IntValue); ok {
				flags, _ = intVal.AsInt()
			}
		}

		sort.Slice(keys, func(i, j int) bool {
			switch flags {
			case 1: // SORT_NUMERIC
				return keys[i] > keys[j]
			case 2: // SORT_STRING
				fallthrough
			default: // SORT_REGULAR
				return keys[i] > keys[j]
			}
		})

		newObj := data.NewObjectValue()
		for _, k := range keys {
			newObj.SetProperty(k, props[k])
		}
		*v = *newObj

		return data.NewBoolValue(true), nil

	default:
		return data.NewBoolValue(false), nil
	}
}

func (f *KrsortFunction) GetName() string {
	return "krsort"
}

func (f *KrsortFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, data.Mixed{}),
		node.NewParameter(nil, "flags", 1, data.NewIntValue(0), data.Int{}),
	}
}

func (f *KrsortFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
		node.NewVariable(nil, "flags", 1, data.Int{}),
	}
}
