package core

import "github.com/php-any/origami/data"

// RestoreExceptionHandlerFunction 实现 restore_exception_handler 函数
//
// 签名：
//
//	restore_exception_handler(): bool
//
// 语义（与 PHP 基本一致）：
//   - 如果当前已经通过 set_exception_handler 注册了回调，则清空它并返回 true
//   - 如果当前没有注册回调，则返回 false
type RestoreExceptionHandlerFunction struct{}

func NewRestoreExceptionHandlerFunction() data.FuncStmt {
	return &RestoreExceptionHandlerFunction{}
}

func (f *RestoreExceptionHandlerFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	vm := ctx.GetVM()

	// 读取当前 handler
	current := getCurrentExceptionHandler(vm)
	if current == nil {
		// 没有可恢复的 handler
		return data.NewBoolValue(false), nil
	}

	handlerVM, ok := vm.(interface {
		SetExceptionHandler(data.Value) data.Value
	})
	if !ok {
		return data.NewBoolValue(false), nil
	}
	handlerVM.SetExceptionHandler(nil)

	return data.NewBoolValue(true), nil
}

func (f *RestoreExceptionHandlerFunction) GetName() string {
	return "restore_exception_handler"
}

func (f *RestoreExceptionHandlerFunction) GetParams() []data.GetValue {
	// 无参数
	return []data.GetValue{}
}

func (f *RestoreExceptionHandlerFunction) GetVariables() []data.Variable {
	return []data.Variable{}
}
