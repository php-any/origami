package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// UnsetParser 解析 unset 语句
type UnsetParser struct {
	*Parser
}

// NewUnsetParser 创建一个新的 unset 解析器
func NewUnsetParser(parser *Parser) StatementParser {
	return &UnsetParser{parser}
}

// Parse 解析 unset 语句
func (p *UnsetParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	// 跳过 unset 关键字
	p.next()

	// 解析参数列表，unset($var1, $var2, ...)
	// 需要解析括号内的参数
	if !p.checkPositionIs(0, token.LPAREN) {
		return nil, data.NewErrorThrow(tracker.EndBefore(), data.NewError(nil, "unset 后面必须跟括号", nil))
	}
	p.next() // 跳过 (

	args := []data.GetValue{}

	// 解析参数列表
	for !p.checkPositionIs(0, token.RPAREN) && !p.isEOF() {
		// 解析一个参数（变量表达式）
		arg, acl := p.parseStatement()
		if acl != nil {
			return nil, acl
		}
		if arg == nil {
			break
		}
		args = append(args, arg)

		// 检查是否有逗号
		if p.checkPositionIs(0, token.COMMA) {
			p.next() // 跳过逗号
		} else {
			break
		}
	}

	// 跳过右括号
	if p.checkPositionIs(0, token.RPAREN) {
		p.next()
	}

	from := tracker.EndBefore()
	// 创建 unset 语句
	return node.NewUnsetStatement(from, args), nil
}
