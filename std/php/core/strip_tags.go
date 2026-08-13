package core

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StripTagsFunction 实现 strip_tags 函数
type StripTagsFunction struct{}

func NewStripTagsFunction() data.FuncStmt {
	return &StripTagsFunction{}
}

func (f *StripTagsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, _ := ctx.GetIndexValue(0)
	if strVal == nil {
		return data.NewStringValue(""), nil
	}

	str := strVal.AsString()
	result := stripTags(str)

	return data.NewStringValue(result), nil
}

func (f *StripTagsFunction) GetName() string {
	return "strip_tags"
}

func (f *StripTagsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, data.String{}),
	}
}

func (f *StripTagsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.String{}),
	}
}

func stripTags(s string) string {
	var result strings.Builder
	inTag := false
	for i := 0; i < len(s); i++ {
		if s[i] == '<' {
			inTag = true
			continue
		}
		if s[i] == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteByte(s[i])
		}
	}
	return result.String()
}
