package lexer

import (
	"unicode"

	"github.com/php-any/origami/token"
)

// HtmlLexer 专门用于 HTML 文件分词的词法分析器
// 与普通 Lexer 不同，HtmlLexer 保留所有字符，包括空格和换行符
type HtmlLexer struct {
	input   string // 输入内容
	pos     int    // 当前位置
	line    int    // 当前行号
	linePos int    // 当前行中的位置
	inTag   bool   // 是否在标签内部（< 和 > 之间）
}

// NewHtmlLexer 创建一个新的 HTML 词法分析器
func NewHtmlLexer() *HtmlLexer {
	return &HtmlLexer{}
}

// Tokenize 将 HTML 输入字符串转换为 token 列表
// 保留所有字符，包括空格、制表符、换行符等
func (h *HtmlLexer) Tokenize(input string) []Token {
	h.input = input
	h.pos = 0
	h.line = 0
	h.linePos = 0

	var tokens []Token

	var tok Token
	var ok bool
	if tok, ok = h.processDoctype(); ok {
		tokens = append(tokens, tok)
	}

	for h.pos < len(h.input) {
		// 按优先级尝试各种处理函数

		if tok, ok = h.processCdata(); ok {
			tokens = append(tokens, tok)
			continue
		}

		if tok, ok = h.processProcessingInstruction(); ok {
			tokens = append(tokens, tok)
			continue
		}

		if tok, ok = h.processHtmlComment(); ok {
			tokens = append(tokens, tok)
			continue
		}

		if tok, ok = h.processDoubleQuotedString(); ok {
			tokens = append(tokens, tok)
			continue
		}

		if tok, ok = h.processSingleQuotedString(); ok {
			tokens = append(tokens, tok)
			continue
		}

		if tok, ok = h.processBacktickString(); ok {
			tokens = append(tokens, tok)
			continue
		}

		if tok, ok = h.processJsServer(); ok {
			tokens = append(tokens, tok)
			continue
		}

		if tok, ok = h.processLt(); ok {
			tokens = append(tokens, tok)
			continue
		}

		if tok, ok = h.processGt(); ok {
			tokens = append(tokens, tok)
			continue
		}

		// 处理 / 符号（用于自闭合标签 /> 和结束标签 </tag>）
		if tok, ok = h.processQuo(); ok {
			tokens = append(tokens, tok)
			continue
		}

		// 处理标识符（标签名、属性名等，在标签内部和结束标签中都需要）
		if tok, ok = h.processIdentifier(); ok {
			tokens = append(tokens, tok)
			continue
		}

		// 在标签内部时，处理标签相关的 token
		if h.inTag {
			if tok, ok = h.processNumber(); ok {
				tokens = append(tokens, tok)
				continue
			}

			if tok, ok = h.processAssign(); ok {
				tokens = append(tokens, tok)
				continue
			}

			// 在标签内部时，跳过空白字符
			if h.pos < len(h.input) && h.isWhitespace(h.input[h.pos]) {
				h.advance()
				continue
			}
		} else {
			// 在标签外部时，优先处理空白字符（避免被 processText 包含）
			if tok, ok = h.processWhitespace(); ok {
				tokens = append(tokens, tok)
				continue
			}

			// 在标签外部时，处理文本内容（作为整体字符串）
			if tok, ok = h.processText(); ok {
				tokens = append(tokens, tok)
				continue
			}
		}

		// 如果所有处理都失败，跳过当前字符（防止无限循环）
		h.advance()
	}
	PrintTokens(tokens, "")
	return tokens
}

// processDoctype 处理 DOCTYPE 声明 <!DOCTYPE
// 只返回 <!DOCTYPE 部分，后续内容（如 html）和 > 由其他处理函数处理
// 返回 token 和是否成功处理
func (h *HtmlLexer) processDoctype() (Token, bool) {
	if h.pos+9 > len(h.input) || h.input[h.pos:h.pos+9] != "<!DOCTYPE" {
		return nil, false
	}

	start := h.pos
	startLine := h.line
	startLinePos := h.linePos

	// 只读取 "<!DOCTYPE" 部分（9个字符）
	for i := 0; i < 9; i++ {
		h.advance()
	}

	return NewWorkerToken(
		token.DOCTYPE,
		h.input[start:h.pos],
		start,
		h.pos,
		startLine,
		startLinePos,
	), true
}

