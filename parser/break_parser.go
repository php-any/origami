package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// BreakParser 表示break语句解析器
type BreakParser struct {
	*Parser
}

// NewBreakParser 创建一个新的break语句解析器
func NewBreakParser(parser *Parser) StatementParser {
	return &BreakParser{
		parser,
	}
}

// Parse 解析break语句
func (p *BreakParser) Parse() (data.GetValue, data.Control) {
	from := p.FromCurrentToken()
	p.next() // 跳过 break

	levels := 1
	if p.current().Type() == token.INT {
		if n, err := parseInt(p.current().Literal()); err == nil {
			levels = n
		}
		p.next()
	}

	return node.NewBreakStatementWithLevel(from, levels), nil
}

func parseInt(s string) (int, error) {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}
