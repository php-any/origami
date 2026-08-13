package php

import (
	"net/url"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewUrldecodeFunction() data.FuncStmt {
	return &UrldecodeFunction{}
}

type UrldecodeFunction struct{}

func (f *UrldecodeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	stringValue, _ := ctx.GetIndexValue(0)

	if stringValue == nil {
		return data.NewStringValue(""), nil
	}

	// 转换为字符串
	str := stringValue.AsString()

	// URL 解码
	decoded, err := url.QueryUnescape(str)
	if err != nil {
		// 解码失败，返回原字符串
		return data.NewStringValue(str), nil
	}

	return data.NewStringValue(decoded), nil
}

func (f *UrldecodeFunction) GetName() string {
	return "urldecode"
}

func (f *UrldecodeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
	}
}

func (f *UrldecodeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
	}
}
