package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// TryParser 表示try语句解析器
type TryParser struct {
	*Parser
}

// NewTryParser 创建一个新的try语句解析器
func NewTryParser(parser *Parser) StatementParser {
	return &TryParser{
		parser,
	}
}

// Parse 解析try语句
func (p *TryParser) Parse() (data.GetValue, data.Control) {
	start := p.GetStart()

	// 跳过try关键字
	p.nextAndCheck(token.TRY)

	// 解析try块
	tryBlock := p.parseBlock()

	// 解析catch块
	var catchBlocks []node.CatchBlock
	var finallyBlock []data.GetValue

	// 解析多个catch块
	for p.checkPositionIs(0, token.CATCH) {
		catchBlock := p.parseCatchBlock()
		if catchBlock != nil {
			catchBlocks = append(catchBlocks, *catchBlock)
		}
	}

	// 解析finally块
	if p.checkPositionIs(0, token.FINALLY) {
		p.next() // 跳过finally关键字
		finallyBlock = p.parseBlock()
	}

	return node.NewTryStatement(
		p.NewTokenFrom(start),
		tryBlock,
		catchBlocks,
		finallyBlock,
	), nil
}

// parseCatchBlock 解析catch块
func (p *TryParser) parseCatchBlock() *node.CatchBlock {
	// 跳过catch关键字
	p.nextAndCheck(token.CATCH)

	// 检查是否有左括号
	if !p.checkPositionIs(0, token.LPAREN) {
		p.addError("Expected '(' after catch")
		return nil
	}
	p.next() // 跳过左括号

	// 解析异常类型
	var exceptionType string
	if p.checkPositionIs(0, token.IDENTIFIER) {
		exceptionType = p.current().Literal
		p.next()
	} else {
		// 如果没有指定异常类型，使用默认的Exception
		exceptionType = "Exception"
	}

	// 检查是否有变量名
	var variable *node.VariableExpression
	if p.checkPositionIs(0, token.VARIABLE) {
		stmt, acl := p.parseStatement()
		if acl != nil {
			p.addControl(acl)
		}
		variable = stmt.(*node.VariableExpression)
		variable.Type = data.NewBaseType(exceptionType)
	} else {
		from := p.NewTokenFrom(p.current().Start)
		name := p.current().Literal
		p.next()
		index := p.scopeManager.CurrentScope().AddVariable(name, data.NewBaseType(exceptionType), from)
		variable = node.NewVariable(from, name, index, data.NewBaseType(exceptionType))
	}

	p.nextAndCheck(token.RPAREN) // 跳过右括号

	// 解析catch块体
	catchBody := p.parseBlock()

	return &node.CatchBlock{
		ExceptionType: exceptionType,
		Variable:      variable,
		Body:          catchBody,
	}
}
