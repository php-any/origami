package php

import (
	"unicode"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewLcfirstFunction() data.FuncStmt {
	return &LcfirstFunction{}
}

type LcfirstFunction struct{}

func (f *LcfirstFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	stringValue, _ := ctx.GetIndexValue(0)

	if stringValue == nil {
		return data.NewStringValue(""), nil
	}

	// 转换为字符串
	str := stringValue.AsString()

	if str == "" {
		return data.NewStringValue(""), nil
	}

	// 将首字母转换为小写
	runes := []rune(str)
	if len(runes) > 0 {
		runes[0] = unicode.ToLower(runes[0])
	}

	return data.NewStringValue(string(runes)), nil
}

func (f *LcfirstFunction) GetName() string {
	return "lcfirst"
}

var lcfirstFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "string", 0, nil, nil),
}

func (f *LcfirstFunction) GetParams() []data.GetValue {
	return lcfirstFunctionGetParams
}

var lcfirstFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
}

func (f *LcfirstFunction) GetVariables() []data.Variable {
	return lcfirstFunctionGetVariables
}
