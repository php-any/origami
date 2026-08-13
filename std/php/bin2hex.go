package php

import (
	"encoding/hex"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// Bin2hexFunction 实现 PHP 内置函数 bin2hex
//
//	bin2hex(string $string): string
//
// 将二进制字符串转为小写十六进制表示。
func NewBin2hexFunction() data.FuncStmt {
	return &Bin2hexFunction{}
}

type Bin2hexFunction struct{}

func (f *Bin2hexFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(hex.EncodeToString([]byte(v.AsString()))), nil
}

func (f *Bin2hexFunction) GetName() string {
	return "bin2hex"
}

func (f *Bin2hexFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
	}
}

func (f *Bin2hexFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
	}
}