// processCdata 处理 CDATA 部分 <![CDATA[...]]>
// 返回 token 和是否成功处理
func (h *HtmlLexer) processCdata() (Token, bool) {
	if h.pos+9 > len(h.input) || h.input[h.pos:h.pos+9] != "<![CDATA[" {
		return nil, false
	}

	start := h.pos
	startLine := h.line
	startLinePos := h.linePos

	// 跳过 "<![CDATA["
	for i := 0; i < 9; i++ {
		h.advance()
	}

	// 读取 CDATA 内容，直到遇到 ]]>
	for h.pos+3 <= len(h.input) {
		if h.input[h.pos:h.pos+3] == "]]>" {
			// 跳过 "]]>"
			for i := 0; i < 3; i++ {
				h.advance()
			}
			break
		}
		h.advance()
	}

	return NewWorkerToken(
		token.STRING,
		h.input[start:h.pos],
		start,
		h.pos,
		startLine,
		startLinePos,
	), true
}

// processProcessingInstruction 处理处理指令 <?...?>（如 PHP 标签 <?php ... ?>）
// 返回 token 和是否成功处理
func (h *HtmlLexer) processProcessingInstruction() (Token, bool) {
	if h.pos+2 > len(h.input) || h.input[h.pos:h.pos+2] != "<?" {
		return nil, false
	}

	start := h.pos
	startLine := h.line
	startLinePos := h.linePos

	// 跳过 "<?"
	for i := 0; i < 2; i++ {
		h.advance()
	}

	// 读取处理指令内容，直到遇到 ?>
	for h.pos+2 <= len(h.input) {
		if h.input[h.pos:h.pos+2] == "?>" {
			// 跳过 "?>"
			for i := 0; i < 2; i++ {
				h.advance()
			}
			break
		}
		h.advance()
	}

	// 检查是否是 PHP 开始标签
	content := h.input[start:h.pos]
	if len(content) >= 5 && content[:5] == "<?php" {
		return NewWorkerToken(
			token.START_TAG,
			content,
			start,
			h.pos,
			startLine,
			startLinePos,
		), true
	}

	return NewWorkerToken(
		token.STRING,
		content,
		start,
		h.pos,
		startLine,
		startLinePos,
	), true
}

// processHtmlComment 处理 HTML 注释 <!-- ... -->
// 返回 token 和是否成功处理
func (h *HtmlLexer) processHtmlComment() (Token, bool) {
	if h.pos+4 > len(h.input) || h.input[h.pos:h.pos+4] != "<!--" {
		return nil, false
	}

	start := h.pos
	startLine := h.line
	startLinePos := h.linePos

	// 跳过 "<!--"
	for i := 0; i < 4; i++ {
		h.advance()
	}

	// 读取注释内容，直到遇到 -->
	for h.pos+3 <= len(h.input) {
		if h.input[h.pos:h.pos+3] == "-->" {
			// 跳过 "-->"
			for i := 0; i < 3; i++ {
				h.advance()
			}
			break
		}
		h.advance()
	}

	return NewWorkerToken(
		token.MULTILINE_COMMENT,
		h.input[start:h.pos],
		start,
		h.pos,
		startLine,
		startLinePos,
	), true
}

