package php

import (
	"unicode/utf8"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StrPadFunction 实现 str_pad(string $string, int $length, string $pad_string = \" \", int $pad_type = STR_PAD_RIGHT): string
// 这里只实现最常用的行为，满足 Symfony Console 对宽度填充的需求：
// - 多字节字符按 rune 计数，不做复杂宽度计算
// - 仅支持 STR_PAD_RIGHT/STR_PAD_LEFT/STR_PAD_BOTH 三种模式（与 PHP 常量值兼容：1/0/2）
type StrPadFunction struct{}

func NewStrPadFunction() data.FuncStmt {
	return &StrPadFunction{}
}

func (f *StrPadFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, _ := ctx.GetIndexValue(0)
	lenVal, _ := ctx.GetIndexValue(1)
	padVal, _ := ctx.GetIndexValue(2)
	typeVal, _ := ctx.GetIndexValue(3)

	if strVal == nil {
		strVal = data.NewStringValue("")
	}
	s := strVal.AsString()

	targetLen := 0
	if lenVal != nil {
		if ai, ok := lenVal.(data.AsInt); ok {
			targetLen, _ = ai.AsInt()
		}
	}
	if targetLen <= 0 {
		return data.NewStringValue(""), nil
	}

	padStr := " "
	if padVal != nil {
		padStr = padVal.AsString()
		if padStr == "" {
			padStr = " "
		}
	}

	padType := 1 // STR_PAD_RIGHT
	if typeVal != nil {
		if ai, ok := typeVal.(data.AsInt); ok {
			padType, _ = ai.AsInt()
		}
	}

	// 近似按 rune 长度计算
	runeCount := utf8.RuneCountInString(s)
	if runeCount >= targetLen {
		return data.NewStringValue(s), nil
	}

	padTotal := targetLen - runeCount

	buildPad := func(n int) string {
		if n <= 0 {
			return ""
		}
		// 按“字符数”而非字节数填充，主要适配 ASCII 使用场景
		var out []rune
		padRunes := []rune(padStr)
		if len(padRunes) == 0 {
			padRunes = []rune{' '}
		}
		for len(out) < n {
			for _, r := range padRunes {
				out = append(out, r)
				if len(out) >= n {
					break
				}
			}
		}
		return string(out)
	}

	var result string
	switch padType {
	case 0: // STR_PAD_LEFT
		left := buildPad(padTotal)
		result = left + s
	case 2: // STR_PAD_BOTH
		leftSize := padTotal / 2
		rightSize := padTotal - leftSize
		left := buildPad(leftSize)
		right := buildPad(rightSize)
		result = left + s + right
	default: // STR_PAD_RIGHT
		right := buildPad(padTotal)
		result = s + right
	}

	return data.NewStringValue(result), nil
}

func (f *StrPadFunction) GetName() string {
	return "str_pad"
}

func (f *StrPadFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, node.NewStringLiteral(nil, ""), nil),
		node.NewParameter(nil, "length", 1, node.NewIntLiteral(nil, "0"), nil),
		node.NewParameter(nil, "pad_string", 2, node.NewStringLiteral(nil, " "), nil),
		node.NewParameter(nil, "pad_type", 3, node.NewIntLiteral(nil, "1"), nil),
	}
}

func (f *StrPadFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "length", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "pad_string", 2, data.NewBaseType("string")),
		node.NewVariable(nil, "pad_type", 3, data.NewBaseType("int")),
	}
}
