package lexer

import (
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/php-any/origami/token"
)

var input = `namespace tests\func;

class A
{
    public function hello()
    {
        $name = 123;

        return "hello world " + $name;
    }
}

$a = new A();
// 无法跳转
echo $a->hello();`

func TestTokenPositions(t *testing.T) {
	// 直接使用测试内容

	// 创建词法分析器
	lexer := NewLexer()
	tokens := lexer.Tokenize(input)

	// 打印所有 Token 信息
	for i, tok := range tokens {
		t.Logf("Token %d: Type=%d, Literal=%q, Start=%d, End=%d, Line=%d, Pos=%d", i, tok.Type, tok.Literal, tok.Start, tok.End, tok.Line, tok.Pos)
	}

	// 验证位置信息的合理性
	for i, tok := range tokens {
		// 检查行号是否在合理范围内
		if tok.Line < 0 {
			t.Errorf("Token %d 行号为负数: %d", i, tok.Line)
		}
		// 检查列号是否在合理范围内
		if tok.Pos < 0 {
			t.Errorf("Token %d 列号为负数: %d", i, tok.Pos)
		}
		// 检查位置范围是否合理
		if tok.Start < 0 || tok.End < 0 || tok.Start >= tok.End {
			t.Errorf("Token %d 位置范围不合理: Start=%d, End=%d", i, tok.Start, tok.End)
		}
	}

	// 验证关键token的位置信息
	// 查找关键token并验证位置
	var namespaceToken, classToken, echoToken, dollarAToken *Token

	for i := range tokens {
		tok := &tokens[i]
		switch tok.Literal {
		case "namespace":
			namespaceToken = tok
		case "class":
			classToken = tok
		case "echo":
			echoToken = tok
		case "$a":
			dollarAToken = tok
		}
	}

	// 验证namespace token位置
	if namespaceToken != nil {
		if namespaceToken.Line != 0 {
			t.Errorf("namespace token 行号错误: 期望 0, 实际 %d", namespaceToken.Line)
		}
		if namespaceToken.Pos != 0 {
			t.Errorf("namespace token 列号错误: 期望 0, 实际 %d", namespaceToken.Pos)
		}
		if namespaceToken.Start != 0 {
			t.Errorf("namespace token Start位置错误: 期望 0, 实际 %d", namespaceToken.Start)
		}
		if namespaceToken.End != 9 {
			t.Errorf("namespace token End位置错误: 期望 9, 实际 %d", namespaceToken.End)
		}
	} else {
		t.Error("未找到 namespace token")
	}

	// 验证class token位置
	if classToken != nil {
		if classToken.Line != 2 {
			t.Errorf("class token 行号错误: 期望 2, 实际 %d", classToken.Line)
		}
		if classToken.Pos != 0 {
			t.Errorf("class token 列号错误: 期望 0, 实际 %d", classToken.Pos)
		}
		if classToken.Start != 23 {
			t.Errorf("class token Start位置错误: 期望 23, 实际 %d", classToken.Start)
		}
		if classToken.End != 28 {
			t.Errorf("class token End位置错误: 期望 28, 实际 %d", classToken.End)
		}
	} else {
		t.Error("未找到 class token")
	}

	// 验证echo token位置
	if echoToken != nil {
		if echoToken.Line != 14 {
			t.Errorf("echo token 行号错误: 期望 14, 实际 %d", echoToken.Line)
		}
		if echoToken.Pos != 0 {
			t.Errorf("echo token 列号错误: 期望 0, 实际 %d", echoToken.Pos)
		}
		if echoToken.Start != 167 {
			t.Errorf("echo token Start位置错误: 期望 167, 实际 %d", echoToken.Start)
		}
		if echoToken.End != 171 {
			t.Errorf("echo token End位置错误: 期望 171, 实际 %d", echoToken.End)
		}
	} else {
		t.Error("未找到 echo token")
	}

	// 验证$a token位置
	if dollarAToken != nil {
		if dollarAToken.Line != 14 {
			t.Errorf("$a token 行号错误: 期望 14, 实际 %d", dollarAToken.Line)
		}
		if dollarAToken.Pos != 6 {
			t.Errorf("$a token 列号错误: 期望 6, 实际 %d", dollarAToken.Pos)
		}
		if dollarAToken.Start != 172 {
			t.Errorf("$a token Start位置错误: 期望 172, 实际 %d", dollarAToken.Start)
		}
		if dollarAToken.End != 174 {
			t.Errorf("$a token End位置错误: 期望 174, 实际 %d", dollarAToken.End)
		}
	} else {
		t.Error("未找到 $a token")
	}
}

