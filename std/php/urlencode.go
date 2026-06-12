package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type UrlencodeFunction struct{}

func NewUrlencodeFunction() data.FuncStmt {
	return &UrlencodeFunction{}
}

func (f *UrlencodeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	stringValue, _ := ctx.GetIndexValue(0)
	if stringValue == nil {
		return data.NewStringValue(""), nil
	}
	str := stringValue.AsString()
	var b strings.Builder
	b.Grow(len(str) * 3 / 2)
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c == ' ' {
			b.WriteByte('+')
		} else if isUrlencodeUnreserved(c) {
			b.WriteByte(c)
		} else {
			b.WriteByte('%')
			b.WriteByte(hextable[c>>4])
			b.WriteByte(hextable[c&0x0f])
		}
	}
	return data.NewStringValue(b.String()), nil
}

// isUrlencodeUnreserved matches PHP urlencode: ALPHA, DIGIT, -, _, .
// Note: ~ IS encoded (%7E) in urlencode, unlike rawurlencode.
func isUrlencodeUnreserved(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') ||
		c == '-' || c == '_' || c == '.'
}

func (f *UrlencodeFunction) GetName() string {
	return "urlencode"
}

func (f *UrlencodeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
	}
}

func (f *UrlencodeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
	}
}
