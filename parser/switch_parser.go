package parser

import (
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// SwitchParser 表示switch语句解析器
type SwitchParser struct {
	*Parser
}

// NewSwitchParser 创建一个新的switch语句解析器
func NewSwitchParser(parser *Parser) StatementParser {
	return &SwitchParser{
		parser,
	}
}

// Parse 解析switch语句
func (p *SwitchParser) Parse() (data.GetValue, data.Control) {
	start := p.GetStart()
	// 跳过switch关键字
	p.next()

	// 解析条件表达式
	condition, acl := p.parseSwitchCondition()
	if acl != nil {
		return nil, acl
	}

	// 解析左大括号
	p.nextAndCheck(token.LBRACE)

	// 解析switch分支
	var cases []node.SwitchCase
	var defaultCase []data.GetValue

	for p.current().Type != token.RBRACE && !p.isEOF() {
		if p.checkPositionIs(0, token.DEFAULT) {
			defaultCase = p.parseDefaultCase()
		} else if p.checkPositionIs(0, token.CASE) {
			caseStmt := p.parseSwitchCase()
			if caseStmt != nil {
				cases = append(cases, *caseStmt)
			}
		} else {
			// 报告错误：期望 case 或 default
			p.addError("switch 语句中期望 'case' 或 'default'")
			return nil, nil
		}
	}

	// 解析右大括号
	p.nextAndCheck(token.RBRACE)

	return node.NewSwitchStatement(
		p.NewTokenFrom(start),
		condition,
		cases,
		defaultCase,
	), nil
}

// parseSwitchCondition 解析switch条件表达式
func (p *SwitchParser) parseSwitchCondition() (data.GetValue, data.Control) {
	// 检查是否是括号形式 switch (condition)
	if p.current().Type == token.LPAREN {
		p.next() // 跳过左括号

		condition, acl := p.parseStatement()
		if acl != nil {
			return nil, acl
		}
		if p.current().Type != token.RPAREN {
			return nil, data.NewErrorThrow(p.newFrom(), errors.New("switch 缺少右括号 ')'"))
		}
		p.next() // 跳过右括号

		return condition, nil
	}

	// 直接解析表达式
	return p.parseStatement()
}

// parseSwitchCase 解析单个switch分支
func (p *SwitchParser) parseSwitchCase() *node.SwitchCase {
	start := p.GetStart()

	// 跳过case关键字
	p.next()

	// 解析case值
	caseValue, acl := p.parseStatement()
	if acl != nil {
		p.addControl(acl)
	}
	// 解析冒号
	p.nextAndCheck(token.COLON)

	// 解析case体（语句列表）
	var statements []data.GetValue

	// 解析直到遇到下一个case、default或右大括号
	for !p.isEOF() && !p.checkPositionIs(0, token.CASE, token.DEFAULT, token.RBRACE) {
		if p.current().Type == token.LBRACE {
			// 这是一个代码块
			statements = p.parseBlock()
			break
		} else if p.checkPositionIs(0, token.BREAK) {
			// 处理break语句
			stmt, acl := p.parseStatement()
			if acl != nil {
				p.addControl(acl)
			}
			if stmt != nil {
				statements = append(statements, stmt)
			}
			// break后通常跟着分号
			if p.checkPositionIs(0, token.SEMICOLON) {
				p.next()
			}
			// break后通常结束当前case
			break
		} else {
			// 这是一个表达式
			stmt, acl := p.parseStatement()
			if acl != nil {
				p.addControl(acl)
			}
			if stmt != nil {
				statements = append(statements, stmt)
			}
		}

		// 跳过分号（可选）
		if p.checkPositionIs(0, token.SEMICOLON) {
			p.next()
		}
	}

	return &node.SwitchCase{
		Node:       node.NewNode(p.NewTokenFrom(start)),
		CaseValue:  caseValue,
		Statements: statements,
	}
}

// parseDefaultCase 解析default分支
func (p *SwitchParser) parseDefaultCase() []data.GetValue {
	// 跳过default关键字
	p.next()

	// 解析冒号
	p.nextAndCheck(token.COLON)

	// 解析default体（语句列表）
	var statements []data.GetValue

	// 解析直到遇到下一个case或右大括号
	for !p.isEOF() && !p.checkPositionIs(0, token.CASE, token.RBRACE) {
		if p.current().Type == token.LBRACE {
			// 这是一个代码块
			statements = p.parseBlock()
			break
		} else if p.checkPositionIs(0, token.BREAK) {
			// 处理break语句
			stmt, acl := p.parseStatement()
			if acl != nil {
				p.addControl(acl)
			}
			if stmt != nil {
				statements = append(statements, stmt)
			}
			// break后通常跟着分号
			if p.checkPositionIs(0, token.SEMICOLON) {
				p.next()
			}
			// break后通常结束当前default
			break
		} else {
			// 这是一个表达式
			stmt, acl := p.parseStatement()
			if acl != nil {
				p.addControl(acl)
			}
			if stmt != nil {
				statements = append(statements, stmt)
			}
		}

		// 跳过分号（可选）
		if p.checkPositionIs(0, token.SEMICOLON) {
			p.next()
		}
	}

	return statements
}
