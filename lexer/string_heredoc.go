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

	// 查找结束标识符（支持前导空格/制表符；结束标记须独占一行）
	tryCloseMarker := func(at int) (int, bool) {
		markerStart := at
		for markerStart < len(input) && (input[markerStart] == ' ' || input[markerStart] == '\t') {
			markerStart++
		}
		if markerStart+len(identifier) > len(input) {
			return 0, false
		}
		if input[markerStart:markerStart+len(identifier)] != identifier {
			return 0, false
		}
		after := markerStart + len(identifier)
		for after < len(input) && (input[after] == ' ' || input[after] == '\t') {
			after++
		}
		if after >= len(input) || input[after] == '\n' || input[after] == '\r' || input[after] == ';' {
			return markerStart + len(identifier), true
		}
		return 0, false
	}

	emit := func(endPos int) (SpecialToken, int, bool) {
		literal := input[start:endPos]
		tt := token.HEREDOC
		if _, isNowdoc := HeredocTokenType(literal); isNowdoc {
			tt = token.NOWDOC
		}
		return SpecialToken{
			Type:    tt,
			Literal: literal,
			Length:  endPos - start,
		}, endPos, true
	}

	searchStart := pos
	if endPos, ok := tryCloseMarker(searchStart); ok {
		return emit(endPos)
	}

	for searchStart < len(input) {
		newlinePos := strings.IndexByte(input[searchStart:], '\n')
		if newlinePos == -1 {
			break
		}
		newlinePos += searchStart
		markerStart := newlinePos + 1
		if endPos, ok := tryCloseMarker(markerStart); ok {
			return emit(endPos)
		}
		searchStart = newlinePos + 1
	}

	return SpecialToken{}, start, false
}
