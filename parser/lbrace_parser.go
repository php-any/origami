package parser

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

type LbraceParser struct {
	*Parser
}

func NewLbraceParser(parser *Parser) StatementParser {
	return &LbraceParser{
		parser,
	}
}

func (ep *LbraceParser) Parse() (data.GetValue, data.Control) {
	tracker := ep.StartTracking()
	ep.next()

	// 检查是否为空对象 {}
	if ep.current().Type == token.RBRACE {
		ep.next()
		v := map[data.GetValue]data.GetValue{}
		from := tracker.EndBefore()
		return node.NewKv(from, v), nil
	}

	var expr data.GetValue
	var acl data.Control
	if ep.peek(1).Type == token.COLON {
		// 先匹配 key: 符号，如果是这个格式，那么是关键字也能作为 key
		expr = data.NewStringValue(ep.current().Literal)
		ep.next()
	} else {
		expr, acl = ep.parseStatement()
		if acl != nil {
			return nil, acl
		}
	}
	switch ep.current().Type {
	case token.ARRAY_KEY_VALUE: // => 对象定义
		v := map[data.GetValue]data.GetValue{}
		ep.next() // =>
		v[expr], acl = ep.parseStatement()
		if acl != nil {
			return nil, acl
		}
		for ep.current().Type != token.RBRACE {
			key, acl := ep.parseStatement()
			if acl != nil {
				return nil, acl
			}
			ep.next()
			val, acl := ep.parseStatement()
			if acl != nil {
				return nil, acl
			}
			v[key] = val
		}
		ep.next()
		from := tracker.EndBefore()
		return node.NewKv(from, v), nil
	case token.COLON: // : JSON 定义
		oldIdentTryString := ep.identTryString
		ep.identTryString = true

		v := map[data.GetValue]data.GetValue{}
		ep.nextAndCheck(token.COLON) // :
		v[expr], acl = ep.parseStatement()
		if acl != nil {
			return nil, acl
		}
		for ep.current().Type != token.RBRACE {
			if ep.checkPositionIs(0, token.COMMA) && ep.checkPositionIs(1, token.RBRACE) {
				ep.next()
				break
			}
			ep.nextAndCheck(token.COMMA)

			var key data.GetValue
			if ep.peek(1).Type == token.COLON {
				// 先匹配 key: 符号，如果是这个格式，那么是关键字也能作为 key
				key = data.NewStringValue(ep.current().Literal)
				ep.next()
			} else {
				key, acl = ep.parseStatement()
				if acl != nil {
					return nil, acl
				}
			}
			ep.nextAndCheck(token.COLON) // :

			val, acl := ep.parseStatement()
			if acl != nil {
				return nil, acl
			}
			v[key] = val
		}

		ep.identTryString = oldIdentTryString

		ep.nextAndCheck(token.RBRACE)
		from := tracker.EndBefore()
		return node.NewKv(from, v), nil
	case token.RBRACE:
		ep.next()
		v := map[data.GetValue]data.GetValue{}
		from := tracker.EndBefore()
		return node.NewKv(from, v), nil
	default:
		return nil, data.NewErrorThrow(ep.FromCurrentToken(), errors.New("TODO: 语法错误"+ep.current().Literal))
	}
}
