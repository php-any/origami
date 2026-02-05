package php

import (
	"unicode/utf8"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StrSplitFunction 实现 PHP 的 str_split
//
// str_split(string $string, int $length = 1): array
// 这里按 UTF-8 rune 进行切分，近似 PHP 行为。
type StrSplitFunction struct{}

func NewStrSplitFunction() data.FuncStmt {
	return &StrSplitFunction{}
}

func (f *StrSplitFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, _ := ctx.GetIndexValue(0)
	lengthVal, _ := ctx.GetIndexValue(1)

	if strVal == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	s := strVal.AsString()

	// 解析 length，默认 1；null 表示使用默认值
	chunkLen := 1
	if lengthVal != nil {
		if _, isNull := lengthVal.(*data.NullValue); !isNull {
			if asInt, ok := lengthVal.(data.AsInt); ok {
				if v, err := asInt.AsInt(); err == nil {
					chunkLen = v
				}
			}
		}
	}

	// 非法 length：<=0 返回 false（与 PHP 一致）
	if chunkLen <= 0 {
		return data.NewBoolValue(false), nil
	}

	// 空串：直接返回空数组
	if s == "" {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 使用 rune 级别的切分，保证多字节字符不会被截断
	runes := []rune(s)
	n := len(runes)

	result := make([]data.Value, 0)
	for i := 0; i < n; i += chunkLen {
		end := i + chunkLen
		if end > n {
			end = n
		}
		part := string(runes[i:end])
		// 理论上这里一定是合法 UTF-8；保险起见再校验一次
		if !utf8.ValidString(part) {
			return data.NewBoolValue(false), nil
		}
		result = append(result, data.NewStringValue(part))
	}

	return data.NewArrayValue(result), nil
}

func (f *StrSplitFunction) GetName() string {
	return "str_split"
}

func (f *StrSplitFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "length", 1, node.NewNullLiteral(nil), nil),
	}
}

func (f *StrSplitFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "length", 1, data.NewBaseType("int")),
	}
}
