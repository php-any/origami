package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ContinueParser 表示continue语句解析器
type ContinueParser struct {
	*Parser
}

// NewContinueParser 创建一个新的continue语句解析器
func NewContinueParser(parser *Parser) StatementParser {
	return &ContinueParser{
		parser,
	}
}

// Parse 解析continue语句
func (p *ContinueParser) Parse() (data.GetValue, data.Control) {
	from := p.FromCurrentToken()
	// 跳过continue关键字
	p.next()

	// continue N：可选层级，最内层为 1
	levels := 1
	if p.checkPositionIs(0, token.INT, token.NUMBER) {
		n, err := parseInt(p.current().Literal())
		if err != nil || n < 1 {
			return nil, data.NewCompileFatal(from, "'continue' operator accepts only positive integers")
		}
		levels = n
		p.next()
	} else if !p.checkPositionIs(0, token.SEMICOLON, token.RBRACE) && !p.isEOF() {
		return nil, data.NewCompileFatal(from, "'continue' operator accepts only positive integers")
	}

	return node.NewContinueStatementWithLevel(
		from,
		levels,
	), nil
}
