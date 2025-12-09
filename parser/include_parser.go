package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// IncludeParser 解析 include/require/include_once/require_once 语句
type IncludeParser struct {
	*Parser
}

func NewIncludeParser(p *Parser) StatementParser {
	return &IncludeParser{p}
}

func (p *IncludeParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	kw := p.current().Type()
	p.next() // 跳过关键字

	// 支持 include expr; 也支持 include(expr);
	var expr data.GetValue
	var acl data.Control
	if p.current().Type() == token.LPAREN {
		p.next()
		expr, acl = p.parseStatement()
		if acl != nil {
			return nil, acl
		}
		if acl = p.nextAndCheck(token.RPAREN); acl != nil {
			return nil, acl
		}
	} else {
		expr, acl = p.parseStatement()
		if acl != nil {
			return nil, acl
		}
	}

	from := tracker.EndBefore()
	return node.NewIncludeStatement(
		from,
		expr,
		kw == token.INCLUDE_ONCE || kw == token.REQUIRE_ONCE,
		kw == token.REQUIRE || kw == token.REQUIRE_ONCE,
	), nil
}
