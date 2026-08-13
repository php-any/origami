package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayUdiffFunction 实现 array_udiff
// array_udiff(array $array, array ...$arrays, callable $value_compare_func): array
type ArrayUdiffFunction struct{}

func NewArrayUdiffFunction() data.FuncStmt {
	return &ArrayUdiffFunction{}
}

func (f *ArrayUdiffFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	paramsVal, _ := ctx.GetIndexValue(0)
	params := paramsToValueList(paramsVal)
	if len(params) < 2 {
		if len(params) == 1 {
			return params[0], nil
		}
		return data.NewArrayValue([]data.Value{}), nil
	}

	callback := params[len(params)-1]
	arrays := params[:len(params)-1]
	baseEntries := toKVEntries(arrays[0])
	if len(baseEntries) == 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}

	otherValues := make([][]data.Value, 0, len(arrays)-1)
	for i := 1; i < len(arrays); i++ {
		entries := toKVEntries(arrays[i])
		vals := make([]data.Value, len(entries))
		for j, e := range entries {
			vals[j] = e.value
		}
		otherValues = append(otherValues, vals)
	}

	resultEntries := make([]kvEntry, 0)
	for _, e := range baseEntries {
		found := false
		for _, other := range otherValues {
			for _, ov := range other {
				if compareCallback(ctx, callback, e.value, ov) == 0 {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			resultEntries = append(resultEntries, e)
		}
	}
	return buildResultFromEntries(resultEntries, true), nil
}

func (f *ArrayUdiffFunction) GetName() string { return "array_udiff" }

func (f *ArrayUdiffFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameters(nil, "arrays", 0, nil, nil)}
}

func (f *ArrayUdiffFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "arrays", 0, data.NewBaseType("array"))}
}
