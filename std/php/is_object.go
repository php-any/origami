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

	if _, ok := value.(*data.ObjectValue); ok {
		return data.NewBoolValue(true), nil
	}
	if _, ok := value.(*data.ClassValue); ok {
		return data.NewBoolValue(true), nil
	}
	// $this 也是对象：PHP 中 is_object($this) 为 true。
	if _, ok := value.(*data.ThisValue); ok {
		return data.NewBoolValue(true), nil
	}
	if _, ok := value.(*data.FuncValue); ok {
		return data.NewBoolValue(true), nil
	}
	if _, ok := value.(*data.ThrowValue); ok {
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(false), nil
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
