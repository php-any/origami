package php

import (
	"net/url"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RawurldecodeFunction 实现 rawurldecode 函数
//
// 语义参考 PHP:
//
//	rawurldecode(string $string): string
//
// 使用 Go 的 url.PathUnescape 来还原百分号编码（不会把 '+' 还原成空格）。
type RawurldecodeFunction struct{}

func NewRawurldecodeFunction() data.FuncStmt {
	return &RawurldecodeFunction{}
}

func (f *RawurldecodeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	stringValue, _ := ctx.GetIndexValue(0)

	if stringValue == nil {
		return data.NewStringValue(""), nil
	}

	str := stringValue.AsString()

	decoded, err := url.PathUnescape(str)
	if err != nil {
		// 解码失败时，保持与 urldecode 类似的宽松行为：返回原字符串
		return data.NewStringValue(str), nil
	}

	return data.NewStringValue(decoded), nil
}

func (f *RawurldecodeFunction) GetName() string {
	return "rawurldecode"
}

func (f *RawurldecodeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
	}
}

func (f *RawurldecodeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
	}
}
