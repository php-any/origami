package lexer

import (
	"github.com/php-any/origami/token"
)

// handleSingleQuotedString 处理单引号字符串
// PHP 单引号字符串中，只有 \\ (literal backslash) 和 \' (literal quote) 是转义序列
// \\' 在字符串末尾时表示 literal backslash 后跟结束引号
func handleSingleQuotedString(input string, start int) (SpecialToken, int, bool) {
	pos := start + 1

	for pos < len(input) {
		ch := input[pos]

		// 处理反斜杠转义
		if ch == '\\' && pos+1 < len(input) {
			next := input[pos+1]
			if next == '\\' || next == '\'' {
				// \\ -> literal backslash, 跳过两个字符
				// \' -> literal single quote, 跳过两个字符（这是字符串内容的一部分，不是结束）
				pos += 2
				continue
			}
		}

		// 处理单引号结束
		if ch == '\'' {
			literal := input[start : pos+1]
			return SpecialToken{
				Type:    token.STRING,
				Literal: literal,
				Length:  pos + 1 - start,
			}, pos + 1, true
		}

		pos++
	}

	return SpecialToken{}, start, false
}

// handleDoubleQuotedString 处理双引号字符串
func handleDoubleQuotedString(input string, start int) (SpecialToken, int, bool) {
	pos := start + 1
	escaped := false

	for pos < len(input) {
		if !escaped && input[pos] == '"' {
			literal := input[start : pos+1]
			return SpecialToken{
				Type:    token.STRING,
				Literal: literal,
				Length:  pos + 1 - start,
			}, pos + 1, true
		}

		if input[pos] == '\\' {
			escaped = !escaped
		} else {
			escaped = false
		}
		pos++
	}

	return SpecialToken{}, start, false
}

// handleBacktickString 处理反引号字符串
func handleBacktickString(input string, start int) (SpecialToken, int, bool) {
	pos := start + 1
	for pos < len(input) {
		if input[pos] == '`' {
			literal := input[start : pos+1]
			return SpecialToken{
				Type:    token.STRING,
				Literal: literal,
				Length:  pos + 1 - start,
			}, pos + 1, true
		}
		pos++
	}
	return SpecialToken{}, start, false
}
