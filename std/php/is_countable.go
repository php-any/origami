package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type IsCountableFunction struct{}

func NewIsCountableFunction() data.FuncStmt {
	return &IsCountableFunction{}
}

func (f *IsCountableFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	if value == nil {
		return data.NewBoolValue(false), nil
	}

	if _, ok := value.(*data.ArrayValue); ok {
		return data.NewBoolValue(true), nil
	}
	// 关联数组在 Origami 中为 ObjectValue，PHP 中仍是 countable array
	if _, ok := value.(*data.ObjectValue); ok {
		return data.NewBoolValue(true), nil
	}

	if classVal, ok := value.(*data.ClassValue); ok {
		if _, has := classVal.GetMethod("count"); has {
			return data.NewBoolValue(true), nil
		}
	}

	return data.NewBoolValue(false), nil
}

func (f *IsCountableFunction) GetName() string {
	return "is_countable"
}

var isCountableFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "value", 0, nil, nil),
}

func (f *IsCountableFunction) GetParams() []data.GetValue {
	return isCountableFunctionGetParams
}

var isCountableFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
}

func (f *IsCountableFunction) GetVariables() []data.Variable {
	return isCountableFunctionGetVariables
}
