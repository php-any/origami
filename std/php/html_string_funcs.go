package php

import (
	"html"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// HtmlentitiesFunction 实现 htmlentities 函数
// htmlentities(string $string, int $flags = ENT_QUOTES | ENT_SUBSTITUTE | ENT_HTML401, ?string $encoding = null, bool $double_encode = true): string
type HtmlentitiesFunction struct{}

func NewHtmlentitiesFunction() data.FuncStmt { return &HtmlentitiesFunction{} }

func (f *HtmlentitiesFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strV, _ := ctx.GetIndexValue(0)
	_, _ = ctx.GetIndexValue(1) // flags 暂时忽略
	// encoding 参数暂时忽略
	doubleEncodeV, _ := ctx.GetIndexValue(3)

	if strV == nil {
		return data.NewStringValue(""), nil
	}
	str := strV.AsString()

	doubleEncode := true
	if doubleEncodeV != nil {
		if bv, ok := doubleEncodeV.(*data.BoolValue); ok {
			doubleEncode = bv.Value
		}
	}

	// 使用 html.EscapeString 处理基本的 & < > "
	result := html.EscapeString(str)

	// 额外处理单引号
	result = strings.ReplaceAll(result, "'", "&#039;")

	if !doubleEncode {
		// 不 double-encode 已经编码的实体
		result = strings.ReplaceAll(result, "&amp;", "&")
	}

	return data.NewStringValue(result), nil
}

func (f *HtmlentitiesFunction) GetName() string { return "htmlentities" }
var htmlentitiesFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "string", 0, nil, nil),
	node.NewParameter(nil, "flags", 1, node.NewIntLiteral(nil, "3"), data.NewBaseType("int")),
	node.NewParameter(nil, "encoding", 2, node.NewNullLiteral(nil), nil),
	node.NewParameter(nil, "double_encode", 3, node.NewBooleanLiteral(nil, true), data.NewBaseType("bool")),
}

func (f *HtmlentitiesFunction) GetParams() []data.GetValue {
	return htmlentitiesFunctionGetParams
}
var htmlentitiesFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "string", 0, nil),
	node.NewVariable(nil, "flags", 1, data.NewBaseType("int")),
	node.NewVariable(nil, "encoding", 2, nil),
	node.NewVariable(nil, "double_encode", 3, data.NewBaseType("bool")),
}

func (f *HtmlentitiesFunction) GetVariables() []data.Variable {
	return htmlentitiesFunctionGetVariables
}

// HtmlspecialcharsDecodeFunction 实现 htmlspecialchars_decode 函数
// htmlspecialchars_decode(string $string, int $flags = ENT_QUOTES | ENT_SUBSTITUTE | ENT_HTML401): string
type HtmlspecialcharsDecodeFunction struct{}

func NewHtmlspecialcharsDecodeFunction() data.FuncStmt { return &HtmlspecialcharsDecodeFunction{} }

func (f *HtmlspecialcharsDecodeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strV, _ := ctx.GetIndexValue(0)
	flagsV, _ := ctx.GetIndexValue(1)

	if strV == nil {
		return data.NewStringValue(""), nil
	}
	str := strV.AsString()

	// 总是解码 &amp; &lt; &gt;
	result := html.UnescapeString(str)

	// 根据 flags 决定是否解码 &quot; 和 &#039;
	// ENT_QUOTES (3) 解码单引号和双引号
	// ENT_COMPAT (2) 只解码双引号
	// ENT_NOQUOTES (0) 不解码引号

	// flags 参数检查
	flags := 3
	if flagsV != nil {
		if iv, ok := flagsV.(*data.IntValue); ok {
			flags = iv.Value
		}
	}

	if flags&1 == 1 { // ENT_QUOTES
		result = strings.ReplaceAll(result, "&#039;", "'")
		result = strings.ReplaceAll(result, "&#39;", "'")
	}
	if flags&2 == 2 { // ENT_COMPAT 或 ENT_QUOTES
		result = strings.ReplaceAll(result, "&quot;", "\"")
	}

	return data.NewStringValue(result), nil
}

func (f *HtmlspecialcharsDecodeFunction) GetName() string { return "htmlspecialchars_decode" }
var htmlspecialcharsDecodeFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "string", 0, nil, nil),
	node.NewParameter(nil, "flags", 1, node.NewIntLiteral(nil, "3"), data.NewBaseType("int")),
}

func (f *HtmlspecialcharsDecodeFunction) GetParams() []data.GetValue {
	return htmlspecialcharsDecodeFunctionGetParams
}
var htmlspecialcharsDecodeFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "string", 0, nil),
	node.NewVariable(nil, "flags", 1, data.NewBaseType("int")),
}

