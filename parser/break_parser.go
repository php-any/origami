package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
	"strconv"
	"strings"
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
	if p.checkPositionIs(0, token.INT, token.NUMBER) {
		n, err := parseInt(p.current().Literal())
		if err != nil || n < 1 {
			return nil, data.NewCompileFatal(from, "'break' operator accepts only positive integers")
		}
		levels = n
		p.next()
	} else if !p.checkPositionIs(0, token.SEMICOLON, token.RBRACE) && !p.isEOF() {
		return nil, data.NewCompileFatal(from, "'break' operator accepts only positive integers")
	}

	return node.NewBreakStatementWithLevel(from, levels), nil
}

func parseInt(s string) (int, error) {
	s = strings.ReplaceAll(s, "_", "")
	n, err := strconv.ParseInt(s, 0, strconv.IntSize)
	return int(n), err
}
