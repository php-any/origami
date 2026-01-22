package parser

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// FuncGetArgsParser 表示 func_get_args 关键字解析器
type FuncGetArgsParser struct {
	*Parser
}

// NewFuncGetArgsParser 创建一个新的 func_get_args 解析器
func NewFuncGetArgsParser(parser *Parser) StatementParser {
	return &FuncGetArgsParser{
		parser,
	}
}

// Parse 解析 func_get_args 关键字，支持 func_get_args 和 func_get_args() 两种语法
func (p *FuncGetArgsParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	// 跳过 func_get_args 关键字
	p.next()

	// 检查是否有括号 func_get_args()
	if p.checkPositionIs(0, token.LPAREN) {
		p.next() // 跳过左括号

		// func_get_args 不接受参数，如果右括号前有内容，这是语法错误
		if p.current().Type() != token.RPAREN {
			return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("func_get_args() 不接受参数"))
		}

		p.next() // 跳过右括号
	}

	from := tracker.EndBefore()
	return node.NewFuncGetArgs(from), nil
}
