package parser

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// StaticParser 解析 static 关键字
// 支持三种格式：
// 1. static function() {} - 静态闭包
// 2. static fn() => expr - 静态箭头函数
// 3. static $variable = value - 静态局部变量
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

	// 1) 支持 static::xxx() / static::$prop 这种静态调用方式
	//    与 self::/parent:: 类似，只是关键字不同
	if sp.checkPositionIs(0, token.SCOPE_RESOLUTION) &&
		(sp.checkPositionIs(1, token.IDENTIFIER) || sp.checkPositionIs(1, token.VARIABLE)) {
		// static::xxx / static::$xxx
		sp.next() // 跳过 ::
		isVariable := sp.current().Type() == token.VARIABLE
		memberName := sp.current().Literal()
		sp.next()
		tokenFrom := tracker.EndBefore()

		// 如果是 VARIABLE，去掉 $ 前缀
		if isVariable && len(memberName) > 0 && memberName[0] == '$' {
			memberName = memberName[1:]
		}

		if sp.checkPositionIs(0, token.LPAREN) {
			// 静态方法调用：static::method()
			vp := &VariableParser{sp.Parser}
			expr := node.NewCallStaticKeywordMethod(tokenFrom, memberName)
			return vp.parseSuffix(expr)
		} else {
			// 静态属性访问：static::$property
			vp := &VariableParser{sp.Parser}
			expr := node.NewCallStaticKeywordProperty(tokenFrom, memberName)
			return vp.parseSuffix(expr)
		}
	}

	// 2) 其它情况：static function / static fn / static $var，走原有逻辑
	// 检查 static 后是 function、fn 还是变量
	if sp.checkPositionIs(0, token.FUNC) {
		// static function() {} - 静态闭包
		return sp.parseStaticFunction(tracker)
	} else if sp.checkPositionIs(0, token.FN) {
		// static fn() => expr - 静态箭头函数
		return sp.parseStaticArrowFunction(tracker)
	} else if sp.checkPositionIs(0, token.VARIABLE) {
		// static $variable = value - 静态局部变量
		return sp.parseStaticVariable(tracker)
	} else {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("static 后必须跟 function、fn 或变量声明"))
	}
}

