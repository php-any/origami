package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/token"
)

// MainStatementParser 主语句解析器
type MainStatementParser struct {
	*Parser
}

// NewMainStatementParser 创建一个主语句解析器
func NewMainStatementParser(parser *Parser) *MainStatementParser {
	return &MainStatementParser{
		parser,
	}
}

// Parse 解析语句
func (sp *MainStatementParser) Parse() (data.GetValue, data.Control) {
	// 获取当前词法单元类型
	tokenType := sp.current().Type()

	// 创建对应的解析器
	var parser StatementParser
	switch tokenType {
	case token.NAMESPACE:
		parser = NewNamespaceParser(sp.Parser)
	case token.USE:
		parser = NewUseParser(sp.Parser)
	case token.CLASS:
		parser = NewClassParser(sp.Parser)
	case token.FUNC:
		parser = NewFunctionParser(sp.Parser)
	case token.START_TAG, token.END_TAG, token.SEMICOLON:
		sp.next()
		return nil, nil
	default:
		// 对于其他类型的语句，使用通用解析器
		return sp.parseStatement()
	}

	return parser.Parse()
}
