package parser

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// GlobalParser 解析 global 关键字语句
// 支持格式：global $var1, $var2, ...;
type GlobalParser struct {
	*Parser
}

// NewGlobalParser 创建一个新的 global 语句解析器
func NewGlobalParser(parser *Parser) StatementParser {
	return &GlobalParser{parser}
}

// Parse 解析 global 语句
func (gp *GlobalParser) Parse() (data.GetValue, data.Control) {
	tracker := gp.StartTracking()
	// 跳过 global 关键字
	gp.next()

	var names []string
	var indexes []int

	for {
		if gp.current().Type() != token.VARIABLE {
			return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("global 后需要变量名"))
		}

		varName := gp.current().Literal()
		gp.next()

		// 去掉 $ 前缀
		if len(varName) > 0 && varName[0] == '$' {
			varName = varName[1:]
		}

		// 在当前作用域中添加（或查找）该变量
		v := gp.scopeManager.CurrentScope().AddVariable(varName, nil, tracker.EndBefore())

		names = append(names, varName)
		indexes = append(indexes, v.GetIndex())

		// 支持 global $a, $b, $c 多变量声明
		if gp.current().Type() == token.COMMA {
			gp.next()
			continue
		}
		break
	}

	return node.NewGlobalStatement(
		tracker.EndBefore(),
		names,
		indexes,
	), nil
}
