package parser

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ArrayParser 解析 array(...) 数组字面量
type ArrayParser struct {
	*Parser
}

func NewArrayParser(parser *Parser) StatementParser {
	return &ArrayParser{parser}
}

func (ep *ArrayParser) Parse() (data.GetValue, data.Control) {
	tracker := ep.StartTracking()
	// 跳过 array
	ep.next()
	// 期待 (
	if acl := ep.nextAndCheck(token.LPAREN); acl != nil {
		return nil, acl
	}

	// 空数组：array()
	if ep.current().Type() == token.RPAREN {
		ep.next()
		from := tracker.EndBefore()
		return node.NewArray(from, []data.GetValue{}), nil
	}

	// 解析第一个元素
	expr, acl := ep.parseStatement()
	if acl != nil {
		return nil, acl
	}

	switch ep.current().Type() {
	case token.RPAREN:
		ep.next()
		from := tracker.EndBefore()
		return node.NewArray(from, []data.GetValue{expr}), nil
	case token.COMMA:
		arr := []data.GetValue{expr}
		for ep.current().Type() == token.COMMA {
			ep.next()
			if ep.checkPositionIs(0, token.RPAREN) {
				continue
			}
			stmt, acl := ep.parseStatement()
			if acl != nil {
				return nil, acl
			}
			arr = append(arr, stmt)
		}
		if ep.current().Type() == token.RPAREN {
			ep.next()
		}
		from := tracker.EndBefore()
		return node.NewArray(from, arr), nil
	case token.ARRAY_KEY_VALUE: // => 关联数组
		v := map[data.GetValue]data.GetValue{}
		ep.next() // =>
		v[expr], acl = ep.parseStatement()
		if acl != nil {
			return nil, acl
		}
		ep.nextAndCheckStip(token.COMMA)

		for ep.current().Type() != token.RPAREN {
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
			ep.nextAndCheckStip(token.COMMA)
		}
		ep.next()
		from := tracker.EndBefore()
		return node.NewKv(from, v), nil
	case token.COLON: // JSON 风格键值
		oldIdentTryString := ep.identTryString
		ep.identTryString = true

		v := map[data.GetValue]data.GetValue{}
		ep.next() // :
		v[expr], acl = ep.parseStatement()
		if acl != nil {
			return nil, acl
		}
		for ep.current().Type() != token.RPAREN {
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

		ep.identTryString = oldIdentTryString
		ep.next()
		from := tracker.EndBefore()
		return node.NewKv(from, v), nil
	default:
		return nil, data.NewErrorThrow(ep.FromCurrentToken(), errors.New("array() 语法错误"))
	}
}
