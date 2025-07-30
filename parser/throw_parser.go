package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ThrowParser 表示throw语句解析器
type ThrowParser struct {
	*Parser
}

// NewThrowParser 创建一个新的throw语句解析器
func NewThrowParser(parser *Parser) StatementParser {
	return &ThrowParser{
		Parser: parser,
	}
}

// Parse 解析throw语句
func (p *ThrowParser) Parse() (data.GetValue, data.Control) {
	start := p.GetStart()
	// 跳过throw关键字
	p.next()

	// 解析要抛出的表达式
	var value data.GetValue
	var acl data.Control
	if p.current().Type != token.SEMICOLON {
		exprParser := NewExpressionParser(p.Parser)
		value, acl = exprParser.Parse()
		if acl != nil {
			return nil, acl
		}
	} else {
		// 如果没有表达式，创建一个默认的异常
		value = node.NewStringLiteral(p.NewTokenFrom(start), "Exception")
	}

	return node.NewThrowStatement(
		p.NewTokenFrom(start),
		value,
	), nil
}
