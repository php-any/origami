package iconv

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IconvStrrposFunction 实现 iconv_strrpos
// 仅支持 UTF-8 语义：返回最后一次出现的位置。
// 注意：为兼容 Symfony polyfill 的用法，这里第三个参数视为 encoding，而非 offset。
// 签名：iconv_strrpos(string $haystack, string $needle, ?string $encoding = null): int|false
type IconvStrrposFunction struct{}

func NewIconvStrrposFunction() data.FuncStmt {
	return &IconvStrrposFunction{}
}

func (f *IconvStrrposFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	hayVal, _ := ctx.GetIndexValue(0)
	needleVal, _ := ctx.GetIndexValue(1)

	if hayVal == nil || needleVal == nil {
		return data.NewBoolValue(false), nil
	}

	hayRunes := []rune(hayVal.AsString())
	needleRunes := []rune(needleVal.AsString())

	if len(needleRunes) == 0 {
		return data.NewBoolValue(false), nil
	}

	for i := len(hayRunes) - len(needleRunes); i >= 0; i-- {
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

func (f *IconvStrrposFunction) GetName() string {
	return "iconv_strrpos"
}

func (f *IconvStrrposFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "haystack", 0, nil, nil),
		node.NewParameter(nil, "needle", 1, nil, nil),
		// 第三个参数在 Symfony polyfill 中作为 encoding 使用，此处不做特殊处理
		node.NewParameter(nil, "encoding", 2, node.NewNullLiteral(nil), nil),
	}
}

func (f *IconvStrrposFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "haystack", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "needle", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "encoding", 2, data.NewBaseType("string")),
	}
}
