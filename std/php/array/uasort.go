package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// UasortFunction 实现 uasort 函数
// uasort(array &$array, callable $callback): bool
// 使用用户自定义比较函数排序，并保持键名关联。
type UasortFunction struct{}

func NewUasortFunction() data.FuncStmt {
	return &UasortFunction{}
}

func (f *UasortFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)
	callbackValue, _ := ctx.GetIndexValue(1)

	if arrayValue == nil || callbackValue == nil {
		return data.NewBoolValue(false), nil
	}

	arrayRef, ok := arrayValue.(*data.ArrayValue)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	if arrayRef.Len() <= 1 {
		return data.NewBoolValue(true), nil
	}

	sorted := data.CloneArrayValue(arrayRef)
	control := userSort(ctx, sorted, callbackValue, true)
	data.CowAssign(ctx.GetIndexZVal(0), sorted)
	if control != nil {
		return nil, control
	}
	return data.NewBoolValue(true), nil
}

func compareLess(v data.Value) bool {
	if iv, ok := v.(*data.IntValue); ok {
		return iv.Value < 0
	}
	if fv, ok := v.(*data.FloatValue); ok {
		return fv.Value < 0
	}
	if as, ok := v.(data.AsInt); ok {
		if i, err := as.AsInt(); err == nil {
			return i < 0
		}
	}
	return false
}

func (f *UasortFunction) GetName() string {
	return "uasort"
}

var uasortFunctionGetParams = []data.GetValue{
	node.NewParameterReference(nil, "array", 0, nil, data.Mixed{}),
	node.NewParameter(nil, "callback", 1, nil, data.Mixed{}),
}

func (f *UasortFunction) GetParams() []data.GetValue {
	return uasortFunctionGetParams
}

var uasortFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.Mixed{}),
	node.NewVariable(nil, "callback", 1, data.Mixed{}),
}

func (f *UasortFunction) GetVariables() []data.Variable {
	return uasortFunctionGetVariables
}
