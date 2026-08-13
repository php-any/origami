package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayCountValuesFunction 实现 array_count_values
// array_count_values(array $array): array
type ArrayCountValuesFunction struct{}

func NewArrayCountValuesFunction() data.FuncStmt {
	return &ArrayCountValuesFunction{}
}

func (f *ArrayCountValuesFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayVal, _ := ctx.GetIndexValue(0)
	entries := toKVEntries(arrayVal)
	if len(entries) == 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}

	counts := make(map[string]int)
	for _, e := range entries {
		if _, ok := e.value.(*data.ArrayValue); ok {
			continue
		}
		if _, ok := e.value.(*data.ObjectValue); ok {
			continue
		}
		k := e.value.AsString()
		counts[k]++
	}

	result := data.NewObjectValue()
	for k, c := range counts {
		result.SetProperty(k, data.NewIntValue(c))
	}
	return result, nil
}

func (f *ArrayCountValuesFunction) GetName() string { return "array_count_values" }

func (f *ArrayCountValuesFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "array", 0, nil, nil)}
}

func (f *ArrayCountValuesFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "array", 0, data.NewBaseType("array"))}
}
