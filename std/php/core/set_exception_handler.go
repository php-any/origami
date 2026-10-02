package core

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SetExceptionHandlerFunction 实现 set_exception_handler 函数
// 签名：callable|null $callback(Throwable $exception)
type SetExceptionHandlerFunction struct{}

func NewSetExceptionHandlerFunction() data.FuncStmt {
	return &SetExceptionHandlerFunction{}
}

func (f *SetExceptionHandlerFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取第一个参数：回调
	cb, ok := ctx.GetIndexValue(0)
	if !ok || !ctx.GetIndexZVal(0).Defined {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("set_exception_handler() expects exactly 1 argument, 0 given"), "ArgumentCountError")
	}
	count := len(ctx.GetCallArgs())
	if flat := ctx.GetFlatCallArgs(); flat != nil {
		count = len(flat)
	}
	if count > 1 {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("set_exception_handler() expects exactly 1 argument, %d given", count), "ArgumentCountError")
	}
	var callable data.Value
	if _, isNull := cb.(*data.NullValue); cb == nil || isNull {
		cb = nil
	} else {
		var ctl data.Control
		callable, ctl = node.ResolveCallback(ctx, cb)
		if ctl != nil {
			return nil, ctl
		}
	}

	vm := ctx.GetVM()

	handlerVM, supported := vm.(interface {
		SetExceptionHandler(data.Value) data.Value
	})
	if !supported {
		return data.NewNullValue(), nil
	}
	var old data.Value
	if resolved, ok := vm.(interface {
		SetExceptionHandlerBinding(data.Value, data.Value) data.Value
	}); ok {
		old = resolved.SetExceptionHandlerBinding(cb, callable)
	} else {
		old = handlerVM.SetExceptionHandler(cb)
	}

	if old == nil {
		return data.NewNullValue(), nil
	}
	return old, nil
}

func (f *SetExceptionHandlerFunction) GetName() string {
	return "set_exception_handler"
}

var setExceptionHandlerFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "callback", 0, nil, data.Mixed{}),
}

func (f *SetExceptionHandlerFunction) GetParams() []data.GetValue {
	return setExceptionHandlerFunctionGetParams
}

var setExceptionHandlerFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "callback", 0, data.Mixed{}),
}

func (f *SetExceptionHandlerFunction) GetVariables() []data.Variable {
	return setExceptionHandlerFunctionGetVariables
}
