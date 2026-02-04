package parser

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// CloneParser 解析 clone 表达式
type CloneParser struct {
	*Parser
}

// NewCloneParser 创建一个新的 clone 解析器
func NewCloneParser(parser *Parser) StatementParser {
	return &CloneParser{parser}
}

// Parse 解析 clone 表达式
// 语法：clone <expression>
func (p *CloneParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()

	// 跳过 clone 关键字
	p.next()

	// 使用表达式解析器解析后续表达式，保证和普通表达式一致的优先级/语义
	exprParser := NewExpressionParser(p.Parser)
	expr, acl := exprParser.Parse()
	if acl != nil {
		return nil, acl
	}
	if expr == nil {
		return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("clone 后面必须跟随待克隆的表达式"))
	}

	from := tracker.EndBefore()
	return node.NewCloneExpression(from, expr), nil
}

