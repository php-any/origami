package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type ClassMagicParser struct {
	*Parser
}

func NewClassMagicParser(parser *Parser) StatementParser {
	return &ClassMagicParser{parser}
}

func (p *ClassMagicParser) Parse() (data.GetValue, data.Control) {
	from := p.FromCurrentToken()
	p.next()

	className := p.currentClass
	if className == "" {
		className = ""
	}
	return node.NewStringLiteralByAst(from, className), nil
}
