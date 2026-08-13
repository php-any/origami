package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StrpbrkFunction 实现 strpbrk 函数
// strpbrk(string $string, string $characters): string|false
// 在 $string 中搜索 $characters 中的任意字符，返回从第一次匹配位置开始的子字符串
type StrpbrkFunction struct{}

func NewStrpbrkFunction() data.FuncStmt {
	return &StrpbrkFunction{}
}

func (f *StrpbrkFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	stringValue, _ := ctx.GetIndexValue(0)
	charactersValue, _ := ctx.GetIndexValue(1)

	if stringValue == nil || charactersValue == nil {
		return data.NewBoolValue(false), nil
	}

	str := stringValue.AsString()
	chars := charactersValue.AsString()

	if str == "" || chars == "" {
		return data.NewBoolValue(false), nil
	}

	idx := strings.IndexAny(str, chars)
	if idx == -1 {
		return data.NewBoolValue(false), nil
	}

	return data.NewStringValue(str[idx:]), nil
}

func (f *StrpbrkFunction) GetName() string {
	return "strpbrk"
}

func (f *StrpbrkFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "characters", 1, nil, nil),
	}
}

func (f *StrpbrkFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "characters", 1, data.NewBaseType("string")),
	}
}
