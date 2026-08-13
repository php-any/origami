package php

import (
	"strings"
	"unicode"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewUcwordsFunction() data.FuncStmt {
	return &UcwordsFunction{}
}

type UcwordsFunction struct{}

func (f *UcwordsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	stringValue, _ := ctx.GetIndexValue(0)
	delimitersValue, _ := ctx.GetIndexValue(1)

	if stringValue == nil {
		return data.NewStringValue(""), nil
	}

	// 转换为字符串
	str := stringValue.AsString()

	if str == "" {
		return data.NewStringValue(""), nil
	}

	// 处理分隔符参数
	delimiters := " \t\r\n\f\v"
	if delimitersValue != nil {
		if _, ok := delimitersValue.(*data.NullValue); !ok {
			delimiters = delimitersValue.AsString()
		}
	}

	// 将每个单词的首字母转换为大写
	runes := []rune(str)
	prevWasDelimiter := true

	for i, r := range runes {
		isDelimiter := strings.ContainsRune(delimiters, r)
		if prevWasDelimiter && !isDelimiter {
			// 单词首字母，转换为大写
			runes[i] = unicode.ToUpper(r)
		}
		prevWasDelimiter = isDelimiter
	}

	return data.NewStringValue(string(runes)), nil
}

func (f *UcwordsFunction) GetName() string {
	return "ucwords"
}

func (f *UcwordsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "delimiters", 1, node.NewNullLiteral(nil), nil),
	}
}

func (f *UcwordsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "delimiters", 1, data.NewBaseType("string")),
	}
}
