package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
	"strconv"
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
// strict_types is retained on the compiled unit and declarations.
func (p *DeclareParser) Parse() (data.GetValue, data.Control) {
	strictDirective := false
	p.next() // 跳过 declare

	// 左括号
	if acl := p.nextAndCheck(token.LPAREN); acl != nil {
		return nil, acl
	}

	for !p.isEOF() && !p.checkPositionIs(0, token.RPAREN) {
		name := p.current().Literal()
		p.next()
		if acl := p.nextAndCheck(token.ASSIGN); acl != nil {
			return nil, acl
		}
		if name == "strict_types" {
			literal := p.current().Literal()
			value, err := strconv.ParseInt(literal, 0, 64)
			if (p.current().Type() != token.INT && p.current().Type() != token.NUMBER) || err != nil || (value != 0 && value != 1) {
				return nil, data.NewCompileFatal(p.newFrom(), "strict_types declaration must have 0 or 1 as its value")
			}
			if !p.declarationsOnly {
				return nil, data.NewCompileFatal(p.newFrom(), "strict_types declaration must be the very first statement in the script")
			}
			strictDirective = true
			p.strictTypes = value == 1
		}
		p.next()
		if p.checkPositionIs(0, token.COMMA) {
			p.next()
		}
	}
	if ctl := p.nextAndCheck(token.RPAREN); ctl != nil {
		return nil, ctl
	}
	if strictDirective && !p.checkPositionIs(0, token.SEMICOLON) {
		return nil, data.NewCompileFatal(p.newFrom(), "strict_types declaration must not use block mode")
	}
	if p.checkPositionIs(0, token.SEMICOLON) {
		p.next()
	}

	return node.NewTodo(), nil
}
