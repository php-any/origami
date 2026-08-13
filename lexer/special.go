package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/php-any/origami/token"
)

// SpecialToken 表示特殊符号的处理结果
type SpecialToken struct {
	Type    token.TokenType
	Literal string
	Length  int // 处理掉的字符长度
}

// SpecialTokenResult 表示特殊符号处理的完整结果
type SpecialTokenResult struct {
	Token      SpecialToken
	NewPos     int // 新的位置
	NewLine    int // 新的行号
	NewLinePos int // 新的行内位置
}

// isCommentStart 检查是否是注释开始
func isCommentStart(input string, start int) bool {
	if start+1 >= len(input) {
		return false
	}
	return input[start] == '/' && (input[start+1] == '/' || input[start+1] == '*')
}

// isNumberStart 检查是否是数字开始
func isNumberStart(r rune) bool {
	return unicode.IsDigit(r)
}

// handleComment 处理注释
func handleComment(input string, start int) (SpecialToken, int, bool) {
	if start+1 >= len(input) {
		return SpecialToken{}, start, false
	}

	next := rune(input[start+1])
	if next == '/' {
		// 单行注释
		pos := start + 2
		for pos < len(input) {
			r, size := utf8.DecodeRuneInString(input[pos:])
			if r == '\n' || r == '\r' {
				break
			}
			pos += size
		}
		return SpecialToken{
			Type:    token.COMMENT,
			Literal: input[start:pos],
			Length:  pos - start,
		}, pos, true
	} else if next == '*' {
		// 多行注释
		pos := start + 2
		for pos < len(input)-1 {
			if input[pos] == '*' && input[pos+1] == '/' {
				pos += 2
				return SpecialToken{
					Type:    token.MULTILINE_COMMENT,
					Literal: input[start:pos],
					Length:  pos - start,
				}, pos, true
			}
			pos++
		}
	}

	return SpecialToken{}, start, false
}

