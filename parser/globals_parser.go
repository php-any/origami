package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// GlobalsParser 解析 $GLOBALS 和 $_ENV
type GlobalsParser struct {
	*Parser
}

func NewGlobalsParser(p *Parser) StatementParser {
	return &GlobalsParser{p}
}

func (gp *GlobalsParser) Parse() (data.GetValue, data.Control) {
	tracker := gp.StartTracking()

	name := gp.current().Literal()

	// 跳过关键字
	gp.next()

	expr := node.NewGlobalsNode(tracker.EndBefore(), name)

	p := &VariableParser{
		Parser: gp.Parser,
	}

	return p.parseSuffix(expr)
}
