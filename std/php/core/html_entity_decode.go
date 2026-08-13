package core

import (
	"html"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// HtmlEntityDecodeFunction 实现 html_entity_decode。
type HtmlEntityDecodeFunction struct{}

func NewHtmlEntityDecodeFunction() data.FuncStmt {
	return &HtmlEntityDecodeFunction{}
}

func (f *HtmlEntityDecodeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, _ := ctx.GetIndexValue(0)
	if strVal == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(html.UnescapeString(strVal.AsString())), nil
}

func (f *HtmlEntityDecodeFunction) GetName() string {
	return "html_entity_decode"
}

func (f *HtmlEntityDecodeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, data.String{}),
		node.NewParameter(nil, "flags", 1, node.NewNullLiteral(nil), data.Mixed{}),
		node.NewParameter(nil, "encoding", 2, node.NewStringLiteral(nil, "UTF-8"), data.String{}),
	}
}

func (f *HtmlEntityDecodeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.String{}),
		node.NewVariable(nil, "flags", 1, data.Mixed{}),
		node.NewVariable(nil, "encoding", 2, data.String{}),
	}
}
