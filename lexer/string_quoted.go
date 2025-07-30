package lexer

import (
	"github.com/php-any/origami/token"
)

// handleSingleQuotedString 处理单引号字符串
func handleSingleQuotedString(input string, start int) (SpecialToken, int, bool) {
	pos := start + 1
	escaped := false

	for pos < len(input) {
		if !escaped && input[pos] == '\'' {
			// 找到字符串结束，包含结束引号
			literal := input[start : pos+1]
			return SpecialToken{
				Type:    token.STRING,
				Literal: literal,
				Length:  pos + 1 - start,
			}, pos + 1, true
		}

		// 处理转义字符
		if input[pos] == '\\' {
			escaped = !escaped
		} else {
			escaped = false
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
			// 找到字符串结束，包含结束引号
			literal := input[start : pos+1]
			return SpecialToken{
				Type:    token.STRING,
				Literal: literal,
				Length:  pos + 1 - start,
			}, pos + 1, true
		}

		// 处理转义字符
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
			// 找到字符串结束，包含结束引号
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
