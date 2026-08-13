package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// OrdFunction 实现 PHP 内置函数 ord
//
// 签名：
//
//	ord(string $string): int
//
// 这里简单按字节返回第一个字节的整数值（0-255）。
// 若字符串为空，则返回 0。
type OrdFunction struct{}

func NewOrdFunction() data.FuncStmt {
	return &OrdFunction{}
}

func (f *OrdFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, _ := ctx.GetIndexValue(0)
	if strVal == nil {
		return data.NewIntValue(0), nil
	}
	s := strVal.AsString()
	if len(s) == 0 {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(int(s[0])), nil
}

func (f *OrdFunction) GetName() string {
	return "ord"
}

func (f *OrdFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
	}
}

func (f *OrdFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
	}
}
