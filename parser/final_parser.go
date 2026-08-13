package parser

import (
	"github.com/php-any/origami/data"
)

// FinalParser 表示 final 关键字解析器
// 目前只跳过 final 关键字，不做其他处理
type FinalParser struct {
	*Parser
}

// NewFinalParser 创建一个新的 final 解析器
func NewFinalParser(parser *Parser) StatementParser {
	return &FinalParser{
		parser,
	}
}

// Parse 解析 final 关键字
// 目前只跳过 final 关键字，然后继续解析下一个语句
func (p *FinalParser) Parse() (data.GetValue, data.Control) {
	// 跳过 final 关键字
	p.next()

	// 继续解析下一个语句（通常是 class）
	if parser, ok := parserRouter[p.current().Type()]; ok {
		return parser(p.Parser).Parse()
	}

	// 如果没有找到对应的解析器，使用通用解析
	return p.parseStatement()
}
