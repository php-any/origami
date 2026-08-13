package php

import (
	"math/rand"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StrWordCountFunction 实现 str_word_count 函数
// str_word_count(string $string, int $format = 0, ?string $characters = null): array|int
type StrWordCountFunction struct{}

func NewStrWordCountFunction() data.FuncStmt { return &StrWordCountFunction{} }

func (f *StrWordCountFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strV, _ := ctx.GetIndexValue(0)
	formatV, _ := ctx.GetIndexValue(1)
	charsV, _ := ctx.GetIndexValue(2)

	if strV == nil {
		return data.NewIntValue(0), nil
	}
	str := strV.AsString()

	format := 0
	if formatV != nil {
		if iv, ok := formatV.(*data.IntValue); ok {
			format = iv.Value
		}
	}

	extraChars := ""
	if charsV != nil {
		extraChars = charsV.AsString()
	}

	// 分割单词
	isWordChar := func(c rune) bool {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
			return true
		}
		return strings.ContainsRune(extraChars, c)
	}

	words := make([]string, 0)
	var current strings.Builder
	for _, c := range str {
		if isWordChar(c) {
			current.WriteRune(c)
		} else {
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
		}
	}
	if current.Len() > 0 {
		words = append(words, current.String())
	}

	switch format {
	case 1: // 返回单词起始位置数组（格式："pos" => "word"）
		list := make([]*data.ZVal, 0, len(words))
		pos := 0
		for _, w := range words {
			idx := strings.Index(str[pos:], w)
			if idx >= 0 {
				pos += idx
				list = append(list, data.NewNamedZVal(strconvItoa(pos), data.NewStringValue(w)))
				pos += len(w)
			}
		}
		return &data.ArrayValue{List: list}, nil
	case 2: // 返回位置为键、单词为值的数组
		list := make([]*data.ZVal, 0, len(words))
		pos := 0
		for _, w := range words {
			idx := strings.Index(str[pos:], w)
			if idx >= 0 {
				pos += idx
				list = append(list, data.NewNamedZVal(data.IntArrayKeyName(pos), data.NewStringValue(w)))
				pos += len(w)
			}
		}
		return &data.ArrayValue{List: list}, nil
	default:
		return data.NewIntValue(len(words)), nil
	}
}

func strconvItoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

func (f *StrWordCountFunction) GetName() string { return "str_word_count" }
func (f *StrWordCountFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "format", 1, node.NewIntLiteral(nil, "0"), data.NewBaseType("int")),
		node.NewParameter(nil, "characters", 2, node.NewNullLiteral(nil), nil),
	}
}
func (f *StrWordCountFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, nil),
		node.NewVariable(nil, "format", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "characters", 2, nil),
	}
}

// WordwrapFunction 实现 wordwrap 函数
// wordwrap(string $string, int $width = 75, string $break = "\n", bool $cut_long_words = false): string
type WordwrapFunction struct{}

func NewWordwrapFunction() data.FuncStmt { return &WordwrapFunction{} }

func (f *WordwrapFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strV, _ := ctx.GetIndexValue(0)
	widthV, _ := ctx.GetIndexValue(1)
	breakV, _ := ctx.GetIndexValue(2)
	cutV, _ := ctx.GetIndexValue(3)

	if strV == nil {
		return data.NewStringValue(""), nil
	}
	str := strV.AsString()

	width := 75
	if widthV != nil {
		if iv, ok := widthV.(*data.IntValue); ok {
			width = iv.Value
		}
	}

	brk := "\n"
	if breakV != nil {
		brk = breakV.AsString()
	}

	cutLong := false
	if cutV != nil {
		if bv, ok := cutV.(*data.BoolValue); ok {
			cutLong = bv.Value
		}
	}

	if width <= 0 {
		width = 1
	}

	var result strings.Builder
	runes := []rune(str)
	lineLen := 0
	lastSpace := -1
	lastSpaceResultLen := -1

	for i, r := range runes {
		if r == '\n' {
			result.WriteRune(r)
			lineLen = 0
			lastSpace = -1
			lastSpaceResultLen = -1
			continue
		}

		result.WriteRune(r)
		lineLen++

		if r == ' ' || r == '\t' {
			lastSpace = i
			lastSpaceResultLen = result.Len()
		}

		if lineLen >= width {
			if lastSpace >= 0 && !cutLong {
				// 在最后一个空格处换行
				s := result.String()
				lastSpacePos := strings.LastIndex(s[:lastSpaceResultLen], " ")
				if lastSpacePos < 0 {
					lastSpacePos = strings.LastIndex(s[:lastSpaceResultLen], "\t")
				}
				if lastSpacePos >= 0 {
					prefix := s[:lastSpacePos]
					suffix := s[lastSpacePos:]
					// 去除 suffix 开头的空白
					suffix = strings.TrimLeft(suffix, " \t")
					result.Reset()
					result.WriteString(prefix)
					result.WriteString(brk)
					result.WriteString(suffix)
					lineLen = len([]rune(suffix))
					lastSpace = -1
					lastSpaceResultLen = -1
					continue
				}
			}
			// 强制断行
			result.WriteString(brk)
			lineLen = 0
			lastSpace = -1
			lastSpaceResultLen = -1
		}
	}

	return data.NewStringValue(result.String()), nil
}

