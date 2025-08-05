package parser

import (
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// NewStructParser 表示 new 表达式解析器
type NewStructParser struct {
	*Parser
}

// NewNewParser 创建一个新的 new 表达式解析器
func NewNewParser(parser *Parser) StatementParser {
	return &NewStructParser{
		parser,
	}
}

// Parse 解析 new 表达式
func (p *NewStructParser) Parse() (data.GetValue, data.Control) {
	// 跳过 new 关键字
	p.next()

	// 解析类名
	if !p.checkPositionIs(0, token.IDENTIFIER, token.GENERIC_TYPE) {
		return nil, data.NewErrorThrow(p.newFrom(), errors.New("new关键字后面必须跟类名"))
	}

	isGenerated := false
	generateType := ""
	if p.checkPositionIs(0, token.GENERIC_TYPE) {
		isGenerated = true
		generateType = p.current().Literal
	}

	// 获取完整的类名路径
	className, acl := p.getClassName(true)
	if acl != nil {
		return nil, acl
	}

	// 解析参数列表
	vp := VariableParser{Parser: p.Parser}
	args, acl := vp.parseFunctionCall()
	if acl != nil {
		return nil, acl
	}

	n := node.NewNewExpression(
		p.FromCurrentToken(),
		className,
		args,
	)

	if p.checkPositionIs(0, token.OBJECT_OPERATOR) {
		// 解析链式调用
		return vp.parseSuffix(n)
	}

	if isGenerated {
		return &node.NewGenerated{
			NewExpression: n,
			T:             generateType,
		}, nil
	}

	// 创建 new 表达式节点
	return n, nil
}
