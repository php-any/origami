package lexer

import (
	"testing"

	"github.com/php-any/origami/token"
)

func TestPreprocessor_Process(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []Token
		expected []Token
	}{
		{
			name: "处理字符串插值",
			tokens: []Token{
				NewWorkerToken(token.STRING, `"Hello {$name}"`, 0, 0, 0, 0),
				NewWorkerToken(token.NEWLINE, "\n", 0, 0, 0, 0),
				NewWorkerToken(token.STRING, `"Hello @{func()} ok"`, 0, 0, 0, 0),
			},
			expected: []Token{
				// 第一个字符串插值应该返回 INTERPOLATION_TOKEN，包含子 token
				// 我们检查类型和子 token 数量，不检查展开的 token
				nil, // 占位符，实际会检查 LingToken
				NewWorkerToken(token.SEMICOLON, "\n", 0, 0, 0, 0),
				// 第二个字符串插值
				nil, // 占位符，实际会检查 LingToken
			},
		},
		{
			name: "处理空白符号和注释",
			tokens: []Token{
				NewWorkerToken(token.WHITESPACE, " ", 0, 0, 0, 0),
				NewWorkerToken(token.COMMENT, "// 这是注释", 0, 0, 0, 0),
				NewWorkerToken(token.IDENTIFIER, "test", 0, 0, 0, 0),
				NewWorkerToken(token.WHITESPACE, "\t", 0, 0, 0, 0),
				NewWorkerToken(token.MULTILINE_COMMENT, "/* 多行注释 */", 0, 0, 0, 0),
			},
			expected: []Token{
				NewWorkerToken(token.IDENTIFIER, "test", 0, 0, 0, 0),
			},
		},
		{
			name: "处理换行符和分号",
			tokens: []Token{
				NewWorkerToken(token.IDENTIFIER, "a", 0, 0, 0, 0),
				NewWorkerToken(token.NEWLINE, "\n", 0, 0, 0, 0),
				NewWorkerToken(token.IDENTIFIER, "b", 0, 0, 0, 0),
				NewWorkerToken(token.NEWLINE, "\n", 0, 0, 0, 0),
				NewWorkerToken(token.IDENTIFIER, "c", 0, 0, 0, 0),
			},
			expected: []Token{
				NewWorkerToken(token.IDENTIFIER, "a", 0, 0, 0, 0),
				NewWorkerToken(token.SEMICOLON, "\n", 0, 0, 0, 0),
				NewWorkerToken(token.IDENTIFIER, "b", 0, 0, 0, 0),
				NewWorkerToken(token.SEMICOLON, "\n", 0, 0, 0, 0),
				NewWorkerToken(token.IDENTIFIER, "c", 0, 0, 0, 0),
			},
		},
		{
			name: "处理不需要分号的情况",
			tokens: []Token{
				NewWorkerToken(token.IDENTIFIER, "a", 0, 0, 0, 0),
				NewWorkerToken(token.ADD, "+", 0, 0, 0, 0),
				NewWorkerToken(token.NEWLINE, "\n", 0, 0, 0, 0),
				NewWorkerToken(token.IDENTIFIER, "b", 0, 0, 0, 0),
			},
			expected: []Token{
				NewWorkerToken(token.IDENTIFIER, "a", 0, 0, 0, 0),
				NewWorkerToken(token.ADD, "+", 0, 0, 0, 0),
				NewWorkerToken(token.IDENTIFIER, "b", 0, 0, 0, 0),
			},
		},
		{
			name: "处理 {$ 后跟非标识符的情况",
			tokens: []Token{
				NewWorkerToken(token.STRING, `"Hello {$+test}"`, 0, 0, 0, 0),
			},
			expected: []Token{
				NewWorkerToken(token.STRING, `"Hello {$+test}"`, 0, 0, 0, 0),
			},
		},
		{
			name: "处理 {$ 后跟数字的情况",
			tokens: []Token{
				NewWorkerToken(token.STRING, `"Hello {$123}"`, 0, 0, 0, 0),
			},
			expected: []Token{
				NewWorkerToken(token.STRING, `"Hello {$123}"`, 0, 0, 0, 0),
			},
		},
		{
			name: "处理 {$ 后跟特殊字符的情况",
			tokens: []Token{
				NewWorkerToken(token.STRING, `"Hello {$!test}"`, 0, 0, 0, 0),
			},
			expected: []Token{
				NewWorkerToken(token.STRING, `"Hello {$!test}"`, 0, 0, 0, 0),
			},
		},
		{
			name: "处理 {$} 空变量名的情况",
			tokens: []Token{
				NewWorkerToken(token.STRING, `"Hello {$}"`, 0, 0, 0, 0),
			},
			expected: []Token{
				NewWorkerToken(token.STRING, `"Hello {$}"`, 0, 0, 0, 0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPreprocessor(tt.tokens)
			got := p.Process()

			// 特殊处理字符串插值测试用例
			if tt.name == "处理字符串插值" {
				// 应该返回 3 个 token：两个 INTERPOLATION_TOKEN 和一个 SEMICOLON
				if len(got) != 3 {
					t.Errorf("Process() 返回的token数量 = %v, 期望 3", len(got))
					return
				}

				// 检查第一个 INTERPOLATION_TOKEN
				lingToken1, ok := got[0].(*LingToken)
				if !ok || lingToken1.Type() != token.INTERPOLATION_TOKEN {
					t.Errorf("Process() token[0] 应该是 INTERPOLATION_TOKEN，实际 Type: %v", got[0].Type())
					return
				}
				children1 := lingToken1.Children()
				// 应该包含: "Hello " (STRING), $name (VARIABLE)
				if len(children1) < 2 {
					t.Errorf("Process() 第一个插值token应该有至少2个子token，实际: %v", len(children1))
					return
				}
				// 检查子 token 类型
				if children1[0].Type() != token.STRING {
					t.Errorf("Process() 第一个插值token的第一个子token应该是STRING，实际: %v", children1[0].Type())
				}
				if lingToken2, ok := children1[1].(*LingToken); ok {
					// 如果是 LingToken，检查其子 token
					subChildren := lingToken2.Children()
					if len(subChildren) > 0 && subChildren[0].Type() != token.VARIABLE {
						t.Errorf("Process() 第一个插值token的第二个子token应该是VARIABLE，实际: %v", subChildren[0].Type())
					}
				}

				// 检查 SEMICOLON
				if got[1].Type() != token.SEMICOLON {
					t.Errorf("Process() token[1] 应该是 SEMICOLON，实际: %v", got[1].Type())
				}

				// 检查第二个 INTERPOLATION_TOKEN
				lingToken2, ok := got[2].(*LingToken)
				if !ok || lingToken2.Type() != token.INTERPOLATION_TOKEN {
					t.Errorf("Process() token[2] 应该是 INTERPOLATION_TOKEN，实际 Type: %v", got[2].Type())
					return
				}
				children2 := lingToken2.Children()
				// 应该包含: "Hello " (STRING), func() (INTERPOLATION_VALUE), " ok" (STRING)
				if len(children2) < 3 {
					t.Errorf("Process() 第二个插值token应该有至少3个子token，实际: %v", len(children2))
					return
				}
				// 检查子 token 类型
				if children2[0].Type() != token.STRING {
					t.Errorf("Process() 第二个插值token的第一个子token应该是STRING，实际: %v", children2[0].Type())
				}
				if children2[1].Type() != token.INTERPOLATION_VALUE {
					t.Errorf("Process() 第二个插值token的第二个子token应该是INTERPOLATION_VALUE，实际: %v", children2[1].Type())
				}
				if children2[2].Type() != token.STRING {
					t.Errorf("Process() 第二个插值token的第三个子token应该是STRING，实际: %v", children2[2].Type())
				}
				return
			}

			// 其他测试用例的原有逻辑
			if len(got) != len(tt.expected) {
				t.Errorf("Process() 返回的token数量 = %v, 期望 %v", len(got), len(tt.expected))
				t.Log("实际返回的tokens:")
				for i, kk := range got {
					t.Logf("  [%d] Type: %v, Literal: %q", i, kk.Type(), kk.Literal())
				}
				t.Log("期望的tokens:")
				for i, kk := range tt.expected {
					if kk != nil {
						t.Logf("  [%d] Type: %v, Literal: %q", i, kk.Type(), kk.Literal())
					}
				}
				return
			}

			for i := 0; i < len(got); i++ {
				if tt.expected[i] == nil {
					continue // 跳过占位符
				}
				if got[i].Type() != tt.expected[i].Type() {
					t.Errorf("Process() token[%d].Type = %v, 期望 %v", i, got[i].Type(), tt.expected[i].Type())
					t.Logf("  [%d] 实际: Type: %v, Literal: %q", i, got[i].Type(), got[i].Literal())
					t.Logf("  [%d] 期望: Type: %v, Literal: %q", i, tt.expected[i].Type(), tt.expected[i].Literal())
				}
				if got[i].Literal() != tt.expected[i].Literal() {
					t.Errorf("Process() token[%d].Literal = %v, 期望 %v", i, got[i].Literal(), tt.expected[i].Literal())
					t.Logf("  [%d] 实际: Type: %v, Literal: %q", i, got[i].Type(), got[i].Literal())
					t.Logf("  [%d] 期望: Type: %v, Literal: %q", i, tt.expected[i].Type(), tt.expected[i].Literal())
				}
			}
		})
	}
}

func TestIsNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected bool
	}{
		{"数字", '1', true},
		{"字母", 'a', false},
		{"特殊字符", '@', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNumber(tt.input); got != tt.expected {
				t.Errorf("isNumber() = %v, 期望 %v", got, tt.expected)
			}
		})
	}
}

func TestIsSpecialSymbol(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected bool
	}{
		{"下划线", '_', true},
		{"标点符号", '.', true},
		{"字母", 'a', false},
		{"数字", '1', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSpecialSymbol(tt.input); got != tt.expected {
				t.Errorf("isSpecialSymbol() = %v, 期望 %v", got, tt.expected)
			}
		})
	}
}
