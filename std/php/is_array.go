package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IsArrayFunction 实现 is_array 函数
type IsArrayFunction struct{}

func NewIsArrayFunction() data.FuncStmt {
	return &IsArrayFunction{}
}

func (f *IsArrayFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	if value == nil {
		return data.NewBoolValue(false), nil
	}

	if _, ok := value.(*data.ArrayValue); ok {
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(false), nil
}

func (f *IsArrayFunction) GetName() string {
	return "is_array"
}

func (f *IsArrayFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *IsArrayFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
	}
}
