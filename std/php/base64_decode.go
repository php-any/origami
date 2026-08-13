package php

import (
	"encoding/base64"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewBase64DecodeFunction() data.FuncStmt {
	return &Base64DecodeFunction{}
}

type Base64DecodeFunction struct{}

func (f *Base64DecodeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	dataValue, _ := ctx.GetIndexValue(0)
	strictValue, _ := ctx.GetIndexValue(1)

	if dataValue == nil {
		return data.NewStringValue(""), nil
	}

	// 转换为字符串
	dataStr := dataValue.AsString()

	// 处理 strict 参数（PHP 中此参数已废弃，但保留兼容性）
	_ = strictValue

	// Base64 解码
	decoded, err := base64.StdEncoding.DecodeString(dataStr)
	if err != nil {
		// 解码失败，返回 false
		return data.NewBoolValue(false), nil
	}

	return data.NewStringValue(string(decoded)), nil
}

func (f *Base64DecodeFunction) GetName() string {
	return "base64_decode"
}

func (f *Base64DecodeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "data", 0, nil, nil),
		node.NewParameter(nil, "strict", 1, node.NewNullLiteral(nil), nil),
	}
}

func (f *Base64DecodeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "data", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "strict", 1, data.NewBaseType("bool")),
	}
}
