package php

import (
	"encoding/hex"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// Hex2binFunction 实现 PHP 内置函数 hex2bin
//
//	hex2bin(string $string): string|false
//
// 将十六进制字符串解码为二进制。奇数长度或非法十六进制字符时返回 false。
func NewHex2binFunction() data.FuncStmt {
	return &Hex2binFunction{}
}

type Hex2binFunction struct{}

func (f *Hex2binFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	s := v.AsString()
	if len(s)%2 != 0 {
		return data.NewBoolValue(false), nil
	}
	decoded, err := hex.DecodeString(s)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(string(decoded)), nil
}

func (f *Hex2binFunction) GetName() string {
	return "hex2bin"
}

func (f *Hex2binFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
	}
}

func (f *Hex2binFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
	}
}
