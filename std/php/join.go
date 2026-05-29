package php

import (
	"github.com/php-any/origami/data"
)

// JoinFunction 是 implode 的别名（join(string $separator, array $array)）。
type JoinFunction struct {
	inner *ImplodeFunction
}

func NewJoinFunction() data.FuncStmt {
	return &JoinFunction{inner: &ImplodeFunction{}}
}

func (f *JoinFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return f.inner.Call(ctx)
}

func (f *JoinFunction) GetName() string { return "join" }

func (f *JoinFunction) GetParams() []data.GetValue  { return f.inner.GetParams() }
func (f *JoinFunction) GetVariables() []data.Variable { return f.inner.GetVariables() }
