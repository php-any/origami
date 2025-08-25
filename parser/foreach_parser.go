package parser

import (
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ForeachParser 表示foreach语句解析器
type ForeachParser struct {
	*Parser
}

// NewForeachParser 创建一个新的foreach语句解析器
func NewForeachParser(parser *Parser) StatementParser {
	return &ForeachParser{
		parser,
	}
}

// Parse 解析foreach语句
func (p *ForeachParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	// 跳过foreach关键字
	p.next()

	// 解析左括号
	if p.current().Type != token.LPAREN {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("foreach 缺少左括号"))
	}
	p.next()

	// 解析数组表达式
	exprParser := NewExpressionParser(p.Parser)
	array, acl := exprParser.Parse()
	if acl != nil {
		return nil, acl
	}
	if array == nil {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("foreach 中需要数组表达式"))
	}

	// 解析 as 关键字
	if p.current().Type != token.AS {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("foreach 中需要 'as' 关键字"))
	}
	p.next()

	// 解析键和值变量
	var key data.Variable
	var value data.Variable

	// 检查是否有键变量 (key => value)
	var ok bool
	keyTemp, acl := p.parseStatement()
	if acl != nil {
		return nil, acl
	}
	if key, ok = keyTemp.(data.Variable); !ok {
		if ident, ok := keyTemp.(*node.StringLiteral); ok {
			name := ident.Value
			val := p.scopeManager.CurrentScope().AddVariable(name, nil, tracker.EndBefore())
			key = node.NewVariableWithFirst(tracker.EndBefore(), val)
		} else {
			return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("foreach 中需要变量"))
		}
	}

	if p.current().Type == token.ARRAY_KEY_VALUE {
		p.next() // 跳过 =>
		keyTemp, acl := p.parseStatement()
		if acl != nil {
			return nil, acl
		}
		if value, ok = keyTemp.(data.Variable); !ok {
			if ident, ok := keyTemp.(*node.StringLiteral); ok {
				name := ident.Value
				val := p.scopeManager.CurrentScope().AddVariable(name, nil, tracker.EndBefore())
				value = node.NewVariableWithFirst(tracker.EndBefore(), val)
			} else {
				return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("foreach 中需要变量"))
			}
		}
	} else {
		value = key
		key = nil
	}

	// 解析右括号
	if p.current().Type != token.RPAREN {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("foreach 缺少右括号"))
	}
	p.next()

	// 解析循环体
	body := p.parseBlock()

	from := tracker.EndBefore()
	return node.NewForeachStatement(
		from,
		array,
		key,
		value,
		body,
	), nil
}
