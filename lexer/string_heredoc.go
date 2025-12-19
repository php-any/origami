package lexer

import (
	"strings"

	"github.com/php-any/origami/token"
)

// handleHeredocString 处理 heredoc/nowdoc 字符串
func handleHeredocString(input string, start int, identifier string) (SpecialToken, int, bool) {
	// 跳过 <<<
	pos := start + 3
	for pos < len(input) && (input[pos] == ' ' || input[pos] == '\t') {
		pos++
	}
	// 检查是否是 nowdoc（单引号括起来的标识符）
	if pos < len(input) && input[pos] == '\'' {
		pos++ // 跳过开始的单引号
		pos += len(identifier)
		// 跳过结束的单引号
		if pos < len(input) && input[pos] == '\'' {
			pos++
		}
	} else {
		pos += len(identifier)
	}

	// 跳过换行符
	for pos < len(input) && (input[pos] == '\n' || input[pos] == '\r') {
		pos++
	}

	// 查找结束标识符（支持前导空格/制表符）
	searchStart := pos
	for searchStart < len(input) {
		// 查找换行符
		newlinePos := strings.IndexByte(input[searchStart:], '\n')
		if newlinePos == -1 {
			return SpecialToken{}, start, false
		}
		newlinePos += searchStart

		// 跳过换行符和前导空格/制表符
		markerStart := newlinePos + 1
		for markerStart < len(input) && (input[markerStart] == ' ' || input[markerStart] == '\t') {
			markerStart++
		}

		// 检查是否是结束标识符
		if markerStart+len(identifier) <= len(input) && input[markerStart:markerStart+len(identifier)] == identifier {
			// 找到结束标识符
			endPos := markerStart + len(identifier)
			literal := input[start:endPos]
			return SpecialToken{
				Type:    token.STRING,
				Literal: literal,
				Length:  endPos - start,
			}, endPos, true
		}

		// 继续查找下一个换行符
		searchStart = newlinePos + 1
	}

	return SpecialToken{}, start, false
}
