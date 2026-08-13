package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SetErrorHandlerFunction 实现 set_error_handler
// set_error_handler(?callable $callback, int $error_types = E_ALL): callable|null
type SetErrorHandlerFunction struct{}

func NewSetErrorHandlerFunction() data.FuncStmt {
	return &SetErrorHandlerFunction{}
}

func (f *SetErrorHandlerFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	cb, _ := ctx.GetIndexValue(0)
	vm := ctx.GetVM()

	handlerVM, ok := vm.(interface {
		SetErrorHandler(data.Value) data.Value
	})
	if !ok {
		return data.NewNullValue(), nil
	}
	old := handlerVM.SetErrorHandler(cb)

	if old == nil {
		return data.NewNullValue(), nil
	}
	return old, nil
}

func (f *SetErrorHandlerFunction) GetName() string { return "set_error_handler" }

func (f *SetErrorHandlerFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, nil, nil),
		node.NewParameter(nil, "error_types", 1, data.NewIntValue(32767), data.Int{}),
	}
}

func (f *SetErrorHandlerFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, nil),
		node.NewVariable(nil, "error_types", 1, data.Int{}),
	}
}

// RestoreErrorHandlerFunction 实现 restore_error_handler(): bool
type RestoreErrorHandlerFunction struct{}

func NewRestoreErrorHandlerFunction() data.FuncStmt {
	return &RestoreErrorHandlerFunction{}
}

func (f *RestoreErrorHandlerFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	vm := ctx.GetVM()
	handlerVM, supported := vm.(interface {
		RestoreErrorHandler() bool
	})
	if !supported {
		return data.NewBoolValue(false), nil
	}
	ok := handlerVM.RestoreErrorHandler()
	return data.NewBoolValue(ok), nil
}

func (f *RestoreErrorHandlerFunction) GetName() string { return "restore_error_handler" }

func (f *RestoreErrorHandlerFunction) GetParams() []data.GetValue { return nil }

func (f *RestoreErrorHandlerFunction) GetVariables() []data.Variable { return nil }
