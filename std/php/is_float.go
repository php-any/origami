package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IsFloatFunction 实现 is_float 函数
type IsFloatFunction struct{}

func NewIsFloatFunction() data.FuncStmt {
	return &IsFloatFunction{}
}

func (f *IsFloatFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	if value == nil {
		return data.NewBoolValue(false), nil
	}

	if _, ok := value.(*data.FloatValue); ok {
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(false), nil
}

func (f *IsFloatFunction) GetName() string {
	return "is_float"
}

var isFloatFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "value", 0, nil, nil),
}

func (f *IsFloatFunction) GetParams() []data.GetValue {
	return isFloatFunctionGetParams
}

var isFloatFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
}

func (f *IsFloatFunction) GetVariables() []data.Variable {
	return isFloatFunctionGetVariables
}
