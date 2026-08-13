package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IsScalarFunction 实现 is_scalar（int/float/string/bool 为 true）
type IsScalarFunction struct{}

func NewIsScalarFunction() data.FuncStmt {
	return &IsScalarFunction{}
}

func (f *IsScalarFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	if value == nil {
		return data.NewBoolValue(false), nil
	}
	switch value.(type) {
	case *data.IntValue, *data.FloatValue, *data.StringValue, *data.BoolValue:
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(false), nil
}

func (f *IsScalarFunction) GetName() string {
	return "is_scalar"
}

func (f *IsScalarFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *IsScalarFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
	}
}
