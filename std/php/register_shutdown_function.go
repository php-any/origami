package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
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

	ctx.GetVM().AddShutdownCallback(cb)
	return data.NewNullValue(), nil
}

func (f *RegisterShutdownFunctionFunction) GetName() string {
	return "register_shutdown_function"
}

func (f *RegisterShutdownFunctionFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, nil, nil),
	}
}

func (f *RegisterShutdownFunctionFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, nil),
	}
}