func TestRawTokenPositions(t *testing.T) {
	// 创建词法分析器并获取原始 Token
	lexer := NewLexer()
	tokens := tokenizeRaw(lexer, input)

	// 打印所有原始 Token 信息
	for i, tok := range tokens {
		t.Logf("Raw Token %d: Type=%d, Literal=%q, Start=%d, End=%d, Line=%d, Pos=%d", i, tok.Type, tok.Literal, tok.Start, tok.End, tok.Line, tok.Pos)
	}

	// 验证位置信息的合理性
	for i, tok := range tokens {
		// 检查行号是否在合理范围内
		if tok.Line < 0 {
			t.Errorf("Raw Token %d 行号为负数: %d", i, tok.Line)
		}
		// 检查列号是否在合理范围内
		if tok.Pos < 0 {
			t.Errorf("Raw Token %d 列号为负数: %d", i, tok.Pos)
		}
		// 检查位置范围是否合理
		if tok.Start < 0 || tok.End < 0 || tok.Start >= tok.End {
			t.Errorf("Raw Token %d 位置范围不合理: Start=%d, End=%d", i, tok.Start, tok.End)
		}
	}

	// 验证关键token的位置信息
	// 查找关键token并验证位置
	var namespaceToken, classToken, echoToken, dollarToken, aToken *Token

	for i := range tokens {
		tok := &tokens[i]
		switch tok.Literal {
		case "namespace":
			namespaceToken = tok
		case "class":
			classToken = tok
		case "echo":
			echoToken = tok
		case "$":
			// 找到第14行的$ token
			if tok.Line == 14 && tok.Pos == 5 {
				dollarToken = tok
			}
		case "a":
			// 找到第14行的a token
			if tok.Line == 14 && tok.Pos == 6 {
				aToken = tok
			}
		}
	}

	// 验证namespace token位置
	if namespaceToken != nil {
		if namespaceToken.Line != 0 {
			t.Errorf("namespace token 行号错误: 期望 0, 实际 %d", namespaceToken.Line)
		}
		if namespaceToken.Pos != 0 {
			t.Errorf("namespace token 列号错误: 期望 0, 实际 %d", namespaceToken.Pos)
		}
		if namespaceToken.Start != 0 {
			t.Errorf("namespace token Start位置错误: 期望 0, 实际 %d", namespaceToken.Start)
		}
		if namespaceToken.End != 9 {
			t.Errorf("namespace token End位置错误: 期望 9, 实际 %d", namespaceToken.End)
		}
	} else {
		t.Error("未找到 namespace token")
	}

	// 验证class token位置
	if classToken != nil {
		if classToken.Line != 2 {
			t.Errorf("class token 行号错误: 期望 2, 实际 %d", classToken.Line)
		}
		if classToken.Pos != 0 {
			t.Errorf("class token 列号错误: 期望 0, 实际 %d", classToken.Pos)
		}
		if classToken.Start != 23 {
			t.Errorf("class token Start位置错误: 期望 23, 实际 %d", classToken.Start)
		}
		if classToken.End != 28 {
			t.Errorf("class token End位置错误: 期望 28, 实际 %d", classToken.End)
		}
	} else {
		t.Error("未找到 class token")
	}

	// 验证echo token位置
	if echoToken != nil {
		if echoToken.Line != 14 {
			t.Errorf("echo token 行号错误: 期望 14, 实际 %d", echoToken.Line)
		}
		if echoToken.Pos != 0 {
			t.Errorf("echo token 列号错误: 期望 0, 实际 %d", echoToken.Pos)
		}
		if echoToken.Start != 167 {
			t.Errorf("echo token Start位置错误: 期望 167, 实际 %d", echoToken.Start)
		}
		if echoToken.End != 171 {
			t.Errorf("echo token End位置错误: 期望 171, 实际 %d", echoToken.End)
		}
	} else {
		t.Error("未找到 echo token")
	}

	// 验证$ token位置（第14行）
	if dollarToken != nil {
		if dollarToken.Line != 14 {
			t.Errorf("$ token 行号错误: 期望 14, 实际 %d", dollarToken.Line)
		}
		if dollarToken.Pos != 5 {
			t.Errorf("$ token 列号错误: 期望 5, 实际 %d", dollarToken.Pos)
		}
		if dollarToken.Start != 172 {
			t.Errorf("$ token Start位置错误: 期望 172, 实际 %d", dollarToken.Start)
		}
		if dollarToken.End != 173 {
			t.Errorf("$ token End位置错误: 期望 173, 实际 %d", dollarToken.End)
		}
	} else {
		t.Error("未找到第14行的 $ token")
	}

	// 验证a token位置（第14行）
	if aToken != nil {
		if aToken.Line != 14 {
			t.Errorf("a token 行号错误: 期望 14, 实际 %d", aToken.Line)
		}
		if aToken.Pos != 6 {
			t.Errorf("a token 列号错误: 期望 6, 实际 %d", aToken.Pos)
		}
		if aToken.Start != 173 {
			t.Errorf("a token Start位置错误: 期望 173, 实际 %d", aToken.Start)
		}
		if aToken.End != 174 {
			t.Errorf("a token End位置错误: 期望 174, 实际 %d", aToken.End)
		}
	} else {
		t.Error("未找到第14行的 a token")
	}
}

