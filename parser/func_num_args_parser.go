package parser

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// FuncNumArgsParser 表示 func_num_args 关键字解析器
type FuncNumArgsParser struct {
	*Parser
}

// NewFuncNumArgsParser 创建一个新的 func_num_args 解析器
func NewFuncNumArgsParser(parser *Parser) StatementParser {
	return &FuncNumArgsParser{
		parser,
	}
}

// Parse 解析 func_num_args 关键字
func (p *FuncNumArgsParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	p.next() // 跳过 func_num_args 关键字

	// 检查是否有括号 func_num_args()
	if p.checkPositionIs(0, token.LPAREN) {
		p.next() // 跳过左括号

		if p.current().Type() != token.RPAREN {
			return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("func_num_args() 不接受参数"))
		}

		p.next() // 跳过右括号
	}

	from := tracker.EndBefore()
	return node.NewFuncNumArgs(from), nil
}
