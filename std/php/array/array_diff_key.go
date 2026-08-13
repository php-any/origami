package array

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayDiffKeyFunction implements array_diff_key
type ArrayDiffKeyFunction struct{}

func NewArrayDiffKeyFunction() data.FuncStmt {
	return &ArrayDiffKeyFunction{}
}

func (f *ArrayDiffKeyFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	firstVal, _ := ctx.GetIndexValue(0)
	if firstVal == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 其余数组的 key 并集（index 1 是 variadic 打包）
	exclude := make(map[string]bool)
	if restVal, ok := ctx.GetIndexValue(1); ok && restVal != nil {
		if restArr, ok := restVal.(*data.ArrayValue); ok {
			for _, zv := range restArr.List {
				for k := range extractKeys(zv.Value) {
					exclude[k] = true
				}
			}
		}
	}

	switch first := firstVal.(type) {
	case *data.ArrayValue:
		result := data.NewArrayValue([]data.Value{}).(*data.ArrayValue)
		for idx, zv := range first.List {
			key := zv.Name
			if key == "" {
				key = fmt.Sprintf("%d", idx)
			}
			if !exclude[key] {
				result.List = append(result.List, &data.ZVal{Name: zv.Name, Value: zv.Value})
			}
		}
		return result, nil
	case *data.ObjectValue:
		result := data.NewObjectValue()
		first.RangeProperties(func(key string, val data.Value) bool {
			if !exclude[key] {
				result.SetProperty(key, val)
			}
			return true
		})
		return result, nil
	default:
		return data.NewArrayValue([]data.Value{}), nil
	}
}

func (f *ArrayDiffKeyFunction) GetName() string { return "array_diff_key" }

func (f *ArrayDiffKeyFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, nil),
		node.NewParameters(nil, "arrays", 1, nil, nil),
	}
}

func (f *ArrayDiffKeyFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "arrays", 1, data.NewBaseType("array")),
	}
}
