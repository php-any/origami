package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"strconv"
)

type LineParser struct {
	*Parser
}

func NewLineParser(parser *Parser) StatementParser {
	return &LineParser{
		Parser: parser,
	}
}

func (p *LineParser) Parse() (data.GetValue, data.Control) {
	from := p.FromCurrentToken()

	// 移动到下一个 token
	p.next()

	// 返回目录路径的字符串字面量
	return node.NewIntLiteral(from, strconv.Itoa(from.SourceFrom.Line())), nil
}
