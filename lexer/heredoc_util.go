package lexer

import "strings"

// ExtractHeredocBody 从 <<<IDENT ... IDENT 字面量中提取正文，并判断是否为 nowdoc。
// 供预处理器与 HeredocParser 共用。
func ExtractHeredocBody(literal string) (body string, isNowdoc bool, ok bool) {
	if len(literal) < 3 || literal[:3] != "<<<" {
		return "", false, false
	}
	pos := 3
	if pos < len(literal) && literal[pos] == '\'' {
		isNowdoc = true
		pos++
	}
	for pos < len(literal) && literal[pos] != '\n' && literal[pos] != '\r' {
		pos++
	}
	for pos < len(literal) && (literal[pos] == '\n' || literal[pos] == '\r') {
		if literal[pos] == '\r' && pos+1 < len(literal) && literal[pos+1] == '\n' {
			pos += 2
		} else {
			pos++
		}
	}
	firstContent := pos
	lastNewline := strings.LastIndexByte(literal, '\n')
	if lastNewline < firstContent {
		return "", isNowdoc, true
	}
	body = literal[firstContent:lastNewline]
	body = strings.TrimSuffix(body, "\r")
	return body, isNowdoc, true
}

// HeredocTokenType 根据 heredoc 字面量判断应使用的 token 类型。
func HeredocTokenType(literal string) (heredoc bool, nowdoc bool) {
	if len(literal) < 4 || literal[:3] != "<<<" {
		return false, false
	}
	return true, literal[3] == '\''
}
