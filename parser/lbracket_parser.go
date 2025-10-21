package parser

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

type LbracketParser struct {
	*Parser
}

// NewLbracketParser [ 开头解析
func NewLbracketParser(parser *Parser) StatementParser {
	return &LbracketParser{
		parser,
	}
}

// Parse 解析函数声明
func (ep *LbracketParser) Parse() (data.GetValue, data.Control) {
	tracker := ep.StartTracking()
	ep.next()

	// 检查是否直接遇到右括号（空数组）
	if ep.current().Type == token.RBRACKET {
		ep.next()
		from := tracker.EndBefore()
		return node.NewArray(from, []data.GetValue{}), nil
	}

	// 解析第一个元素
	expr, acl := ep.parseStatement()
	if acl != nil {
		return nil, acl
	}

	if ep.current().Type == token.RBRACKET {
		// 只有一个元素
		ep.next()
		from := tracker.EndBefore()
		return node.NewArray(from, []data.GetValue{expr}), nil
	} else {
		switch ep.current().Type {
		case token.COMMA: // , 数组定义
			arr := make([]data.GetValue, 1)
			arr[0] = expr
			for ep.current().Type == token.COMMA {
				ep.next()
				if ep.checkPositionIs(0, token.RBRACKET) {
					continue
				}
				stmt, acl := ep.parseStatement()
				if acl != nil {
					return nil, acl
				}
				arr = append(arr, stmt)
			}
			if ep.current().Type == token.RBRACKET {
				ep.next()
			}
			from := tracker.EndBefore()
			return node.NewArray(from, arr), nil
		case token.ARRAY_KEY_VALUE: // => 对象定义
			v := map[data.GetValue]data.GetValue{}
			ep.next() // =>
			v[expr], acl = ep.parseStatement()
			if acl != nil {
				return nil, acl
			}
			ep.nextAndCheckStip(token.COMMA)

			for ep.current().Type != token.RBRACKET {
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
		case token.COLON: // : JSON 定义
			oldIdentTryString := ep.identTryString
			ep.identTryString = true

			v := map[data.GetValue]data.GetValue{}
			ep.next() // :
			v[expr], acl = ep.parseStatement()
			if acl != nil {
				return nil, acl
			}
			for ep.current().Type != token.RBRACKET {
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
			return nil, data.NewErrorThrow(ep.FromCurrentToken(), errors.New("TODO: 语法错误"))
		}
	}
}
