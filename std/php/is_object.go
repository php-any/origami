package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IsObjectFunction 实现 is_object 函数
type IsObjectFunction struct{}

func NewIsObjectFunction() data.FuncStmt {
	return &IsObjectFunction{}
}

func (f *IsObjectFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	if value == nil {
		return data.NewBoolValue(false), nil
	}

	return data.NewBoolValue(data.ValueKindOf(value) == data.ValueObject), nil
}

func (f *IsObjectFunction) GetName() string {
	return "is_object"
}

var isObjectFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "value", 0, nil, nil),
}

func (f *IsObjectFunction) GetParams() []data.GetValue {
	return isObjectFunctionGetParams
}

var isObjectFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
}

func (f *IsObjectFunction) GetVariables() []data.Variable {
	return isObjectFunctionGetVariables
}
