package lexer

import (
	"unicode"

	"github.com/php-any/origami/token"
)

// Preprocessor 表示预处理器
type Preprocessor struct {
	tokens []Token
}

// NewPreprocessor 创建一个新的预处理器
func NewPreprocessor(tokens []Token) *Preprocessor {
	return &Preprocessor{
		tokens: tokens,
	}
}

// cannotAddSemicolon 判断前一个token后是否不能补分号
func cannotAddSemicolon(t Token) bool {
	switch t.Type {
	case token.SEMICOLON: // 前一个已经是分号，不用补充
		return true
	case token.COMMA: // 逗号后不用补充
		return true
	case token.DOT: // 点后不用补充
		return true
	case token.RBRACE: // 右花括号后不用补充
		return true
	case token.OBJECT_OPERATOR: // 箭头后不用补充
		return true
	case token.ADD: // 加号后不用补充
		return true
	case token.SUB: // 减号后不用补充
		return true
	case token.MUL: // 乘号后不用补充
		return true
	case token.QUO: // 除号后不用补充
		return true
	case token.REM: // 取模后不用补充
		return true
	case token.BIT_AND: // 按位与后不用补充
		return true
	case token.BIT_OR: // 按位或后不用补充
		return true
	case token.LAND: // 逻辑与后不用补充
		return true
	case token.LOR: // 逻辑或后不用补充
		return true
	case token.LBRACKET: // 左方括号后不用补充
		return true
	case token.LBRACE: // 左花括号后不用补充
		return true
	case token.LPAREN: // 左圆括号后不用补充
		return true
	default:
		return false // 其他情况需要补充分号
	}
}

// cannotAddSemicolonAfter 判断后一个token前是否不能补分号
func cannotAddSemicolonAfter(t Token) bool {
	switch t.Type {
	case token.LBRACKET: // 左方括号前不用补充
		return true
	case token.RBRACKET: // 右方括号前不用补充
		return true
	case token.LBRACE: // 左花括号前不用补充
		return true
	case token.RBRACE: // 右花括号前不用补充
		return true
	case token.LPAREN: // 左圆括号前不用补充
		return true
	case token.RPAREN: // 右圆括号前不用补充
		return true
	case token.ARRAY_KEY_VALUE:
		return true
	case token.OBJECT_OPERATOR:
		return true

	default:
		return false // 其他情况需要补充分号
	}
}

// Process 处理所有token，实现自动补分号、字符串插值、跳过无意义符号
// 识别$标识符为变量
func (p *Preprocessor) Process() []Token {
	var filtered []Token
	// 1. 跳过空白符和注释，处理$标识符
	for i := 0; i < len(p.tokens); i++ {
		t := p.tokens[i]
		switch t.Type {
		case token.WHITESPACE, token.COMMENT, token.MULTILINE_COMMENT:
			continue
		case token.STRING:
			// 2. 字符串插值
			tokens := processStringInterpolation(t)
			filtered = append(filtered, tokens...)
		case token.DOLLAR:
			// 处理$标识符组合
			if i+1 < len(p.tokens) && (p.tokens[i+1].Type == token.IDENTIFIER || (p.tokens[i+1].Type >= token.KEYWORD_START && p.tokens[i+1].Type <= token.KEYWORD_END)) ||
				p.tokens[i+1].Type == token.NULL || // 添加对null的支持
				p.tokens[i+1].Type == token.TRUE || // 添加对true的支持
				p.tokens[i+1].Type == token.FALSE { // 添加对false的支持

				// 将$和标识符合并为一个变量token，保留$符号
				next := p.tokens[i+1]
				filtered = append(filtered, Token{
					Type:    token.VARIABLE,
					Literal: "$" + next.Literal,
					Start:   t.Start,
					End:     next.End,
					Line:    next.Line,
					Pos:     next.Pos,
				})
				i++ // 跳过下一个token
			} else {
				filtered = append(filtered, t)
			}
		default:
			filtered = append(filtered, t)
		}
	}

	// 3. 自动补分号（TS风格）
	var result []Token
	for i := 0; i < len(filtered); i++ {
		t := filtered[i]
		if t.Type == token.NEWLINE {
			// 检查前一个token是否需要补分号
			if i > 0 && !cannotAddSemicolon(filtered[i-1]) {
				// 检查后一个token是否需要补分号
				if i+1 < len(filtered) && !cannotAddSemicolonAfter(filtered[i+1]) {
					// 将换行符替换为分号，保持原有位置信息但不修改 Literal
					semicolon := Token{
						Type:    token.SEMICOLON,
						Literal: t.Literal, // 保持原始 Literal 值（换行符）
						Start:   t.Start,
						End:     t.End,
						Line:    t.Line,
						Pos:     t.Pos,
					}
					result = append(result, semicolon)
				}
			}
			// 跳过换行符
			continue
		} else {
			result = append(result, t)
		}
	}

	// 检查标识符是否是变量
	for i, t := range result {
		if t.Type == token.IDENTIFIER {
			if (i+1) < len(result) && result[i+1].Type == token.ASSIGN {
				if i > 2 {
					check := result[i-1].Type
					for _, temp := range []token.TokenType{
						token.LBRACKET,  // [
						token.LBRACE,    // {
						token.LPAREN,    // (
						token.SEMICOLON, // ;
						token.COMMA,     // ,
					} {
						if check == temp {
							result[i].Type = token.VARIABLE
						}
					}
				}
			}
		}
	}

	return result
}

