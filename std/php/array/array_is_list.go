package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayIsListFunction 实现 array_is_list 函数 (PHP 8.1+)
// array_is_list(array $array): bool
type ArrayIsListFunction struct{}

func NewArrayIsListFunction() data.FuncStmt {
	return &ArrayIsListFunction{}
}

func (f *ArrayIsListFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrValue, _ := ctx.GetIndexValue(0)
	if arrValue == nil {
		return data.NewBoolValue(false), nil
	}

	arr, ok := arrValue.(*data.ArrayValue)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	// 空数组是 list
	if arr.Len() == 0 {
		return data.NewBoolValue(true), nil
	}
	for arraySlots94, position := arr.View(), 0; position < arraySlots94.Len(); position++ {
		z := arraySlots94.At(position)
		key, ok := z.PHPArrayKey(position).(*data.IntValue)
		if !ok || key.Value != position {
			return data.NewBoolValue(false), nil
		}
	}

	return data.NewBoolValue(true), nil
}

func (f *ArrayIsListFunction) GetName() string {
	return "array_is_list"
}

var arrayIsListFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "array", 0, nil, nil),
}

func (f *ArrayIsListFunction) GetParams() []data.GetValue {
	return arrayIsListFunctionGetParams
}

var arrayIsListFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
}

func (f *ArrayIsListFunction) GetVariables() []data.Variable {
	return arrayIsListFunctionGetVariables
}
