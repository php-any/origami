package php

import (
	"net/url"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewUrlencodeFunction() data.FuncStmt {
	return &UrlencodeFunction{}
}

type UrlencodeFunction struct{}

func (f *UrlencodeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	stringValue, _ := ctx.GetIndexValue(0)

	if stringValue == nil {
		return data.NewStringValue(""), nil
	}

	// 转换为字符串
	str := stringValue.AsString()

	// URL 编码
	encoded := url.QueryEscape(str)

	return data.NewStringValue(encoded), nil
}

func (f *UrlencodeFunction) GetName() string {
	return "urlencode"
}

func (f *UrlencodeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
	}
}

func (f *UrlencodeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
	}
}
