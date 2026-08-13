package parser

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// NullableParser 表示可空类型解析器
type NullableParser struct {
	*Parser
}

// NewNullableParser 创建一个新的可空类型解析器
func NewNullableParser(parser *Parser) StatementParser {
	return &NullableParser{
		parser,
	}
}

// Parse 解析可空类型声明
func (p *NullableParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()

	// 跳过 ? 符号
	p.next()

	// 检查下一个token是否是类型标识符
	if !isIdentOrTypeToken(p.current().Type()) {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("可空类型声明需要类型标识符"))
	}

	typeName := p.current().Literal()
	p.next()

	// 检查是否有变量名
	if p.current().Type() != token.VARIABLE {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("可空类型声明需要变量名"))
	}

	varName := p.current().Literal()
	p.next()

	// 创建可空类型
	nullableType := data.NewNullableType(data.NewBaseType(typeName))

	from := tracker.EndBefore()
	// 在作用域中添加变量
	val := p.scopeManager.CurrentScope().AddVariable(varName, nullableType, from)

	expr := node.NewVariableWithFirst(from, val)
	// 解析后续操作（函数调用、数组访问等）
	vp := &VariableParser{p.Parser}
	return vp.parseSuffix(expr)
}
