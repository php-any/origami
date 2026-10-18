package lexer

import "github.com/php-any/origami/token"

func endsCompilerHalt(tokens []Token) bool {
	expected := []token.TokenType{token.SEMICOLON, token.RPAREN, token.LPAREN, token.HALT_COMPILER}
	i := len(tokens) - 1
	for _, kind := range expected {
		for i >= 0 && (tokens[i].Type() == token.NEWLINE || tokens[i].Type() == token.COMMENT || tokens[i].Type() == token.MULTILINE_COMMENT || tokens[i].Type() == token.WHITESPACE) {
			i--
		}
		if i < 0 || tokens[i].Type() != kind {
			return false
		}
		i--
	}
	return true
}
