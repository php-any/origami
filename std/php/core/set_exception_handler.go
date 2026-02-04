package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
)

// SetExceptionHandlerFunction 实现 set_exception_handler 函数
// 目前主要支持闭包/匿名函数形式的回调，签名：callable $callback(Throwable $exception)
type SetExceptionHandlerFunction struct{}

func NewSetExceptionHandlerFunction() data.FuncStmt {
	return &SetExceptionHandlerFunction{}
}

func (f *SetExceptionHandlerFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取第一个参数：回调
	cb, ok := ctx.GetIndexValue(0)
	if !ok {
		// 未传入回调时，直接返回当前回调（与 PHP 行为略有差异，但更安全）
		if old := getCurrentExceptionHandler(ctx.GetVM()); old != nil {
			return old, nil
		}
		return data.NewNullValue(), nil
	}

	vm := ctx.GetVM()

	// 根据 VM 实现类型设置异常处理回调
	var old data.Value
	switch v := vm.(type) {
	case *runtime.VM:
		old = v.SetExceptionHandler(cb)
	case *runtime.TempVM:
		old = v.SetExceptionHandler(cb)
	default:
		// 其他 VM 实现（如 LspVM）忽略设置，仅返回 null
		return data.NewNullValue(), nil
	}

	if old == nil {
		return data.NewNullValue(), nil
	}
	return old, nil
}

func (f *SetExceptionHandlerFunction) GetName() string {
	return "set_exception_handler"
}

func (f *SetExceptionHandlerFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, nil, data.Mixed{}),
	}
}

func (f *SetExceptionHandlerFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, data.Mixed{}),
	}
}

// getCurrentExceptionHandler 辅助函数，从 VM 中获取当前异常处理回调
func getCurrentExceptionHandler(vm data.VM) data.Value {
	switch v := vm.(type) {
	case *runtime.VM:
		return v.GetExceptionHandler()
	case *runtime.TempVM:
		return v.GetExceptionHandler()
	default:
		return nil
	}
}

