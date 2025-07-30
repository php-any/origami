package lexer

import (
	"testing"

	"github.com/php-any/origami/token"
)

// 验证测试用例中的期望值是否正确
func validateTestExpectations(t *testing.T, input string, expected []Token) {
	for i, tk := range expected {
		// 检查 Start 和 End 是否在输入字符串范围内
		if tk.Start < 0 || tk.Start >= len(input) {
			t.Errorf("测试用例错误: tk[%d] 的 Start 位置 %d 超出输入字符串范围 [0, %d]",
				i, tk.Start, len(input)-1)
		}
		if tk.End < 0 || tk.End > len(input) {
			t.Errorf("测试用例错误: tk[%d] 的 End 位置 %d 超出输入字符串范围 [0, %d]",
				i, tk.End, len(input))
		}
		if tk.Start >= tk.End {
			t.Errorf("测试用例错误: tk[%d] 的 Start(%d) 大于等于 End(%d)",
				i, tk.Start, tk.End)
		}

		// 检查字面值是否与输入字符串匹配
		expectedLiteral := input[tk.Start:tk.End]
		if tk.Literal != expectedLiteral {
			t.Errorf("测试用例错误: tk[%d] 的字面值不匹配\n\t期望: %q\n\t实际: %q\n\t位置: [%d:%d]",
				i, expectedLiteral, tk.Literal, tk.Start, tk.End)
		} else {
			//fmt.Printf("字符串 [%v] == [%v]\n", tk.Literal, expectedLiteral)
		}
	}
}