func (f *HtmlspecialcharsDecodeFunction) GetVariables() []data.Variable {
	return htmlspecialcharsDecodeFunctionGetVariables
}

// AddcslashesFunction 实现 addcslashes 函数
// addcslashes(string $string, string $characters): string
type AddcslashesFunction struct{}

func NewAddcslashesFunction() data.FuncStmt { return &AddcslashesFunction{} }

func (f *AddcslashesFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strV, _ := ctx.GetIndexValue(0)
	charsV, _ := ctx.GetIndexValue(1)

	if strV == nil || charsV == nil {
		return data.NewStringValue(""), nil
	}
	str := strV.AsString()
	chars := charsV.AsString()

	var result strings.Builder
	for _, c := range str {
		if strings.ContainsRune(chars, c) {
			result.WriteByte('\\')
		}
		result.WriteRune(c)
	}
	return data.NewStringValue(result.String()), nil
}

func (f *AddcslashesFunction) GetName() string { return "addcslashes" }
var addcslashesFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "string", 0, nil, nil),
	node.NewParameter(nil, "characters", 1, nil, nil),
}

func (f *AddcslashesFunction) GetParams() []data.GetValue {
	return addcslashesFunctionGetParams
}
var addcslashesFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "string", 0, nil),
	node.NewVariable(nil, "characters", 1, nil),
}

func (f *AddcslashesFunction) GetVariables() []data.Variable {
	return addcslashesFunctionGetVariables
}

// MetaphoneFunction 实现 metaphone 函数
// metaphone(string $string, int $max_phonemes = 0): string
type MetaphoneFunction struct{}

func NewMetaphoneFunction() data.FuncStmt { return &MetaphoneFunction{} }

func (f *MetaphoneFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strV, _ := ctx.GetIndexValue(0)
	maxV, _ := ctx.GetIndexValue(1)

	if strV == nil {
		return data.NewStringValue(""), nil
	}
	str := strV.AsString()

	maxPhonemes := 0
	if maxV != nil {
		if iv, ok := maxV.(*data.IntValue); ok {
			maxPhonemes = iv.Value
		}
	}

	// 简化的 metaphone 算法实现
	metaphone := strings.ToUpper(str)
	metaphone = strings.ReplaceAll(metaphone, "PH", "F")
	metaphone = strings.ReplaceAll(metaphone, "KN", "N")
	metaphone = strings.ReplaceAll(metaphone, "GN", "N")
	metaphone = strings.ReplaceAll(metaphone, "MB", "M")
	metaphone = strings.ReplaceAll(metaphone, "WR", "R")
	metaphone = strings.ReplaceAll(metaphone, "X", "KS")
	metaphone = strings.ReplaceAll(metaphone, "CI", "S")
	metaphone = strings.ReplaceAll(metaphone, "CE", "S")
	metaphone = strings.ReplaceAll(metaphone, "CY", "S")
	metaphone = strings.ReplaceAll(metaphone, "GI", "J")
	metaphone = strings.ReplaceAll(metaphone, "GE", "J")
	metaphone = strings.ReplaceAll(metaphone, "GY", "J")
	metaphone = strings.ReplaceAll(metaphone, "CK", "K")
	metaphone = strings.ReplaceAll(metaphone, "DG", "J")
	metaphone = strings.ReplaceAll(metaphone, "GH", "")
	metaphone = strings.ReplaceAll(metaphone, "QU", "K")
	metaphone = strings.ReplaceAll(metaphone, "TH", "0")
	metaphone = strings.ReplaceAll(metaphone, "TCH", "CH")

	// 移除元音（除非在开头）
	var result strings.Builder
	for i := 0; i < len(metaphone); i++ {
		c := metaphone[i]
		isVowel := c == 'A' || c == 'E' || c == 'I' || c == 'O' || c == 'U'
		if isVowel && i > 0 {
			continue
		}
		if c == 'H' && i > 0 && i < len(metaphone)-1 {
			continue
		}
		if c == 'W' || c == 'Y' {
			if i > 0 && i < len(metaphone)-1 {
				continue
			}
		}
		result.WriteByte(c)
	}

	if maxPhonemes > 0 && result.Len() > maxPhonemes {
		return data.NewStringValue(result.String()[:maxPhonemes]), nil
	}
	return data.NewStringValue(result.String()), nil
}

func (f *MetaphoneFunction) GetName() string { return "metaphone" }
var metaphoneFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "string", 0, nil, nil),
	node.NewParameter(nil, "max_phonemes", 1, node.NewIntLiteral(nil, "0"), data.NewBaseType("int")),
}

