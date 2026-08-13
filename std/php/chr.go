package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ChrFunction 实现 PHP 内置函数 chr
//
//	chr(int $codepoint): string
//
// 返回指定字节（0-255）对应的单字符字符串；常用于生成 ANSI 转义等。
func NewChrFunction() data.FuncStmt {
	return &ChrFunction{}
}

type ChrFunction struct{}

func (f *ChrFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewStringValue(""), nil
	}
	var code int
	if asInt, ok := v.(data.AsInt); ok {
		c, err := asInt.AsInt()
		if err != nil {
			return data.NewStringValue(""), nil
		}
		code = c
	} else {
		code = 0
	}
	if code < 0 || code > 255 {
		code = code & 0xFF
	}
	return data.NewStringValue(string(byte(code))), nil
}

func (f *ChrFunction) GetName() string {
	return "chr"
}

func (f *ChrFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "codepoint", 0, nil, nil),
	}
}

func (f *ChrFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "codepoint", 0, data.NewBaseType("int")),
	}
}
