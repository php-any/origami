package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// IssetParser 解析 isset 语句
type IssetParser struct {
	*Parser
}

// NewIssetParser 创建一个新的 isset 解析器
func NewIssetParser(parser *Parser) StatementParser {
	return &IssetParser{parser}
}

// Parse 解析 isset 语句
func (p *IssetParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	// 跳过 isset 关键字
	p.next()

	// 解析参数列表，isset($var1, $var2, ...)
	// 需要解析括号内的参数
	if !p.checkPositionIs(0, token.LPAREN) {
		return nil, data.NewErrorThrow(tracker.EndBefore(), data.NewError(nil, "isset 后面必须跟括号", nil))
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
	// 创建 isset 语句
	return node.NewIssetStatement(from, args), nil
}
