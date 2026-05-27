package parser

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/lexer"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// HeredocParser 解析 HEREDOC / NOWDOC 字面量（词法阶段已吞掉 <<< 至结束定界符）。
type HeredocParser struct {
	*Parser
}

func NewHeredocParser(parser *Parser) *HeredocParser {
	return &HeredocParser{Parser: parser}
}

// ParseLiteral 将 HEREDOC 或 NOWDOC token 解析为字符串字面量 AST。
func (hp *HeredocParser) ParseLiteral() (data.GetValue, data.Control) {
	tracker := hp.StartTracking()
	tok := hp.current()
	literal := tok.Literal()
	tt := tok.Type()

	if tt != token.HEREDOC && tt != token.NOWDOC {
		return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("期望 heredoc/nowdoc token，实际 %v", tt))
	}

	body, isNowdoc, ok := lexer.ExtractHeredocBody(literal)
	if !ok {
		return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("无效的 heredoc 字面量"))
	}
	if tt == token.NOWDOC {
		isNowdoc = true
	}

	hp.next()
	from := tracker.EndBefore()
	if isNowdoc {
		return node.NewNowdocLiteral(from, body), nil
	}
	return node.NewHeredocLiteral(from, body), nil
}
