package parser

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
	"strings"
)

func (p *ClassParser) parsePropertyHooks(name string, propertyType data.Types) (get, set data.Method, virtual bool, ctl data.Control) {
	p.next()
	backed := false
	for !p.currentIsTypeOrEOF(token.RBRACE) {
		if p.checkPositionIs(0, token.SEMICOLON) {
			p.next()
			continue
		}
		kind := p.current().Literal()
		if kind != "get" && kind != "set" {
			return nil, nil, false, data.NewCompileFatal(p.newFrom(), "Expected get or set property hook")
		}
		if kind == "get" && get != nil || kind == "set" && set != nil {
			return nil, nil, false, data.NewCompileFatal(p.newFrom(), "Duplicate property hook")
		}
		from := p.StartTracking()
		p.next()
		p.scopeManager.NewScope(false)
		var params []data.GetValue
		if kind == "set" {
			if p.checkPositionIs(0, token.LPAREN) {
				params, ctl = p.ParseParameters()
			} else {
				variable := p.scopeManager.CurrentScope().AddVariable("value", propertyType, p.newFrom())
				params = []data.GetValue{node.NewParameter(p.newFrom(), "value", variable.GetIndex(), nil, propertyType)}
			}
		}
		if ctl != nil {
			p.scopeManager.PopScope()
			return nil, nil, false, ctl
		}
		start := p.position
		var body []data.GetValue
		if p.checkPositionIs(0, token.ARRAY_KEY_VALUE) {
			p.next()
			expr, err := NewExpressionParser(p.Parser).Parse()
			ctl = err
			if kind == "get" {
				body = []data.GetValue{node.NewReturnStatement(from.EndBefore(), expr)}
			} else {
				backed = true
				body = []data.GetValue{node.NewBinaryAssign(from.EndBefore(), node.NewObjectProperty(from.EndBefore(), node.NewThis(from.EndBefore()), strings.TrimPrefix(name, "$")), expr)}
			}
			if ctl == nil {
				ctl = p.nextAndCheck(token.SEMICOLON)
			}
		} else if p.checkPositionIs(0, token.LBRACE) {
			body, ctl = p.ParseFunctionBody()
		} else {
			ctl = data.NewCompileFatal(p.newFrom(), "Expected property hook body")
		}
		for i := start; i+2 < p.position; i++ {
			if p.tokens[i].Literal() == "$this" && p.tokens[i+1].Type() == token.OBJECT_OPERATOR && p.tokens[i+2].Literal() == strings.TrimPrefix(name, "$") {
				backed = true
			}
		}
		variables := p.GetVariables()
		p.scopeManager.PopScope()
		if ctl != nil {
			return nil, nil, false, ctl
		}
		var ret data.Types
		if kind == "get" {
			ret = propertyType
		}
		method := node.NewMethod(from.EndBefore(), fmt.Sprintf("$%s::%s", strings.TrimPrefix(name, "$"), kind), "public", false, params, body, variables, ret, p.strictTypes)
		if kind == "get" {
			get = method
		} else {
			set = method
		}
	}
	ctl = p.nextAndCheck(token.RBRACE)
	return get, set, !backed, ctl
}
