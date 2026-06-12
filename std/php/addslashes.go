package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type AddslashesFunction struct{}

func NewAddslashesFunction() data.FuncStmt {
	return &AddslashesFunction{}
}

func (f *AddslashesFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	stringValue, _ := ctx.GetIndexValue(0)
	if stringValue == nil {
		return data.NewStringValue(""), nil
	}
	str := stringValue.AsString()
	var b strings.Builder
	b.Grow(len(str) + 4)
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c == '\'' || c == '"' || c == '\\' {
			b.WriteByte('\\')
		}
		b.WriteByte(c)
	}
	return data.NewStringValue(b.String()), nil
}

func (f *AddslashesFunction) GetName() string {
	return "addslashes"
}

func (f *AddslashesFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
	}
}

func (f *AddslashesFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
	}
}
