package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// UsortFunction 实现 usort 函数
// usort(array &$array, callable $callback): bool
// 使用用户自定义的比较函数对数组进行排序
type UsortFunction struct{}

func NewUsortFunction() data.FuncStmt {
	return &UsortFunction{}
}

func (f *UsortFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)
	callbackValue, _ := ctx.GetIndexValue(1)

	if arrayValue == nil || callbackValue == nil {
		return data.NewBoolValue(false), nil
	}

	arrayRef, ok := arrayValue.(*data.ArrayValue)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	if arrayRef.Len() == 0 {
		return data.NewBoolValue(true), nil
	}

	// PHP sorts a duplicate so callback captures see the unsorted argument.
	sorted := data.CloneArrayValue(arrayRef)
	control := userSort(ctx, sorted, callbackValue, false)
	data.CowAssign(ctx.GetIndexZVal(0), sorted)
	if control != nil {
		return nil, control
	}
	return data.NewBoolValue(true), nil
}

func (f *UsortFunction) GetName() string {
	return "usort"
}

var usortFunctionGetParams = []data.GetValue{
	node.NewParameterReference(nil, "array", 0, nil, data.Mixed{}),
	node.NewParameter(nil, "callback", 1, nil, data.Mixed{}),
}

func (f *UsortFunction) GetParams() []data.GetValue {
	return usortFunctionGetParams
}

var usortFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.Mixed{}),
	node.NewVariable(nil, "callback", 1, data.Mixed{}),
}

func (f *UsortFunction) GetVariables() []data.Variable {
	return usortFunctionGetVariables
}
