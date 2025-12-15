package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IsStringFunction 实现 is_string 函数
type IsStringFunction struct{}

func NewIsStringFunction() data.FuncStmt {
	return &IsStringFunction{}
}

func (f *IsStringFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	if value == nil {
		return data.NewBoolValue(false), nil
	}

	if _, ok := value.(*data.StringValue); ok {
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(false), nil
}

func (f *IsStringFunction) GetName() string {
	return "is_string"
}

func (f *IsStringFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *IsStringFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
	}
}