// tokenizeRaw 返回未经预处理器处理的原始 Token
func tokenizeRaw(lexer *Lexer, input string) []Token {
	var tokens []Token
	pos := 0
	line := 0    // 从0开始
	linePos := 0 // 从0开始

	for pos < len(input) {
		// 跳过空白符但记录位置
		if unicode.IsSpace(rune(input[pos])) {
			if input[pos] == '\n' {
				line++
				linePos = 0 // 从0开始
			} else {
				linePos++
			}
			pos++
			continue
		}

		// 尝试处理特殊 Token
		if specialToken, found := HandleSpecialToken(input, pos, 1, pos); found {
			tok := Token{
				Type:    specialToken.Token.Type,
				Literal: specialToken.Token.Literal,
				Start:   pos,
				End:     pos + len(specialToken.Token.Literal),
				Line:    line,
				Pos:     linePos,
			}
			tokens = append(tokens, tok)

			// 更新位置
			pos = specialToken.NewPos
			line = specialToken.NewLine
			linePos = specialToken.NewLinePos
			continue
		}

		// 尝试匹配最长 Token
		if tokenDef, newPos, found := lexer.matchLongestToken(input, pos); found && newPos > pos {
			tok := Token{
				Type:    tokenDef.Type,
				Literal: input[pos:newPos],
				Start:   pos,
				End:     newPos,
				Line:    line,
				Pos:     linePos,
			}
			tokens = append(tokens, tok)

			// 更新位置
			for i := pos; i < newPos; i++ {
				if input[i] == '\n' {
					line++
					linePos = 0 // 从0开始
				} else {
					linePos++
				}
			}
			pos = newPos
			continue
		}

		// 处理标识符（字母、下划线、中文字符开头）
		if pos < len(input) {
			r, size := utf8.DecodeRuneInString(input[pos:])
			if r != utf8.RuneError && (unicode.IsLetter(r) || r == '_' || r >= 0x4e00) {
				start := pos
				startLinePos := linePos
				pos += size
				linePos += size

				// 继续读取标识符字符
				for pos < len(input) {
					r, size := utf8.DecodeRuneInString(input[pos:])
					if r == utf8.RuneError {
						break
					}
					if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '\\' && r < 0x4e00 {
						break
					}
					pos += size
					linePos += size
				}

				tokens = append(tokens, Token{
					Type:    token.IDENTIFIER,
					Literal: input[start:pos],
					Start:   start,
					End:     pos,
					Line:    line,
					Pos:     startLinePos,
				})
				continue
			}
		}

		// 处理单个字符的 token（如操作符、分隔符等）
		if pos < len(input) {
			r, size := utf8.DecodeRuneInString(input[pos:])
			if r != utf8.RuneError {
				tokens = append(tokens, Token{
					Type:    token.UNKNOWN,
					Literal: string(r),
					Start:   pos,
					End:     pos + size,
					Line:    line,
					Pos:     linePos,
				})
				pos += size
				linePos += size
				continue
			}
		}

		// 如果无法匹配，跳过一个字符
		if input[pos] == '\n' {
			line++
			linePos = 0 // 从0开始
		} else {
			linePos++
		}
		pos++
	}

	return tokens
}
