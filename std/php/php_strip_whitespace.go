package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"strings"
)

type PHPStripWhitespaceFunction struct{}

func (*PHPStripWhitespaceFunction) GetName() string { return "php_strip_whitespace" }
func (*PHPStripWhitespaceFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "filename", 0, nil, data.String{})}
}
func (*PHPStripWhitespaceFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "filename", 0, data.String{})}
}
func (*PHPStripWhitespaceFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	filename, _ := ctx.GetIndexValue(0)
	source, err := data.ReadPHPFile(ctx.GetVM(), filename.AsString())
	if err != nil {
		return data.NewStringValue(""), nil
	}
	var out strings.Builder
	whitespace := false
	for _, token := range tokenizePhp(string(source)) {
		if fields, ok := token.(*data.ArrayValue); ok {
			values := fields.ToValueList()
			id := values[0].(*data.IntValue).Value
			switch id {
			case T_COMMENT, T_DOC_COMMENT:
				continue
			case T_WHITESPACE:
				if !whitespace {
					out.WriteByte(' ')
				}
				whitespace = true
			default:
				out.WriteString(values[1].AsString())
				whitespace = false
			}
		} else {
			out.WriteString(token.AsString())
			whitespace = false
		}
	}
	return data.NewStringValue(out.String()), nil
}
