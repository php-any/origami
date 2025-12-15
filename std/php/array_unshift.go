package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayUnshiftFunction 实现 array_unshift 函数
type ArrayUnshiftFunction struct{}

func NewArrayUnshiftFunction() data.FuncStmt {
	return &ArrayUnshiftFunction{}
}

func (f *ArrayUnshiftFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)

	// Get variadic arguments
	valuesValue, _ := ctx.GetIndexValue(1)
	var values []data.Value

	if valuesValue != nil {
		if paramsArray, ok := valuesValue.(*data.ArrayValue); ok {
			values = paramsArray.Value
		}
	}

	if len(values) == 0 {
		// No values to prepend, just return count
		if arr, ok := arrayValue.(*data.ArrayValue); ok {
			return data.NewIntValue(len(arr.Value)), nil
		}
		return data.NewIntValue(0), nil
	}

	if arr, ok := arrayValue.(*data.ArrayValue); ok {
		// Prepend values
		// array_unshift prepends passed elements to the front of the array.
		// Note that the list of elements is prepended as a whole, so that the prepended elements stay in the same order.

		newArr := make([]data.Value, len(values)+len(arr.Value))
		copy(newArr, values)
		copy(newArr[len(values):], arr.Value)

		arr.Value = newArr

		return data.NewIntValue(len(arr.Value)), nil
	}

	// Warning: array_unshift() expects parameter 1 to be array
	return data.NewIntValue(0), nil
}

func (f *ArrayUnshiftFunction) GetName() string {
	return "array_unshift"
}

func (f *ArrayUnshiftFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, nil),
		node.NewParameters(nil, "values", 1, nil, nil), // Variadic
	}
}

func (f *ArrayUnshiftFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "values", 1, data.NewBaseType("mixed")), // How to mark variadic?
	}
}
