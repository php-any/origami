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

func (f *IsCountableFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *IsCountableFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
	}
}