// parseStaticFunction 解析 static function() {} 格式
// 支持 static function() use ($var) {} 格式
func (sp *StaticParser) parseStaticFunction(tracker *PositionTracker) (data.GetValue, data.Control) {
	// 跳过 function
	sp.next()

	// 期待匿名函数格式: function (...) { ... } 或 function (...) use (...) { ... }
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

	// 解析 use 捕获列表（可选）：function () use ($a, $b) 或 use (&$a, $b)
	captures, acl := fp.parseClosureUse()
	if acl != nil {
		return nil, acl
	}

	// 解析返回类型
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

	// 按引用捕获: 对 use(&$var) 的变量，将其替换为 VariableReference
	if len(captures) > 0 {
		for _, c := range captures {
			if !c.IsReference {
				continue
			}
			if childVar, ok := sp.scopeManager.CurrentScope().GetVariable(c.Name); ok {
				sp.scopeManager.CurrentScope().SetVariable(
					c.Name,
					node.NewVariableReference(
						sp.FromCurrentToken(),
						childVar.GetName(),
						childVar.GetIndex(),
						childVar.GetType(),
					),
				)
			}
		}
		// 更新 vars，确保其中的引用变量已经变成 VariableReference
		vars = sp.scopeManager.CurrentScope().GetVariables()
	}

	// 弹出函数作用域
	sp.scopeManager.PopScope()

	// 如果有 use 子句，使用 LambdaExpression 并构建 parent 映射
	if len(captures) > 0 {
		// 构建 parent 映射，仅捕获 use 声明的变量
		parent := make(map[int]int)
		for _, outer := range sp.scopeManager.CurrentScope().GetVariables() {
			for _, child := range vars {
				if child.GetName() == outer.GetName() {
					for _, c := range captures {
						if c.Name == child.GetName() {
							parent[child.GetIndex()] = outer.GetIndex()
						}
					}
				}
			}
		}

		// 静态闭包使用 LambdaExpression（支持 use 子句）
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

	// 没有 use 子句，使用 FunctionStatement
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

// parseStaticArrowFunction 解析 static fn(...) => ... 格式
func (sp *StaticParser) parseStaticArrowFunction(tracker *PositionTracker) (data.GetValue, data.Control) {
	// 跳过 fn
	sp.next()

	// 期待箭头函数格式: fn (...) => ...
	if !sp.checkPositionIs(0, token.LPAREN) {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("static fn 语法错误，缺少参数列表"))
	}

	// 复用 FunctionParser 处理参数/返回类型
	fp := &FunctionParser{sp.Parser}

	// 创建新的函数作用域（箭头函数是 lambda，自动捕获外部变量）
	sp.scopeManager.NewScope(true)

	// 解析参数列表
	params, acl := fp.parseParameters()
	if acl != nil {
		return nil, acl
	}

	// 解析返回类型（可选）
	ret, acl := fp.parserReturnType()
	if acl != nil {
		return nil, acl
	}

	// 检查是否有 => 符号
	if !sp.checkPositionIs(0, token.ARRAY_KEY_VALUE) {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("static fn 语法错误，缺少 => 符号"))
	}
	sp.next() // 跳过 =>

	// 解析函数体（箭头函数是单个表达式）
	body, acl := fp.parseBlock()
	if acl != nil {
		return nil, acl
	}

	vars := sp.scopeManager.CurrentScope().GetVariables()
	// 弹出函数作用域
	sp.scopeManager.PopScope()

	// 构建 parent 映射，自动捕获外部变量
	parent := make(map[int]int)
	for _, parentVariable := range sp.scopeManager.CurrentScope().GetVariables() {
		for _, childVariable := range vars {
			if childVariable.GetName() == parentVariable.GetName() {
				parent[childVariable.GetIndex()] = parentVariable.GetIndex()
			}
		}
	}

	// 静态箭头函数创建为 Lambda 表达式（无绑定 $this）
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

// parseStaticVariable 解析 static $variable = value 格式
func (sp *StaticParser) parseStaticVariable(tracker *PositionTracker) (data.GetValue, data.Control) {
	// 解析类型声明（可选，如 static int $count = 0）
	var varType data.Types
	if isIdentOrTypeToken(sp.current().Type()) {
		typeName := sp.current().Literal()
		sp.next()
		varType = data.NewBaseType(typeName)
	} else if sp.checkPositionIs(0, token.TERNARY) && isIdentOrTypeToken(sp.peek(1).Type()) {
		// ?int 方式
		sp.next()
		base := data.NewBaseType(sp.current().Literal())
		sp.next()
		varType = data.NewNullableType(base)
	}

	// 解析变量名
	if sp.current().Type() != token.VARIABLE {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("static 后需要变量名"))
	}

	name := sp.current().Literal()
	sp.next()

	// 解析初始化表达式
	var initializer data.GetValue
	if sp.current().Type() == token.ASSIGN {
		sp.next() // 跳过等号
		exprParser := NewExpressionParser(sp.Parser)
		var acl data.Control
		initializer, acl = exprParser.Parse()
		if acl != nil {
			return nil, acl
		}
	}

	// 在作用域中注册变量（静态局部变量）
	// 去掉变量名的 $ 前缀
	varName := name
	if len(varName) > 0 && varName[0] == '$' {
		varName = varName[1:]
	}
	sp.scopeManager.CurrentScope().AddVariable(varName, varType, tracker.EndBefore())

	// 创建静态局部变量声明语句
	return node.NewStaticVarStatement(
		tracker.EndBefore(),
		name,
		initializer,
	), nil
}
