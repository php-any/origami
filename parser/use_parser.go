package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
	"strings"
)

// UseParser 表示use语句解析器
type UseParser struct {
	*Parser
}

// NewUseParser 创建一个新的use语句解析器
func NewUseParser(parser *Parser) StatementParser {
	return &UseParser{
		parser,
	}
}

// Parse 解析use语句
func (p *UseParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	// 跳过 use 关键字
	p.next()

	// 解析命名空间
	if p.current().Type != token.IDENTIFIER {
		p.addError("Expected namespace after 'use'")
		return nil, nil
	}

	// 获取完整的命名空间路径
	namespace := p.current().Literal
	p.next()

	// 检查是否有 as 关键字
	var alias string
	if p.current().Type == token.AS {
		p.next()
		if p.current().Type != token.IDENTIFIER {
			p.addError("Expected identifier after 'as'")
			return nil, nil
		}
		alias = p.current().Literal
		p.next()
	} else {
		// 如果没有 as 关键字，从完整路径中提取最后一个部分作为 alias
		// 例如：从 "a\b\c" 中提取 "c"
		parts := strings.Split(namespace, "\\")
		alias = parts[len(parts)-1]
	}

	// 检查分号
	if p.current().Type != token.SEMICOLON {
		p.addError("Expected ';' after use statement")
		return nil, nil
	}
	p.next()

	from := tracker.EndBefore()
	// 创建 use 语句节点
	return node.NewUseStatement(
		from,
		namespace,
		alias,
	), nil
}
