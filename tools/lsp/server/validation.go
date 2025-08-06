package server

import (
	"strings"
)

// validateDocument 验证文档语法
// 检查 Origami 语言文档的语法错误并返回诊断信息
func (s *Server) validateDocument(uri, content string) {
	// 创建诊断信息数组
	diagnostics := []map[string]interface{}{}

	// 按行分割文档内容进行逐行检查
	lines := strings.Split(content, "\n")

	// 遍历每一行进行语法检查
	for lineNumber, line := range lines {
		// 去除行首尾空白字符
		trimmedLine := strings.TrimSpace(line)

		// 跳过空行和注释行
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") || strings.HasPrefix(trimmedLine, "/*") {
			continue
		}

		// 检查各种语法错误
		lineDiagnostics := s.checkLineSyntax(line, lineNumber)
		diagnostics = append(diagnostics, lineDiagnostics...)
	}

	// 发送诊断信息到客户端
	s.sendDiagnostics(uri, diagnostics)
}

// checkLineSyntax 检查单行语法
// 返回该行发现的所有语法错误
func (s *Server) checkLineSyntax(line string, lineNumber int) []map[string]interface{} {
	diagnostics := []map[string]interface{}{}
	trimmedLine := strings.TrimSpace(line)

	// 检查语句是否以分号结尾（排除控制结构）
	if s.shouldEndWithSemicolon(trimmedLine) && !strings.HasSuffix(trimmedLine, ";") {
		diagnostics = append(diagnostics, s.createDiagnostic(
			lineNumber,
			len(line)-len(strings.TrimRightFunc(line, func(r rune) bool { return r == ' ' || r == '\t' })),
			len(line),
			"语句应该以分号结尾",
			2, // Warning
		))
	}

	// 检查括号匹配
	if !s.checkBracketBalance(trimmedLine) {
		diagnostics = append(diagnostics, s.createDiagnostic(
			lineNumber,
			0,
			len(line),
			"括号不匹配",
			1, // Error
		))
	}

	// 检查关键字拼写
	if misspelledKeyword := s.checkKeywordSpelling(trimmedLine); misspelledKeyword != "" {
		// 找到拼写错误关键字的位置
		start := strings.Index(line, misspelledKeyword)
		if start != -1 {
			diagnostics = append(diagnostics, s.createDiagnostic(
				lineNumber,
				start,
				start+len(misspelledKeyword),
				"可能的关键字拼写错误: "+misspelledKeyword,
				2, // Warning
			))
		}
	}

	return diagnostics
}

// shouldEndWithSemicolon 检查语句是否应该以分号结尾
// 控制结构（if、for、while等）的开始行不需要分号
func (s *Server) shouldEndWithSemicolon(line string) bool {
	// 控制结构关键字列表
	controlKeywords := []string{
		"if", "else", "for", "foreach", "while", "switch", "case", "default",
		"function", "class", "try", "catch", "finally",
	}

	// 检查是否以控制结构关键字开始
	for _, keyword := range controlKeywords {
		if strings.HasPrefix(line, keyword+" ") || strings.HasPrefix(line, keyword+"(") {
			return false
		}
	}

	// 检查是否以大括号结尾（代码块）
	if strings.HasSuffix(line, "{") || strings.HasSuffix(line, "}") {
		return false
	}

	// 其他语句应该以分号结尾
	return true
}

// checkBracketBalance 检查括号是否平衡
// 检查圆括号、方括号和大括号的匹配
func (s *Server) checkBracketBalance(line string) bool {
	stack := []rune{}
	brackets := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}

	for _, char := range line {
		switch char {
		case '(', '[', '{':
			stack = append(stack, char)
		case ')', ']', '}':
			if len(stack) == 0 {
				return false // 右括号多于左括号
			}
			if stack[len(stack)-1] != brackets[char] {
				return false // 括号类型不匹配
			}
			stack = stack[:len(stack)-1] // 弹出匹配的左括号
		}
	}

	// 如果栈不为空，说明有未匹配的左括号
	return len(stack) == 0
}

// checkKeywordSpelling 检查关键字拼写
// 返回可能拼写错误的关键字，如果没有错误返回空字符串
func (s *Server) checkKeywordSpelling(line string) string {
	// 正确的关键字列表
	correctKeywords := map[string]bool{
		"if": true, "else": true, "for": true, "foreach": true, "while": true,
		"function": true, "class": true, "return": true, "break": true, "continue": true,
		"switch": true, "case": true, "default": true, "try": true, "catch": true,
		"finally": true, "echo": true, "print": true, "var": true, "let": true,
		"const": true, "public": true, "private": true, "protected": true,
	}

	// 常见的拼写错误映射
	commonMisspellings := map[string]string{
		"fi":       "if",
		"esle":     "else",
		"fro":      "for",
		"whiel":    "while",
		"fucntion": "function",
		"calss":    "class",
		"retrun":   "return",
		"braek":    "break",
		"contine":  "continue",
		"swtich":   "switch",
		"defualt":  "default",
		"tyr":      "try",
		"cathc":    "catch",
		"ehco":     "echo",
	}

	// 提取行中的单词
	words := strings.Fields(line)
	for _, word := range words {
		// 移除标点符号
		cleanWord := strings.Trim(word, "(){}[];,.")

		// 检查是否是常见的拼写错误
		if _, exists := commonMisspellings[cleanWord]; exists {
			return cleanWord
		}

		// 检查是否是看起来像关键字但拼写错误的单词
		if s.looksLikeKeyword(cleanWord) && !correctKeywords[cleanWord] {
			return cleanWord
		}
	}

	return ""
}

// looksLikeKeyword 检查单词是否看起来像关键字
// 基于长度和字符组成进行简单判断
func (s *Server) looksLikeKeyword(word string) bool {
	// 关键字通常是小写字母组成的短单词
	if len(word) < 2 || len(word) > 10 {
		return false
	}

	for _, char := range word {
		if char < 'a' || char > 'z' {
			return false
		}
	}

	return true
}

// createDiagnostic 创建诊断信息
// 构建符合 LSP 规范的诊断信息对象
func (s *Server) createDiagnostic(line, startChar, endChar int, message string, severity int) map[string]interface{} {
	return map[string]interface{}{
		"range": map[string]interface{}{
			"start": map[string]interface{}{
				"line":      line,
				"character": startChar,
			},
			"end": map[string]interface{}{
				"line":      line,
				"character": endChar,
			},
		},
		"severity": severity, // 1 = Error, 2 = Warning, 3 = Information, 4 = Hint
		"message":  message,
		"source":   "origami-lsp",
	}
}

// sendDiagnostics 发送诊断信息到客户端
// 使用 LSP 的 textDocument/publishDiagnostics 通知
func (s *Server) sendDiagnostics(uri string, diagnostics []map[string]interface{}) {
	notification := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "textDocument/publishDiagnostics",
		"params": map[string]interface{}{
			"uri":         uri,
			"diagnostics": diagnostics,
		},
	}

	// 发送通知（不需要响应）
	s.sendMessage(notification)
}
