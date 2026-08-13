package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// FunctionMagicParser 处理 __FUNCTION__ 魔术常量
type FunctionMagicParser struct {
	*Parser
}

func NewFunctionMagicParser(parser *Parser) StatementParser {
	return &FunctionMagicParser{Parser: parser}
}

func (p *FunctionMagicParser) Parse() (data.GetValue, data.Control) {
	from := p.FromCurrentToken()
	p.next()
	return node.NewStringLiteralByAst(from, p.currentFunction), nil
}
