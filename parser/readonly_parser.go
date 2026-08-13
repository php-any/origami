package parser

import (
	"github.com/php-any/origami/data"
)

// ReadonlyParser 解析 readonly 关键字（PHP 8.1 属性 / PHP 8.2 readonly class）。
// 当前跳过修饰符后继续解析后续语句（通常是 class / final class）。
type ReadonlyParser struct {
	*Parser
}

func NewReadonlyParser(parser *Parser) StatementParser {
	return &ReadonlyParser{parser}
}

func (p *ReadonlyParser) Parse() (data.GetValue, data.Control) {
	p.next() // skip readonly

	if parser, ok := parserRouter[p.current().Type()]; ok {
		return parser(p.Parser).Parse()
	}

	return p.parseStatement()
}
