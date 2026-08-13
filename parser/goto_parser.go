package parser

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// GotoParser 解析 goto 语句
type GotoParser struct {
	*Parser
}

func NewGotoParser(p *Parser) StatementParser {
	return &GotoParser{p}
}

func (p *GotoParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	// 跳过 goto 关键字
	p.next()

	if p.current().Type() != token.IDENTIFIER {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("goto 语句后面需要跟随标签名"))
	}

	label := p.current().Literal()
	p.next()

	from := tracker.EndBefore()
	return node.NewGotoStatement(from, label), nil
}
