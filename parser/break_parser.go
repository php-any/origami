package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// BreakParser 表示break语句解析器
type BreakParser struct {
	*Parser
}

// NewBreakParser 创建一个新的break语句解析器
func NewBreakParser(parser *Parser) StatementParser {
	return &BreakParser{
		parser,
	}
}

// Parse 解析break语句
func (p *BreakParser) Parse() (data.GetValue, data.Control) {
	start := p.GetStart()
	// 跳过break关键字
	p.next()

	return node.NewBreakStatement(
		p.NewTokenFrom(start),
	), nil
}
