package parser

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
	"path/filepath"
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
	var name string = np.current().Literal()
	if name == "" {
		return nil, data.NewErrorThrow(np.newFrom(), fmt.Errorf("命名空间名称不能为空"))
	}

	np.AddScanNamespace(name, filepath.Dir(*np.source))

	// 解析命名空间体
	statements := make([]data.GetValue, 0)
	from := tracker.EndBefore()
	np.namespace = node.NewNamespace(from, name, statements)

	if np.current().Type() == token.LBRACE {
		np.next() // 跳过 {

		// 创建主语句解析器
		stmtParser := NewMainStatementParser(np.Parser)

		// 解析命名空间内的语句
		for np.current().Type() != token.RBRACE && !np.isEOF() {
			stmt, acl := stmtParser.Parse()
			if acl != nil {
				return nil, acl
			}
			if stmt != nil {
				statements = append(statements, stmt)
			}
		}

		if np.current().Type() == token.RBRACE {
			np.next() // 跳过 }
		} else {
			return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("期望 }"))
		}
	}

	return np.namespace, nil
}
