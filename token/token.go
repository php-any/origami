package token

import (
	"sync"
)

// WordType 表示 token 的单词类型
type WordType int

const (
	// KEYWORD 表示关键字
	KEYWORD WordType = iota
	// OPERATOR 表示运算符
	OPERATOR
	// SYMBOL 表示符号
	SYMBOL
	// LITERAL 表示字面量
	LITERAL
	// HTML 表示 HTML 标签
	HTML
)

// TokenDefinition 表示一个 token 的定义
type TokenDefinition struct {
	// token 的类型
	Type TokenType
	// token 的字面值
	Literal string
	// WordType 表示 token 的单词类型
	WordType WordType
}

// TokenDefinitions 定义了所有的 token
var TokenDefinitions = []TokenDefinition{
	// 关键字
	{Type: IF, Literal: "if", WordType: KEYWORD},
	{Type: ELSE, Literal: "else", WordType: KEYWORD},
	{Type: ELSE_IF, Literal: "elseif", WordType: KEYWORD},
	{Type: WHILE, Literal: "while", WordType: KEYWORD},
	{Type: FOR, Literal: "for", WordType: KEYWORD},
	{Type: FOREACH, Literal: "foreach", WordType: KEYWORD},
	{Type: DO, Literal: "do", WordType: KEYWORD},
	{Type: SWITCH, Literal: "switch", WordType: KEYWORD},
	{Type: CASE, Literal: "case", WordType: KEYWORD},
	{Type: BREAK, Literal: "break", WordType: KEYWORD},
	{Type: CONTINUE, Literal: "continue", WordType: KEYWORD},
	{Type: RETURN, Literal: "return", WordType: KEYWORD},
	{Type: FUNC, Literal: "function", WordType: KEYWORD},
	{Type: CLASS, Literal: "class", WordType: KEYWORD},
	{Type: PUBLIC, Literal: "public", WordType: KEYWORD},
	{Type: PRIVATE, Literal: "private", WordType: KEYWORD},
	{Type: PROTECTED, Literal: "protected", WordType: KEYWORD},
	{Type: STATIC, Literal: "static", WordType: KEYWORD},
	{Type: FINAL, Literal: "final", WordType: KEYWORD},
	{Type: ABSTRACT, Literal: "abstract", WordType: KEYWORD},
	{Type: INTERFACE, Literal: "interface", WordType: KEYWORD},
	{Type: TRAIT, Literal: "trait", WordType: KEYWORD},
	{Type: NAMESPACE, Literal: "namespace", WordType: KEYWORD},
	{Type: USE, Literal: "use", WordType: KEYWORD},
	{Type: AS, Literal: "as", WordType: KEYWORD},
	{Type: NEW, Literal: "new", WordType: KEYWORD},
	{Type: INSTANCEOF, Literal: "instanceof", WordType: KEYWORD},
	{Type: LIKE, Literal: "like", WordType: KEYWORD},
	{Type: CONST, Literal: "const", WordType: KEYWORD},
	{Type: VAR, Literal: "var", WordType: KEYWORD},
	{Type: ECHO, Literal: "echo", WordType: KEYWORD},
	{Type: THROW, Literal: "throw", WordType: KEYWORD},
	{Type: TRY, Literal: "try", WordType: KEYWORD},
	{Type: CATCH, Literal: "catch", WordType: KEYWORD},
	{Type: FINALLY, Literal: "finally", WordType: KEYWORD},
	{Type: CLONE, Literal: "clone", WordType: KEYWORD},
	{Type: YIELD, Literal: "yield", WordType: KEYWORD},
	{Type: FROM, Literal: "from", WordType: KEYWORD},
	{Type: INSTEAD_OF, Literal: "insteadof", WordType: KEYWORD},
	{Type: EXTENDS, Literal: "extends", WordType: KEYWORD},
	{Type: IMPLEMENTS, Literal: "implements", WordType: KEYWORD},
	{Type: LIST, Literal: "list", WordType: KEYWORD},
	{Type: HALT_COMPILER, Literal: "__halt_compiler", WordType: KEYWORD},
	{Type: MATCH, Literal: "match", WordType: KEYWORD},
	{Type: ENUM, Literal: "enum", WordType: KEYWORD},
	{Type: READONLY, Literal: "readonly", WordType: KEYWORD},
	{Type: FN, Literal: "fn", WordType: KEYWORD},
	{Type: SPAWN, Literal: "spawn", WordType: KEYWORD},
	{Type: THIS, Literal: "$this", WordType: KEYWORD},
	{Type: PARENT, Literal: "parent", WordType: KEYWORD},
	{Type: IN, Literal: "in", WordType: KEYWORD},
	{Type: DEFAULT, Literal: "default", WordType: KEYWORD},
	{Type: UNUSED, Literal: "_", WordType: KEYWORD},
	{Type: DIR, Literal: "__DIR__", WordType: KEYWORD},

	// 运算符
	{Type: ADD, Literal: "+", WordType: OPERATOR},
	{Type: SUB, Literal: "-", WordType: OPERATOR},
	{Type: MUL, Literal: "*", WordType: OPERATOR},
	{Type: QUO, Literal: "/", WordType: OPERATOR},
	{Type: REM, Literal: "%", WordType: OPERATOR},
	{Type: ASSIGN, Literal: "=", WordType: OPERATOR},
	{Type: EQ, Literal: "==", WordType: OPERATOR},
	{Type: NE, Literal: "!=", WordType: OPERATOR},
	{Type: EQ_STRICT, Literal: "===", WordType: OPERATOR},
	{Type: NE_STRICT, Literal: "!==", WordType: OPERATOR},
	{Type: LT, Literal: "<", WordType: OPERATOR},
	{Type: GT, Literal: ">", WordType: OPERATOR},
	{Type: LE, Literal: "<=", WordType: OPERATOR},
	{Type: GE, Literal: ">=", WordType: OPERATOR},
	{Type: LAND, Literal: "&&", WordType: OPERATOR},
	{Type: LOR, Literal: "||", WordType: OPERATOR},
	{Type: NOT, Literal: "!", WordType: OPERATOR},
	{Type: BIT_AND, Literal: "&", WordType: OPERATOR},
	{Type: BIT_OR, Literal: "|", WordType: OPERATOR},
	{Type: BIT_XOR, Literal: "^", WordType: OPERATOR},
	{Type: BIT_NOT, Literal: "~", WordType: OPERATOR},
	{Type: SHL, Literal: "<<", WordType: OPERATOR},
	{Type: SHR, Literal: ">>", WordType: OPERATOR},
	{Type: INCR, Literal: "++", WordType: OPERATOR},
	{Type: DECR, Literal: "--", WordType: OPERATOR},
	{Type: OBJECT_OPERATOR, Literal: "->", WordType: OPERATOR},
	{Type: ARRAY_KEY_VALUE, Literal: "=>", WordType: OPERATOR},
	{Type: TERNARY, Literal: "?", WordType: OPERATOR},
	{Type: COLON, Literal: ":", WordType: OPERATOR},
	{Type: SCOPE_RESOLUTION, Literal: "::", WordType: OPERATOR},
	{Type: AT, Literal: "@", WordType: OPERATOR},
	{Type: DOLLAR, Literal: "$", WordType: SYMBOL},
	{Type: COMMA, Literal: ",", WordType: OPERATOR},
	{Type: SEMICOLON, Literal: ";", WordType: OPERATOR},
	{Type: LPAREN, Literal: "(", WordType: OPERATOR},
	{Type: RPAREN, Literal: ")", WordType: OPERATOR},
	{Type: LBRACE, Literal: "{", WordType: OPERATOR},
	{Type: RBRACE, Literal: "}", WordType: OPERATOR},
	{Type: LBRACKET, Literal: "[", WordType: OPERATOR},
	{Type: RBRACKET, Literal: "]", WordType: OPERATOR},
	{Type: SPACESHIP, Literal: "<=>", WordType: OPERATOR},
	{Type: NULLSAFE_CALL, Literal: "??->", WordType: OPERATOR},
	{Type: NULL_COALESCE, Literal: "??", WordType: OPERATOR},
	{Type: POWER, Literal: "**", WordType: OPERATOR},
	{Type: POWER_EQ, Literal: "**=", WordType: OPERATOR},
	{Type: ADD_EQ, Literal: "+=", WordType: OPERATOR},
	{Type: SUB_EQ, Literal: "-=", WordType: OPERATOR},
	{Type: MUL_EQ, Literal: "*=", WordType: OPERATOR},
	{Type: QUO_EQ, Literal: "/=", WordType: OPERATOR},
	{Type: REM_EQ, Literal: "%=", WordType: OPERATOR},
	{Type: CONCAT_EQ, Literal: ".=", WordType: OPERATOR},
	{Type: BIT_AND_EQ, Literal: "&=", WordType: OPERATOR},
	{Type: BIT_OR_EQ, Literal: "|=", WordType: OPERATOR},
	{Type: BIT_XOR_EQ, Literal: "^=", WordType: OPERATOR},
	{Type: SHL_EQ, Literal: "<<=", WordType: OPERATOR},
	{Type: SHR_EQ, Literal: ">>=", WordType: OPERATOR},
	{Type: NAMESPACE_SEPARATOR, Literal: "\\", WordType: OPERATOR},
	{Type: DOT, Literal: ".", WordType: OPERATOR},
	{Type: ELLIPSIS, Literal: "...", WordType: OPERATOR},
	{Type: DOUBLE_DOT, Literal: "..", WordType: OPERATOR},

	// 字面量
	{Type: NULL, Literal: "null", WordType: LITERAL},
	{Type: TRUE, Literal: "true", WordType: LITERAL},
	{Type: FALSE, Literal: "false", WordType: LITERAL},
	{Type: BOOL, Literal: "bool", WordType: KEYWORD},
	{Type: NEWLINE, Literal: "\n", WordType: SYMBOL},

	{Type: START_TAG, Literal: "<?php", WordType: OPERATOR},
	{Type: END_TAG, Literal: "?>", WordType: OPERATOR},
}

var initTokenDefinitions = false
var tree = make(map[TokenType][]TokenDefinition)
var once sync.Once

// GetTokenDefinitions 根据类型获取单词
func GetTokenDefinitions() map[TokenType][]TokenDefinition {
	once.Do(func() {
		// 初始化代码只会执行一次，即使有多个goroutine同时调用
		for _, definition := range TokenDefinitions {
			if _, ok := tree[definition.Type]; !ok {
				tree[definition.Type] = []TokenDefinition{definition}
			} else {
				tree[definition.Type] = append(tree[definition.Type], definition)
			}
		}
		initTokenDefinitions = true
	})
	return tree
}
