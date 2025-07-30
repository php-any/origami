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
	start := p.GetStart()
	// 跳过foreach关键字
	p.next()

	// 解析左括号
	if p.current().Type != token.LPAREN {
		p.addError("foreach 后面需要 '('")
		return nil, nil
	}
	p.next()

	// 解析数组表达式
	exprParser := NewExpressionParser(p.Parser)
	array, acl := exprParser.Parse()
	if acl != nil {
		return nil, acl
	}
	if array == nil {
		return nil, data.NewErrorThrow(p.NewTokenFrom(start), errors.New("foreach 中需要数组表达式"))
	}

	// 解析 as 关键字
	if p.current().Type != token.AS {
		return nil, data.NewErrorThrow(p.NewTokenFrom(start), errors.New("foreach 中需要 'as' 关键字"))
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
			index := p.scopeManager.CurrentScope().AddVariable(name, nil, p.NewTokenFrom(start))
			key = node.NewVariable(p.NewTokenFrom(start), name, index, nil)
		} else {
			p.addError("foreach 无法解析变量 key")
			return nil, nil
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
				index := p.scopeManager.CurrentScope().AddVariable(name, nil, p.NewTokenFrom(start))
				value = node.NewVariable(p.NewTokenFrom(start), name, index, nil)
			} else {
				p.addError("foreach 无法解析变量 key")
				return nil, nil
			}
		}
	} else {
		value = key
		key = nil
	}

	// 解析右括号
	if p.current().Type != token.RPAREN {
		p.addError("foreach 变量后面需要 ')'")
		return nil, nil
	}
	p.next()

	// 解析循环体
	body := p.parseBlock()

	return node.NewForeachStatement(
		p.NewTokenFrom(start),
		array,
		key,
		value,
		body,
	), nil
}
