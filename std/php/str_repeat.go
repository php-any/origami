package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StrRepeatFunction 实现 PHP 内置函数 str_repeat
//
// 签名：
//
//	str_repeat(string $string, int $times): string
//
// 近似行为：
//   - $times <= 0 时返回空字符串（PHP 会发 warning 并返回空串，这里忽略 warning）
type StrRepeatFunction struct{}

func NewStrRepeatFunction() data.FuncStmt {
	return &StrRepeatFunction{}
}

func (f *StrRepeatFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, _ := ctx.GetIndexValue(0)
	timesVal, _ := ctx.GetIndexValue(1)

	if strVal == nil || timesVal == nil {
		return data.NewStringValue(""), nil
	}

	s := strVal.AsString()
	times := 0
	if asInt, ok := timesVal.(data.AsInt); ok {
		if v, err := asInt.AsInt(); err == nil {
			times = v
		}
	}
	if times <= 0 || s == "" {
		return data.NewStringValue(""), nil
	}

	return data.NewStringValue(strings.Repeat(s, times)), nil
}

func (f *StrRepeatFunction) GetName() string {
	return "str_repeat"
}

func (f *StrRepeatFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "times", 1, nil, nil),
	}
}

func (f *StrRepeatFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "times", 1, data.NewBaseType("int")),
	}
}
