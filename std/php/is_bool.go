package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IsBoolFunction 实现 is_bool 函数
type IsBoolFunction struct{}

func NewIsBoolFunction() data.FuncStmt {
	return &IsBoolFunction{}
}

func (f *IsBoolFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	if value == nil {
		return data.NewBoolValue(false), nil
	}

	if _, ok := value.(*data.BoolValue); ok {
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(false), nil
}

func (f *IsBoolFunction) GetName() string {
	return "is_bool"
}

func (f *IsBoolFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *IsBoolFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
	}
}
