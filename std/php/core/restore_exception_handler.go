package core

import "github.com/php-any/origami/data"

// RestoreExceptionHandlerFunction 实现 restore_exception_handler 函数
//
// 签名：
//
//	restore_exception_handler(): bool
//
// 恢复上一次注册的处理器；空栈也返回 true。
type RestoreExceptionHandlerFunction struct{}

func NewRestoreExceptionHandlerFunction() data.FuncStmt {
	return &RestoreExceptionHandlerFunction{}
}

func (f *RestoreExceptionHandlerFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	vm := ctx.GetVM()

	handlerVM, ok := vm.(interface {
		RestoreExceptionHandler() bool
	})
	if !ok {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(handlerVM.RestoreExceptionHandler()), nil
}

func (f *RestoreExceptionHandlerFunction) GetName() string {
	return "restore_exception_handler"
}

var restoreExceptionHandlerFunctionGetParams = []data.GetValue{}

func (f *RestoreExceptionHandlerFunction) GetParams() []data.GetValue {
	// 无参数
	return restoreExceptionHandlerFunctionGetParams
}

var restoreExceptionHandlerFunctionGetVariables = []data.Variable{}

func (f *RestoreExceptionHandlerFunction) GetVariables() []data.Variable {
	return restoreExceptionHandlerFunctionGetVariables
}
