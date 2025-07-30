package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// VarParser 表示var语句解析器
type VarParser struct {
	*Parser
}

// NewVarParser 创建一个新的var语句解析器
func NewVarParser(parser *Parser) StatementParser {
	return &VarParser{
		parser,
	}
}

// Parse 解析var语句
func (p *VarParser) Parse() (data.GetValue, data.Control) {
	start := p.GetStart()
	// 跳过var关键字
	p.next()

	// 解析变量名
	if p.current().Type != token.IDENTIFIER {
		p.addError("Expected variable name after var")
		return nil, nil
	}
	name := p.current().Literal
	p.next()

	// 解析初始化表达式
	var initializer data.GetValue
	if p.current().Type == token.ASSIGN {
		p.next() // 跳过等号
		exprParser := NewExpressionParser(p.Parser)
		var acl data.Control
		initializer, acl = exprParser.Parse()
		if acl != nil {
			return nil, acl
		}
	}

	return node.NewVarStatement(
		p.NewTokenFrom(start),
		name,
		initializer,
	), nil
}
