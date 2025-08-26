package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewFunctionExistsFunction() data.FuncStmt {
	return &FunctionExistsFunction{}
}

type FunctionExistsFunction struct{}

func (f *FunctionExistsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, ctl := ctx.GetVariableValue(node.NewVariable(nil, "function_name", 0, data.String{}))
	if ctl != nil {
		return data.NewBoolValue(false), ctl
	}
	funcName := v.(data.AsString).AsString()
	_, exist := ctx.GetVM().GetFunc(funcName)
	return data.NewBoolValue(exist), nil
}

func (f *FunctionExistsFunction) GetName() string {
	return "function_exists"
}

func (f *FunctionExistsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "function_name", 0, nil, data.String{}),
	}
}

func (f *FunctionExistsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "function_name", 0, data.String{}),
	}
}
