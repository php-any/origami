package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewStrlenFunction() data.FuncStmt {
	return &StrlenFunction{
		params: []data.GetValue{
			node.NewParameter(nil, "string", 0, nil, nil),
		},
		vars: []data.Variable{
			node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		},
	}
}

type StrlenFunction struct {
	params []data.GetValue
	vars   []data.Variable
}

func (f *StrlenFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewIntValue(0), nil
	}

	// 检查是否为 NullValue
	if _, ok := v.(*data.NullValue); ok {
		return data.NewIntValue(0), nil
	}

	// 转换为字符串
	var str string
	if strVal, ok := v.(data.AsString); ok {
		str = strVal.AsString()
	} else {
		str = v.AsString()
	}

	return data.NewIntValue(len(str)), nil
}

func (f *StrlenFunction) GetName() string {
	return "strlen"
}

func (f *StrlenFunction) GetParams() []data.GetValue {
	return f.params
}

func (f *StrlenFunction) GetVariables() []data.Variable {
	return f.vars
}
