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
	for p.current().Type() != token.RBRACE && !p.isEOF() {
		if p.checkPositionIs(0, token.DEFAULT) {
			p.next()
			p.nextAndCheck(token.ARRAY_KEY_VALUE)
			if p.current().Type() == token.LBRACE {
				def, acl = p.parseBlock()
				if acl != nil {
					return nil, acl
				}
			} else {
				stmt, acl := p.parseStatement()
				if acl != nil {
					return nil, acl
				}
				def = []data.GetValue{stmt}
			}
			if p.checkPositionIs(0, token.SEMICOLON, token.COMMA) {
				p.next()
			}
		} else {
			arm, acl := p.parseMatchArm()
			if acl != nil {
				return nil, acl
			}
			if arm != nil {
				arms = append(arms, *arm)
			} else {
				return nil, data.NewErrorThrow(p.FromCurrentToken(), errors.New("match 语句中期望匹配分支或 'default'"))
			}
		}
	}

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
	if p.current().Type() == token.LPAREN {
		p.next()

		condition, acl := p.parseStatement()
		if acl != nil {
			return nil, acl
		}
		if p.current().Type() != token.RPAREN {
			return nil, data.NewErrorThrow(p.FromCurrentToken(), errors.New("match 缺少右括号 ')'"))
		}
		p.next()
		return condition, nil
	}
	return p.parseStatement()
}

// parseMatchArm 解析单个match分支
func (p *MatchParser) parseMatchArm() (*node.MatchArm, data.Control) {
	tracker := p.StartTracking()

	var conditions []data.GetValue
	for !p.checkPositionIs(0, token.EOF, token.ARRAY_KEY_VALUE) {
		if p.checkPositionIs(0, token.ARRAY_KEY_VALUE) {
			break
		}

		// 统一使用 expressionParser 解析完整条件表达式，支持 instanceof + && + || 等组合
		condition, acl := p.expressionParser.Parse()
		if acl != nil {
			return nil, acl
		}
		if condition == nil {
			return nil, data.NewErrorThrow(p.FromCurrentToken(), errors.New("match 左边表达式不能为空"))
		}
		conditions = append(conditions, condition)

		if p.current().Type() == token.COMMA {
			p.next()
			continue
		}

		if p.current().Type() == token.ARRAY_KEY_VALUE {
			break
		}
	}

	p.nextAndCheck(token.ARRAY_KEY_VALUE)

	var expression data.GetValue
	var statements []data.GetValue

	if p.current().Type() == token.LBRACE {
		var acl data.Control
		statements, acl = p.parseBlock()
		if acl != nil {
			return nil, acl
		}
	} else {
		var acl data.Control
		expression, acl = p.parseStatement()
		if acl != nil {
			return nil, acl
		}
	}

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
