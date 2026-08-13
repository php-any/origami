package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/php-any/origami/token"
)

// IsDelimiter 判断是否是分割符号, 对单词参数中断的影响
// 分割符号包括：
// - 括号类: () [] {} <>
// - 运算符: + - * / % = ! & | ^ ~ ? :
// - 标点符号: , ; : ' " ` @ # $
// - 空白字符: 空格 \t \n \r \f \v
func IsDelimiter(r rune) bool {
	// 空白字符: 空格 \t \n \r \f \v
	if unicode.IsSpace(r) {
		return true
	}

	// 获取所有 token 定义
	tokenDefs := token.GetTokenDefinitions()

	// 定义会分割标识符或数字的 token 类型
	delimiterTypes := []token.TokenType{
		// 括号类: () [] {} <>
		token.LPAREN, token.RPAREN, // ()
		token.LBRACKET, token.RBRACKET, // []
		token.LBRACE, token.RBRACE, // {}
		token.LT, token.GT, // <>

		// 运算符: + - * / % = ! & | ^ ~ ? :
		token.ADD, token.SUB, token.MUL, token.QUO, token.REM, // + - * / %
		token.ASSIGN, token.NOT, token.BIT_AND, token.BIT_OR, token.BIT_XOR, token.BIT_NOT, // = ! & | ^ ~
		token.TERNARY, token.COLON, // ? :

		// 标点符号: , ; : ' " ` @ # $
		token.COMMA, token.SEMICOLON, token.COLON, // , ; :
		token.AT, token.DOLLAR, // @ $

		// 其他分割符
		token.DOT, // . (注意：\ 不在分割符列表中，因为它在命名空间中是有意义的)
	}

	// 检查字符是否匹配任何分割符类型
	for _, tokenType := range delimiterTypes {
		if definitions, exists := tokenDefs[tokenType]; exists {
			for _, definition := range definitions {
				// 使用正确的 UTF-8 解码方式
				if len(definition.Literal) > 0 {
					// 对于单字节字符，直接比较
					if len(definition.Literal) == 1 {
						if rune(definition.Literal[0]) == r {
							return true
						}
					} else {
						// 对于多字节字符，使用 UTF-8 解码
						defRune, _ := utf8.DecodeRuneInString(definition.Literal)
						if defRune == r {
							return true
						}
					}
				}
			}
		}
	}

	return false
}

func RuneIsToken(r rune, t token.TokenType) bool {
	tokenDefs := token.GetTokenDefinitions()
	if definitions, exists := tokenDefs[t]; exists {
		for _, definition := range definitions {
			// 使用正确的 UTF-8 解码方式
			if len(definition.Literal) > 0 {
				// 对于单字节字符，直接比较
				if len(definition.Literal) == 1 {
					if rune(definition.Literal[0]) == r {
						return true
					}
				} else {
					// 对于多字节字符，使用 UTF-8 解码
					defRune, _ := utf8.DecodeRuneInString(definition.Literal)
					if defRune == r {
						return true
					}
				}
			}
		}
	}
	return false
}