// processDoubleQuotedString 处理双引号字符串 "..."（用于 HTML 属性值）
// 返回 token 和是否成功处理
func (h *HtmlLexer) processDoubleQuotedString() (Token, bool) {
	if h.pos >= len(h.input) || h.input[h.pos] != '"' {
		return nil, false
	}

	start := h.pos
	startLine := h.line
	startLinePos := h.linePos

	// 读取开始引号
	h.advance()

	// 读取字符串内容，直到遇到未转义的结束引号
	escaped := false
	for h.pos < len(h.input) {
		if !escaped && h.input[h.pos] == '"' {
			// 找到字符串结束，包含结束引号
			h.advance()
			return NewWorkerToken(
				token.STRING,
				h.input[start:h.pos],
				start,
				h.pos,
				startLine,
				startLinePos,
			), true
		}

		// 处理转义字符
		if h.input[h.pos] == '\\' {
			escaped = !escaped
		} else {
			escaped = false
		}
		h.advance()
	}

	// 如果没有找到结束引号，返回已读取的部分
	return NewWorkerToken(
		token.STRING,
		h.input[start:h.pos],
		start,
		h.pos,
		startLine,
		startLinePos,
	), true
}

// processSingleQuotedString 处理单引号字符串 '...'（用于 HTML 属性值）
// 返回 token 和是否成功处理
func (h *HtmlLexer) processSingleQuotedString() (Token, bool) {
	if h.pos >= len(h.input) || h.input[h.pos] != '\'' {
		return nil, false
	}

	start := h.pos
	startLine := h.line
	startLinePos := h.linePos

	// 读取开始引号
	h.advance()

	// 读取字符串内容，直到遇到未转义的结束引号
	escaped := false
	for h.pos < len(h.input) {
		if !escaped && h.input[h.pos] == '\'' {
			// 找到字符串结束，包含结束引号
			h.advance()
			return NewWorkerToken(
				token.STRING,
				h.input[start:h.pos],
				start,
				h.pos,
				startLine,
				startLinePos,
			), true
		}

		// 处理转义字符
		if h.input[h.pos] == '\\' {
			escaped = !escaped
		} else {
			escaped = false
		}
		h.advance()
	}

	// 如果没有找到结束引号，返回已读取的部分
	return NewWorkerToken(
		token.STRING,
		h.input[start:h.pos],
		start,
		h.pos,
		startLine,
		startLinePos,
	), true
}

// processBacktickString 处理反引号字符串 `...`（用于 HTML 属性值）
// 返回 token 和是否成功处理
func (h *HtmlLexer) processBacktickString() (Token, bool) {
	if h.pos >= len(h.input) || h.input[h.pos] != '`' {
		return nil, false
	}

	start := h.pos
	startLine := h.line
	startLinePos := h.linePos

	// 读取开始反引号
	h.advance()

	// 读取字符串内容，直到遇到结束反引号
	for h.pos < len(h.input) {
		if h.input[h.pos] == '`' {
			// 找到字符串结束，包含结束反引号
			h.advance()
			return NewWorkerToken(
				token.STRING,
				h.input[start:h.pos],
				start,
				h.pos,
				startLine,
				startLinePos,
			), true
		}
		h.advance()
	}

	// 如果没有找到结束反引号，返回已读取的部分
	return NewWorkerToken(
		token.STRING,
		h.input[start:h.pos],
		start,
		h.pos,
		startLine,
		startLinePos,
	), true
}

// processJsServer 处理 $.SERVER 关键字
// 返回 token 和是否成功处理
func (h *HtmlLexer) processJsServer() (Token, bool) {
	if h.pos+8 > len(h.input) || h.input[h.pos:h.pos+8] != "$.SERVER" {
		return nil, false
	}

	start := h.pos
	startLine := h.line
	startLinePos := h.linePos

	// 跳过 "$.SERVER"（8个字符）
	for i := 0; i < 8; i++ {
		h.advance()
	}

	return NewWorkerToken(
		token.JS_SERVER,
		h.input[start:h.pos],
		start,
		h.pos,
		startLine,
		startLinePos,
	), true
}

