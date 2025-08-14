package parser

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

type SpawnParser struct {
	*Parser
}

func NewSpawnParser(parser *Parser) StatementParser {
	return &SpawnParser{
		parser,
	}
}

func (sp *SpawnParser) Parse() (data.GetValue, data.Control) {
	tracking := sp.StartTracking()

	// 跳过 spawn 关键字
	sp.next()

	// 解析 spawn 后面的表达式（可能是变量、函数调用等）
	expr, acl := sp.expressionParser.Parse()
	if acl != nil {
		return nil, acl
	}
	if expr == nil {
		return nil, data.NewErrorThrow(tracking.EndBefore(), fmt.Errorf("spawn 后面需要有效的表达式"))
	}

	switch call := expr.(type) {
	case *node.LambdaExpression:
		return node.NewSpawnStatement(tracking.EndBefore(), call), nil
	case *node.ForStatement:
		return node.NewSpawnStatement(tracking.EndBefore(), call), nil
	case *node.ForeachStatement:
		return node.NewSpawnStatement(tracking.EndBefore(), call), nil
	default:
		// 使用 VariableParser 的 ParseSuffix 方法来解析后续的方法调用
		vp := &VariableParser{sp.Parser}
		callExpr, acl := vp.parseSuffix(expr)
		if acl != nil {
			return nil, acl
		}
		// 检查是否有分号
		sp.nextAndCheckStip(token.SEMICOLON)

		return node.NewSpawnStatement(tracking.EndBefore(), callExpr), nil
	}
}
