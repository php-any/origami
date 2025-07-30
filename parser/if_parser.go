package parser

import (
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// IfParser 表示if语句解析器
type IfParser struct {
	*Parser
}

// NewIfParser 创建一个新的if语句解析器
func NewIfParser(parser *Parser) StatementParser {
	return &IfParser{
		parser,
	}
}

// Parse 解析if语句
func (p *IfParser) Parse() (data.GetValue, data.Control) {
	start := p.GetStart()
	// 跳过if关键字
	p.next()

	// 解析条件部分
	condition, acl := p.parseIfCondition()
	if acl != nil {
		return nil, acl
	}
	// 解析then分支
	thenBranch := p.parseBlock()

	// 解析 else if 和 else 分支
	var elseIfBranches []node.ElseIfBranch
	var elseBranch []data.GetValue

	if (p.checkPositionIs(0, token.ELSE) && p.checkPositionIs(1, token.IF)) || p.checkPositionIs(0, token.ELSE_IF) {
		// 解析多个 else if 分支
		for {
			if (p.checkPositionIs(0, token.ELSE) && p.checkPositionIs(1, token.IF)) || p.checkPositionIs(0, token.ELSE_IF) {
				if p.checkPositionIs(0, token.ELSE_IF) {
					p.next()
				} else {
					p.next()
					p.next()
				}
				// 解析 else if 的条件部分
				elseIfCondition, acl := p.parseIfCondition()
				if acl != nil {
					return nil, acl
				}
				// 解析 else if 的 then 分支
				elseIfThenBranch := p.parseBlock()

				elseIfBranches = append(elseIfBranches, node.ElseIfBranch{
					Condition:  elseIfCondition,
					ThenBranch: elseIfThenBranch,
				})

			} else {
				break
			}
		}
	}
	if p.checkPositionIs(0, token.ELSE) {
		// 这是 else 分支，不是 else if
		p.next()
		elseBranch = p.parseBlock()
	}

	return node.NewIfStatement(
		p.NewTokenFrom(start),
		condition,
		thenBranch,
		elseIfBranches,
		elseBranch,
	), nil
}

// parseIfCondition 解析if条件，支持多种语法形式
func (p *IfParser) parseIfCondition() (data.GetValue, data.Control) {
	// 检查是否是括号形式 if (condition)
	if p.current().Type == token.LPAREN {
		p.next() // 跳过左括号

		condition, acl := p.parseStatement()
		if acl != nil {
			p.addControl(acl)
		}
		if p.current().Type != token.RPAREN {
			return nil, data.NewErrorThrow(p.newFrom(), errors.New("if 缺少右括号 ')'"))
		}
		p.next() // 跳过右括号

		return condition, nil
	}

	// 检查是否是分号分隔的形式 if init; condition; increment {}
	// 这种形式类似于 for 循环的语法
	return p.parseSemicolonCondition()
}

// parseSemicolonCondition 解析分号分隔的条件形式
func (p *IfParser) parseSemicolonCondition() (data.GetValue, data.Control) {
	exprParser := NewExpressionParser(p.Parser)

	// 解析第一个表达式（初始化或条件）
	firstExpr, acl := exprParser.Parse()
	if acl != nil {
		return nil, acl
	}
	if firstExpr == nil {
		return nil, data.NewErrorThrow(p.newFrom(), errors.New("if 条件中缺少表达式"))
	}

	// 检查是否有分号
	if p.current().Type == token.SEMICOLON {
		// 这是 if init; condition; increment {} 形式
		p.next() // 跳过第一个分号

		// 解析条件表达式
		condition, acl := exprParser.Parse()
		if acl != nil {
			return nil, acl
		}
		if condition == nil {
			return nil, data.NewErrorThrow(p.newFrom(), errors.New("if 分号后缺少条件表达式"))
		}

		// 检查是否有第二个分号
		if p.current().Type == token.SEMICOLON {
			p.next() // 跳过第二个分号

			// 解析增量表达式（可选） TODO
			_, acl = exprParser.Parse()
			if acl != nil {
				return nil, acl
			}
			// 这里我们可以选择忽略增量表达式，或者将其作为条件的一部分
			// 为了简化，我们只返回条件表达式
			return condition, nil
		}

		// 只有两个表达式：if init; condition {}
		return condition, nil
	}

	// 只有一个表达式：if condition {}
	return firstExpr, nil
}
