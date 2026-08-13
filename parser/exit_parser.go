package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ExitParser 解析 exit / die 语句。
type ExitParser struct {
	*Parser
}

func NewExitParser(parser *Parser) StatementParser {
	return &ExitParser{parser}
}

func (p *ExitParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	p.next() // skip exit/die

	var value data.GetValue
	switch p.current().Type() {
	case token.SEMICOLON, token.RBRACE, token.RPAREN:
		// bare exit;
	case token.LPAREN:
		p.next()
		if p.current().Type() != token.RPAREN {
			expr, acl := NewExpressionParser(p.Parser).Parse()
			if acl != nil {
				return nil, acl
			}
			value = expr
		}
		if p.current().Type() == token.RPAREN {
			p.next()
		}
	default:
		expr, acl := NewExpressionParser(p.Parser).Parse()
		if acl != nil {
			return nil, acl
		}
		value = expr
	}

	return node.NewExitStatement(tracker.EndBefore(), value), nil
}
