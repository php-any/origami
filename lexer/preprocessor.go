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
	switch t.Type() {
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
	switch t.Type() {
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
		switch t.Type() {
		case token.WHITESPACE, token.COMMENT, token.MULTILINE_COMMENT:
			continue
		case token.STRING:
			// 2. 字符串插值
			filtered = append(filtered, processStringInterpolation(t))
		case token.DOLLAR:
			// 处理$标识符组合
			if i+1 < len(p.tokens) && (p.tokens[i+1].Type() == token.IDENTIFIER || (p.tokens[i+1].Type() >= token.KEYWORD_START && p.tokens[i+1].Type() <= token.KEYWORD_END)) ||
				p.tokens[i+1].Type() == token.NULL || // 添加对null的支持
				p.tokens[i+1].Type() == token.TRUE || // 添加对true的支持
				p.tokens[i+1].Type() == token.FALSE { // 添加对false的支持

				// 将$和标识符合并为一个变量token，保留$符号
				next := p.tokens[i+1]
				filtered = append(filtered, NewWorkerToken(
					token.VARIABLE,
					"$"+next.Literal(),
					t.Start(),
					next.End(),
					next.Line(),
					next.Pos(),
				))
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
		if t.Type() == token.NEWLINE {
			// 检查前一个token是否需要补分号
			if i > 0 && !cannotAddSemicolon(filtered[i-1]) {
				// 检查后一个token是否需要补分号
				if i+1 < len(filtered) && !cannotAddSemicolonAfter(filtered[i+1]) {
					// 将换行符替换为分号，保持原有位置信息但不修改 Literal
					semicolon := NewWorkerToken(
						token.SEMICOLON,
						t.Literal(), // 保持原始 Literal 值（换行符）
						t.Start(),
						t.End(),
						t.Line(),
						t.Pos(),
					)
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
		if t.Type() == token.IDENTIFIER {
			if (i+1) < len(result) && result[i+1].Type() == token.ASSIGN {
				if i > 2 {
					check := result[i-1].Type()
					for _, temp := range []token.TokenType{
						token.LBRACKET,  // [
						token.LBRACE,    // {
						token.LPAREN,    // (
						token.SEMICOLON, // ;
						token.COMMA,     // ,
					} {
						if check == temp {
							// 创建一个新的 WorkerToken 替换原来的 token
							result[i] = NewWorkerToken(
								token.VARIABLE,
								t.Literal(),
								t.Start(),
								t.End(),
								t.Line(),
								t.Pos(),
							)
						}
					}
				}
			}
		}
	}

	// 调试：打印所有处理后的 tokens
	PrintTokens(result, "分词后的 Token 列表")

	return result
}

// processStringInterpolation 处理字符串插值，如果有插值返回LingToken，否则返回WorkerToken
func processStringInterpolation(t Token) Token {
	literal := t.Literal()
	if len(literal) < 2 {
		return t
	}
	quote := literal[0]
	if literal[len(literal)-1] != quote {
		return t
	}
	content := literal[1 : len(literal)-1]
	var children []Token // 用于存储插值块内的子 token
	var currentStr []rune
	runes := []rune(content)
	hasInterpolation := false // 标记是否有插值

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
				// 添加当前字符串
				children = append(children, NewWorkerToken(
					token.STRING,
					string(quote)+string(currentStr)+string(quote),
					t.Start(),
					t.End(),
					t.Line(),
					t.Pos(),
				))
				currentStr = nil
			}

			// 如果还没有任何 children，添加空字符串
			if len(children) == 0 {
				children = append(children, NewWorkerToken(
					token.STRING,
					"",
					t.Start(),
					t.Start(),
					t.Line(),
					t.Pos(),
				))
			}

			// 收集{$...}中的完整表达式内容（支持方法调用等复杂表达式）
			start := i + 2
			j := start
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
				exprContent := string(runes[start:j])

				// 复杂表达式，需要重新分词
				code := "$" + exprContent
				l := NewLexer()
				codeTokens := l.Tokenize(code)
				// 将分词结果添加到children中，并调整位置信息
				baseStart := t.Start() + start - 1
				values := make([]Token, 0)
				for _, codeToken := range codeTokens {
					// 创建新的 WorkerToken 并调整位置
					values = append(values, NewWorkerToken(
						codeToken.Type(),
						codeToken.Literal(),
						codeToken.Start()+baseStart,
						codeToken.End()+baseStart,
						t.Line(),
						t.Pos()+start-1+(codeToken.Start()),
					))
				}
				children = append(children, NewLingToken(
					token.INTERPOLATION_VALUE,
					code,
					t.Start()+j,
					t.Start()+j,
					t.Line(),
					t.Pos()+j,
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
				// 添加当前字符串
				children = append(children, NewWorkerToken(
					token.STRING,
					string(quote)+string(currentStr)+string(quote),
					t.Start(),
					t.End(),
					t.Line(),
					t.Pos(),
				))
				currentStr = nil
			}

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
				// 将分词结果添加到children中，并调整位置信息
				baseStart := t.Start() + start
				values := make([]Token, 0)
				for _, codeToken := range codeTokens {
					// 创建新的 WorkerToken 并调整位置
					values = append(values, NewWorkerToken(
						codeToken.Type(),
						codeToken.Literal(),
						codeToken.Start()+baseStart,
						codeToken.End()+baseStart,
						t.Line(),
						t.Pos()+start+codeToken.Start(),
					))
				}
				children = append(children, NewLingToken(
					token.INTERPOLATION_VALUE,
					code,
					t.Start()+j,
					t.Start()+j,
					t.Line(),
					t.Pos()+j,
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
			// 如果有插值，添加到children中
			children = append(children, NewWorkerToken(
				token.STRING,
				string(quote)+string(currentStr)+string(quote),
				t.Start(),
				t.End(),
				t.Line(),
				t.Pos(),
			))
		}
		// 如果没有插值，currentStr 会在后面处理
	}

	// 如果有插值，创建 LingToken
	if hasInterpolation {
		// 如果没有任何children，说明是空字符串，添加一个空字符串token
		if len(children) == 0 {
			children = append(children, NewWorkerToken(
				token.STRING,
				string(quote)+string(quote),
				t.Start(),
				t.End(),
				t.Line(),
				t.Pos(),
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
			literal,
			t.Start(),
			t.End(),
			t.Line(),
			t.Pos(),
			children,
		)
	}

	// 如果没有插值，返回普通字符串token
	if len(currentStr) > 0 {
		return NewWorkerToken(
			token.STRING,
			string(quote)+string(currentStr)+string(quote),
			t.Start(),
			t.End(),
			t.Line(),
			t.Pos(),
		)
	}

	// 空字符串
	return NewWorkerToken(
		token.STRING,
		string(quote)+string(quote),
		t.Start(),
		t.End(),
		t.Line(),
		t.Pos(),
	)
}

// isValidVarChar 检查是否是有效的变量名字符
func isValidVarChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || ('\u4e00' <= r && r <= '\u9fff') // 常见中文 Unicode 范围
}

// isNumber 检查是否是数字
func isNumber(r rune) bool {
	return unicode.IsDigit(r)
}

// isSpecialSymbol 检查是否是特殊符号
func isSpecialSymbol(r rune) bool {
	return r == '_' || unicode.IsPunct(r)
}
