package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewArrayDiffFunction() data.FuncStmt {
	return &ArrayDiffFunction{}
}

type ArrayDiffFunction struct{}

func (fn *ArrayDiffFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	paramsVal, _ := ctx.GetIndexValue(0)
	if paramsVal == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	paramsArr, ok := paramsVal.(*data.ArrayValue)
	if !ok {
		return data.NewArrayValue([]data.Value{}), nil
	}

	arrays := make([][]data.Value, 0)
	for _, z := range paramsArr.List {
		if z.Value == nil {
			continue
		}
		switch v := z.Value.(type) {
		case *data.ArrayValue:
			arrays = append(arrays, v.ToValueList())
		case *data.ObjectValue:
			vals := make([]data.Value, 0)
			v.RangeProperties(func(key string, value data.Value) bool {
				vals = append(vals, value)
				return true
			})
			arrays = append(arrays, vals)
		}
	}

	if len(arrays) < 2 {
		if len(arrays) == 1 {
			return data.NewArrayValue(arrays[0]), nil
		}
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 以第一个数组为基准，移除在其他数组中出现的值
	first := arrays[0]
	others := arrays[1:]
	result := make([]data.Value, 0)
	for _, item := range first {
		found := false
		for _, other := range others {
			for _, o := range other {
				if item.AsString() == o.AsString() {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			result = append(result, item)
		}
	}

	return data.NewArrayValue(result), nil
}

func (fn *ArrayDiffFunction) GetName() string {
	return "array_diff"
}

func (fn *ArrayDiffFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "arrays", 0, nil, nil),
	}
}

func (fn *ArrayDiffFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "arrays", 0, data.NewBaseType("array")),
	}
}
