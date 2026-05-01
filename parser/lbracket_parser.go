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

	// 处理跳过的元素: [, $a] 或 [,, $a]
	if ep.current().Type() == token.COMMA {
		// 第一个元素被跳过，创建空占位并进入逗号分隔数组模式
		arr := []data.GetValue{data.NewNullValue()}
		for ep.current().Type() == token.COMMA {
			ep.next()
			if ep.checkPositionIs(0, token.RBRACKET) {
				arr = append(arr, data.NewNullValue())
				break
			}
			if ep.current().Type() == token.COMMA {
				arr = append(arr, data.NewNullValue())
				continue
			}
			if ep.current().Type() == token.ELLIPSIS {
				ep.next()
				spreadExpr, acl := ep.parseStatement()
				if acl != nil {
					return nil, acl
				}
				arr = append(arr, node.NewArraySpread(spreadExpr))
			} else {
				stmt, acl := ep.parseStatement()
				if acl != nil {
					return nil, acl
				}
				arr = append(arr, stmt)
			}
		}
		if ep.current().Type() == token.RBRACKET {
			ep.next()
		}
		from := tracker.EndBefore()
		return node.NewArray(from, arr), nil
	}

	// 检查第一个元素是否是展开运算符
	var expr data.GetValue
	var acl data.Control
	if ep.current().Type() == token.ELLIPSIS {
		ep.next() // 跳过 ...
		// 解析要展开的表达式
		spreadExpr, acl := ep.parseStatement()
		if acl != nil {
			return nil, acl
		}
		if spreadExpr == nil {
			return nil, data.NewErrorThrow(ep.FromCurrentToken(), errors.New("展开运算符后需要一个表达式"))
		}
		expr = node.NewArraySpread(spreadExpr)
	} else {
		// 解析第一个元素
		expr, acl = ep.parseStatement()
		if acl != nil {
			return nil, acl
		}
		if expr == nil {
			return nil, data.NewErrorThrow(ep.FromCurrentToken(), errors.New("数组第一个元素解析失败，无法识别当前 token"))
		}
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

			for ep.current().Type() == token.COMMA {
				ep.next()
				if ep.checkPositionIs(0, token.RBRACKET) {
					continue
				}
				// 检查是否是展开运算符
				if ep.current().Type() == token.ELLIPSIS {
					ep.next() // 跳过 ...
					// 解析要展开的表达式
					spreadExpr, acl := ep.parseStatement()
					if acl != nil {
						return nil, acl
					}
					if spreadExpr == nil {
						return nil, data.NewErrorThrow(ep.FromCurrentToken(), errors.New("展开运算符后需要一个表达式"))
					}
					// 创建展开节点
					arr = append(arr, node.NewArraySpread(spreadExpr))
				} else {
					stmt, acl := ep.parseStatement()
					if acl != nil {
						return nil, acl
					}
					if stmt == nil {
						return nil, data.NewErrorThrow(ep.FromCurrentToken(), fmt.Errorf("数组元素解析失败，无法识别当前 token: %s (类型: %d)", ep.current().Literal(), ep.current().Type()))
					}
					// 支持混合 key => value: [$val, $key => $val2]
					if ep.checkPositionIs(0, token.ARRAY_KEY_VALUE) {
						kv := make([]node.KvPair, 0, len(arr)+1)
						for i, v := range arr {
							kv = append(kv, node.KvPair{Key: data.NewIntValue(i), Value: v})
						}
						ep.next() // =>
						val, acl := ep.parseStatement()
						if acl != nil {
							return nil, acl
						}
						kv = append(kv, node.KvPair{Key: stmt, Value: val})
						// 后续全部按 KvPair 解析
						ep.nextAndCheckStip(token.COMMA)
						for ep.current().Type() != token.RBRACKET {
							key, acl := ep.parseStatement()
							if acl != nil {
								return nil, acl
							}
							ep.next()
							nextVal, acl := ep.parseStatement()
							if acl != nil {
								return nil, acl
							}
							kv = append(kv, node.KvPair{Key: key, Value: nextVal})
							ep.nextAndCheckStip(token.COMMA)
						}
						ep.next()
						from := tracker.EndBefore()
						return node.NewKv(from, kv), nil
					}
					arr = append(arr, stmt)
				}
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
