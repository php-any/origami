package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// UnusedParser 表示unused语句解析器
type UnusedParser struct {
	*Parser
}

// NewUnusedParser 创建一个新的unused语句解析器
func NewUnusedParser(parser *Parser) StatementParser {
	return &UnusedParser{
		parser,
	}
}

// Parse 解析unused语句
func (p *UnusedParser) Parse() (data.GetValue, data.Control) {
	from := p.NewTokenFrom(p.GetStart())
	// 跳过 _ 符号
	name := "_"
	p.next()
	index := p.scopeManager.CurrentScope().AddVariable(name, nil, from)
	// 创建变量声明语句
	return node.NewVariable(
		from,
		name,
		index,
		nil,
	), nil
}
