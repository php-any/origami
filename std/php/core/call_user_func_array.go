package core

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type CallUserFuncArrayFunction struct{ CallUserFuncFunction }

func NewCallUserFuncArrayFunction() data.FuncStmt               { return &CallUserFuncArrayFunction{} }
func (f *CallUserFuncArrayFunction) GetName() string            { return "call_user_func_array" }
func (f *CallUserFuncArrayFunction) GetParams() []data.GetValue { return callUserFuncArrayParams }
func (f *CallUserFuncArrayFunction) GetVariables() []data.Variable {
	return callUserFuncFunctionGetVariables
}

var callUserFuncArrayParams = []data.GetValue{
	node.NewParameter(nil, "callback", 0, nil, nil),
	node.NewParameter(nil, "args", 1, nil, data.NewBaseType("array")),
}

func (f *CallUserFuncArrayFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	callback, _ := ctx.GetIndexValue(0)
	arguments, _ := ctx.GetIndexValue(1)
	array, ok := arguments.(*data.ArrayValue)
	if !ok {
		return nil, data.NewTypeError(nil, errors.New("call_user_func_array(): args must be an array"))
	}
	function, ctl := f.resolveCallback(ctx, callback)
	if ctl != nil {
		return nil, ctl
	}
	if function == nil {
		return data.NewBoolValue(false), nil
	}
	var callable data.GetValue = function
	if bound, ok := callback.(*data.BoundFuncValue); ok {
		callable = bound
	}
	return node.NewCallMethod(nil, callable, []data.GetValue{node.NewSpreadArgument(nil, array)}).GetValue(ctx)
}
