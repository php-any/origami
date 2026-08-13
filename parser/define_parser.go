package parser

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// DefineParser 解析 define(...) 语句
type DefineParser struct {
	*Parser
}

func NewDefineParser(p *Parser) StatementParser {
	return &DefineParser{p}
}

// Parse 解析 define 语句并立即注册一个常量变量跟踪
// 语法：define(<string>, <expr>[, <case_insensitive>])
func (p *DefineParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	p.next() // 跳过 define
	if acl := p.nextAndCheck(token.LPAREN); acl != nil {
		return nil, acl
	}

	// 解析常量名：必须是字符串字面量
	if p.current().Type() != token.STRING {
		return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("define 第一个参数必须是字符串常量名"))
	}
	constName := p.current().Literal()
	p.next() // 消费常量名

	// 逗号
	if acl := p.nextAndCheck(token.COMMA); acl != nil {
		return nil, acl
	}

	// 解析值表达式
	valueExpr, acl := p.parseStatement()
	if acl != nil {
		return nil, acl
	}

	// 可选的第三个参数，忽略其值
	if p.checkPositionIs(0, token.COMMA) {
		p.next() // 逗号
		if _, acl := p.parseStatement(); acl != nil {
			return nil, acl
		}
	}

	// 右括号
	if acl := p.nextAndCheck(token.RPAREN); acl != nil {
		return nil, acl
	}

	// 立即在当前作用域创建变量跟踪（类型使用 data.Const 与 const 保持一致）
	ty := data.Const{}
	val := p.scopeManager.CurrentScope().AddVariable(constName, ty, tracker.EndBefore())

	return node.NewConstStatement(
		tracker.EndBefore(),
		node.NewVariableWithFirst(tracker.EndBefore(), val),
		valueExpr,
	), nil
}
