package parser

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// DoWhileParser 表示 do-while 语句解析器
type DoWhileParser struct {
	*Parser
}

// NewDoWhileParser 创建一个新的 do-while 语句解析器
func NewDoWhileParser(parser *Parser) StatementParser {
	return &DoWhileParser{
		parser,
	}
}

// Parse 解析 do-while 语句
func (p *DoWhileParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	// 跳过 do 关键字
	p.next()

	// 解析循环体
	body, acl := p.parseBlock()
	if acl != nil {
		return nil, acl
	}

	// 检查是否有 while 关键字
	if p.current().Type() != token.WHILE {
		return nil, data.NewErrorThrow(p.FromCurrentToken(), errors.New("do-while 循环缺少 while 关键字"))
	}
	p.next() // 跳过 while

	// 解析条件表达式
	if p.current().Type() == token.LPAREN {
		p.next() // 跳过 (
	}

	condition, acl := p.parseStatement()
	if acl != nil {
		return nil, acl
	}

	if p.current().Type() == token.RPAREN {
		p.next() // 跳过 )
	}

	// 检查是否有分号
	if p.current().Type() == token.SEMICOLON {
		p.next() // 跳过 ;
	}

	from := tracker.EndBefore()

	return node.NewDoWhileStatement(
		from,
		condition,
		body,
	), nil
}
