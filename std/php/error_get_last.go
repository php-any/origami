package php

import (
	"github.com/php-any/origami/data"
)

// ErrorGetLastFunction 实现 error_get_last(): ?array
type ErrorGetLastFunction struct{}

func NewErrorGetLastFunction() data.FuncStmt {
	return &ErrorGetLastFunction{}
}

func (f *ErrorGetLastFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	info := getLastError()
	if info == nil {
		return data.NewNullValue(), nil
	}
	return lastErrorToArray(info), nil
}

func (f *ErrorGetLastFunction) GetName() string { return "error_get_last" }

func (f *ErrorGetLastFunction) GetParams() []data.GetValue { return nil }

func (f *ErrorGetLastFunction) GetVariables() []data.Variable { return nil }
