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
// 处理 $.SERVER($variable) 格式，将其转换为 JavaScript 值格式的字符串
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

	// 跳过左括号
	p.next()

	// 解析变量名
	if !p.checkPositionIs(0, token.VARIABLE) {
		return nil, data.NewErrorThrow(tracker.EndBefore(), data.NewError(tracker.EndBefore(), "$.SERVER() 中必须包含变量", nil))
	}

	varName := p.current().Literal()
	p.next()

	// 检查右括号
	if !p.checkPositionIs(0, token.RPAREN) {
		return nil, data.NewErrorThrow(tracker.EndBefore(), data.NewError(tracker.EndBefore(), "$.SERVER() 缺少右括号", nil))
	}

	// 跳过右括号
	p.next()

	// 创建 JS_SERVER 表达式节点
	// 这将被渲染为 JavaScript 值格式的字符串
	expr := node.NewJsServerExpression(tracker.EndBefore(), varName)

	return expr, nil
}
