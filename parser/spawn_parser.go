package parser

import (
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
	from := sp.NewTokenFrom(sp.GetStart())

	// 跳过 spawn 关键字
	sp.next()

	// 解析 spawn 后面的表达式（可能是变量、函数调用等）
	expr, acl := sp.expressionParser.Parse()
	if acl != nil {
		return nil, acl
	}
	if expr == nil {
		sp.addError("spawn 后面需要有效的表达式")
		return nil, nil
	}

	switch call := expr.(type) {
	case *node.LambdaExpression:
		return node.NewSpawnStatement(from, call), nil
	case *node.ForStatement:
		return node.NewSpawnStatement(from, call), nil
	case *node.ForeachStatement:
		return node.NewSpawnStatement(from, call), nil
	default:
		// 使用 VariableParser 的 ParseSuffix 方法来解析后续的方法调用
		vp := &VariableParser{sp.Parser}
		callExpr, acl := vp.parseSuffix(expr)
		if acl != nil {
			return nil, acl
		}
		// 检查是否有分号
		sp.nextAndCheckStip(token.SEMICOLON)

		return node.NewSpawnStatement(from, callExpr), nil
	}
}
