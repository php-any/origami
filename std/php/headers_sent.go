package php

import (
	"github.com/php-any/origami/data"
)

// HeadersSentFunction 实现 headers_sent 函数
type HeadersSentFunction struct{}

func NewHeadersSentFunction() data.FuncStmt { return &HeadersSentFunction{} }

func (f *HeadersSentFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(false), nil
}

func (f *HeadersSentFunction) GetName() string               { return "headers_sent" }
func (f *HeadersSentFunction) GetParams() []data.GetValue    { return nil }
func (f *HeadersSentFunction) GetVariables() []data.Variable { return nil }