// processStringInterpolation 处理字符串插值，返回拆分后的token列表
func processStringInterpolation(t Token) []Token {
	if len(t.Literal) < 2 {
		return []Token{t}
	}
	quote := t.Literal[0]
	if t.Literal[len(t.Literal)-1] != quote {
		return []Token{t}
	}
	content := t.Literal[1 : len(t.Literal)-1]
	var tokens []Token
	var currentStr []rune
	runes := []rune(content)

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

			// 处理变量插值
			if len(currentStr) > 0 {
				if len(tokens) > 0 {
					tokens = append(tokens, Token{
						Type:    token.ADD,
						Literal: "+",
						Start:   t.Start + i,
						End:     t.Start + i + 1,
						Line:    t.Line,
						Pos:     t.Pos + i,
					})
				}
				// 添加当前字符串
				tokens = append(tokens, Token{
					Type:    token.STRING,
					Literal: string(quote) + string(currentStr) + string(quote),
					Start:   t.Start,
					End:     t.End,
					Line:    t.Line,
					Pos:     t.Pos,
				})
				currentStr = nil
			}

			// 添加加号
			if len(tokens) == 0 {
				tokens = append(tokens, Token{
					Type:    token.STRING,
					Literal: "",
					Start:   t.Start,
					End:     t.Start,
					Line:    t.Line,
					Pos:     t.Pos,
				})
			}
			tokens = append(tokens, Token{
				Type:    token.ADD,
				Literal: "+",
				Start:   t.Start + i,
				End:     t.Start + i + 1,
				Line:    t.Line,
				Pos:     t.Pos + i,
			})

			// 收集变量名
			start := i + 2
			j := start
			for j < len(runes) && isValidVarChar(runes[j]) {
				j++
			}
			if j < len(runes) && runes[j] == '}' {
				// 添加变量token，包含$前缀
				tokens = append(tokens, Token{
					Type:    token.VARIABLE,
					Literal: "$" + string(runes[start:j]),
					Start:   t.Start + start - 1, // -1 是因为要包含$符号
					End:     t.Start + j,
					Line:    t.Line,
					Pos:     t.Pos + start - 1,
				})
				i = j
				continue
			}
			// 如果没有找到匹配的 }，将 { 和 $ 作为普通字符处理
			currentStr = append(currentStr, r)
			currentStr = append(currentStr, runes[i+1])
			i++ // 跳过 $ 字符
			continue
		} else if r == '@' && i+2 < len(runes) && runes[i+1] == '{' {
			// 处理函数插值
			if len(currentStr) > 0 {
				if len(tokens) > 0 {
					tokens = append(tokens, Token{
						Type:    token.ADD,
						Literal: "+",
						Start:   t.Start + i,
						End:     t.Start + i + 1,
						Line:    t.Line,
						Pos:     t.Pos + i,
					})
				}
				// 添加当前字符串
				tokens = append(tokens, Token{
					Type:    token.STRING,
					Literal: string(quote) + string(currentStr) + string(quote),
					Start:   t.Start,
					End:     t.End,
					Line:    t.Line,
					Pos:     t.Pos,
				})
				currentStr = nil
			}

			// 添加加号
			tokens = append(tokens, Token{
				Type:    token.ADD,
				Literal: "+",
				Start:   t.Start + i,
				End:     t.Start + i + 1,
				Line:    t.Line,
				Pos:     t.Pos + i,
			})

			// 收集@{...}中的内容
			start := i + 2
			j := start
			parenCount := 0
			for j < len(runes) {
				if runes[j] == '{' {
					parenCount++
				} else if runes[j] == '}' {
					if parenCount == 0 {
						break
					}
					parenCount--
				}
				j++
			}
			if j < len(runes) && runes[j] == '}' {
				// 对@{...}中的内容进行重新分词
				code := string(runes[start:j])
				l := NewLexer()
				codeTokens := l.Tokenize(code)
				// 将分词结果添加到tokens中，并调整位置信息
				for _, codeToken := range codeTokens {
					codeToken.Start += t.Start + start
					codeToken.End += t.Start + start
					codeToken.Line = t.Line
					codeToken.Pos = t.Pos + start + (codeToken.Start - (t.Start + start))
					tokens = append(tokens, codeToken)
				}
				i = j
				continue
			}
		} else {
			currentStr = append(currentStr, r)
		}
	}

	// 添加剩余的字符串
	if len(currentStr) > 0 {
		if len(tokens) == 2 && tokens[0].Literal == "" && tokens[1].Type == token.ADD {
			// "{$data}"
			tokens = []Token{}
		} else {
			if len(tokens) > 1 && token.ADD != tokens[len(tokens)-1].Type {
				tokens = append(tokens, Token{
					Type:    token.ADD,
					Literal: "+",
					Start:   t.Start,
					End:     t.End,
					Line:    t.Line,
					Pos:     t.Pos,
				})
			}
		}
		tokens = append(tokens, Token{
			Type:    token.STRING,
			Literal: string(quote) + string(currentStr) + string(quote),
			Start:   t.Start,
			End:     t.End,
			Line:    t.Line,
			Pos:     t.Pos,
		})
	}

	// 如果没有任何token，说明是空字符串，添加一个空字符串token
	if len(tokens) == 0 {
		tokens = append(tokens, Token{
			Type:    token.STRING,
			Literal: string(quote) + string(quote),
			Start:   t.Start,
			End:     t.End,
			Line:    t.Line,
			Pos:     t.Pos,
		})
	}

	return tokens
}

// isValidVarChar 检查是否是有效的变量名字符
func isValidVarChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || ('\u4e00' <= r && r <= '\u9fff') // 常见中文 Unicode 范围
}

// isValidFuncChar 检查是否是有效的函数名字符
func isValidFuncChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

// isNumber 检查是否是数字
func isNumber(r rune) bool {
	return unicode.IsDigit(r)
}

// isSpecialSymbol 检查是否是特殊符号
func isSpecialSymbol(r rune) bool {
	return r == '_' || unicode.IsPunct(r)
}
