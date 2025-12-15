package php

import (
	"strconv"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IsNumericFunction 实现 is_numeric 函数
type IsNumericFunction struct{}

func NewIsNumericFunction() data.FuncStmt {
	return &IsNumericFunction{}
}

func (f *IsNumericFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	if value == nil {
		return data.NewBoolValue(false), nil
	}

	if _, ok := value.(*data.IntValue); ok {
		return data.NewBoolValue(true), nil
	}
	if _, ok := value.(*data.FloatValue); ok {
		return data.NewBoolValue(true), nil
	}
	if str, ok := value.(*data.StringValue); ok {
		// Check if string is numeric
		s := str.AsString()
		if _, err := strconv.ParseFloat(s, 64); err == nil {
			return data.NewBoolValue(true), nil
		}
	}

	return data.NewBoolValue(false), nil
}

func (f *IsNumericFunction) GetName() string {
	return "is_numeric"
}

func (f *IsNumericFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *IsNumericFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
	}
}
