package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// YieldParser 表示yield语句解析器
type YieldParser struct {
	*Parser
}

// NewYieldParser 创建一个新的yield语句解析器
func NewYieldParser(parser *Parser) StatementParser {
	return &YieldParser{
		parser,
	}
}

// Parse 解析yield语句
func (p *YieldParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	// 跳过yield关键字
	p.next()

	// 检查是否是 yield from
	if p.current().Type() == token.FROM {
		// yield from expression
		p.next() // 跳过 from

		// 解析表达式
		exprParser := NewExpressionParser(p.Parser)
		expr, acl := exprParser.Parse()
		if acl != nil {
			return nil, acl
		}

		from := tracker.EndBefore()
		return node.NewYieldFromStatement(from, expr), nil
	}

	// 普通 yield 语句
	// yield key => value 或 yield value
	var key data.GetValue
	var value data.GetValue

	// 解析第一个表达式（可能是key或value）
	exprParser := NewExpressionParser(p.Parser)
	firstExpr, acl := exprParser.Parse()
	if acl != nil {
		return nil, acl
	}

	// 检查是否有 => 操作符（键值对）
	if p.checkPositionIs(0, token.ARRAY_KEY_VALUE) {
		// yield key => value 格式
		key = firstExpr
		p.next() // 跳过 =>
		value, acl = exprParser.Parse()
		if acl != nil {
			return nil, acl
		}
	} else {
		// yield value 格式
		value = firstExpr
	}

	from := tracker.EndBefore()
	return node.NewYieldStatement(from, key, value), nil
}
