package parser

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

type ListParser struct{ *Parser }

func NewListParser(p *Parser) StatementParser { return &ListParser{p} }
func (p *ListParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	p.next()
	if ctl := p.nextAndCheck(token.LPAREN); ctl != nil {
		return nil, ctl
	}
	var targets []data.GetValue
	for !p.checkPositionIs(0, token.RPAREN) {
		if p.checkPositionIs(0, token.EOF) {
			return nil, data.NewErrorThrow(p.FromCurrentToken(), fmt.Errorf("Unclosed list()"))
		}
		if p.checkPositionIs(0, token.COMMA) {
			targets = append(targets, data.NewNullValue())
			p.next()
			continue
		}
		start := p.position
		target, ctl := p.expressionParser.parseTernary()
		if ctl != nil {
			return nil, ctl
		}
		if target == nil || p.position == start {
			return nil, data.NewErrorThrow(p.FromCurrentToken(), fmt.Errorf("Invalid list target"))
		}
		targets = append(targets, target)
		if p.checkPositionIs(0, token.RPAREN) {
			break
		}
		if ctl := p.nextAndCheck(token.COMMA); ctl != nil {
			return nil, ctl
		}
	}
	p.next()
	return node.NewArray(tracker.EndBefore(), targets), nil
}
