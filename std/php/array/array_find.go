package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayFindKeyFunction 实现 PHP 8.4 array_find_key(array $array, callable $callback): int|string|null
// 返回第一个使 callback($value, $key) 为真的键；找不到返回 null。
type ArrayFindKeyFunction struct{}

func NewArrayFindKeyFunction() data.FuncStmt {
	return &ArrayFindKeyFunction{}
}

func (f *ArrayFindKeyFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrVal, _ := ctx.GetIndexValue(0)
	cbVal, _ := ctx.GetIndexValue(1)
	if arrVal == nil || cbVal == nil {
		return data.NewNullValue(), nil
	}

	for _, e := range toKVEntries(arrVal) {
		ret, ctl := invokeCallback(ctx, cbVal, []data.Value{e.value, e.key})
		if ctl != nil {
			return nil, ctl
		}
		if isCallbackTruthy(ret) {
			return e.key, nil
		}
	}
	return data.NewNullValue(), nil
}

func (f *ArrayFindKeyFunction) GetName() string { return "array_find_key" }

var arrayFindKeyFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "array", 0, nil, data.Arrays{}),
	node.NewParameter(nil, "callback", 1, nil, nil),
}

func (f *ArrayFindKeyFunction) GetParams() []data.GetValue {
	return arrayFindKeyFunctionGetParams
}

var arrayFindKeyFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.Arrays{}),
	node.NewVariable(nil, "callback", 1, data.Mixed{}),
}

func (f *ArrayFindKeyFunction) GetVariables() []data.Variable {
	return arrayFindKeyFunctionGetVariables
}

// ArrayFindFunction 实现 PHP 8.4 array_find(array $array, callable $callback): mixed
// 返回第一个使 callback($value, $key) 为真的值；找不到返回 null。
type ArrayFindFunction struct{}

func NewArrayFindFunction() data.FuncStmt {
	return &ArrayFindFunction{}
}

func (f *ArrayFindFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrVal, _ := ctx.GetIndexValue(0)
	cbVal, _ := ctx.GetIndexValue(1)
	if arrVal == nil || cbVal == nil {
		return data.NewNullValue(), nil
	}

	for _, e := range toKVEntries(arrVal) {
		ret, ctl := invokeCallback(ctx, cbVal, []data.Value{e.value, e.key})
		if ctl != nil {
			return nil, ctl
		}
		if isCallbackTruthy(ret) {
			return e.value, nil
		}
	}
	return data.NewNullValue(), nil
}

func (f *ArrayFindFunction) GetName() string { return "array_find" }

var arrayFindFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "array", 0, nil, data.Arrays{}),
	node.NewParameter(nil, "callback", 1, nil, nil),
}

func (f *ArrayFindFunction) GetParams() []data.GetValue {
	return arrayFindFunctionGetParams
}

var arrayFindFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.Arrays{}),
	node.NewVariable(nil, "callback", 1, data.Mixed{}),
}

func (f *ArrayFindFunction) GetVariables() []data.Variable {
	return arrayFindFunctionGetVariables
}

// ArrayAnyFunction 实现 PHP 8.4 array_any(array $array, callable $callback): bool
type ArrayAnyFunction struct{}

func NewArrayAnyFunction() data.FuncStmt { return &ArrayAnyFunction{} }

func (f *ArrayAnyFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrVal, _ := ctx.GetIndexValue(0)
	cbVal, _ := ctx.GetIndexValue(1)
	if arrVal == nil || cbVal == nil {
		return data.NewBoolValue(false), nil
	}
	for _, e := range toKVEntries(arrVal) {
		ret, ctl := invokeCallback(ctx, cbVal, []data.Value{e.value, e.key})
		if ctl != nil {
			return nil, ctl
		}
		if isCallbackTruthy(ret) {
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}

func (f *ArrayAnyFunction) GetName() string { return "array_any" }
var arrayAnyFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "array", 0, nil, data.Arrays{}),
	node.NewParameter(nil, "callback", 1, nil, nil),
}

func (f *ArrayAnyFunction) GetParams() []data.GetValue {
	return arrayAnyFunctionGetParams
}
var arrayAnyFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.Arrays{}),
	node.NewVariable(nil, "callback", 1, data.Mixed{}),
}

func (f *ArrayAnyFunction) GetVariables() []data.Variable {
	return arrayAnyFunctionGetVariables
}

// ArrayAllFunction 实现 PHP 8.4 array_all(array $array, callable $callback): bool
type ArrayAllFunction struct{}

func NewArrayAllFunction() data.FuncStmt { return &ArrayAllFunction{} }

func (f *ArrayAllFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrVal, _ := ctx.GetIndexValue(0)
	cbVal, _ := ctx.GetIndexValue(1)
	if arrVal == nil || cbVal == nil {
		return data.NewBoolValue(true), nil
	}
	for _, e := range toKVEntries(arrVal) {
		ret, ctl := invokeCallback(ctx, cbVal, []data.Value{e.value, e.key})
		if ctl != nil {
			return nil, ctl
		}
		if !isCallbackTruthy(ret) {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func (f *ArrayAllFunction) GetName() string { return "array_all" }
var arrayAllFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "array", 0, nil, data.Arrays{}),
	node.NewParameter(nil, "callback", 1, nil, nil),
}

func (f *ArrayAllFunction) GetParams() []data.GetValue {
	return arrayAllFunctionGetParams
}
var arrayAllFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.Arrays{}),
	node.NewVariable(nil, "callback", 1, data.Mixed{}),
}

func (f *ArrayAllFunction) GetVariables() []data.Variable {
	return arrayAllFunctionGetVariables
}

func isCallbackTruthy(v data.Value) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(data.AsBool); ok {
		okv, err := b.AsBool()
		return err == nil && okv
	}
	return false
}
