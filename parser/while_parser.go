package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// WhileParser 表示while语句解析器
type WhileParser struct {
	*Parser
}

// NewWhileParser 创建一个新的while语句解析器
func NewWhileParser(parser *Parser) StatementParser {
	return &WhileParser{
		parser,
	}
}

// Parse 解析while语句
func (p *WhileParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	// 跳过while关键字
	p.next()

	if p.current().Type == token.LPAREN {
		p.next()
	}

	condition, acl := p.parseStatement()
	if acl != nil {
		return nil, acl
	}
	if p.current().Type == token.RPAREN {
		p.next()
	}

	// 解析循环体
	body, acl := p.parseBlock()
	if acl != nil {
		return nil, acl
	}
	from := tracker.EndBefore()

	return node.NewWhileStatement(
		from,
		condition,
		body,
	), nil
}
