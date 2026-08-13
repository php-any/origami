package php

import (
	"strings"
	"unicode"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	MB_CASE_UPPER        = 0
	MB_CASE_LOWER        = 1
	MB_CASE_TITLE        = 2
	MB_CASE_FOLD         = 3
	MB_CASE_UPPER_SIMPLE = 4
	MB_CASE_LOWER_SIMPLE = 5
	MB_CASE_TITLE_SIMPLE = 6
	MB_CASE_FOLD_SIMPLE  = 7
)

// MbConvertCaseFunction 实现 mb_convert_case 函数
type MbConvertCaseFunction struct{}

func NewMbConvertCaseFunction() data.FuncStmt {
	return &MbConvertCaseFunction{}
}

func (f *MbConvertCaseFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	stringValue, _ := ctx.GetIndexValue(0)
	modeValue, _ := ctx.GetIndexValue(1)

	if stringValue == nil {
		return data.NewStringValue(""), nil
	}

	str := stringValue.AsString()
	mode := 0
	if modeValue != nil {
		if iv, ok := modeValue.(*data.IntValue); ok {
			mode = iv.Value
		}
	}

	var result string
	switch mode {
	case MB_CASE_UPPER, MB_CASE_UPPER_SIMPLE:
		result = strings.ToUpper(str)
	case MB_CASE_LOWER, MB_CASE_LOWER_SIMPLE:
		result = strings.ToLower(str)
	case MB_CASE_TITLE, MB_CASE_TITLE_SIMPLE:
		result = toTitleCase(str)
	case MB_CASE_FOLD, MB_CASE_FOLD_SIMPLE:
		result = caseFold(str)
	default:
		result = str
	}

	return data.NewStringValue(result), nil
}

func toTitleCase(s string) string {
	runes := []rune(s)
	prevWasLetter := false
	for i, r := range runes {
		if prevWasLetter {
			runes[i] = unicode.ToLower(r)
		} else {
			runes[i] = unicode.ToTitle(r)
		}
		prevWasLetter = unicode.IsLetter(r)
	}
	return string(runes)
}

func caseFold(s string) string {
	return strings.ToLower(s)
}

func (f *MbConvertCaseFunction) GetName() string {
	return "mb_convert_case"
}

func (f *MbConvertCaseFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "mode", 1, nil, nil),
		node.NewParameter(nil, "encoding", 2, node.NewNullLiteral(nil), nil),
	}
}

func (f *MbConvertCaseFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "mode", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "encoding", 2, data.NewNullableType(data.NewBaseType("string"))),
	}
}
