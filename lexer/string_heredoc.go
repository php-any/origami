package lexer

import (
	"strings"

	"github.com/php-any/origami/token"
)

// handleHeredocString 处理 heredoc 字符串
func handleHeredocString(input string, start int, identifier string) (SpecialToken, int, bool) {
	// 跳过 <<< 和标识符
	pos := start + 3
	for pos < len(input) && (input[pos] == ' ' || input[pos] == '\t') {
		pos++
	}
	pos += len(identifier)

	// 跳过换行符
	for pos < len(input) && (input[pos] == '\n' || input[pos] == '\r') {
		pos++
	}

	// 查找结束标识符
	endMarker := "\n" + identifier
	if pos+len(endMarker) >= len(input) {
		return SpecialToken{}, start, false
	}

	endPos := strings.Index(input[pos:], endMarker)
	if endPos == -1 {
		return SpecialToken{}, start, false
	}

	// 计算结束位置（包含结束标识符）
	endPos = pos + endPos + len(endMarker)

	// 创建 token
	literal := input[start:endPos]
	return SpecialToken{
		Type:    token.STRING,
		Literal: literal,
		Length:  endPos - start,
	}, endPos, true
}
