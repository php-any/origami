package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IsNullFunction 实现 is_null 函数
type IsNullFunction struct{}

func NewIsNullFunction() data.FuncStmt {
	return &IsNullFunction{}
}

func (f *IsNullFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	if value == nil {
		return data.NewBoolValue(true), nil
	}

	if _, ok := value.(*data.NullValue); ok {
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(false), nil
}

func (f *IsNullFunction) GetName() string {
	return "is_null"
}

func (f *IsNullFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *IsNullFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
	}
}
