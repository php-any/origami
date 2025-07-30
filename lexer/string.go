package lexer

// isStringStart 检查是否是字符串开始
func isStringStart(input string, start int) (string, bool) {
	if start >= len(input) {
		return "", false
	}

	// 检查单引号、双引号或反引号
	if input[start] == '\'' || input[start] == '"' || input[start] == '`' {
		return string(input[start]), true
	}

	// 检查 heredoc 语法
	if start+2 < len(input) && input[start:start+3] == "<<<" {
		// 找到 heredoc 标识符的结束位置
		pos := start + 3
		for pos < len(input) && (input[pos] == ' ' || input[pos] == '\t') {
			pos++
		}
		idStart := pos
		for pos < len(input) && (input[pos] != '\n' && input[pos] != '\r') {
			pos++
		}
		if pos > idStart {
			return input[idStart:pos], true
		}
	}

	return "", false
}

// HandleString 处理字符串
func HandleString(input string, start int) (SpecialToken, int, bool) {
	if start >= len(input) {
		return SpecialToken{}, start, false
	}

	// 检查字符串开始
	if quote, ok := isStringStart(input, start); ok {
		if len(quote) == 1 {
			// 根据引号类型选择处理函数
			if quote[0] == '\'' {
				return handleSingleQuotedString(input, start)
			} else if quote[0] == '`' {
				return handleBacktickString(input, start)
			} else {
				return handleDoubleQuotedString(input, start)
			}
		} else {
			// heredoc 字符串
			return handleHeredocString(input, start, quote)
		}
	}

	return SpecialToken{}, start, false
}
