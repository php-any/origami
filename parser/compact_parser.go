package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// CompactParser 解析 compact 语句
type CompactParser struct {
	*Parser
}

// NewCompactParser 创建一个新的 compact 解析器
func NewCompactParser(parser *Parser) StatementParser {
	return &CompactParser{parser}
}

// Parse 解析 compact 语句
func (p *CompactParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	// 跳过 compact 关键字
	p.next()

	// 解析参数列表，compact('var1', 'var2', ...)
	// 需要解析括号内的参数
	if !p.checkPositionIs(0, token.LPAREN) {
		return nil, data.NewErrorThrow(tracker.EndBefore(), data.NewError(nil, "compact 后面必须跟括号", nil))
	}
	p.next() // 跳过 (

	varNames := []data.GetValue{}

	// 解析参数列表
	for !p.checkPositionIs(0, token.RPAREN) && !p.isEOF() {
		// 解析一个参数（应该是字符串字面量）
		arg, acl := p.parseStatement()
		if acl != nil {
			return nil, acl
		}
		if arg == nil {
			break
		}
		// 如果 arg 是字符串，自动转变量
		if argv, ok := arg.(*node.StringLiteral); ok {
			arg = p.scopeManager.LookupVariable(argv.Value)
		}
		varNames = append(varNames, arg)

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
	// 创建 compact 语句
	return node.NewCompactStatement(from, varNames), nil
}
