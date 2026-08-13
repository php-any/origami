package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// MethodMagicParser 处理 __METHOD__ 魔术常量
type MethodMagicParser struct {
	*Parser
}

func NewMethodMagicParser(parser *Parser) StatementParser {
	return &MethodMagicParser{Parser: parser}
}

func (p *MethodMagicParser) Parse() (data.GetValue, data.Control) {
	from := p.FromCurrentToken()
	p.next()
	method := p.currentFunction
	if p.currentClass != "" && method != "" {
		method = p.currentClass + "::" + method
	}
	return node.NewStringLiteralByAst(from, method), nil
}
