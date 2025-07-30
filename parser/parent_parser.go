package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

type ParentParser struct {
	*VariableParser
}

func NewParentParser(parser *Parser) StatementParser {
	return &ParentParser{
		VariableParser: &VariableParser{
			parser,
		},
	}
}

func (pp *ParentParser) Parse() (data.GetValue, data.Control) {
	// 获取变量名
	tokenFrom := pp.NewTokenFrom(pp.GetStart())
	pp.next()

	// 检查是否是 parent:: 语法
	if pp.checkPositionIs(0, token.SCOPE_RESOLUTION) && pp.checkPositionIs(1, token.IDENTIFIER) {
		pp.next() // 跳过 ::
		methodName := pp.current().Literal
		pp.next()

		if pp.checkPositionIs(0, token.LPAREN) {
			// 创建静态方法调用表达式
			vp := &VariableParser{pp.Parser}
			expr := node.NewCallParentMethod(tokenFrom, methodName)
			return vp.parseSuffix(expr)
		} else {
			// 创建静态属性访问表达式
			vp := &VariableParser{pp.Parser}
			expr := node.NewCallParentProperty(tokenFrom, methodName)
			return vp.parseSuffix(expr)
		}
	}

	// 创建变量表达式
	expr := node.NewParent(tokenFrom)

	// 解析后续操作（函数调用、数组访问等）
	return pp.parseSuffix(expr)
}
