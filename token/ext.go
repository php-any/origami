package token

func NewKeyword(word string, t TokenType) {
	TokenDefinitions = append(TokenDefinitions, TokenDefinition{
		Type:     t,
		Literal:  word,
		WordType: KEYWORD,
	})
}

func NewOperator(word string, t TokenType) {
	TokenDefinitions = append(TokenDefinitions, TokenDefinition{
		Type:     t,
		Literal:  word,
		WordType: OPERATOR,
	})
}
