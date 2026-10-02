package php

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
	"strings"
)

// RegisterShutdownFunctionFunction 实现 register_shutdown_function 函数
// 注册一个在脚本执行结束时调用的回调，回调在 main 中统一执行
type RegisterShutdownFunctionFunction struct{}

func NewRegisterShutdownFunctionFunction() data.FuncStmt {
	return &RegisterShutdownFunctionFunction{}
}

func (f *RegisterShutdownFunctionFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	cb, ok := ctx.GetIndexValue(0)
	if !ok || cb == nil {
		return data.NewNullValue(), nil
	}

	if object, ok := cb.(*data.ClassValue); ok {
		cb = data.NewArrayValue([]data.Value{object, data.NewStringValue("__invoke")})
	} else if name, ok := cb.(*data.StringValue); ok {
		if class, method, found := strings.Cut(name.Value, "::"); found {
			cb = data.NewArrayValue([]data.Value{data.NewStringValue(class), data.NewStringValue(method)})
		}
	}
	check := core.NewIsCallableFunction()
	checkCtx := ctx.CreateContext(check.GetVariables())
	checkCtx.SetIndexZVal(0, data.NewZVal(cb))
	valid, ctl := check.Call(checkCtx)
	if ctl != nil {
		return nil, ctl
	}
	if b, ok := valid.(*data.BoolValue); !ok || !b.Value {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("register_shutdown_function(): Argument #1 ($callback) must be a valid callback"), "TypeError")
	}
	var args []data.Value
	if values, ok := ctx.GetIndexValue(1); ok {
		if array, ok := values.(*data.ArrayValue); ok {
			args = array.ToValueList()
		}
	}
	ctx.GetVM().AddShutdownCallback(data.NewFuncValue(&shutdownInvocation{callback: cb, args: args}))
	return data.NewNullValue(), nil
}

func (f *RegisterShutdownFunctionFunction) GetName() string {
	return "register_shutdown_function"
}

var registerShutdownFunctionFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "callback", 0, nil, nil),
	node.NewParameters(nil, "args", 1, nil, nil),
}

func (f *RegisterShutdownFunctionFunction) GetParams() []data.GetValue {
	return registerShutdownFunctionFunctionGetParams
}

var registerShutdownFunctionFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "callback", 0, nil),
	node.NewVariable(nil, "args", 1, nil),
}

func (f *RegisterShutdownFunctionFunction) GetVariables() []data.Variable {
	return registerShutdownFunctionFunctionGetVariables
}

type shutdownInvocation struct {
	callback data.Value
	args     []data.Value
}

func (*shutdownInvocation) GetName() string               { return "shutdown_callback" }
func (*shutdownInvocation) GetParams() []data.GetValue    { return nil }
func (*shutdownInvocation) GetVariables() []data.Variable { return nil }
func (f *shutdownInvocation) Call(ctx data.Context) (data.GetValue, data.Control) {
	call := core.NewCallUserFuncFunction()
	callCtx := ctx.CreateContext(call.GetVariables())
	callCtx.SetIndexZVal(0, data.NewZVal(f.callback))
	callCtx.SetIndexZVal(1, data.NewZVal(data.NewArrayValue(f.args)))
	return call.Call(callCtx)
}