// handleNumber 处理数字
func handleNumber(input string, start int) (SpecialToken, int, bool) {
	pos := start
	literal := ""

	// 检查是否以数字开始（包括负号）
	if pos >= len(input) {
		return SpecialToken{}, start, false
	}

	// 处理负号
	if pos < len(input) && input[pos] == '-' {
		pos++
		// 检查负号后是否有数字
		if pos >= len(input) || !unicode.IsDigit(rune(input[pos])) {
			return SpecialToken{}, start, false
		}
	} else if !unicode.IsDigit(rune(input[pos])) {
		// 如果不是负号，必须是以数字开始
		return SpecialToken{}, start, false
	}

	// 匹配到分隔符才停止
	for pos < len(input) {
		r, size := utf8.DecodeRuneInString(input[pos:])
		if r == utf8.RuneError {
			return SpecialToken{}, start, false
		}

		// 科学计数法中的+和-不应被视为分隔符
		if IsDelimiter(r) && r != '.' && !(r == '+' || r == '-') {
			break
		}
		if (r == '+' || r == '-') && pos > start {
			prev := input[pos-1]
			if prev != 'e' && prev != 'E' {
				break
			}
		}

		if pos+2 < len(input) && input[pos+1] == '.' && input[pos+2] == '.' {
			pos += size
			break
		}

		if r == 'e' || r == 'E' {
			pos += size
			if pos < len(input) && (input[pos] == '+' || input[pos] == '-') {
				pos++
			}
			continue
		}

		pos += size
	}

	if pos <= start {
		return SpecialToken{}, start, false
	}

	literal = input[start:pos]

	// 判断是否包含非数字字符
	hasNonDigit := false
	for i := 0; i < len(literal); i++ {
		r := rune(literal[i])
		if !unicode.IsDigit(r) && r != '.' && r != 'e' && r != 'E' && r != '+' && r != '-' && r != 'x' && r != 'X' && r != 'b' && r != 'B' {
			hasNonDigit = true
			break
		}
	}

	// 如果包含非数字字符，返回 NUMBER 类型（而不是 STRING）
	if hasNonDigit {
		return SpecialToken{
			Type:    token.NUMBER,
			Literal: literal,
			Length:  pos - start,
		}, pos, true
	}

	// 判断数字类型
	// 1. 检查是否是十六进制
	if len(literal) > 2 && literal[0] == '0' && (literal[1] == 'x' || literal[1] == 'X') {
		// 验证十六进制格式
		for i := 2; i < len(literal); i++ {
			r := rune(literal[i])
			if !unicode.IsDigit(r) && (r < 'a' || r > 'f') && (r < 'A' || r > 'F') {
				return SpecialToken{
					Type:    token.NUMBER,
					Literal: literal,
					Length:  pos - start,
				}, pos, true
			}
		}
		// 十六进制格式正确，返回 NUMBER 类型
		return SpecialToken{
			Type:    token.NUMBER,
			Literal: literal,
			Length:  pos - start,
		}, pos, true
	}

	// 2. 检查是否是二进制
	if len(literal) > 2 && literal[0] == '0' && (literal[1] == 'b' || literal[1] == 'B') {
		// 验证二进制格式
		for i := 2; i < len(literal); i++ {
			if literal[i] != '0' && literal[i] != '1' {
				return SpecialToken{
					Type:    token.NUMBER,
					Literal: literal,
					Length:  pos - start,
				}, pos, true
			}
		}
		return SpecialToken{
			Type:    token.NUMBER,
			Literal: literal,
			Length:  pos - start,
		}, pos, true
	}

	// 3. 检查是否是浮点数或整数（先检查是否有小数点，避免将 0.0 误判为八进制）
	hasDot := false
	hasExp := false
	for i := 0; i < len(literal); i++ {
		r := rune(literal[i])
		if r == '.' {
			if hasDot || hasExp {
				return SpecialToken{
					Type:    token.NUMBER,
					Literal: literal,
					Length:  pos - start,
				}, pos, true
			}
			hasDot = true
		} else if r == 'e' || r == 'E' {
			if hasExp {
				return SpecialToken{
					Type:    token.NUMBER,
					Literal: literal,
					Length:  pos - start,
				}, pos, true
			}
			hasExp = true
			// 科学计数法中的 + 和 - 是有效的
			if i+1 < len(literal) && (literal[i+1] == '+' || literal[i+1] == '-') {
				i++
			}
		} else if !unicode.IsDigit(r) && r != '-' && r != '+' {
			return SpecialToken{
				Type:    token.NUMBER,
				Literal: literal,
				Length:  pos - start,
			}, pos, true
		}
	}

	// 科学计数法应该返回 NUMBER 类型，而不是 FLOAT
	if hasExp {
		return SpecialToken{
			Type:    token.NUMBER,
			Literal: literal,
			Length:  pos - start,
		}, pos, true
	}

	// 4. 如果包含小数点，返回 FLOAT 类型
	if hasDot {
		return SpecialToken{
			Type:    token.FLOAT,
			Literal: literal,
			Length:  pos - start,
		}, pos, true
	}

	// 5. 检查是否是八进制（只有在没有小数点的情况下才检查）
	if len(literal) > 1 && literal[0] == '0' {
		// 验证八进制格式
		for i := 1; i < len(literal); i++ {
			if literal[i] < '0' || literal[i] > '7' {
				return SpecialToken{
					Type:    token.NUMBER,
					Literal: literal,
					Length:  pos - start,
				}, pos, true
			}
		}
		return SpecialToken{
			Type:    token.NUMBER,
			Literal: literal,
			Length:  pos - start,
		}, pos, true
	}

	// 6. 默认返回整数类型
	return SpecialToken{
		Type:    token.INT,
		Literal: literal,
		Length:  pos - start,
	}, pos, true
}

// handleByte 处理字节字面量
func handleByte(input string, start int) (SpecialToken, int, bool) {
	if start+2 >= len(input) || input[start] != 'b' || input[start+1] != '\'' {
		return SpecialToken{}, start, false
	}

	pos := start + 2
	for pos < len(input) {
		if input[pos] == '\'' {
			pos++
			return SpecialToken{
				Type:    token.BYTE,
				Literal: input[start:pos],
				Length:  pos - start,
			}, pos, true
		}
		if input[pos] == '\\' {
			pos++
			if pos < len(input) {
				pos++
			}
		} else {
			pos++
		}
	}

	return SpecialToken{}, start, false
}