func TestLexer_Tokenize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "简单英文代码",
			input: "if (x > 0) { echo 'hello'; }",
			expected: []Token{
				{Type: token.IF, Literal: "if", Start: 0, End: 2},
				{Type: token.LPAREN, Literal: "(", Start: 3, End: 4},
				{Type: token.IDENTIFIER, Literal: "x", Start: 4, End: 5},
				{Type: token.GT, Literal: ">", Start: 6, End: 7},
				{Type: token.INT, Literal: "0", Start: 8, End: 9},
				{Type: token.RPAREN, Literal: ")", Start: 9, End: 10},
				{Type: token.RBRACE, Literal: "{", Start: 11, End: 12},
				{Type: token.ECHO, Literal: "echo", Start: 13, End: 17},
				{Type: token.STRING, Literal: "'hello'", Start: 18, End: 25},
				{Type: token.SEMICOLON, Literal: ";", Start: 25, End: 26},
				{Type: token.RBRACE, Literal: "}", Start: 27, End: 28},
			},
		},
		{
			name:  "中英文混合代码",
			input: "如果 (x > 0) { 输出 '你好'; }",
			expected: []Token{
				{Type: token.IDENTIFIER, Literal: "如果", Start: 0, End: 6},
				{Type: token.LPAREN, Literal: "(", Start: 7, End: 8},
				{Type: token.IDENTIFIER, Literal: "x", Start: 8, End: 9},
				{Type: token.LT, Literal: ">", Start: 10, End: 11},
				{Type: token.INT, Literal: "0", Start: 12, End: 13},
				{Type: token.RPAREN, Literal: ")", Start: 13, End: 14},
				{Type: token.RBRACE, Literal: "{", Start: 15, End: 16},
				{Type: token.IDENTIFIER, Literal: "输出", Start: 17, End: 23},
				{Type: token.STRING, Literal: "'你好'", Start: 24, End: 32},
				{Type: token.SEMICOLON, Literal: ";", Start: 32, End: 33},
				{Type: token.RBRACE, Literal: "}", Start: 34, End: 35},
			},
		},
		{
			name:  "中文标识符",
			input: "变量 = 100;",
			expected: []Token{
				{Type: token.IDENTIFIER, Literal: "变量", Start: 0, End: 6},
				{Type: token.ASSIGN, Literal: "=", Start: 7, End: 8},
				{Type: token.INT, Literal: "100", Start: 9, End: 12},
				{Type: token.SEMICOLON, Literal: ";", Start: 12, End: 13},
			},
		},
		{
			name:  "混合运算符",
			input: "x += 你好;",
			expected: []Token{
				{Type: token.IDENTIFIER, Literal: "x", Start: 0, End: 1},
				{Type: token.ADD_EQ, Literal: "+=", Start: 2, End: 4},
				{Type: token.IDENTIFIER, Literal: "你好", Start: 5, End: 11},
				{Type: token.SEMICOLON, Literal: ";", Start: 11, End: 12},
			},
		},
		{
			name:  "中文空格",
			input: "x　y", // 使用全角空格
			expected: []Token{
				{Type: token.IDENTIFIER, Literal: "x", Start: 0, End: 1},
				{Type: token.IDENTIFIER, Literal: "y", Start: 4, End: 5},
			},
		},
		{
			name:  "混合空格",
			input: "x \t\n\n　y", // 混合使用半角和全角空格
			expected: []Token{
				{Type: token.IDENTIFIER, Literal: "x", Start: 0, End: 1},
				{Type: token.NEWLINE, Literal: "\n", Start: 3, End: 4}, // 多个换行只保留一个
				{Type: token.IDENTIFIER, Literal: "y", Start: 8, End: 9},
			},
		},
		{
			name:  "变量赋值",
			input: "$test = 123;",
			expected: []Token{
				{Type: token.DOLLAR, Literal: "$", Start: 0, End: 1},
				{Type: token.IDENTIFIER, Literal: "test", Start: 1, End: 5},
				{Type: token.ASSIGN, Literal: "=", Start: 6, End: 7},
				{Type: token.INT, Literal: "123", Start: 8, End: 11},
				{Type: token.SEMICOLON, Literal: ";", Start: 11, End: 12},
			},
		},
		{
			name: "多语句和换行",
			input: `$a = 1;
$b = "hello";
if ($a > 0) {
    echo $b;
}`,
			expected: []Token{
				{Type: token.DOLLAR, Literal: "$", Start: 0, End: 1},
				{Type: token.IDENTIFIER, Literal: "a", Start: 1, End: 2},
				{Type: token.ASSIGN, Literal: "=", Start: 3, End: 4},
				{Type: token.INT, Literal: "1", Start: 5, End: 6},
				{Type: token.SEMICOLON, Literal: ";", Start: 6, End: 7},
				{Type: token.NEWLINE, Literal: "\n", Start: 7, End: 8},
				{Type: token.DOLLAR, Literal: "$", Start: 8, End: 9},
				{Type: token.IDENTIFIER, Literal: "b", Start: 9, End: 10},
				{Type: token.ASSIGN, Literal: "=", Start: 11, End: 12},
				{Type: token.STRING, Literal: "\"hello\"", Start: 13, End: 20},
				{Type: token.SEMICOLON, Literal: ";", Start: 20, End: 21},
				{Type: token.NEWLINE, Literal: "\n", Start: 21, End: 22},
				{Type: token.IF, Literal: "if", Start: 22, End: 24},
				{Type: token.LPAREN, Literal: "(", Start: 25, End: 26},
				{Type: token.DOLLAR, Literal: "$", Start: 26, End: 27},
				{Type: token.IDENTIFIER, Literal: "a", Start: 27, End: 28},
				{Type: token.GT, Literal: ">", Start: 29, End: 30},
				{Type: token.INT, Literal: "0", Start: 31, End: 32},
				{Type: token.RPAREN, Literal: ")", Start: 32, End: 33},
				{Type: token.RBRACE, Literal: "{", Start: 34, End: 35},
				{Type: token.NEWLINE, Literal: "\n", Start: 35, End: 36},
				{Type: token.ECHO, Literal: "echo", Start: 40, End: 44},
				{Type: token.DOLLAR, Literal: "$", Start: 45, End: 46},
				{Type: token.IDENTIFIER, Literal: "b", Start: 46, End: 47},
				{Type: token.SEMICOLON, Literal: ";", Start: 47, End: 48},
				{Type: token.NEWLINE, Literal: "\n", Start: 48, End: 49},
				{Type: token.RBRACE, Literal: "}", Start: 49, End: 50},
			},
		},
		{
			name: "连续换行",
			input: `$a = 1;


$b = 2;`,
			expected: []Token{
				{Type: token.DOLLAR, Literal: "$", Start: 0, End: 1},
				{Type: token.IDENTIFIER, Literal: "a", Start: 1, End: 2},
				{Type: token.ASSIGN, Literal: "=", Start: 3, End: 4},
				{Type: token.INT, Literal: "1", Start: 5, End: 6},
				{Type: token.SEMICOLON, Literal: ";", Start: 6, End: 7},
				{Type: token.NEWLINE, Literal: "\n", Start: 7, End: 8},
				{Type: token.DOLLAR, Literal: "$", Start: 10, End: 11},
				{Type: token.IDENTIFIER, Literal: "b", Start: 11, End: 12},
				{Type: token.ASSIGN, Literal: "=", Start: 13, End: 14},
				{Type: token.INT, Literal: "2", Start: 15, End: 16},
				{Type: token.SEMICOLON, Literal: ";", Start: 16, End: 17},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 首先验证测试用例的期望值是否正确
			validateTestExpectations(t, tt.input, tt.expected)

			lexer := NewLexer()
			tokens := lexer.Tokenize(tt.input)

			if len(tokens) != len(tt.expected) {
				t.Errorf("期望 %d 个 token，实际得到 %d 个\n期望: %#v\n实际: %#v", len(tt.expected), len(tokens), tt.expected, tokens)
				for i := 0; i < len(tokens) && i < len(tt.expected); i++ {
					t.Errorf("token[%d] 期望: %+v，实际: %+v", i, tt.expected[i], tokens[i])
				}
				return
			}

			for i, tk := range tokens {
				expected := tt.expected[i]
				if tk.Type != expected.Type || tk.Literal != expected.Literal || tk.Start != expected.Start || tk.End != expected.End {
					t.Errorf("token[%d] 不匹配:\n\t期望: %+v\n\t实际: %+v", i, expected, tk)
				}
			}
		})
	}
}

