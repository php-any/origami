package parser

import (
	"errors"
	"fmt"

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
	if ep.current().Type() == token.RBRACKET {
		ep.next()
		from := tracker.EndBefore()
		return node.NewArray(from, []data.GetValue{}), nil
	}

	// 解析第一个元素
	expr, acl := ep.parseStatement()
	if acl != nil {
		return nil, acl
	}
	if expr == nil {
		return nil, data.NewErrorThrow(ep.FromCurrentToken(), errors.New("数组第一个元素解析失败，无法识别当前 token"))
	}

	if ep.current().Type() == token.RBRACKET {
		// 只有一个元素
		ep.next()
		from := tracker.EndBefore()
		return node.NewArray(from, []data.GetValue{expr}), nil
	} else {
		switch ep.current().Type() {
		case token.COMMA: // , 数组定义
			arr := make([]data.GetValue, 1)
			arr[0] = expr

			var test = 0
			for ep.current().Type() == token.COMMA {
				ep.next()
				if ep.checkPositionIs(0, token.RBRACKET) {
					continue
				}
				if *ep.source == "/Users/lvluo/Desktop/github.com/php-any/laravel/vendor/nunomaduro/collision/src/Adapters/Phpunit/Subscribers/EnsurePrinterIsRegisteredSubscriber.php" {
					test++
					if ep.position == 424 {
						ep.printRemaining()
						ttt := ep.current()
						fmt.Println(ttt)
					}
				}
				stmt, acl := ep.parseStatement()
				if acl != nil {
					return nil, acl
				}
				if stmt == nil {
					return nil, data.NewErrorThrow(ep.FromCurrentToken(), fmt.Errorf("数组元素解析失败，无法识别当前 token: %s (类型: %d)", ep.current().Literal(), ep.current().Type()))
				}
				arr = append(arr, stmt)
			}
			if ep.current().Type() == token.RBRACKET {
				ep.next()
			}
			from := tracker.EndBefore()
			return node.NewArray(from, arr), nil
		case token.ARRAY_KEY_VALUE: // => 对象定义
			v := []node.KvPair{}
			ep.next() // =>
			firstVal, acl := ep.parseStatement()
			if acl != nil {
				return nil, acl
			}
			v = append(v, node.KvPair{Key: expr, Value: firstVal})
			ep.nextAndCheckStip(token.COMMA)

			for ep.current().Type() != token.RBRACKET {
				key, acl := ep.parseStatement()
				if acl != nil {
					return nil, acl
				}
				ep.next()
				val, acl := ep.parseStatement()
				if acl != nil {
					return nil, acl
				}
				v = append(v, node.KvPair{Key: key, Value: val})
				ep.nextAndCheckStip(token.COMMA)
			}
			ep.next()
			from := tracker.EndBefore()
			return node.NewKv(from, v), nil
		case token.COLON: // : JSON 定义
			oldIdentTryString := ep.identTryString
			ep.identTryString = true

			v := []node.KvPair{}
			ep.next() // :
			firstVal, acl := ep.parseStatement()
			if acl != nil {
				return nil, acl
			}
			v = append(v, node.KvPair{Key: expr, Value: firstVal})
			for ep.current().Type() != token.RBRACKET {
				key, acl := ep.parseStatement()
				if acl != nil {
					return nil, acl
				}
				ep.next()
				val, acl := ep.parseStatement()
				if acl != nil {
					return nil, acl
				}
				v = append(v, node.KvPair{Key: key, Value: val})
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