// HandleSpecialToken 处理特殊符号
func HandleSpecialToken(input string, start int, currentLine, currentLinePos int) (SpecialTokenResult, bool) {
	if start >= len(input) {
		return SpecialTokenResult{}, false
	}

	// 检查字符串开始
	if tk, newPos, ok := HandleString(input, start); ok {
		// 计算字符串中的换行符数量
		newlineCount := 0
		lastNewlinePos := -1
		for i := start; i < newPos; i++ {
			if input[i] == '\n' {
				newlineCount++
				lastNewlinePos = i
			}
		}

		newLine := currentLine + newlineCount
		var newLinePos int
		if lastNewlinePos >= 0 {
			newLinePos = newPos - lastNewlinePos - 1
		} else {
			newLinePos = currentLinePos + (newPos - start)
		}

		return SpecialTokenResult{
			Token:      tk,
			NewPos:     newPos,
			NewLine:    newLine,
			NewLinePos: newLinePos,
		}, true
	}

	// 检查字节字面量
	if tk, newPos, ok := handleByte(input, start); ok {
		// 字节字面量通常不包含换行符
		newLinePos := currentLinePos + (newPos - start)
		return SpecialTokenResult{
			Token:      tk,
			NewPos:     newPos,
			NewLine:    currentLine,
			NewLinePos: newLinePos,
		}, true
	}

	r, _ := utf8.DecodeRuneInString(input[start:])
	switch {
	case isCommentStart(input, start):
		return handleCommentWithLineInfo(input, start, currentLine, currentLinePos)
	case isNumberStart(r) || (r == '-' && start+1 < len(input) && unicode.IsDigit(rune(input[start+1]))):
		return handleNumberWithLineInfo(input, start, currentLine, currentLinePos)
	}

	return SpecialTokenResult{}, false
}

// handleCommentWithLineInfo 处理注释并返回行号信息
func handleCommentWithLineInfo(input string, start int, currentLine, currentLinePos int) (SpecialTokenResult, bool) {
	if start+1 >= len(input) {
		return SpecialTokenResult{}, false
	}

	next := rune(input[start+1])
	if next == '/' {
		// 单行注释
		pos := start + 2
		newlineCount := 0
		for pos < len(input) {
			r, size := utf8.DecodeRuneInString(input[pos:])
			if r == '\n' || r == '\r' {
				newlineCount++
				pos += size
				break
			}
			pos += size
		}

		newLine := currentLine + newlineCount
		var newLinePos int
		if newlineCount > 0 {
			newLinePos = 0 // 换行后从行首开始
		} else {
			newLinePos = currentLinePos + (pos - start)
		}

		return SpecialTokenResult{
			Token: SpecialToken{
				Type:    token.COMMENT,
				Literal: input[start:pos],
				Length:  pos - start,
			},
			NewPos:     pos,
			NewLine:    newLine,
			NewLinePos: newLinePos,
		}, true
	} else if next == '*' {
		// 多行注释
		pos := start + 2
		newlineCount := 0
		for pos < len(input)-1 {
			if input[pos] == '*' && input[pos+1] == '/' {
				pos += 2
				break
			}
			if input[pos] == '\n' {
				newlineCount++
			}
			pos++
		}

		newLine := currentLine + newlineCount
		var newLinePos int
		if newlineCount > 0 {
			// 计算最后一个换行符后的位置
			lastNewlinePos := -1
			for i := start; i < pos; i++ {
				if input[i] == '\n' {
					lastNewlinePos = i
				}
			}
			if lastNewlinePos >= 0 {
				newLinePos = pos - lastNewlinePos - 1
			} else {
				newLinePos = currentLinePos + (pos - start)
			}
		} else {
			newLinePos = currentLinePos + (pos - start)
		}

		return SpecialTokenResult{
			Token: SpecialToken{
				Type:    token.MULTILINE_COMMENT,
				Literal: input[start:pos],
				Length:  pos - start,
			},
			NewPos:     pos,
			NewLine:    newLine,
			NewLinePos: newLinePos,
		}, true
	}

	return SpecialTokenResult{}, false
}

// handleNumberWithLineInfo 处理数字并返回行号信息
func handleNumberWithLineInfo(input string, start int, currentLine, currentLinePos int) (SpecialTokenResult, bool) {
	// 数字通常不包含换行符，所以行号不变
	if tk, newPos, ok := handleNumber(input, start); ok {
		newLinePos := currentLinePos + (newPos - start)
		return SpecialTokenResult{
			Token:      tk,
			NewPos:     newPos,
			NewLine:    currentLine,
			NewLinePos: newLinePos,
		}, true
	}
	return SpecialTokenResult{}, false
}
