package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayKeyLastFunction 实现 array_key_last
// array_key_last(array $array): int|string|null
type ArrayKeyLastFunction struct{}

func NewArrayKeyLastFunction() data.FuncStmt {
	return &ArrayKeyLastFunction{}
}

func (f *ArrayKeyLastFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	val, _ := ctx.GetIndexValue(0)
	if val == nil {
		return data.NewNullValue(), nil
	}

	switch v := val.(type) {
	case *data.ArrayValue:
		if v.Len() == 0 {
			return data.NewNullValue(), nil
		}
		position := v.Len() - 1
		return v.At(position).PHPArrayKey(position), nil

	default:
		return data.NewNullValue(), nil
	}
}

func (f *ArrayKeyLastFunction) GetName() string { return "array_key_last" }

var arrayKeyLastFunctionGetParams = []data.GetValue{node.NewParameter(nil, "array", 0, nil, nil)}

func (f *ArrayKeyLastFunction) GetParams() []data.GetValue {
	return arrayKeyLastFunctionGetParams
}

var arrayKeyLastFunctionGetVariables = []data.Variable{node.NewVariable(nil, "array", 0, data.NewBaseType("array"))}

func (f *ArrayKeyLastFunction) GetVariables() []data.Variable {
	return arrayKeyLastFunctionGetVariables
}
