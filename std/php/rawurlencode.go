package php

import (
	"net/url"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RawurlencodeFunction 实现 rawurlencode 函数
//
// 语义参考 PHP:
//
//	rawurlencode(string $string): string
//
// 使用 Go 的 url.PathEscape，与 PHP 的 rawurlencode 行为接近：
//   - 空格编码为 %20，而不是 '+'
type RawurlencodeFunction struct{}

func NewRawurlencodeFunction() data.FuncStmt {
	return &RawurlencodeFunction{}
}

func (f *RawurlencodeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	stringValue, _ := ctx.GetIndexValue(0)

	if stringValue == nil {
		return data.NewStringValue(""), nil
	}

	str := stringValue.AsString()
	encoded := url.PathEscape(str)

	return data.NewStringValue(encoded), nil
}

func (f *RawurlencodeFunction) GetName() string {
	return "rawurlencode"
}

func (f *RawurlencodeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
	}
}

func (f *RawurlencodeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
	}
}
