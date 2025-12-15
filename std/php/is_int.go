package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IsIntFunction 实现 is_int 函数
type IsIntFunction struct{}

func NewIsIntFunction() data.FuncStmt {
	return &IsIntFunction{}
}

func (f *IsIntFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	if value == nil {
		return data.NewBoolValue(false), nil
	}

	if _, ok := value.(*data.IntValue); ok {
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(false), nil
}

func (f *IsIntFunction) GetName() string {
	return "is_int"
}

func (f *IsIntFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *IsIntFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
	}
}
