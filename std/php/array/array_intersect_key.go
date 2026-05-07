package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayIntersectKeyFunction implements array_intersect_key
type ArrayIntersectKeyFunction struct{}

func NewArrayIntersectKeyFunction() data.FuncStmt {
	return &ArrayIntersectKeyFunction{}
}

func (f *ArrayIntersectKeyFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	paramsValue, _ := ctx.GetIndexValue(0)
	if paramsValue == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}
	paramsArray, ok := paramsValue.(*data.ArrayValue)
	if !ok {
		return data.NewArrayValue([]data.Value{}), nil
	}
	arrays := paramsArray.ToValueList()
	if len(arrays) < 2 {
		return arrays[0], nil
	}

	var firstArr *data.ArrayValue
	switch v := arrays[0].(type) {
	case *data.ArrayValue:
		firstArr = v
	default:
		return data.NewArrayValue([]data.Value{}), nil
	}

	allKeys := make([]map[string]bool, 0)
	for i := 1; i < len(arrays); i++ {
		keys := make(map[string]bool)
		if arr, ok := arrays[i].(*data.ArrayValue); ok {
			for idx, zv := range arr.List {
				if zv.Name != "" {
					keys[zv.Name] = true
				} else {
					keys[data.NewIntValue(idx).AsString()] = true
				}
			}
		}
		allKeys = append(allKeys, keys)
	}

	resultList := make([]data.Value, 0)
	resultNames := make([]string, 0)
	for idx, zv := range firstArr.List {
		key := zv.Name
		if key == "" {
			key = data.NewIntValue(idx).AsString()
		}
		inAll := true
		for _, ks := range allKeys {
			if !ks[key] {
				inAll = false
				break
			}
		}
		if inAll {
			resultList = append(resultList, zv.Value)
			resultNames = append(resultNames, zv.Name)
		}
	}

	arr := data.NewArrayValue(resultList)
	for i, name := range resultNames {
		if i < len(arr.(*data.ArrayValue).List) {
			arr.(*data.ArrayValue).List[i].Name = name
		}
	}
	return arr, nil
}

func (f *ArrayIntersectKeyFunction) GetName() string { return "array_intersect_key" }

func (f *ArrayIntersectKeyFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "arrays", 0, nil, nil),
	}
}

func (f *ArrayIntersectKeyFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "arrays", 0, data.NewBaseType("array")),
	}
}
