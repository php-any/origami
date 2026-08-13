package iconv

import (
	"unicode/utf8"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IconvSubstrFunction 实现 iconv_substr
// 仅支持 UTF-8 语义：按字符偏移和长度截取子串。
// 签名：iconv_substr(string $str, int $offset, ?int $length = null, ?string $encoding = null): string|false
type IconvSubstrFunction struct{}

func NewIconvSubstrFunction() data.FuncStmt {
	return &IconvSubstrFunction{}
}

func (f *IconvSubstrFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, _ := ctx.GetIndexValue(0)
	offsetVal, _ := ctx.GetIndexValue(1)
	lengthVal, _ := ctx.GetIndexValue(2)

	if strVal == nil || offsetVal == nil {
		return data.NewBoolValue(false), nil
	}

	s := strVal.AsString()
	runes := []rune(s)
	n := len(runes)

	// 解析 offset
	offset := 0
	if asInt, ok := offsetVal.(data.AsInt); ok {
		v, _ := asInt.AsInt()
		offset = v
	}
	if offset < 0 {
		offset = n + offset
	}
	if offset < 0 {
		offset = 0
	}
	if offset >= n {
		return data.NewStringValue(""), nil
	}

	// 解析 length：null 表示直到末尾
	length := n - offset
	if lengthVal != nil {
		if _, isNull := lengthVal.(*data.NullValue); !isNull {
			if asInt, ok := lengthVal.(data.AsInt); ok {
				v, _ := asInt.AsInt()
				length = v
			}
		}
	}
	if length < 0 {
		length = n + length - offset
	}
	if length <= 0 {
		return data.NewStringValue(""), nil
	}

	end := offset + length
	if end > n {
		end = n
	}

	sub := string(runes[offset:end])
	// 确保返回仍是有效 UTF-8
	if !utf8.ValidString(sub) {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(sub), nil
}

func (f *IconvSubstrFunction) GetName() string {
	return "iconv_substr"
}

func (f *IconvSubstrFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "str", 0, nil, nil),
		node.NewParameter(nil, "offset", 1, nil, nil),
		node.NewParameter(nil, "length", 2, node.NewNullLiteral(nil), nil),
		node.NewParameter(nil, "encoding", 3, node.NewNullLiteral(nil), nil),
	}
}

func (f *IconvSubstrFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "str", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "offset", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "length", 2, data.NewBaseType("int")),
		node.NewVariable(nil, "encoding", 3, data.NewBaseType("string")),
	}
}