func (f *WordwrapFunction) GetName() string { return "wordwrap" }
func (f *WordwrapFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "width", 1, node.NewIntLiteral(nil, "75"), data.NewBaseType("int")),
		node.NewParameter(nil, "break", 2, node.NewStringLiteral(nil, `"\n"`), nil),
		node.NewParameter(nil, "cut_long_words", 3, node.NewBooleanLiteral(nil, false), data.NewBaseType("bool")),
	}
}
func (f *WordwrapFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, nil),
		node.NewVariable(nil, "width", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "break", 2, nil),
		node.NewVariable(nil, "cut_long_words", 3, data.NewBaseType("bool")),
	}
}

// StrShuffleFunction 实现 str_shuffle 函数
// str_shuffle(string $string): string
type StrShuffleFunction struct{}

func NewStrShuffleFunction() data.FuncStmt { return &StrShuffleFunction{} }

func (f *StrShuffleFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strV, _ := ctx.GetIndexValue(0)
	if strV == nil {
		return data.NewStringValue(""), nil
	}
	str := strV.AsString()
	runes := []rune(str)
	rand.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})
	return data.NewStringValue(string(runes)), nil
}

func (f *StrShuffleFunction) GetName() string { return "str_shuffle" }
func (f *StrShuffleFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "string", 0, nil, nil)}
}
func (f *StrShuffleFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "string", 0, nil)}
}

// StrRot13Function 实现 str_rot13 函数
// str_rot13(string $string): string
type StrRot13Function struct{}

func NewStrRot13Function() data.FuncStmt { return &StrRot13Function{} }

func (f *StrRot13Function) Call(ctx data.Context) (data.GetValue, data.Control) {
	strV, _ := ctx.GetIndexValue(0)
	if strV == nil {
		return data.NewStringValue(""), nil
	}
	str := strV.AsString()
	var result strings.Builder
	for _, c := range str {
		switch {
		case c >= 'a' && c <= 'z':
			result.WriteRune('a' + (c-'a'+13)%26)
		case c >= 'A' && c <= 'Z':
			result.WriteRune('A' + (c-'A'+13)%26)
		default:
			result.WriteRune(c)
		}
	}
	return data.NewStringValue(result.String()), nil
}

func (f *StrRot13Function) GetName() string { return "str_rot13" }
func (f *StrRot13Function) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "string", 0, nil, nil)}
}
func (f *StrRot13Function) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "string", 0, nil)}
}

// AddslashesFunction 实现 addslashes 函数
// addslashes(string $string): string
type AddslashesFunction struct{}

func NewAddslashesFunction() data.FuncStmt { return &AddslashesFunction{} }

func (f *AddslashesFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strV, _ := ctx.GetIndexValue(0)
	if strV == nil {
		return data.NewStringValue(""), nil
	}
	str := strV.AsString()
	var result strings.Builder
	for _, c := range str {
		switch c {
		case '\'', '"', '\\':
			result.WriteByte('\\')
			result.WriteRune(c)
		case '\x00':
			result.WriteString("\\0")
		default:
			result.WriteRune(c)
		}
	}
	return data.NewStringValue(result.String()), nil
}

func (f *AddslashesFunction) GetName() string { return "addslashes" }
func (f *AddslashesFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "string", 0, nil, nil)}
}
func (f *AddslashesFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "string", 0, nil)}
}

// ChopFunction 实现 chop 函数（rtrim 的别名）
// chop(string $string, string $characters = " \n\r\t\v\x00"): string
type ChopFunction struct{}

func NewChopFunction() data.FuncStmt { return &ChopFunction{} }

func (f *ChopFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strV, _ := ctx.GetIndexValue(0)
	charsV, _ := ctx.GetIndexValue(1)

	if strV == nil {
		return data.NewStringValue(""), nil
	}
	str := strV.AsString()
	if charsV == nil {
		return data.NewStringValue(strings.TrimRight(str, " \n\r\t\v\x00")), nil
	}
	// 检查是否为 NullValue
	if _, isNull := charsV.(*data.NullValue); isNull {
		return data.NewStringValue(strings.TrimRight(str, " \n\r\t\v\x00")), nil
	}
	return data.NewStringValue(strings.TrimRight(str, charsV.AsString())), nil
}

func (f *ChopFunction) GetName() string { return "chop" }
func (f *ChopFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "characters", 1, node.NewNullLiteral(nil), nil),
	}
}
func (f *ChopFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, nil),
		node.NewVariable(nil, "characters", 1, nil),
	}
}

// QuotemetaFunction 实现 quotemeta 函数
// quotemeta(string $string): string
type QuotemetaFunction struct{}

func NewQuotemetaFunction() data.FuncStmt { return &QuotemetaFunction{} }

