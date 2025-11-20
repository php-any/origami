package php

import (
	"encoding/base64"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewBase64EncodeFunction() data.FuncStmt {
	return &Base64EncodeFunction{}
}

type Base64EncodeFunction struct{}

func (f *Base64EncodeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	dataValue, _ := ctx.GetIndexValue(0)

	if dataValue == nil {
		return data.NewStringValue(""), nil
	}

	// 转换为字符串
	dataStr := dataValue.AsString()

	// Base64 编码
	encoded := base64.StdEncoding.EncodeToString([]byte(dataStr))

	return data.NewStringValue(encoded), nil
}

func (f *Base64EncodeFunction) GetName() string {
	return "base64_encode"
}

func (f *Base64EncodeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "data", 0, nil, nil),
	}
}

func (f *Base64EncodeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "data", 0, data.NewBaseType("string")),
	}
}
