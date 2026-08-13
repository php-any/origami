package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArraySumFunction 实现 array_sum
// array_sum(array $array): int|float
type ArraySumFunction struct{}

func NewArraySumFunction() data.FuncStmt {
	return &ArraySumFunction{}
}

func (f *ArraySumFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayVal, _ := ctx.GetIndexValue(0)
	entries := toKVEntries(arrayVal)
	sum := 0.0
	hasFloat := false
	for _, e := range entries {
		if _, ok := e.value.(*data.FloatValue); ok {
			hasFloat = true
		}
		sum += toFloatValue(e.value)
	}
	if hasFloat {
		return data.NewFloatValue(sum), nil
	}
	return data.NewIntValue(int(sum)), nil
}

func (f *ArraySumFunction) GetName() string { return "array_sum" }

func (f *ArraySumFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "array", 0, nil, nil)}
}

func (f *ArraySumFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "array", 0, data.NewBaseType("array"))}
}