func TestLexer_Whitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "中文空格",
			input: "x　y", // 使用全角空格
			expected: []Token{
				{Type: token.IDENTIFIER, Literal: "x", Start: 0, End: 1},
				{Type: token.IDENTIFIER, Literal: "y", Start: 4, End: 5},
			},
		},
		{
			name:  "混合空格",
			input: "x \t\n　y", // 混合使用半角和全角空格
			expected: []Token{
				{Type: token.IDENTIFIER, Literal: "x", Start: 0, End: 1},
				{Type: token.NEWLINE, Literal: "\n", Start: 3, End: 4},
				{Type: token.IDENTIFIER, Literal: "y", Start: 7, End: 8},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 首先验证测试用例的期望值是否正确
			validateTestExpectations(t, tt.input, tt.expected)

			lexer := NewLexer()
			tokens := lexer.Tokenize(tt.input)

			if len(tokens) != len(tt.expected) {
				t.Errorf("期望 %d 个 token，实际得到 %d 个", len(tt.expected), len(tokens))
				return
			}

			for i, tk := range tokens {
				expected := tt.expected[i]
				if tk.Type != expected.Type {
					t.Errorf("token[%d] 类型不匹配: 期望 %v，实际 %v", i, expected.Type, tk.Type)
				}
				if tk.Literal != expected.Literal {
					t.Errorf("token[%d] 字面值不匹配: 期望 %q，实际 %q", i, expected.Literal, tk.Literal)
				}
				if tk.Start != expected.Start {
					t.Errorf("token[%d] 起始位置不匹配: 期望 %d，实际 %d", i, expected.Start, tk.Start)
				}
				if tk.End != expected.End {
					t.Errorf("token[%d] 结束位置不匹配: 期望 %d，实际 %d", i, expected.End, tk.End)
				}
			}
		})
	}
}