func (f *QuotemetaFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strV, _ := ctx.GetIndexValue(0)
	if strV == nil {
		return data.NewStringValue(""), nil
	}
	str := strV.AsString()
	var result strings.Builder
	metaChars := ".\\+*?[^]$(){}=!<>|:-"
	for _, c := range str {
		if strings.ContainsRune(metaChars, c) {
			result.WriteByte('\\')
		}
		result.WriteRune(c)
	}
	return data.NewStringValue(result.String()), nil
}

func (f *QuotemetaFunction) GetName() string { return "quotemeta" }
func (f *QuotemetaFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "string", 0, nil, nil)}
}
func (f *QuotemetaFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "string", 0, nil)}
}

// Nl2brFunction 实现 nl2br 函数
// nl2br(string $string, bool $use_xhtml = true): string
type Nl2brFunction struct{}

func NewNl2brFunction() data.FuncStmt { return &Nl2brFunction{} }

func (f *Nl2brFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strV, _ := ctx.GetIndexValue(0)
	xhtmlV, _ := ctx.GetIndexValue(1)

	if strV == nil {
		return data.NewStringValue(""), nil
	}
	str := strV.AsString()
	useXhtml := true
	if xhtmlV != nil {
		if bv, ok := xhtmlV.(*data.BoolValue); ok {
			useXhtml = bv.Value
		}
	}
	br := "<br />"
	if !useXhtml {
		br = "<br>"
	}
	// 将 \n 替换为 <br />\n，将 \r\n 处理为 <br />\r\n
	result := strings.ReplaceAll(str, "\r\n", br+"\r\n")
	result = strings.ReplaceAll(result, "\n", br+"\n")
	return data.NewStringValue(result), nil
}

func (f *Nl2brFunction) GetName() string { return "nl2br" }
func (f *Nl2brFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "use_xhtml", 1, node.NewBooleanLiteral(nil, true), data.NewBaseType("bool")),
	}
}
func (f *Nl2brFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, nil),
		node.NewVariable(nil, "use_xhtml", 1, data.NewBaseType("bool")),
	}
}

// SubstrCompareFunction 实现 substr_compare 函数
// substr_compare(string $haystack, string $needle, int $offset, ?int $length = null, bool $case_insensitive = false): int
type SubstrCompareFunction struct{}

func NewSubstrCompareFunction() data.FuncStmt { return &SubstrCompareFunction{} }

func (f *SubstrCompareFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	hayV, _ := ctx.GetIndexValue(0)
	needleV, _ := ctx.GetIndexValue(1)
	offsetV, _ := ctx.GetIndexValue(2)
	lengthV, _ := ctx.GetIndexValue(3)
	caseInsensitiveV, _ := ctx.GetIndexValue(4)

	if hayV == nil || needleV == nil {
		return data.NewIntValue(-1), nil
	}
	haystack := hayV.AsString()
	needle := needleV.AsString()

	offset := 0
	if offsetV != nil {
		if iv, ok := offsetV.(*data.IntValue); ok {
			offset = iv.Value
		}
	}
	if offset < 0 {
		offset = len(haystack) + offset
	}
	if offset < 0 || offset > len(haystack) {
		return data.NewIntValue(-1), nil
	}

	length := len(haystack) - offset
	if lengthV != nil {
		if _, isNull := lengthV.(*data.NullValue); isNull {
			// 未指定 length 时，使用 needle 的长度
			length = len(needle)
		} else if iv, ok := lengthV.(*data.IntValue); ok {
			length = iv.Value
		}
	}
	if length < 0 {
		length = len(haystack) + length - offset
	}
	if length > len(haystack)-offset {
		length = len(haystack) - offset
	}
	if length < 0 {
		length = 0
	}

	substr := haystack[offset : offset+length]

	caseInsensitive := false
	if caseInsensitiveV != nil {
		if bv, ok := caseInsensitiveV.(*data.BoolValue); ok {
			caseInsensitive = bv.Value
		}
	}

	if caseInsensitive {
		return data.NewIntValue(strings.Compare(strings.ToLower(substr), strings.ToLower(needle))), nil
	}
	return data.NewIntValue(strings.Compare(substr, needle)), nil
}

func (f *SubstrCompareFunction) GetName() string { return "substr_compare" }
func (f *SubstrCompareFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "haystack", 0, nil, nil),
		node.NewParameter(nil, "needle", 1, nil, nil),
		node.NewParameter(nil, "offset", 2, nil, data.NewBaseType("int")),
		node.NewParameter(nil, "length", 3, node.NewNullLiteral(nil), nil),
		node.NewParameter(nil, "case_insensitive", 4, node.NewBooleanLiteral(nil, false), data.NewBaseType("bool")),
	}
}
func (f *SubstrCompareFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "haystack", 0, nil),
		node.NewVariable(nil, "needle", 1, nil),
		node.NewVariable(nil, "offset", 2, data.NewBaseType("int")),
		node.NewVariable(nil, "length", 3, nil),
		node.NewVariable(nil, "case_insensitive", 4, data.NewBaseType("bool")),
	}
}
