package php

import (
	"github.com/php-any/origami/data"
)

func NewFuncNumArgsFunction() data.FuncStmt {
	return &FuncNumArgsFunction{}
}

type FuncNumArgsFunction struct{}

func (fn *FuncNumArgsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	args := ctx.GetCallArgs()
	if args == nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(len(args)), nil
}

func (fn *FuncNumArgsFunction) GetName() string { return "func_num_args" }

func (fn *FuncNumArgsFunction) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (fn *FuncNumArgsFunction) GetVariables() []data.Variable {
	return []data.Variable{}
}
