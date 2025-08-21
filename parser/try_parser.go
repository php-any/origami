package parser

import (
	"fmt"
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
	tracker := p.StartTracking()

	// 跳过try关键字
	p.nextAndCheck(token.TRY)

	// 解析try块
	tryBlock := p.parseBlock()

	// 解析catch块
	var catchBlocks []node.CatchBlock
	var finallyBlock []data.GetValue

	// 解析多个catch块
	for p.checkPositionIs(0, token.CATCH) {
		catchBlock, acl := p.parseCatchBlock(tracker)
		if acl != nil {
			return nil, acl
		}
		if catchBlock != nil {
			catchBlocks = append(catchBlocks, *catchBlock)
		}
	}

	// 解析finally块
	if p.checkPositionIs(0, token.FINALLY) {
		p.next() // 跳过finally关键字
		finallyBlock = p.parseBlock()
	}

	from := tracker.EndBefore()
	return node.NewTryStatement(
		from,
		tryBlock,
		catchBlocks,
		finallyBlock,
	), nil
}

// parseCatchBlock 解析catch块
func (p *TryParser) parseCatchBlock(tracker *PositionTracker) (*node.CatchBlock, data.Control) {
	// 跳过catch关键字
	p.nextAndCheck(token.CATCH)

	// 检查是否有左括号
	if !p.checkPositionIs(0, token.LPAREN) {
		return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("Expected '(' after catch"))
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
	var variable data.Variable
	if p.checkPositionIs(0, token.VARIABLE) {
		stmt, acl := p.parseStatement()
		if acl != nil {
			p.addControl(acl)
		}
		variable1 := stmt.(*node.VariableExpression)
		variable1.Type = data.NewBaseType(exceptionType)
		variable = variable1
	} else {
		name := p.current().Literal
		p.next()
		val := p.scopeManager.CurrentScope().AddVariable(name, data.NewBaseType(exceptionType), tracker.EndBefore())
		variable = node.NewVariableWithFirst(tracker.EndBefore(), val)
	}

	p.nextAndCheck(token.RPAREN) // 跳过右括号

	// 解析catch块体
	catchBody := p.parseBlock()

	return &node.CatchBlock{
		ExceptionType: exceptionType,
		Variable:      variable,
		Body:          catchBody,
	}, nil
}
