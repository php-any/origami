package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// DeclareParser 解析 declare(...) 语句
type DeclareParser struct {
	*Parser
}

func NewDeclareParser(p *Parser) StatementParser {
	return &DeclareParser{p}
}

// Parse 解析 declare 语句
// 语法：declare(strict_types=1);
// 目前只做解析，不实现具体逻辑
func (p *DeclareParser) Parse() (data.GetValue, data.Control) {
	p.next() // 跳过 declare

	// 左括号
	if acl := p.nextAndCheck(token.LPAREN); acl != nil {
		return nil, acl
	}

	for !p.isEOF() && !p.checkPositionIs(0, token.RPAREN) {
		p.next()
	}
	p.next()
	if p.checkPositionIs(0, token.SEMICOLON) {
		p.next()
	}

	return node.NewTodo(), nil
}
