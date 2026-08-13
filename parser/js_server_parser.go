package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// JsServerParser 表示 JS_SERVER 解析器
type JsServerParser struct {
	*Parser
}

// NewJsServerParser 创建一个新的 JS_SERVER 解析器
func NewJsServerParser(parser *Parser) StatementParser {
	return &JsServerParser{
		Parser: parser,
	}
}

// Parse 解析 JS_SERVER 表达式
// 处理 $.SERVER($variable) 格式，将其解析为函数调用表达式
// 当前 token 应该是 JS_SERVER ($.SERVER)
func (p *JsServerParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()

	// 跳过 $.SERVER (JS_SERVER)
	if p.current().Type() != token.JS_SERVER {
		return nil, data.NewErrorThrow(tracker.EndBefore(), data.NewError(tracker.EndBefore(), "$.SERVER 格式错误", nil))
	}
	p.next()

	// 检查是否有左括号
	if !p.checkPositionIs(0, token.LPAREN) {
		return nil, data.NewErrorThrow(tracker.EndBefore(), data.NewError(tracker.EndBefore(), "$.SERVER 后面必须跟左括号", nil))
	}

	// 解析函数调用参数
	vp := &VariableParser{p.Parser}
	args, acl := vp.parseFunctionCall()
	if acl != nil {
		return nil, acl
	}

	// 创建 JS_SERVER 表达式节点，传入参数列表
	expr := node.NewJsServerExpression(tracker.EndBefore(), args)

	return expr, nil
}
