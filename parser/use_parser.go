package parser

import (
	"fmt"
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
		return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("use 命名空间名称不能为空"))
	}

	// 获取完整的命名空间路径
	namespace := p.current().Literal
	p.next()

	// 检查是否有 as 关键字
	var alias string
	if p.current().Type == token.AS {
		p.next()
		if p.current().Type != token.IDENTIFIER {
			return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("as 关键字后需要变量名"))
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
		return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("use 语句缺少分号"))
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
