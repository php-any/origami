package php

import (
	"github.com/php-any/origami/data"
)

// ErrorClearLastFunction 实现 error_clear_last(): void
// 清除 error_get_last() 可返回的最近错误。
type ErrorClearLastFunction struct{}

func NewErrorClearLastFunction() data.FuncStmt {
	return &ErrorClearLastFunction{}
}

func (f *ErrorClearLastFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	clearLastError()
	return data.NewNullValue(), nil
}

func (f *ErrorClearLastFunction) GetName() string { return "error_clear_last" }

func (f *ErrorClearLastFunction) GetParams() []data.GetValue { return nil }

func (f *ErrorClearLastFunction) GetVariables() []data.Variable { return nil }
