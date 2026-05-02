package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StrtokFunction 实现 strtok 函数
// strtok(string $string, string $token): string|false
type StrtokFunction struct{}

func NewStrtokFunction() data.FuncStmt {
	return &StrtokFunction{}
}

func (f *StrtokFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	stringValue, _ := ctx.GetIndexValue(0)
	tokenValue, _ := ctx.GetIndexValue(1)

	if stringValue == nil || tokenValue == nil {
		return data.NewBoolValue(false), nil
	}

	str := stringValue.AsString()
	token := tokenValue.AsString()

	if str == "" || token == "" {
		return data.NewBoolValue(false), nil
	}

	idx := strings.Index(str, token)
	if idx == -1 {
		return data.NewStringValue(str), nil
	}

	if idx == 0 {
		return data.NewStringValue(""), nil
	}

	return data.NewStringValue(str[:idx]), nil
}

func (f *StrtokFunction) GetName() string {
	return "strtok"
}

func (f *StrtokFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "token", 1, nil, nil),
	}
}

func (f *StrtokFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "token", 1, data.NewBaseType("string")),
	}
}
