package parser

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// FnParser 解析 fn 关键字（箭头函数）
//
// 使用 fn 关键字创建箭头函数（短闭包）
//
// 语法: fn(参数列表): 返回类型 => 表达式
//
// 特性:
//   - 箭头函数自动捕获外部作用域的变量（按值捕获）
//   - 箭头函数是 lambda 表达式，不能绑定 $this
//   - 函数体是单个表达式，不需要花括号
//
// 使用示例:
//
//	// 基本用法
//	$fn = fn($x) => $x * 2;
//
//	// 带类型声明
//	$fn = fn(int $x): int => $x * 2;
//
//	// 多个参数
//	$fn = fn($a, $b) => $a + $b;
//
//	// 捕获外部变量
//	$multiplier = 10;
//	$fn = fn($x) => $x * $multiplier;
//
//	// 在数组方法中使用
//	$numbers = [1, 2, 3];
//	$doubled = $numbers->map(fn($n) => $n * 2);
type FnParser struct {
	*Parser
}

func NewFnParser(parser *Parser) StatementParser {
	return &FnParser{parser}
}

func (fp *FnParser) Parse() (data.GetValue, data.Control) {
	tracker := fp.StartTracking()
	// 跳过 fn
	fp.next()

	// 期待箭头函数格式: fn (...) => ...
	if !fp.checkPositionIs(0, token.LPAREN) {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("fn 语法错误，缺少参数列表"))
	}

	// 复用 FunctionParser 处理参数/返回类型
	fpHelper := &FunctionParser{fp.Parser}

	// 创建新的函数作用域（箭头函数是 lambda，自动捕获外部变量）
	fp.scopeManager.NewScope(true)

	// 解析参数列表
	params, acl := fpHelper.parseParameters()
	if acl != nil {
		return nil, acl
	}

	// 解析返回类型（可选）
	ret, acl := fpHelper.parserReturnType()
	if acl != nil {
		return nil, acl
	}

	// 检查是否有 => 符号
	if !fp.checkPositionIs(0, token.ARRAY_KEY_VALUE) {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("fn 语法错误，缺少 => 符号"))
	}
	fp.next() // 跳过 =>

	// 解析函数体（箭头函数是单个表达式）
	body, acl := fpHelper.parseBlock()
	if acl != nil {
		return nil, acl
	}

	vars := fp.scopeManager.CurrentScope().GetVariables()
	// 弹出函数作用域
	fp.scopeManager.PopScope()

	// 构建 parent 映射，自动捕获外部变量
	parent := make(map[int]int)
	for _, childVariable := range vars {
		for _, parentVariable := range fp.scopeManager.CurrentScope().GetVariables() {
			if childVariable.GetName() == parentVariable.GetName() {
				// 形参由调用方传入，不应从父作用域捕获，否则体内会误读外层同名变量
				if isParameterName(params, childVariable.GetName()) {
					continue
				}
				parent[childVariable.GetIndex()] = parentVariable.GetIndex()
			}
		}
	}

	// 箭头函数创建为 Lambda 表达式（无绑定 $this）
	fn := node.NewLambdaExpression(
		tracker.EndBefore(),
		params,
		body,
		vars,
		parent,
	)

	// 设置返回类型（如果指定了）
	if ret != nil {
		fn.FunctionStatement.Ret = ret
	}

	return fn, nil
}

// isParameterName 判断 name 是否为 params 中某个形参的名称（形参不应从父作用域捕获）
func isParameterName(params []data.GetValue, name string) bool {
	for _, p := range params {
		if v, ok := p.(interface{ GetName() string }); ok && v.GetName() == name {
			return true
		}
	}
	return false
}