func TestLexer_Literals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "整数和浮点数字面量",
			input: "123 3.14 -42 0x1A 0b1010 0777",
			expected: []Token{
				{Type: token.INT, Literal: "123", Start: 0, End: 3},
				{Type: token.FLOAT, Literal: "3.14", Start: 4, End: 8},
				{Type: token.INT, Literal: "-42", Start: 9, End: 12},
				{Type: token.NUMBER, Literal: "0x1A", Start: 13, End: 17},
				{Type: token.NUMBER, Literal: "0b1010", Start: 18, End: 24},
				{Type: token.NUMBER, Literal: "0777", Start: 25, End: 29},
			},
		},
		{
			name:  "字符串字面量",
			input: "'hello' \"world\" `backtick`",
			expected: []Token{
				{Type: token.STRING, Literal: "'hello'", Start: 0, End: 7},
				{Type: token.STRING, Literal: "\"world\"", Start: 8, End: 15},
				{Type: token.STRING, Literal: "`backtick`", Start: 16, End: 26},
			},
		},
		{
			name:  "布尔值和null字面量",
			input: "true false null",
			expected: []Token{
				{Type: token.TRUE, Literal: "true", Start: 0, End: 4},
				{Type: token.FALSE, Literal: "false", Start: 5, End: 10},
				{Type: token.NULL, Literal: "null", Start: 11, End: 15},
			},
		},
		{
			name:  "字节字面量",
			input: "b'a' b'\\n' b'\\t'",
			expected: []Token{
				{Type: token.BYTE, Literal: "b'a'", Start: 0, End: 4},
				{Type: token.BYTE, Literal: "b'\\n'", Start: 5, End: 10},
				{Type: token.BYTE, Literal: "b'\\t'", Start: 11, End: 16},
			},
		},
		{
			name:  "混合字面量",
			input: "123 'hello' true b'x' null 3.14",
			expected: []Token{
				{Type: token.INT, Literal: "123", Start: 0, End: 3},
				{Type: token.STRING, Literal: "'hello'", Start: 4, End: 11},
				{Type: token.TRUE, Literal: "true", Start: 12, End: 16},
				{Type: token.BYTE, Literal: "b'x'", Start: 17, End: 21},
				{Type: token.NULL, Literal: "null", Start: 22, End: 26},
				{Type: token.FLOAT, Literal: "3.14", Start: 27, End: 31},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 首先验证测试用例的期望值是否正确
			validateTestExpectations(t, tt.input, tt.expected)

			lexer := NewLexer()
			tokens := lexer.Tokenize(tt.input)

			if len(tokens) != len(tt.expected) {
				t.Errorf("期望 %d 个 token，实际得到 %d 个", len(tt.expected), len(tokens))
				return
			}

			for i, tk := range tokens {
				expected := tt.expected[i]
				if tk.Type != expected.Type {
					t.Errorf("token[%d] 类型不匹配: 期望 %v，实际 %v", i, expected.Type, tk.Type)
				}
				if tk.Literal != expected.Literal {
					t.Errorf("token[%d] 字面值不匹配: 期望 %q，实际 %q", i, expected.Literal, tk.Literal)
				}
				if tk.Start != expected.Start {
					t.Errorf("token[%d] 起始位置不匹配: 期望 %d，实际 %d", i, expected.Start, tk.Start)
				}
				if tk.End != expected.End {
					t.Errorf("token[%d] 结束位置不匹配: 期望 %d，实际 %d", i, expected.End, tk.End)
				}
			}
		})
	}
}

