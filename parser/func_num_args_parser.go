package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

type FuncNumArgsParser struct {
	*Parser
}

func NewFuncNumArgsParser(parser *Parser) StatementParser {
	return &FuncNumArgsParser{parser}
}

func (p *FuncNumArgsParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	p.next() // 跳过 func_num_args 关键字

	// 期望左括号
	if p.current().Type() != token.LPAREN {
		return nil, data.NewErrorThrow(tracker.EndBefore(), nil)
	}
	p.next()

	// 期望右括号
	if p.current().Type() != token.RPAREN {
		return nil, data.NewErrorThrow(tracker.EndBefore(), nil)
	}
	p.next()

	return node.NewFuncNumArgs(tracker.EndBefore()), nil
}
