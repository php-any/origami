package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArraySearchFunction 实现 array_search 函数
// array_search(mixed $needle, array $haystack, bool $strict = false): int|string|false
type ArraySearchFunction struct{}

func NewArraySearchFunction() data.FuncStmt {
	return &ArraySearchFunction{}
}

func (f *ArraySearchFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	needleValue, _ := ctx.GetIndexValue(0)
	haystackValue, _ := ctx.GetIndexValue(1)
	strictValue, _ := ctx.GetIndexValue(2)

	if needleValue == nil || haystackValue == nil {
		return data.NewBoolValue(false), nil
	}

	arr, ok := haystackValue.(*data.ArrayValue)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	strict := false
	if strictValue != nil {
		if b, ok := strictValue.(*data.BoolValue); ok {
			strict = b.Value
		}
	}

	needleStr := needleValue.AsString()
	for i, z := range arr.List {
		v := z.Value
		match := false
		if strict {
			match = (v.AsString() == needleStr)
		} else {
			match = (v.AsString() == needleStr)
		}
		if match {
			if z.Name != "" {
				return data.NewStringValue(z.Name), nil
			}
			return data.NewIntValue(i), nil
		}
	}

	return data.NewBoolValue(false), nil
}

func (f *ArraySearchFunction) GetName() string {
	return "array_search"
}

func (f *ArraySearchFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "needle", 0, nil, nil),
		node.NewParameter(nil, "haystack", 1, nil, nil),
		node.NewParameter(nil, "strict", 2, data.NewBoolValue(false), nil),
	}
}

func (f *ArraySearchFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "needle", 0, nil),
		node.NewVariable(nil, "haystack", 1, data.NewBaseType("array")),
		node.NewVariable(nil, "strict", 2, data.NewBaseType("bool")),
	}
}
