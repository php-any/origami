package parser

import (
    "github.com/php-any/origami/data"
    "github.com/php-any/origami/node"
    "github.com/php-any/origami/token"
)

// BoolParser 表示bool类型声明解析器
type BoolParser struct {
    *Parser
}

// NewBoolParser 创建一个新的bool类型声明解析器
func NewBoolParser(parser *Parser) StatementParser {
    return &BoolParser{
        parser,
    }
}

// Parse 解析bool类型声明
func (p *BoolParser) Parse() (data.GetValue, data.Control) {
    // 开始位置跟踪
    tracker := p.StartTracking()

    // 跳过bool关键字
    p.next()

    // 检查下一个token是否是变量
    if !p.checkPositionIs(0, token.VARIABLE, token.IDENTIFIER) {
        from := tracker.End()
        return nil, data.NewErrorThrow(from, data.NewError(from, "bool类型声明需要变量名", nil))
    }

    // 获取变量名
    varName := p.current().Literal
    p.next()

    // 结束位置跟踪，获取准确的From信息
    from := tracker.EndBefore()

    // 在作用域中添加变量
    val := p.scopeManager.CurrentScope().AddVariable(varName, data.Bool{}, from)

    // 创建变量表达式
    expr := node.NewVariableWithFirst(from, val)

    // 解析后续操作（函数调用、数组访问等）
    vp := &VariableParser{p.Parser}
    return vp.parseSuffix(expr)
}