func TestLexer_Operators(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "算术运算符",
			input: "1 + 2 - 3 * 4 / 5 % 6",
			expected: []Token{
				{Type: token.INT, Literal: "1", Start: 0, End: 1},
				{Type: token.ADD, Literal: "+", Start: 2, End: 3},
				{Type: token.INT, Literal: "2", Start: 4, End: 5},
				{Type: token.SUB, Literal: "-", Start: 6, End: 7},
				{Type: token.INT, Literal: "3", Start: 8, End: 9},
				{Type: token.MUL, Literal: "*", Start: 10, End: 11},
				{Type: token.INT, Literal: "4", Start: 12, End: 13},
				{Type: token.QUO, Literal: "/", Start: 14, End: 15},
				{Type: token.INT, Literal: "5", Start: 16, End: 17},
				{Type: token.REM, Literal: "%", Start: 18, End: 19},
				{Type: token.INT, Literal: "6", Start: 20, End: 21},
			},
		},
		{
			name:  "比较运算符",
			input: "a == b != c > d < e >= f <= g",
			expected: []Token{
				{Type: token.IDENTIFIER, Literal: "a", Start: 0, End: 1},
				{Type: token.EQ, Literal: "==", Start: 2, End: 4},
				{Type: token.IDENTIFIER, Literal: "b", Start: 5, End: 6},
				{Type: token.NE, Literal: "!=", Start: 7, End: 9},
				{Type: token.IDENTIFIER, Literal: "c", Start: 10, End: 11},
				{Type: token.GT, Literal: ">", Start: 12, End: 13},
				{Type: token.IDENTIFIER, Literal: "d", Start: 14, End: 15},
				{Type: token.LT, Literal: "<", Start: 16, End: 17},
				{Type: token.IDENTIFIER, Literal: "e", Start: 18, End: 19},
				{Type: token.GE, Literal: ">=", Start: 20, End: 22},
				{Type: token.IDENTIFIER, Literal: "f", Start: 23, End: 24},
				{Type: token.LE, Literal: "<=", Start: 25, End: 27},
				{Type: token.IDENTIFIER, Literal: "g", Start: 28, End: 29},
			},
		},
		{
			name:  "逻辑运算符",
			input: "a && b || c !d",
			expected: []Token{
				{Type: token.IDENTIFIER, Literal: "a", Start: 0, End: 1},
				{Type: token.LAND, Literal: "&&", Start: 2, End: 4},
				{Type: token.IDENTIFIER, Literal: "b", Start: 5, End: 6},
				{Type: token.LOR, Literal: "||", Start: 7, End: 9},
				{Type: token.IDENTIFIER, Literal: "c", Start: 10, End: 11},
				{Type: token.NOT, Literal: "!", Start: 12, End: 13},
				{Type: token.IDENTIFIER, Literal: "d", Start: 13, End: 14},
			},
		},
		{
			name:  "赋值运算符",
			input: "a = b += c -= d *= e /= f %= g",
			expected: []Token{
				{Type: token.IDENTIFIER, Literal: "a", Start: 0, End: 1},
				{Type: token.ASSIGN, Literal: "=", Start: 2, End: 3},
				{Type: token.IDENTIFIER, Literal: "b", Start: 4, End: 5},
				{Type: token.ADD_EQ, Literal: "+=", Start: 6, End: 8},
				{Type: token.IDENTIFIER, Literal: "c", Start: 9, End: 10},
				{Type: token.SUB_EQ, Literal: "-=", Start: 11, End: 13},
				{Type: token.IDENTIFIER, Literal: "d", Start: 14, End: 15},
				{Type: token.MUL_EQ, Literal: "*=", Start: 16, End: 18},
				{Type: token.IDENTIFIER, Literal: "e", Start: 19, End: 20},
				{Type: token.QUO_EQ, Literal: "/=", Start: 21, End: 23},
				{Type: token.IDENTIFIER, Literal: "f", Start: 24, End: 25},
				{Type: token.REM_EQ, Literal: "%=", Start: 26, End: 28},
				{Type: token.IDENTIFIER, Literal: "g", Start: 29, End: 30},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validateTestExpectations(t, tt.input, tt.expected)

			lexer := NewLexer()
			tokens := lexer.Tokenize(tt.input)

			if len(tokens) != len(tt.expected) {
				t.Errorf("期望 %d 个 token，实际得到 %d 个", len(tt.expected), len(tokens))
				return
			}

			for i, tk := range tokens {
				expected := tt.expected[i]
				if tk.Type != expected.Type {
					t.Errorf("token[%d] 类型不匹配: 期望 %v，实际 %v", i, expected.Type, tk.Type)
				}
				if tk.Literal != expected.Literal {
					t.Errorf("token[%d] 字面值不匹配: 期望 %q，实际 %q", i, expected.Literal, tk.Literal)
				}
				if tk.Start != expected.Start {
					t.Errorf("token[%d] 起始位置不匹配: 期望 %d，实际 %d", i, expected.Start, tk.Start)
				}
				if tk.End != expected.End {
					t.Errorf("token[%d] 结束位置不匹配: 期望 %d，实际 %d", i, expected.End, tk.End)
				}
			}
		})
	}
}
