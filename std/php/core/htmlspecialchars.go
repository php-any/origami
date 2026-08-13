package core

import (
	"html"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// HtmlspecialcharsFunction 实现 htmlspecialchars 函数
type HtmlspecialcharsFunction struct{}

func NewHtmlspecialcharsFunction() data.FuncStmt {
	return &HtmlspecialcharsFunction{}
}

func (f *HtmlspecialcharsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, _ := ctx.GetIndexValue(0)
	if strVal == nil {
		return data.NewStringValue(""), nil
	}

	str := strVal.AsString()
	result := html.EscapeString(str)

	return data.NewStringValue(result), nil
}

func (f *HtmlspecialcharsFunction) GetName() string {
	return "htmlspecialchars"
}

func (f *HtmlspecialcharsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, data.String{}),
	}
}

func (f *HtmlspecialcharsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.String{}),
	}
}