// processNumber 处理数字字面量（整数或浮点数）
// 返回 token 和是否成功处理
func (h *HtmlLexer) processNumber() (Token, bool) {
	if h.pos >= len(h.input) {
		return nil, false
	}

	ch := h.input[h.pos]
	// 检查是否是数字或负号（用于负数）
	if !h.isDigit(ch) && ch != '-' && ch != '+' {
		return nil, false
	}

	start := h.pos
	startLine := h.line
	startLinePos := h.linePos

	// 处理符号
	if ch == '-' || ch == '+' {
		h.advance()
		// 符号后必须是数字
		if h.pos >= len(h.input) || !h.isDigit(h.input[h.pos]) {
			// 回退
			h.pos = start
			h.line = startLine
			h.linePos = startLinePos
			return nil, false
		}
	}

	// 读取整数部分
	hasDigit := false
	for h.pos < len(h.input) && h.isDigit(h.input[h.pos]) {
		hasDigit = true
		h.advance()
	}

	// 检查是否有小数点（浮点数）
	if h.pos < len(h.input) && h.input[h.pos] == '.' {
		h.advance()
		// 读取小数部分
		for h.pos < len(h.input) && h.isDigit(h.input[h.pos]) {
			hasDigit = true
			h.advance()
		}
	}

	// 检查是否有指数部分（科学计数法）
	if h.pos < len(h.input) && (h.input[h.pos] == 'e' || h.input[h.pos] == 'E') {
		h.advance()
		// 指数可以有符号
		if h.pos < len(h.input) && (h.input[h.pos] == '+' || h.input[h.pos] == '-') {
			h.advance()
		}
		// 读取指数数字
		for h.pos < len(h.input) && h.isDigit(h.input[h.pos]) {
			hasDigit = true
			h.advance()
		}
	}

	// 如果收集到了有效的数字
	if hasDigit {
		// 确定 token 类型（整数或浮点数）
		tokenType := token.INT
		literal := h.input[start:h.pos]
		// 检查是否包含小数点或指数
		for i := 0; i < len(literal); i++ {
			if literal[i] == '.' || literal[i] == 'e' || literal[i] == 'E' {
				tokenType = token.FLOAT
				break
			}
		}

		return NewWorkerToken(
			tokenType,
			literal,
			start,
			h.pos,
			startLine,
			startLinePos,
		), true
	}

	// 回退
	h.pos = start
	h.line = startLine
	h.linePos = startLinePos
	return nil, false
}

// processGt 处理单独的 > 符号
// 返回 token 和是否成功处理
func (h *HtmlLexer) processGt() (Token, bool) {
	if h.pos >= len(h.input) || h.input[h.pos] != '>' {
		return nil, false
	}

	start := h.pos
	startLine := h.line
	startLinePos := h.linePos

	h.advance()
	// 标记离开标签内部
	h.inTag = false

	return NewWorkerToken(
		token.GT,
		h.input[start:h.pos],
		start,
		h.pos,
		startLine,
		startLinePos,
	), true
}

// processLt 处理单独的 < 符号
// 返回 token 和是否成功处理
func (h *HtmlLexer) processLt() (Token, bool) {
	if h.pos >= len(h.input) || h.input[h.pos] != '<' {
		return nil, false
	}

	// 检查是否是已知的特殊结构
	if h.pos+1 < len(h.input) {
		next := h.input[h.pos+1]
		// 如果是 !、?，应该由其他函数处理
		if next == '!' || next == '?' {
			return nil, false
		}
		// 如果是 /，这是结束标签，需要处理
		// 结束标签的 < 也需要被识别为 LT token
	}

	// 处理单独的 < 符号（包括标签开头的 < 和结束标签的 <）
	start := h.pos
	startLine := h.line
	startLinePos := h.linePos

	h.advance()

	// 检查是否是结束标签 </tag>
	if h.pos < len(h.input) && h.input[h.pos] == '/' {
		// 结束标签，不设置 inTag = true
		// / 符号会由 processQuo 处理
	} else {
		// 开始标签，标记进入标签内部
		h.inTag = true
	}

	return NewWorkerToken(
		token.LT,
		h.input[start:h.pos],
		start,
		h.pos,
		startLine,
		startLinePos,
	), true
}

