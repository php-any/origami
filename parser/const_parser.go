package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ConstParser 表示const语句解析器
type ConstParser struct {
	*Parser
}

// NewConstParser 创建一个新的const语句解析器
func NewConstParser(parser *Parser) StatementParser {
	return &ConstParser{
		parser,
	}
}

// Parse 解析const语句
func (p *ConstParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	// 跳过const关键字
	p.next()

	// 解析常量名
	if p.current().Type != token.IDENTIFIER {
		p.addError("Expected constant name after const")
		return nil, nil
	}
	name := p.current().Literal
	p.next()

	// 解析初始化表达式（常量必须初始化）
	if p.current().Type != token.ASSIGN {
		p.addError("Expected '=' after constant name")
		return nil, nil
	}
	p.next() // 跳过等号
	initializer, acl := p.parseStatement()
	from := tracker.EndBefore()

	return node.NewConstStatement(
		from,
		name,
		initializer,
	), acl
}
