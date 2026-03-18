package php

import (
	"github.com/php-any/origami/data"
)

// NewSetErrorHandlerFunction 创建 set_error_handler 函数。
// PHP 语义（简化版）：
//
//	set_error_handler(callable $callback, int $error_types = E_ALL): callable|null
//
// - 设置用户自定义的错误处理函数；
// - 返回之前的错误处理函数（如果有），否则返回 null。
func NewSetErrorHandlerFunction() data.FuncStmt {
	return &SetErrorHandlerFunction{}
}

type SetErrorHandlerFunction struct {
	data.Function
}

func (f *SetErrorHandlerFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 简化实现：不实际存储错误处理函数，只返回 null
	// 真实实现需要在 VM 或上下文中存储错误处理回调
	return data.NewNullValue(), nil
}

func (f *SetErrorHandlerFunction) GetName() string {
	return "set_error_handler"
}

func (f *SetErrorHandlerFunction) GetParams() []data.GetValue {
	return nil
}

func (f *SetErrorHandlerFunction) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("callback", 0, nil),
		data.NewVariable("error_types", 1, data.NewBaseType("int")),
	}
}
