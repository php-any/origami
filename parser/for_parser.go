package parser

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ForParser 表示for语句解析器
type ForParser struct {
	*Parser
}

// NewForParser 创建一个新的for语句解析器
func NewForParser(parser *Parser) StatementParser {
	return &ForParser{
		parser,
	}
}

// Parse 解析for语句
func (p *ForParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	// 跳过for关键字
	p.next()
	var acl data.Control
	// for $v in $arr {}
	if p.checkPositionIs(1, token.IN) || p.checkPositionIs(2, token.IN) || p.checkPositionIs(3, token.IN) || p.checkPositionIs(4, token.IN) {
		// 解析成 foreach 执行
		// 解析键和值变量
		var key data.Variable
		var value data.Variable
		// 解析初始化表达式
		hasLparen := false
		if p.current().Type == token.LPAREN {
			p.nextAndCheck(token.LPAREN)
			hasLparen = true
		}
		exprParser := NewMainStatementParser(p.Parser)

		var ok bool
		initializer, acl := exprParser.Parse()
		if acl != nil {
			return initializer, acl
		}
		if key, ok = initializer.(data.Variable); !ok {
			// 声明为变量
			if ident, ok := initializer.(*node.StringLiteral); ok {
				name := ident.Value
				val := p.scopeManager.CurrentScope().AddVariable(name, nil, tracker.EndBefore())
				key = node.NewVariableWithFirst(tracker.EndBefore(), val)
			} else {
				return nil, data.NewErrorThrow(tracker.EndBefore(), nil)
			}
		}
		if p.current().Type == token.COMMA {
			p.next()
			expr, acl := exprParser.Parse()
			if acl != nil {
				return expr, acl
			}
			if value, ok = expr.(data.Variable); !ok {
				// 声明为变量
				if ident, ok := expr.(*node.StringLiteral); ok {
					name := ident.Value
					varInfo := p.scopeManager.LookupVariable(name)
					if varInfo != nil {
						// 解析后续操作（函数调用、数组访问等）
						vp := &VariableParser{p.Parser}
						return vp.parseSuffix(varInfo)
					} else {
						val := p.scopeManager.CurrentScope().AddVariable(name, nil, tracker.EndBefore())
						value = node.NewVariableWithFirst(tracker.EndBefore(), val)
					}
				} else {
					return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("for in 无法解析变量 value"))
				}
			}
		} else {
			name := "_"
			val := p.scopeManager.CurrentScope().AddVariable(name, nil, tracker.EndBefore())
			value = node.NewVariableWithFirst(tracker.EndBefore(), val)
		}

		p.nextAndCheck(token.IN)

		array, acl := exprParser.Parse()

		if hasLparen {
			p.nextAndCheck(token.RPAREN)
		}

		// 解析循环体
		body := p.parseBlock()

		return node.NewForeachStatement(
			tracker.EndBefore(),
			array,
			key,
			value,
			body,
		), acl
	} else {
		// 解析初始化表达式
		hasLparen := false
		if p.current().Type == token.LPAREN {
			p.nextAndCheck(token.LPAREN)
			hasLparen = true
		}
		exprParser := NewMainStatementParser(p.Parser)
		var initializer node.Statement
		if p.checkPositionIs(0, token.SEMICOLON) {
			p.nextAndCheckStip(token.SEMICOLON) // 跳过分号
		} else if !p.checkPositionIs(0, token.LBRACE) {
			initializer, acl = exprParser.Parse()
			p.nextAndCheckStip(token.SEMICOLON) // 跳过分号
		}

		// 解析条件表达式
		var condition data.GetValue
		if p.checkPositionIs(0, token.SEMICOLON) {
			p.nextAndCheckStip(token.SEMICOLON) // 跳过分号
		} else if !p.checkPositionIs(0, token.LBRACE) {
			condition, acl = exprParser.Parse()
			p.nextAndCheckStip(token.SEMICOLON) // 跳过分号
		}

		// 解析递增表达式
		var increment data.GetValue
		if p.checkPositionIs(0, token.SEMICOLON) {
			p.nextAndCheckStip(token.SEMICOLON) // 跳过分号
		} else if !p.checkPositionIs(0, token.LBRACE, token.RPAREN) {
			increment, acl = exprParser.Parse()
			p.nextAndCheckStip(token.SEMICOLON) // 跳过分号
		}

		if hasLparen {
			p.nextAndCheck(token.RPAREN)
		}

		// 解析循环体
		body := p.parseBlock()

		return node.NewForStatement(
			tracker.EndBefore(),
			initializer,
			condition,
			increment,
			body,
		), acl
	}
}
