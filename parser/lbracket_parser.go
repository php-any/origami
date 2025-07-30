package parser

import (
	"errors"
	"github.com/php-any/origami/node"
)
import "github.com/php-any/origami/token"
import "github.com/php-any/origami/data"

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
	ep.next()
	expr, acl := ep.parseStatement()
	if acl != nil {
		return nil, acl
	}
	if ep.current().Type == token.RBRACKET {
		// 直接结束了
		ep.next()
		if expr == nil {
			return node.NewArray(ep.NewTokenFrom(ep.position), []node.Statement{}), nil
		}
		return node.NewArray(ep.NewTokenFrom(ep.position), []node.Statement{expr}), nil
	} else {
		switch ep.current().Type {
		case token.COMMA: // , 数组定义
			arr := make([]node.Statement, 1)
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
			return node.NewArray(ep.NewTokenFrom(ep.position), arr), nil
		case token.ARRAY_KEY_VALUE: // => 对象定义
			v := map[node.Statement]node.Statement{}
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
			return node.NewKv(ep.NewTokenFrom(ep.position), v), nil
		case token.COLON: // : JSON 定义
			oldIdentTryString := ep.identTryString
			ep.identTryString = true

			v := map[node.Statement]node.Statement{}
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
			return node.NewKv(ep.NewTokenFrom(ep.position), v), nil
		default:
			return nil, data.NewErrorThrow(ep.NewTokenFrom(ep.position), errors.New("TODO: 语法错误"))
		}
	}
}
