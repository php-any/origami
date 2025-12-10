package parser

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// StaticParser 解析 static function() {} 闭包
type StaticParser struct {
	*Parser
}

func NewStaticParser(parser *Parser) StatementParser {
	return &StaticParser{parser}
}

func (sp *StaticParser) Parse() (data.GetValue, data.Control) {
	tracker := sp.StartTracking()
	// 跳过 static
	sp.next()

	// static 后必须跟 function 关键字
	if !sp.checkPositionIs(0, token.FUNC) {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("static 后必须跟 function 定义闭包"))
	}

	// 跳过 function
	sp.next()

	// 期待匿名函数格式: function (...) { ... }
	if !sp.checkPositionIs(0, token.LPAREN) {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("static function 语法错误，缺少参数列表"))
	}

	// 复用 FunctionParser 处理参数/返回类型/函数体
	fp := &FunctionParser{sp.Parser}

	// 创建新的函数作用域
	sp.scopeManager.NewScope(false)

	// 解析参数列表
	params, acl := fp.parseParameters()
	if acl != nil {
		return nil, acl
	}
	ret, acl := fp.parserReturnType()
	if acl != nil {
		return nil, acl
	}

	// 解析函数体
	body, acl := fp.parseBlock()
	if acl != nil {
		return nil, acl
	}

	vars := sp.scopeManager.CurrentScope().GetVariables()
	// 弹出函数作用域
	sp.scopeManager.PopScope()

	// 静态闭包目前语义等同普通闭包（无绑定 $this）
	fn := data.NewFuncValue(node.NewFunctionStatement(
		tracker.EndBefore(),
		"",
		params,
		body,
		vars,
		ret,
	))

	return fn, nil
}
