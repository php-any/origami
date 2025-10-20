package parser

import (
	"errors"
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

type LparenParser struct {
	*Parser
}

// NewLparenParser 类型转换 or lambda or 括号表达式
func NewLparenParser(parser *Parser) StatementParser {
	return &LparenParser{
		parser,
	}
}

// Parse 解析左括号后的内容
func (ep *LparenParser) Parse() (data.GetValue, data.Control) {
	tracking := ep.StartTracking()
	// 检查是否是类型转换: (string) $data
	if ep.isTypeCast() {
		ep.nextAndCheck(token.LPAREN) // 跳过左括号
		return ep.parseTypeCast(tracking)
	}

	// 检查是否是 Lambda 表达式: (a, b) => {}
	if ep.isLambdaExpression() {
		return ep.parseLambdaExpression(tracking)
	}

	// 检查是否是括号表达式: (a + b)
	ep.nextAndCheck(token.LPAREN) // 跳过左括号
	return ep.parseParenthesizedExpression(tracking)
}

// isTypeCast 检查是否是类型转换
func (ep *LparenParser) isTypeCast() bool {
	// 检查模式: (IDENTIFIER) EXPRESSION
	if ep.checkPositionIs(1, token.IDENTIFIER) &&
		ep.checkPositionIs(2, token.RPAREN) &&
		!ep.checkPositionIs(3, token.ARRAY_KEY_VALUE) &&
		ep.checkPositionIs(3, token.IDENTIFIER, token.VARIABLE, token.LPAREN, token.INT, token.FLOAT, token.STRING, token.NULL, token.TRUE, token.FALSE) {
		return true
	}
	return false
}

// parseTypeCast 解析类型转换
func (ep *LparenParser) parseTypeCast(tracking *PositionTracker) (data.GetValue, data.Control) {
	typeName := ep.current().Literal
	ep.next()                     // 跳过类型名
	ep.nextAndCheck(token.RPAREN) // 跳过右括号

	val, acl := ep.parseStatement()
	if acl != nil {
		return nil, acl
	}
	fn, ok := ep.vm.GetFunc(typeName)
	if !ok {
		return nil, data.NewErrorThrow(tracking.EndBefore(), errors.New("未定义的函数:"+typeName))
	}
	return node.NewCallExpression(tracking.EndBefore(), typeName, []data.GetValue{val}, fn), nil
}

// isLambdaExpression 检查是否是 Lambda 表达式
func (ep *LparenParser) isLambdaExpression() bool {
	// 检查是否包含 => 符号，需要正确处理括号嵌套
	pos := 1 // 从 ( 后面开始检查
	parenCount := 0

	for pos < len(ep.tokens)-ep.position {
		tokenType := ep.tokens[ep.position+pos].Type

		if tokenType == token.LPAREN {
			parenCount++
		} else if tokenType == token.RPAREN {
			parenCount--
			if parenCount < 0 {
				// 找到了匹配的右括号，检查后面是否有 =>
				if pos+1 < len(ep.tokens)-ep.position &&
					ep.tokens[ep.position+pos+1].Type == token.ARRAY_KEY_VALUE {
					return true
				}
				break
			}
		} else if tokenType == token.ARRAY_KEY_VALUE && parenCount == 0 {
			// 在括号平衡的情况下找到了 =>
			return true
		}
		pos++
	}
	return false
}

// parseLambdaExpression 解析 Lambda 表达式
func (ep *LparenParser) parseLambdaExpression(tracking *PositionTracker) (data.GetValue, data.Control) {
	fp := &FunctionParser{
		ep.Parser,
	}

	// 创建新的函数作用域
	fp.scopeManager.NewScope(true)

	// 解析参数列表
	params, acl := fp.parseParameters()
	if acl != nil {
		return nil, acl
	}
	fp.nextAndCheck(token.ARRAY_KEY_VALUE)

	// 解析函数体
	body, acl := fp.parseBlock()
	if acl != nil {
		return nil, acl
	}
	vars := fp.scopeManager.CurrentScope().GetVariables()

	// 弹出函数作用域
	fp.scopeManager.PopScope()

	parent := make(map[int]int)
	for _, parentVariable := range fp.scopeManager.CurrentScope().GetVariables() {
		for _, childVariable := range vars {
			if childVariable.GetName() == parentVariable.GetName() {
				parent[childVariable.GetIndex()] = parentVariable.GetIndex()
			}
		}
	}

	return node.NewLambdaExpression(
		tracking.EndBefore(),
		params,
		body,
		vars,
		parent,
	), nil
}

// parseParenthesizedExpression 解析括号表达式
func (ep *LparenParser) parseParenthesizedExpression(tracking *PositionTracker) (data.GetValue, data.Control) {
	// 解析括号内的表达式
	expr, acl := ep.parseStatement()
	if acl != nil {
		return nil, acl
	}
	// 检查是否有右括号
	if ep.current().Type != token.RPAREN {
		return nil, data.NewErrorThrow(tracking.EndBefore(), fmt.Errorf("缺少右括号 ')'"))
	}
	ep.next() // 跳过右括号

	// 将表达式包装为语句
	return expr, nil
}
