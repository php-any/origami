package iconv

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IconvStrposFunction 实现 iconv_strpos
// 仅支持 UTF-8 语义：按字符偏移返回第一次出现的位置。
// 签名：iconv_strpos(string $haystack, string $needle, int $offset = 0, ?string $encoding = null): int|false
type IconvStrposFunction struct{}

func NewIconvStrposFunction() data.FuncStmt {
	return &IconvStrposFunction{}
}

func (f *IconvStrposFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	hayVal, _ := ctx.GetIndexValue(0)
	needleVal, _ := ctx.GetIndexValue(1)
	offsetVal, _ := ctx.GetIndexValue(2)

	if hayVal == nil || needleVal == nil {
		return data.NewBoolValue(false), nil
	}

	hayRunes := []rune(hayVal.AsString())
	needleRunes := []rune(needleVal.AsString())

	if len(needleRunes) == 0 {
		return data.NewBoolValue(false), nil
	}

	// 解析 offset
	offset := 0
	if offsetVal != nil {
		if asInt, ok := offsetVal.(data.AsInt); ok {
			v, _ := asInt.AsInt()
			offset = v
		}
	}
	if offset < 0 || offset >= len(hayRunes) {
		return data.NewBoolValue(false), nil
	}

	for i := offset; i+len(needleRunes) <= len(hayRunes); i++ {
		match := true
		for j := 0; j < len(needleRunes); j++ {
			if hayRunes[i+j] != needleRunes[j] {
				match = false
				break
			}
		}
		if match {
			return data.NewIntValue(i), nil
		}
	}

	return data.NewBoolValue(false), nil
}

func (f *IconvStrposFunction) GetName() string {
	return "iconv_strpos"
}

func (f *IconvStrposFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "haystack", 0, nil, nil),
		node.NewParameter(nil, "needle", 1, nil, nil),
		node.NewParameter(nil, "offset", 2, node.NewNullLiteral(nil), nil),
		node.NewParameter(nil, "encoding", 3, node.NewNullLiteral(nil), nil),
	}
}

func (f *IconvStrposFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "haystack", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "needle", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "offset", 2, data.NewBaseType("int")),
		node.NewVariable(nil, "encoding", 3, data.NewBaseType("string")),
	}
}
