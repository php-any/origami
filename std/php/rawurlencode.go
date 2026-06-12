package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type RawurlencodeFunction struct{}

func NewRawurlencodeFunction() data.FuncStmt {
	return &RawurlencodeFunction{}
}

func (f *RawurlencodeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	stringValue, _ := ctx.GetIndexValue(0)
	if stringValue == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(rawurlencode(stringValue.AsString())), nil
}

func rawurlencode(s string) string {
	var b strings.Builder
	b.Grow(len(s) * 3 / 2)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if isUnreserved(c) {
			b.WriteByte(c)
		} else {
			b.WriteByte('%')
			b.WriteByte(hextable[c>>4])
			b.WriteByte(hextable[c&0x0f])
		}
	}
	return b.String()
}

func isUnreserved(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') ||
		c == '-' || c == '_' || c == '.' || c == '~'
}

const hextable = "0123456789ABCDEF"

func (f *RawurlencodeFunction) GetName() string {
	return "rawurlencode"
}

func (f *RawurlencodeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
	}
}

func (f *RawurlencodeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
	}
}
