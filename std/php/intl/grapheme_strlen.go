package intl

import (
	"unicode/utf8"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// GraphemeStrlenFunction 实现 grapheme_strlen（intl 扩展）
// 返回 UTF-8 字符串的“字素簇”数量；此处用 rune 数量近似，以满足 Symfony Console Helper::length 等调用
func NewGraphemeStrlenFunction() data.FuncStmt {
	return &GraphemeStrlenFunction{}
}

type GraphemeStrlenFunction struct{}

func (f *GraphemeStrlenFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewIntValue(0), nil
	}
	var str string
	if s, ok := v.(data.AsString); ok {
		str = s.AsString()
	} else {
		str = v.AsString()
	}
	if !utf8.ValidString(str) {
		return data.NewBoolValue(false), nil
	}
	return data.NewIntValue(utf8.RuneCountInString(str)), nil
}

func (f *GraphemeStrlenFunction) GetName() string {
	return "grapheme_strlen"
}

func (f *GraphemeStrlenFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
	}
}

func (f *GraphemeStrlenFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
	}
}
