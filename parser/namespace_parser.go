package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// NamespaceParser 命名空间解析器
type NamespaceParser struct {
	*Parser
}

// NewNamespaceParser 创建一个命名空间解析器
func NewNamespaceParser(parser *Parser) StatementParser {
	return &NamespaceParser{
		parser,
	}
}

// Parse 解析命名空间语句
func (np *NamespaceParser) Parse() (data.GetValue, data.Control) {
	// 跳过 namespace 关键字
	tracker := np.StartTracking()
	np.next()

	// 解析命名空间名称
	var name string
	if np.current().Type == token.IDENTIFIER {
		name = np.current().Literal
		np.next()
		for np.current().Type != token.SEMICOLON && np.current().Type != token.LBRACE && !np.isEOF() {
			np.next()
			if np.current().Type != token.SEMICOLON && np.current().Type != token.LBRACE {
				name = name + "\\" + np.current().Literal
			}
		}

	} else {
		np.addError("期望命名空间名称")
		return nil, nil
	}

	// 解析命名空间体
	statements := make([]node.Statement, 0)
	from := tracker.EndBefore()
	np.namespace = node.NewNamespace(from, name, statements)

	if np.current().Type == token.LBRACE {
		np.next() // 跳过 {

		// 创建主语句解析器
		stmtParser := NewMainStatementParser(np.Parser)

		// 解析命名空间内的语句
		for np.current().Type != token.RBRACE && !np.isEOF() {
			stmt, acl := stmtParser.Parse()
			if acl != nil {
				return nil, acl
			}
			if stmt != nil {
				statements = append(statements, stmt)
			}
		}

		if np.current().Type == token.RBRACE {
			np.next() // 跳过 }
		} else {
			np.addError("期望 }")
		}
	}

	return np.namespace, nil
}
