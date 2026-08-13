package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewStrtoupperFunction() data.FuncStmt {
	return &StrtoupperFunction{}
}

type StrtoupperFunction struct{}

func (f *StrtoupperFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	stringValue, _ := ctx.GetIndexValue(0)

	if stringValue == nil {
		return data.NewStringValue(""), nil
	}

	// 转换为字符串
	str := stringValue.AsString()

	// 转换为大写
	upper := strings.ToUpper(str)

	return data.NewStringValue(upper), nil
}

func (f *StrtoupperFunction) GetName() string {
	return "strtoupper"
}

func (f *StrtoupperFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
	}
}

func (f *StrtoupperFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
	}
}
