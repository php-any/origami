package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type ThisParser struct {
	*VariableParser
}

func NewThisParser(parser *Parser) StatementParser {
	return &ThisParser{
		VariableParser: &VariableParser{
			parser,
		},
	}
}

func (vp *ThisParser) Parse() (data.GetValue, data.Control) {
	// 获取变量名
	tokenFrom := vp.FromCurrentToken()
	vp.next()

	// 创建变量表达式
	expr := node.NewThis(tokenFrom)

	// 解析后续操作（函数调用、数组访问等）
	return vp.parseSuffix(expr)
}
