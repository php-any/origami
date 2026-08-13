package intl

import (
	"unicode/utf8"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// GraphemeSubstrFunction 实现 grapheme_substr（intl 扩展）
// 按 UTF-8 字素（此处用 rune）截取子串，供 Symfony String 等使用
func NewGraphemeSubstrFunction() data.FuncStmt {
	return &GraphemeSubstrFunction{}
}

type GraphemeSubstrFunction struct{}

func (f *GraphemeSubstrFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	stringValue, _ := ctx.GetIndexValue(0)
	startValue, _ := ctx.GetIndexValue(1)
	lengthValue, _ := ctx.GetIndexValue(2)

	if stringValue == nil {
		return data.NewBoolValue(false), nil
	}
	str := stringValue.AsString()
	if !utf8.ValidString(str) {
		return data.NewBoolValue(false), nil
	}

	runes := []rune(str)
	runeLen := len(runes)
	start := 0
	if startValue != nil {
		if s, ok := startValue.(data.AsInt); ok {
			if n, err := s.AsInt(); err == nil {
				start = n
			}
		}
	}
	if start < 0 {
		start = runeLen + start
		if start < 0 {
			start = 0
		}
	}
	if start >= runeLen {
		return data.NewStringValue(""), nil
	}

	length := runeLen - start
	if lengthValue != nil {
		if l, ok := lengthValue.(data.AsInt); ok {
			if n, err := l.AsInt(); err == nil && n >= 0 {
				length = n
			}
		}
	}
	if start+length > runeLen {
		length = runeLen - start
	}
	return data.NewStringValue(string(runes[start : start+length])), nil
}

func (f *GraphemeSubstrFunction) GetName() string {
	return "grapheme_substr"
}

func (f *GraphemeSubstrFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "offset", 1, nil, nil),
		node.NewParameter(nil, "length", 2, nil, nil),
	}
}

func (f *GraphemeSubstrFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "offset", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "length", 2, data.NewBaseType("int")),
	}
}
