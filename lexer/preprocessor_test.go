package lexer

import (
	"github.com/php-any/origami/token"
	"testing"
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
				{Type: token.STRING, Literal: `"Hello {$name}"`},
				{Type: token.NEWLINE, Literal: "\n"},
				{Type: token.STRING, Literal: `"Hello @{func()} ok"`},
			},
			expected: []Token{
				{Type: token.STRING, Literal: `"Hello "`},
				{Type: token.ADD, Literal: `+`},
				{Type: token.VARIABLE, Literal: `$name`},
				{Type: token.SEMICOLON, Literal: "\n"},
				{Type: token.STRING, Literal: `"Hello "`},
				{Type: token.ADD, Literal: `+`},
				{Type: token.IDENTIFIER, Literal: `func`},
				{Type: token.LPAREN, Literal: `(`},
				{Type: token.RPAREN, Literal: `)`},
				{Type: token.ADD, Literal: `+`},
				{Type: token.STRING, Literal: `" ok"`},
			},
		},
		{
			name: "处理空白符号和注释",
			tokens: []Token{
				{Type: token.WHITESPACE, Literal: " "},
				{Type: token.COMMENT, Literal: "// 这是注释"},
				{Type: token.IDENTIFIER, Literal: "test"},
				{Type: token.WHITESPACE, Literal: "\t"},
				{Type: token.MULTILINE_COMMENT, Literal: "/* 多行注释 */"},
			},
			expected: []Token{
				{Type: token.IDENTIFIER, Literal: "test"},
			},
		},
		{
			name: "处理换行符和分号",
			tokens: []Token{
				{Type: token.IDENTIFIER, Literal: "a"},
				{Type: token.NEWLINE, Literal: "\n"},
				{Type: token.IDENTIFIER, Literal: "b"},
				{Type: token.NEWLINE, Literal: "\n"},
				{Type: token.IDENTIFIER, Literal: "c"},
			},
			expected: []Token{
				{Type: token.IDENTIFIER, Literal: "a"},
				{Type: token.SEMICOLON, Literal: "\n"},
				{Type: token.IDENTIFIER, Literal: "b"},
				{Type: token.SEMICOLON, Literal: "\n"},
				{Type: token.IDENTIFIER, Literal: "c"},
			},
		},
		{
			name: "处理不需要分号的情况",
			tokens: []Token{
				{Type: token.IDENTIFIER, Literal: "a"},
				{Type: token.ADD, Literal: "+"},
				{Type: token.NEWLINE, Literal: "\n"},
				{Type: token.IDENTIFIER, Literal: "b"},
			},
			expected: []Token{
				{Type: token.IDENTIFIER, Literal: "a"},
				{Type: token.ADD, Literal: "+"},
				{Type: token.IDENTIFIER, Literal: "b"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPreprocessor(tt.tokens)
			got := p.Process()

			if len(got) != len(tt.expected) {
				t.Errorf("Process() 返回的token数量 = %v, 期望 %v", len(got), len(tt.expected))
				t.Log("实际返回的tokens:")
				for i, token := range got {
					t.Logf("  [%d] Type: %v, Literal: %q", i, token.Type, token.Literal)
				}
				t.Log("期望的tokens:")
				for i, token := range tt.expected {
					t.Logf("  [%d] Type: %v, Literal: %q", i, token.Type, token.Literal)
				}
				return
			}

			for i := 0; i < len(got); i++ {
				if got[i].Type != tt.expected[i].Type {
					t.Errorf("Process() token[%d].Type = %v, 期望 %v", i, got[i].Type, tt.expected[i].Type)
					t.Logf("  [%d] 实际: Type: %v, Literal: %q", i, got[i].Type, got[i].Literal)
					t.Logf("  [%d] 期望: Type: %v, Literal: %q", i, tt.expected[i].Type, tt.expected[i].Literal)
				}
				if got[i].Literal != tt.expected[i].Literal {
					t.Errorf("Process() token[%d].Literal = %v, 期望 %v", i, got[i].Literal, tt.expected[i].Literal)
					t.Logf("  [%d] 实际: Type: %v, Literal: %q", i, got[i].Type, got[i].Literal)
					t.Logf("  [%d] 期望: Type: %v, Literal: %q", i, tt.expected[i].Type, tt.expected[i].Literal)
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
