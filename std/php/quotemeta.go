package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type QuotemetaFunction struct{}

func NewQuotemetaFunction() data.FuncStmt {
	return &QuotemetaFunction{}
}

func (f *QuotemetaFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	stringValue, _ := ctx.GetIndexValue(0)
	if stringValue == nil {
		return data.NewStringValue(""), nil
	}
	str := stringValue.AsString()
	var b strings.Builder
	b.Grow(len(str) * 2)
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c == '.' || c == '\\' || c == '+' || c == '*' || c == '?' || c == '[' || c == '^' || c == ']' || c == '(' || c == ')' || c == '$' {
			b.WriteByte('\\')
		}
		b.WriteByte(c)
	}
	return data.NewStringValue(b.String()), nil
}

func (f *QuotemetaFunction) GetName() string {
	return "quotemeta"
}

func (f *QuotemetaFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
	}
}

func (f *QuotemetaFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
	}
}