func (f *MetaphoneFunction) GetParams() []data.GetValue {
	return metaphoneFunctionGetParams
}
var metaphoneFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "string", 0, nil),
	node.NewVariable(nil, "max_phonemes", 1, data.NewBaseType("int")),
}

func (f *MetaphoneFunction) GetVariables() []data.Variable {
	return metaphoneFunctionGetVariables
}

// SoundexFunction 实现 soundex 函数
// soundex(string $string): string
type SoundexFunction struct{}

func NewSoundexFunction() data.FuncStmt { return &SoundexFunction{} }

func (f *SoundexFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strV, _ := ctx.GetIndexValue(0)
	if strV == nil {
		return data.NewStringValue(""), nil
	}
	str := strings.ToUpper(strV.AsString())
	if len(str) == 0 {
		return data.NewStringValue(""), nil
	}

	// Soundex 编码表
	codeMap := map[rune]byte{
		'B': '1', 'F': '1', 'P': '1', 'V': '1',
		'C': '2', 'G': '2', 'J': '2', 'K': '2', 'Q': '2', 'S': '2', 'X': '2', 'Z': '2',
		'D': '3', 'T': '3',
		'L': '4',
		'M': '5', 'N': '5',
		'R': '6',
	}

	runes := []rune(str)
	result := make([]byte, 0, 4)
	result = append(result, byte(runes[0]))

	prevCode := byte(0)
	if code, ok := codeMap[runes[0]]; ok {
		prevCode = code
	}

	for i := 1; i < len(runes) && len(result) < 4; i++ {
		code, ok := codeMap[runes[i]]
		if !ok {
			prevCode = 0
			continue
		}
		if code != prevCode {
			result = append(result, code)
			prevCode = code
		}
	}

	// 补齐到4位
	for len(result) < 4 {
		result = append(result, '0')
	}

	return data.NewStringValue(string(result)), nil
}

func (f *SoundexFunction) GetName() string { return "soundex" }
var soundexFunctionGetParams = []data.GetValue{node.NewParameter(nil, "string", 0, nil, nil)}

func (f *SoundexFunction) GetParams() []data.GetValue {
	return soundexFunctionGetParams
}
var soundexFunctionGetVariables = []data.Variable{node.NewVariable(nil, "string", 0, nil)}

func (f *SoundexFunction) GetVariables() []data.Variable {
	return soundexFunctionGetVariables
}

// SimilarTextFunction 实现 similar_text 函数
// similar_text(string $string1, string $string2, float &$percent = null): int
type SimilarTextFunction struct{}

func NewSimilarTextFunction() data.FuncStmt { return &SimilarTextFunction{} }

func (f *SimilarTextFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	str1V, _ := ctx.GetIndexValue(0)
	str2V, _ := ctx.GetIndexValue(1)
	percentRef, _ := ctx.GetIndexValue(2)

	if str1V == nil || str2V == nil {
		return data.NewIntValue(0), nil
	}
	s1 := str1V.AsString()
	s2 := str2V.AsString()

	// 计算最长公共子序列长度
	similar := lcsLength(s1, s2)

	// 计算百分比
	total := len(s1) + len(s2)
	if total > 0 {
		percent := float64(similar * 200) / float64(total)
		// 设置 percent 引用参数
		if percentRef != nil {
			if iv, ok := percentRef.(*data.FloatValue); ok {
				iv.Value = percent
			} else if iv, ok := percentRef.(*data.IntValue); ok {
				iv.Value = int(percent)
			}
		}
	}

	return data.NewIntValue(similar), nil
}

func lcsLength(a, b string) int {
	m, n := len(a), len(b)
	if m == 0 || n == 0 {
		return 0
	}
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				if dp[i-1][j] > dp[i][j-1] {
					dp[i][j] = dp[i-1][j]
				} else {
					dp[i][j] = dp[i][j-1]
				}
			}
		}
	}
	return dp[m][n]
}

func (f *SimilarTextFunction) GetName() string { return "similar_text" }
var similarTextFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "string1", 0, nil, nil),
	node.NewParameter(nil, "string2", 1, nil, nil),
	node.NewParameterReference(nil, "percent", 2, node.NewNullLiteral(nil), data.NewBaseType("float")),
}

func (f *SimilarTextFunction) GetParams() []data.GetValue {
	return similarTextFunctionGetParams
}
var similarTextFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "string1", 0, nil),
	node.NewVariable(nil, "string2", 1, nil),
	node.NewVariable(nil, "percent", 2, data.NewBaseType("float")),
}

func (f *SimilarTextFunction) GetVariables() []data.Variable {
	return similarTextFunctionGetVariables
}
