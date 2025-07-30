package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ReturnParser 表示return语句解析器
type ReturnParser struct {
	*Parser
}

// NewReturnParser 创建一个新的return语句解析器
func NewReturnParser(parser *Parser) StatementParser {
	return &ReturnParser{
		parser,
	}
}

// Parse 解析return语句
func (p *ReturnParser) Parse() (data.GetValue, data.Control) {
	start := p.GetStart()
	// 跳过return关键字
	p.next()

	// 解析返回值表达式
	var value data.GetValue
	if p.current().Type != token.SEMICOLON {
		exprParser := NewExpressionParser(p.Parser)
		var acl data.Control
		value, acl = exprParser.Parse()
		if acl != nil {
			return nil, acl
		}
		if p.checkPositionIs(0, token.COMMA) {
			values := []data.GetValue{value}
			for p.checkPositionIs(0, token.COMMA) {
				p.next()
				next, acl := exprParser.Parse()
				if acl != nil {
					return nil, acl
				}
				values = append(values, next)
			}
			return node.NewReturnsStatement(
				p.NewTokenFrom(start),
				values,
			), nil
		}
	}

	return node.NewReturnStatement(
		p.NewTokenFrom(start),
		value,
	), nil
}
