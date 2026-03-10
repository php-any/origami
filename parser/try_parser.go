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
	tryBlock, acl := p.parseBlock()
	if acl != nil {
		return nil, acl
	}
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
		finallyBlock, acl = p.parseBlock()
		if acl != nil {
			return nil, acl
		}
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

	// 解析异常类型（支持多类型 catch (A | B | C $e)，用 NewUnionType 表示）
	typeNames := make([]string, 0)
	for {
		if !p.checkPositionIs(0, token.IDENTIFIER) {
			if len(typeNames) == 0 {
				typeNames = append(typeNames, "Exception")
			}
			break
		}
		name := p.current().Literal()
		p.next()
		// 与参数/返回值类型解析保持一致：优先解析为完整类名（考虑命名空间与 use）
		if data.ISBaseType(name) {
			typeNames = append(typeNames, name)
		} else if full, ok := p.findFullClassNameByNamespace(name); ok {
			typeNames = append(typeNames, full)
		} else {
			typeNames = append(typeNames, name)
		}
		if !p.checkPositionIs(0, token.BIT_OR) {
			break
		}
		p.next() // 跳过 |
	}

	// 构建 data.Types：单类型用 BaseType，多类型用 NewUnionType
	var exceptionType data.Types
	types := make([]data.Types, 0, len(typeNames))
	for _, n := range typeNames {
		types = append(types, data.NewBaseType(n))
	}
	if len(types) == 1 {
		exceptionType = types[0]
	} else {
		exceptionType = data.NewUnionType(types)
	}

	// 检查是否有变量名
	var variable data.Variable
	if p.checkPositionIs(0, token.VARIABLE) {
		stmt, acl := p.parseStatement()
		if acl != nil {
			return nil, acl
		}
		variable1 := stmt.(*node.VariableExpression)
		variable1.Type = exceptionType
		variable = variable1

		acl = p.nextAndCheck(token.RPAREN) // 跳过右括号
		if acl != nil {
			return nil, acl
		}
	} else {
		name := p.current().Literal()
		p.next()
		val := p.scopeManager.CurrentScope().AddVariable(name, exceptionType, tracker.EndBefore())
		variable = node.NewVariableWithFirst(tracker.EndBefore(), val)
	}

	// 解析catch块体
	catchBody, acl := p.parseBlock()
	if acl != nil {
		return nil, acl
	}
	return &node.CatchBlock{
		ExceptionType: exceptionType,
		Variable:      variable,
		Body:          catchBody,
	}, nil
}
