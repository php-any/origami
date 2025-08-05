package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// EchoParser 解析 echo 语句
type EchoParser struct {
	*Parser
}

// NewEchoParser 创建一个新的 echo 解析器
func NewEchoParser(parser *Parser) StatementParser {
	return &EchoParser{parser}
}

// Parse 解析 echo 语句
func (p *EchoParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	// 跳过 echo 关键字
	p.next()

	exprs := make([]data.GetValue, 1)

	// 解析表达式
	var acl data.Control
	exprs[0], acl = p.parseStatement()
	if acl != nil {
		return nil, acl
	}
	for p.current().Type == token.COMMA {
		p.next() // ,
		stmt, acl := p.parseStatement()
		if acl != nil {
			return nil, acl
		}
		exprs = append(exprs, stmt)
	}

	from := tracker.EndBefore()
	// 创建 echo 语句
	return node.NewEchoStatement(from, exprs), nil
}