// processAssign 处理 = 符号
// 返回 token 和是否成功处理
func (h *HtmlLexer) processAssign() (Token, bool) {
	if h.pos >= len(h.input) || h.input[h.pos] != '=' {
		return nil, false
	}

	start := h.pos
	startLine := h.line
	startLinePos := h.linePos

	h.advance()

	return NewWorkerToken(
		token.ASSIGN,
		h.input[start:h.pos],
		start,
		h.pos,
		startLine,
		startLinePos,
	), true
}

// processQuo 处理 / 符号
// 返回 token 和是否成功处理
func (h *HtmlLexer) processQuo() (Token, bool) {
	if h.pos >= len(h.input) || h.input[h.pos] != '/' {
		return nil, false
	}

	start := h.pos
	startLine := h.line
	startLinePos := h.linePos

	h.advance()

	return NewWorkerToken(
		token.QUO,
		h.input[start:h.pos],
		start,
		h.pos,
		startLine,
		startLinePos,
	), true
}

// processIdentifier 处理标识符（标签名、属性名等）
// 返回 token 和是否成功处理
func (h *HtmlLexer) processIdentifier() (Token, bool) {
	if h.pos >= len(h.input) {
		return nil, false
	}

	ch := h.input[h.pos]
	// 检查是否是字母、下划线或中文字符开头
	if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_') {
		return nil, false
	}

	start := h.pos
	startLine := h.line
	startLinePos := h.linePos

	// 读取标识符字符
	for h.pos < len(h.input) {
		ch := h.input[h.pos]
		// 标识符可以包含字母、数字、下划线、连字符、冒号（用于 XML 命名空间）
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') || ch == '_' || ch == '-' || ch == ':' {
			h.advance()
		} else {
			break
		}
	}

	return NewWorkerToken(
		token.IDENTIFIER,
		h.input[start:h.pos],
		start,
		h.pos,
		startLine,
		startLinePos,
	), true
}

// processWhitespace 处理空白字符（空格、制表符、换行符等）
// 返回 token 和是否成功处理
func (h *HtmlLexer) processWhitespace() (Token, bool) {
	if h.pos >= len(h.input) || !h.isWhitespace(h.input[h.pos]) {
		return nil, false
	}

	start := h.pos
	startLine := h.line
	startLinePos := h.linePos

	// 收集连续的空白字符
	for h.pos < len(h.input) && h.isWhitespace(h.input[h.pos]) {
		h.advance()
	}

	return NewWorkerToken(
		token.WHITESPACE,
		h.input[start:h.pos],
		start,
		h.pos,
		startLine,
		startLinePos,
	), true
}

// processText 处理文本内容（标签之间的内容，支持插值）
// 返回 token 和是否成功处理
func (h *HtmlLexer) processText() (Token, bool) {
	// 只在标签外部处理文本内容
	if h.inTag {
		return nil, false
	}

	if h.pos >= len(h.input) {
		return nil, false
	}

	start := h.pos
	startLine := h.line
	startLinePos := h.linePos

	// 读取直到遇到 <（标签开始），将所有内容作为整体字符串
	// 不停止于 >、引号等，因为这些在标签外部都是文本内容的一部分
	for h.pos < len(h.input) && h.input[h.pos] != '<' {
		h.advance()
	}

	// 如果收集到了文本内容
	if start < h.pos {
		text := h.input[start:h.pos]
		// 检查是否只包含空白字符
		if h.isOnlyWhitespace(text) {
			return NewWorkerToken(
				token.WHITESPACE,
				text,
				start,
				h.pos,
				startLine,
				startLinePos,
			), true
		}
		// 处理文本插值
		return h.processTextInterpolation(text, start, startLine, startLinePos), true
	}

	return nil, false
}

