package parser

import (
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// MatchParser 表示match语句解析器
type MatchParser struct {
	*Parser
}

// NewMatchParser 创建一个新的match语句解析器
func NewMatchParser(parser *Parser) StatementParser {
	return &MatchParser{
		parser,
	}
}

// Parse 解析match语句
func (p *MatchParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()

	// 跳过match关键字
	p.next()

	// 解析条件表达式
	condition, acl := p.parseMatchCondition()
	if acl != nil {
		return nil, acl
	}
	// 解析左大括号
	p.nextAndCheck(token.LBRACE)

	// 解析match分支
	var arms []node.MatchArm
	var def []data.GetValue
	for p.current().Type != token.RBRACE && !p.isEOF() {
		if p.checkPositionIs(0, token.DEFAULT) {
			p.next()
			p.nextAndCheck(token.ARRAY_KEY_VALUE)
			if p.current().Type == token.LBRACE {
				// 手动解析代码块，不消耗右大括号
				def = p.parseBlock()
			} else {
				// 这是一个表达式
				stmt, acl := p.parseStatement()
				if acl != nil {
					return nil, acl
				}
				def = []data.GetValue{stmt}
			}
		} else {
			arm, acl := p.parseMatchArm()
			if acl != nil {
				return nil, acl
			}
			if arm != nil {
				arms = append(arms, *arm)
			} else {
				// 报告错误：期望 match arm 或 default
				p.addError("match 语句中期望匹配分支或 'default'")
				return nil, nil
			}
		}
	}

	// 解析右大括号
	p.nextAndCheck(token.RBRACE)

	return node.NewMatchStatement(
		tracker.EndBefore(),
		condition,
		arms,
		def,
	), nil
}

// parseMatchCondition 解析match条件表达式
func (p *MatchParser) parseMatchCondition() (data.GetValue, data.Control) {
	// 检查是否是括号形式 match (condition)
	if p.current().Type == token.LPAREN {
		p.next() // 跳过左括号

		condition, acl := p.parseStatement()
		if acl != nil {
			return nil, acl
		}
		if p.current().Type != token.RPAREN {
			return nil, data.NewErrorThrow(p.FromCurrentToken(), errors.New("match 缺少右括号 ')'"))
		}
		p.next() // 跳过右括号

		return condition, nil
	}

	// 直接解析表达式
	return p.parseStatement()
}

// parseMatchArm 解析单个match分支
func (p *MatchParser) parseMatchArm() (*node.MatchArm, data.Control) {
	tracker := p.StartTracking()

	// 解析条件部分（可以是多个条件，用逗号分隔）
	var conditions []data.GetValue
	for !p.checkPositionIs(0, token.EOF, token.ARRAY_KEY_VALUE) {
		condition, ok := p.parseValue()
		if !ok {
			return nil, data.NewErrorThrow(p.FromCurrentToken(), errors.New("match 左边必须是值"))
		}
		conditions = append(conditions, condition)
	}

	// 解析箭头 =>
	p.nextAndCheck(token.ARRAY_KEY_VALUE)

	// 解析表达式或代码块
	var expression data.GetValue
	var statements []data.GetValue

	if p.current().Type == token.LBRACE {
		// 手动解析代码块，不消耗右大括号
		statements = p.parseBlock()
	} else {
		// 这是一个表达式
		var acl data.Control
		expression, acl = p.parseStatement()
		if acl != nil {
			return nil, acl
		}
	}

	// 解析分号（可选）
	if p.checkPositionIs(0, token.SEMICOLON, token.COMMA) {
		p.next()
	}

	return &node.MatchArm{
		Node:       node.NewNode(tracker.EndBefore()),
		Conditions: conditions,
		Expression: expression,
		Statements: statements,
	}, nil
}
