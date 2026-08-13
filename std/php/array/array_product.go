package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayProductFunction 实现 array_product
// array_product(array $array): int|float
type ArrayProductFunction struct{}

func NewArrayProductFunction() data.FuncStmt {
	return &ArrayProductFunction{}
}

func (f *ArrayProductFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayVal, _ := ctx.GetIndexValue(0)
	entries := toKVEntries(arrayVal)
	if len(entries) == 0 {
		return data.NewIntValue(1), nil
	}
	product := 1.0
	hasFloat := false
	for _, e := range entries {
		if _, ok := e.value.(*data.FloatValue); ok {
			hasFloat = true
		}
		product *= toFloatValue(e.value)
	}
	if hasFloat {
		return data.NewFloatValue(product), nil
	}
	return data.NewIntValue(int(product)), nil
}

func (f *ArrayProductFunction) GetName() string { return "array_product" }

func (f *ArrayProductFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "array", 0, nil, nil)}
}

func (f *ArrayProductFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "array", 0, data.NewBaseType("array"))}
}
