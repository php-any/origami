package iconv

import (
	"unicode/utf8"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IconvStrlenFunction 实现 iconv_strlen
// 在 UTF-8 环境中等价于按字符数统计长度。
// 签名：iconv_strlen(string $str, ?string $encoding = null): int|false
type IconvStrlenFunction struct{}

func NewIconvStrlenFunction() data.FuncStmt {
	return &IconvStrlenFunction{}
}

func (f *IconvStrlenFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, _ := ctx.GetIndexValue(0)
	if strVal == nil {
		return data.NewBoolValue(false), nil
	}
	s := strVal.AsString()

	// 这里直接按 UTF-8 rune 计数，忽略编码参数
	l := utf8.RuneCountInString(s)
	return data.NewIntValue(l), nil
}

func (f *IconvStrlenFunction) GetName() string {
	return "iconv_strlen"
}

func (f *IconvStrlenFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "str", 0, nil, nil),
		// encoding 可选，这里不强制类型
		node.NewParameter(nil, "encoding", 1, node.NewNullLiteral(nil), nil),
	}
}

func (f *IconvStrlenFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "str", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "encoding", 1, data.NewBaseType("string")),
	}
}
