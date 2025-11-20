package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewStrtolowerFunction() data.FuncStmt {
	return &StrtolowerFunction{}
}

type StrtolowerFunction struct{}

func (f *StrtolowerFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	stringValue, _ := ctx.GetIndexValue(0)

	if stringValue == nil {
		return data.NewStringValue(""), nil
	}

	// 转换为字符串
	str := stringValue.AsString()

	// 转换为小写
	lower := strings.ToLower(str)

	return data.NewStringValue(lower), nil
}

func (f *StrtolowerFunction) GetName() string {
	return "strtolower"
}

func (f *StrtolowerFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
	}
}

func (f *StrtolowerFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
	}
}
