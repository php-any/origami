package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
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

	return node.NewContinueStatement(
		from,
	), nil
}