// processTextInterpolation 处理文本中的插值（{$...} 和 @{...}）
// 参考 preprocessor.go:192 的 processStringInterpolation 实现，但不处理引号
func (h *HtmlLexer) processTextInterpolation(text string, start, startLine, startLinePos int) Token {
	var children []Token // 用于存储插值块内的子 token
	var currentStr []rune
	runes := []rune(text)
	hasInterpolation := false // 标记是否有插值

	// 辅助函数：将 rune 索引转换为字节位置
	runeToBytePos := func(runeIdx int) int {
		if runeIdx <= 0 {
			return 0
		}
		if runeIdx >= len(runes) {
			return len(text)
		}
		return len(string(runes[:runeIdx]))
	}

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '{' && i+2 < len(runes) && runes[i+1] == '$' {
			// 检查 $ 后面是否是有效的变量名起始字符
			nextChar := runes[i+2]
			// 变量名必须以字母、下划线或中文字符开头，不能是数字或特殊符号
			if !isValidVarChar(nextChar) || unicode.IsDigit(nextChar) {
				// 如果 $ 后面不是有效的变量名起始字符，将 { 和 $ 都作为普通字符处理
				currentStr = append(currentStr, r)
				currentStr = append(currentStr, runes[i+1])
				i++ // 跳过 $ 字符，下次循环会处理 $ 后面的字符
				continue
			}

			hasInterpolation = true
			// 处理变量插值
			if len(currentStr) > 0 {
				// 添加当前字符串（不带引号）
				strStart := i - len(currentStr)
				children = append(children, NewWorkerToken(
					token.STRING,
					string(currentStr),
					start+runeToBytePos(strStart),
					start+runeToBytePos(i),
					startLine,
					startLinePos+runeToBytePos(strStart),
				))
				currentStr = nil
			}

			// 如果还没有任何 children，添加空字符串
			if len(children) == 0 {
				children = append(children, NewWorkerToken(
					token.STRING,
					"",
					start,
					start,
					startLine,
					startLinePos,
				))
			}

			// 收集{$...}中的完整表达式内容（支持方法调用等复杂表达式）
			exprStart := i + 2
			j := exprStart
			braceDepth := 1 // 从 { 开始，深度为1
			parenDepth := 0
			bracketDepth := 0

			for j < len(runes) {
				if runes[j] == '{' {
					braceDepth++
				} else if runes[j] == '}' {
					braceDepth--
					if braceDepth == 0 {
						break
					}
				} else if runes[j] == '(' {
					parenDepth++
				} else if runes[j] == ')' {
					parenDepth--
				} else if runes[j] == '[' {
					bracketDepth++
				} else if runes[j] == ']' {
					bracketDepth--
				}
				j++
			}

			if j < len(runes) && runes[j] == '}' && braceDepth == 0 {
				// 找到了匹配的 }，提取表达式内容
				exprContent := string(runes[exprStart:j])

				// 复杂表达式，需要重新分词
				code := "$" + exprContent
				l := NewLexer()
				codeTokens := l.Tokenize(code)
				// 将分词结果添加到children中，并调整位置信息
				baseStart := start + runeToBytePos(i+1) // { 的位置 + 1 是 $ 的位置
				values := make([]Token, 0)
				for _, codeToken := range codeTokens {
					// 创建新的 WorkerToken 并调整位置
					values = append(values, NewWorkerToken(
						codeToken.Type(),
						codeToken.Literal(),
						codeToken.Start()+baseStart,
						codeToken.End()+baseStart,
						startLine,
						startLinePos+runeToBytePos(i+1)+codeToken.Start(),
					))
				}
				children = append(children, NewLingToken(
					token.INTERPOLATION_VALUE,
					code,
					start+runeToBytePos(j),
					start+runeToBytePos(j+1),
					startLine,
					startLinePos+runeToBytePos(j),
					values,
				))
				i = j
				continue
			}
			// 如果没有找到匹配的 }，将 { 和 $ 作为普通字符处理
			currentStr = append(currentStr, r)
			currentStr = append(currentStr, runes[i+1])
			i++ // 跳过 $ 字符
			continue
		} else if r == '@' && i+2 < len(runes) && runes[i+1] == '{' {
			hasInterpolation = true
			// 处理函数插值
			if len(currentStr) > 0 {
				// 添加当前字符串（不带引号）
				strStart := i - len(currentStr)
				children = append(children, NewWorkerToken(
					token.STRING,
					string(currentStr),
					start+runeToBytePos(strStart),
					start+runeToBytePos(i),
					startLine,
					startLinePos+runeToBytePos(strStart),
				))
				currentStr = nil
			}

			// 收集@{...}中的内容
			exprStart := i + 2
			j := exprStart
			braceCount := 0
			for j < len(runes) {
				if runes[j] == '{' {
					braceCount++
				} else if runes[j] == '}' {
					if braceCount == 0 {
						break
					}
					braceCount--
				}
				j++
			}
			if j < len(runes) && runes[j] == '}' {
				// 对@{...}中的内容进行重新分词
				code := string(runes[exprStart:j])
				l := NewLexer()
				codeTokens := l.Tokenize(code)
				// 将分词结果添加到children中，并调整位置信息
				baseStart := start + runeToBytePos(i+2) // @{ 之后的位置
				values := make([]Token, 0)
				for _, codeToken := range codeTokens {
					// 创建新的 WorkerToken 并调整位置
					values = append(values, NewWorkerToken(
						codeToken.Type(),
						codeToken.Literal(),
						codeToken.Start()+baseStart,
						codeToken.End()+baseStart,
						startLine,
						startLinePos+runeToBytePos(i+2)+codeToken.Start(),
					))
				}
				children = append(children, NewLingToken(
					token.INTERPOLATION_VALUE,
					code,
					start+runeToBytePos(j),
					start+runeToBytePos(j+1),
					startLine,
					startLinePos+runeToBytePos(j),
					values,
				))
				i = j
				continue
			}
		} else {
			currentStr = append(currentStr, r)
		}
	}

	// 添加剩余的字符串
	if len(currentStr) > 0 {
		if hasInterpolation {
			// 如果有插值，添加到children中（不带引号）
			strStart := len(runes) - len(currentStr)
			children = append(children, NewWorkerToken(
				token.STRING,
				string(currentStr),
				start+runeToBytePos(strStart),
				start+len(text),
				startLine,
				startLinePos+runeToBytePos(strStart),
			))
		}
	}

	// 如果有插值，创建 LingToken
	if hasInterpolation {
		// 如果没有任何children，说明是空字符串，添加一个空字符串token
		if len(children) == 0 {
			children = append(children, NewWorkerToken(
				token.STRING,
				"",
				start,
				start,
				startLine,
				startLinePos,
			))
		} else {
			// 处理特殊情况："{$data}/other" -> 移除开头的空字符串
			if len(children) >= 1 && children[0].Literal() == "" {
				children = children[1:]
			}
		}
		// 创建 LingToken 包含所有子 token
		return NewLingToken(
			token.INTERPOLATION_TOKEN,
			text,
			start,
			start+len(text),
			startLine,
			startLinePos,
			children,
		)
	}

	// 如果没有插值，返回普通字符串token
	if len(currentStr) > 0 {
		return NewWorkerToken(
			token.STRING,
			string(currentStr),
			start,
			start+len(text),
			startLine,
			startLinePos,
		)
	}

	// 空字符串
	return NewWorkerToken(
		token.STRING,
		"",
		start,
		start,
		startLine,
		startLinePos,
	)
}

// advance 前进一个字符，并更新行号和位置信息
func (h *HtmlLexer) advance() {
	if h.pos >= len(h.input) {
		return
	}

	if h.input[h.pos] == '\n' {
		h.line++
		h.linePos = 0
	} else {
		h.linePos++
	}
	h.pos++
}

// isWhitespace 检查字符是否是空白字符（包括换行符）
func (h *HtmlLexer) isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n'
}

// isOnlyWhitespace 检查字符串是否只包含空白字符
func (h *HtmlLexer) isOnlyWhitespace(s string) bool {
	for i := 0; i < len(s); i++ {
		if !h.isWhitespace(s[i]) {
			// 检查是否是全角空格
			if i+2 < len(s) && s[i] == 0xe3 && s[i+1] == 0x80 && s[i+2] == 0x80 {
				i += 2
				continue
			}
			return false
		}
	}
	return true
}

// isDigit 检查字符是否是数字
func (h *HtmlLexer) isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
