package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayUdiffAssocFunction 实现 array_udiff_assoc
// array_udiff_assoc(array $array, array ...$arrays, callable $value_compare_func): array
// 与 array_udiff 不同：键必须相同才参与值比较（Symfony UrlGenerator 用它算 query extra）。
type ArrayUdiffAssocFunction struct{}

func NewArrayUdiffAssocFunction() data.FuncStmt {
	return &ArrayUdiffAssocFunction{}
}

func (f *ArrayUdiffAssocFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
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

	others := make([][]kvEntry, 0, len(arrays)-1)
	for i := 1; i < len(arrays); i++ {
		others = append(others, toKVEntries(arrays[i]))
	}

	resultEntries := make([]kvEntry, 0)
	for _, e := range baseEntries {
		found := false
		for _, other := range others {
			for _, oe := range other {
				if oe.keyStr != e.keyStr {
					continue
				}
				if compareCallback(ctx, callback, e.value, oe.value) == 0 {
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

func (f *ArrayUdiffAssocFunction) GetName() string { return "array_udiff_assoc" }

func (f *ArrayUdiffAssocFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameters(nil, "arrays", 0, nil, nil)}
}

func (f *ArrayUdiffAssocFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "arrays", 0, data.NewBaseType("array"))}
}
