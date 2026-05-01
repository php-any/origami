package parser

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// VariableParser 表示变量解析器
type VariableParser struct {
	*Parser
}

// NewVariableParser 创建一个新的变量解析器
func NewVariableParser(parser *Parser) StatementParser {
	return &VariableParser{
		parser,
	}
}

// Parse 解析变量表达式
func (vp *VariableParser) Parse() (data.GetValue, data.Control) {
	// 获取变量名
	expr := vp.parseVariable()

	// 解析后续操作（函数调用、数组访问等）
	return vp.parseSuffix(expr)
}

func (vp *VariableParser) parseVariable() data.Variable {
	tracker := vp.StartTracking()

	// 获取变量名
	name := vp.current().Literal()
	vp.next()

	// 特殊变量：CLI 超全局 $argv / $argc 以及其它 PHP 超全局
	switch name {
	case "$argv":
		// 首先判断当前是否处于全局作用域（无父作用域）
		isGlobalScope := vp.scopeManager.CurrentScope() != nil && vp.scopeManager.CurrentScope().GetParent() == nil
		if !isGlobalScope {
			// 查找变量索引
			varInfo := vp.scopeManager.LookupVariable(name)
			if varInfo == nil {
				// 如果变量不存在，在当前作用域中创建它
				val := vp.scopeManager.CurrentScope().AddVariable(name, nil, tracker.EndBefore())
				varInfo = node.NewVariableWithFirst(tracker.EndBefore(), val)
			}

			// 创建变量表达式
			return node.NewVariableWithFirst(tracker.EndBefore(), varInfo)
		}
		// 在运行时由 node/argv_variable.go 中的 ArgvVariable 完成取值
		return node.NewArgvVariable(tracker.EndBefore())
	case "$argc":
		isGlobalScope := vp.scopeManager.CurrentScope() != nil && vp.scopeManager.CurrentScope().GetParent() == nil
		if !isGlobalScope {
			// 查找变量索引
			varInfo := vp.scopeManager.LookupVariable(name)
			if varInfo == nil {
				// 如果变量不存在，在当前作用域中创建它
				val := vp.scopeManager.CurrentScope().AddVariable(name, nil, tracker.EndBefore())
				varInfo = node.NewVariableWithFirst(tracker.EndBefore(), val)
			}

			// 创建变量表达式
			return node.NewVariableWithFirst(tracker.EndBefore(), varInfo)
		}
		// 在运行时由 node/argv_variable.go 中的 ArgcVariable 完成取值
		return node.NewArgcVariable(tracker.EndBefore())
	case "$GLOBALS":
		return node.NewGlobalsArrayVariable(tracker.EndBefore())
	case "$_ENV":
		return node.NewEnvVariable(tracker.EndBefore())
	case "$_SERVER":
		return node.NewServerVariable(tracker.EndBefore())
	case "$_GET":
		return node.NewGetVariable(tracker.EndBefore())
	case "$_POST":
		return node.NewPostVariable(tracker.EndBefore())
	case "$_COOKIE":
		return node.NewCookieVariable(tracker.EndBefore())
	case "$_SESSION":
		return node.NewSessionVariable(tracker.EndBefore())
	case "$_FILES":
		return node.NewFilesVariable(tracker.EndBefore())
	case "$_REQUEST":
		return node.NewRequestVariable(tracker.EndBefore())

	}

	// 查找变量索引
	varInfo := vp.scopeManager.LookupVariable(name)
	if varInfo == nil {
		// 如果变量不存在，在当前作用域中创建它
		val := vp.scopeManager.CurrentScope().AddVariable(name, nil, tracker.EndBefore())
		varInfo = node.NewVariableWithFirst(tracker.EndBefore(), val)
	}

	// 创建变量表达式
	return node.NewVariableWithFirst(tracker.EndBefore(), varInfo)
}

// parseSuffix 解析变量后缀操作
func (vp *VariableParser) parseSuffix(expr data.GetValue) (data.GetValue, data.Control) {
	var acl data.Control
	for {
		switch vp.current().Type() {
		case token.LPAREN:
			// 在解析函数调用之前记录位置
			tracker := vp.StartTracking()
			stmt, acl := vp.parseFunctionCall()
			if acl != nil {
				return nil, acl
			}
			from := tracker.EndBefore()
			expr = node.NewCallMethod(from, expr, stmt)
		case token.LBRACKET:
			expr, acl = vp.parseArrayAccess(expr)
			if acl != nil {
				return nil, acl
			}
		case token.DOT:
			expr, acl = vp.parsePropertyAccess(expr)
			if acl != nil {
				return nil, acl
			}
		case token.TERNARY:
			// 检查是否是链式空安全调用：?-> (PHP 8.0+)
			if vp.checkPositionIs(1, token.OBJECT_OPERATOR) {
				// 解析链式空安全调用 ?->
				vp.next() // 跳过 ?
				vp.next() // 跳过 ->
				// 手动解析方法/属性调用（因为我们已经跳过了 ->）
				tracker2 := vp.StartTracking()
				var callExpr data.GetValue
				// 先处理花括号动态属性：$obj->{$name} 或 $obj->{expr}
				if vp.current().Type() == token.LBRACE {
					vp.next() // 跳过 {
					// 支持任意表达式作为属性名：$obj?->{expr}
					nameExpr, acl := vp.expressionParser.Parse()
					if acl != nil {
						return nil, acl
					}
					vp.nextAndCheck(token.RBRACE)
					callExpr = node.NewCallObjectDynamicProperty(tracker2.EndBefore(), expr, nameExpr)
				} else if !(vp.checkPositionIs(0, token.IDENTIFIER, token.VARIABLE) || (vp.current().Type() > token.KEYWORD_START && vp.current().Type() < token.VALUE_START)) {
					return nil, data.NewErrorThrow(tracker2.EndBefore(), errors.New("符号'?->'后面需要跟随单词"))
				} else if vp.current().Type() == token.VARIABLE {
					nameExpr := vp.parseVariable()
					callExpr = node.NewCallObjectDynamicProperty(tracker2.EndBefore(), expr, nameExpr)
				} else {
					method := vp.current().Literal()
					vp.next()
					if vp.current().Type() == token.LPAREN {
						stmt, acl := vp.parseFunctionCall()
						if acl != nil {
							return nil, acl
						}
						callExpr = node.NewObjectMethod(tracker2.EndBefore(), expr, method, stmt)
					} else {
						callExpr = node.NewObjectProperty(tracker2.EndBefore(), expr, method)
					}
				}
				// 继续解析链式调用（支持 ?->method()?->property）
				callExpr, acl := vp.parseSuffix(callExpr)
				if acl != nil {
					return nil, acl
				}
				// 包装为空安全调用节点
				expr = node.NewNullsafeCall(vp.FromCurrentToken(), expr, callExpr)
				continue
			}
			// 不是 ?->，返回当前表达式
			return expr, nil
		case token.OBJECT_OPERATOR:
			expr, acl = vp.parseMethodCall(expr)
			if acl != nil {
				return nil, acl
			}
		case token.SCOPE_RESOLUTION:
			vp.next() // 跳过 ::
			tracker := vp.StartTracking()
			if vp.checkPositionIs(0, token.CLASS) {
				// 处理 $var::class 语法
				vp.next() // 跳过 class
				// 生成 ClassConstant 节点
				return node.NewClassConstant(tracker.EndBefore(), expr), nil
			} else if vp.checkPositionIs(1, token.LPAREN) {
				// 处理 $var::method() 静态方法调用
				fnName := vp.current().Literal()
				vp.next() // 跳过方法名
				// 创建静态方法调用表达式（运行时会动态解析类名）
				vp := &VariableParser{vp.Parser}
				expr := node.NewCallStaticMethod(tracker.EndBefore(), expr, fnName)
				return vp.parseSuffix(expr)
			} else {
				// 处理 $var::PROPERTY 静态属性访问
				attrName := vp.current().Literal()
				vp.next()
				vp := &VariableParser{vp.Parser}
				expr := node.NewCallStaticProperty(tracker.EndBefore(), expr, attrName)
				return vp.parseSuffix(expr)
			}
		default:
			return expr, nil
		}
	}
}

// parseFunctionCall 解析函数调用
func (vp *VariableParser) parseFunctionCall() ([]data.GetValue, data.Control) {
	vp.nextAndCheck(token.LPAREN) // 跳过左括号

	args := make([]data.GetValue, 0)
	if vp.current().Type() != token.RPAREN {
		for {
			// 优先检查命名参数（标识符或关键字后跟冒号，PHP 8.0+）
			if (vp.checkPositionIs(0, token.IDENTIFIER, token.DEFAULT) ||
				(vp.current().Type() > token.KEYWORD_START && vp.current().Type() < token.VALUE_START && vp.current().Type() != token.RPAREN)) &&
				vp.checkPositionIs(1, token.COLON) {
				tracker := vp.StartTracking()
				name := vp.current().Literal()
				vp.next()
				vp.next()
				from := tracker.EndBefore()

				value, acl := vp.expressionParser.Parse()
				if acl != nil {
					return nil, acl
				}
				value, acl = vp.parseSuffix(value)
				if acl != nil {
					return nil, acl
				}
				args = append(args, node.NewNamedArgument(from, name, value))
				if vp.current().Type() != token.COMMA {
					break
				}
				vp.next()
			} else if vp.current().Type() == token.ELLIPSIS {
				tracker := vp.StartTracking()
				vp.next() // 跳过 ...
				// 如果 ... 后紧跟 ) 或 , 则是 first-class callable（如 $fn = $obj->method(...)）
				if vp.current().Type() == token.RPAREN || vp.current().Type() == token.COMMA {
					// 使用 ToClosure 标记，让外层 CallObjectMethod 处理
					from := tracker.EndBefore()
					args = append(args, node.NewSpreadArgument(from, nil))
					if vp.current().Type() != token.COMMA {
						break
					}
					vp.next()
					continue
				}
				// 展开实参 ...expr
				expr, acl := vp.parseStatement()
				if acl != nil {
					return nil, acl
				}
				expr, acl = vp.parseSuffix(expr)
				if acl != nil {
					return nil, acl
				}
				from := tracker.EndBefore()
				args = append(args, node.NewSpreadArgument(from, expr))

				if vp.current().Type() != token.COMMA {
					break
				}
				vp.next()
			} else {
				// 解析单个参数表达式，但不消费逗号（避免将多个参数解析为一个 VariableList）
				// 使用 parseTernary 而不是 parseStatement/Parse，因为后者会消费逗号
				expr, acl := vp.expressionParser.parseTernary()
				if acl != nil {
					return nil, acl
				}
				if list, ok := expr.(*node.VariableList); ok {
					for _, expr = range list.Vars {
						args = append(args, expr)
					}
				} else {
					expr, acl = vp.parseSuffix(expr)
					if acl != nil {
						return nil, acl
					}
					args = append(args, expr)
				}

				if vp.current().Type() != token.COMMA {
					break
				}
				vp.next()
			}
		}
	}

	return args, vp.nextAndCheck(token.RPAREN)
}

// parseArrayAccess 解析数组访问
func (vp *VariableParser) parseArrayAccess(array data.GetValue) (data.GetValue, data.Control) {
	tracker := vp.StartTracking()
	vp.next() // 跳过左方括号
	from := tracker.EndBefore()

	if vp.current().Type() == token.DOUBLE_DOT {
		// arr[..1]
		vp.next()
		if vp.current().Type() == token.RBRACKET {
			// arr[..]
			vp.next()
			return node.NewRange(
				from,
				array,
				nil,
				nil,
			), nil
		}
		stop, acl := vp.expressionParser.Parse()
		if acl != nil {
			return nil, acl
		}
		vp.nextAndCheck(token.RBRACKET)
		return node.NewRange(
			from,
			array,
			nil,
			stop,
		), nil
	}

	var index data.GetValue
	var acl data.Control

	if vp.current().Type() == token.RBRACKET {
		// arr[]
		vp.next()
		index = node.NewObjectProperty(from, array, "length")
		return node.NewIndexExpression(
			from,
			array,
			index,
		), nil
	} else {
		index, acl = vp.expressionParser.Parse()
		if acl != nil {
			return nil, acl
		}
	}

	if vp.current().Type() == token.RBRACKET {
		// arr[1]
		vp.next()

		return node.NewIndexExpression(
			from,
			array,
			index,
		), nil
	}

	if vp.current().Type() == token.DOUBLE_DOT {
		start := index
		vp.next()
		if vp.current().Type() == token.RBRACKET {
			// arr[1..]
			vp.next()
			return node.NewRange(
				from,
				array,
				start,
				nil,
			), nil
		}
		stop, acl := vp.expressionParser.Parse()
		if acl != nil {
			return nil, acl
		}
		vp.nextAndCheck(token.RBRACKET)

		return node.NewRange(
			from,
			array,
			start,
			stop,
		), nil
	}

	vp.nextAndCheck(token.RBRACKET)

	return node.NewIndexExpression(
		from,
		array,
		index,
	), nil
}

// parsePropertyAccess 解析属性访问或者拼接
func (vp *VariableParser) parsePropertyAccess(object data.GetValue) (data.GetValue, data.Control) {
	tracker := vp.StartTracking()
	vp.next() // 跳过点号

	//if vp.checkPositionIs(0, token.LPAREN) {
	//	// data . () 只能是拼接
	//	parser := NewLparenParser(vp.Parser)
	//	property, acl := parser.Parse()
	//	if acl != nil {
	//		return nil, acl
	//	}
	//	from := tracker.EndBefore()
	//	return node.NewBinaryAdd(from, object, property), nil
	//} else if vp.checkPositionIs(0, token.IDENTIFIER) || (vp.current().Type() > token.KEYWORD_START && vp.current().Type() < token.VALUE_START) {
	//	property := vp.current().Literal()
	//	vp.next()
	//
	//	if vp.checkPositionIs(0, token.LPAREN) {
	//		stmt, acl := vp.parseFunctionCall()
	//		if acl != nil {
	//			return nil, acl
	//		}
	//		from := tracker.EndBefore()
	//		return node.NewObjectMethod(
	//			from,
	//			object,
	//			property,
	//			stmt,
	//		), nil
	//	} else {
	//		from := tracker.EndBefore()
	//		return node.NewObjectProperty(
	//			from,
	//			object,
	//			property,
	//		), nil
	//	}
	//}

	// 尝试兼容 php . 符号作为字符串链接
	property, acl := vp.parseStatement()
	from := tracker.EndBefore()
	return node.NewBinaryAdd(
		from,
		object,
		property,
	), acl
}

func (vp *VariableParser) parseMethodCall(object data.GetValue) (data.GetValue, data.Control) {
	vp.next() // 跳过箭头
	tracker := vp.StartTracking()

	// 先处理花括号动态属性：$obj->{$name}
	if vp.current().Type() == token.LBRACE {
		// 语法：$obj->{$name} 或更通用的 $obj->{expr}
		vp.next() // 跳过 {

		// 允许任意表达式作为动态属性名：$obj->{expr}
		// 例如 $obj->{$name}、$obj->{getName()}、$obj->{$a.$b} 等
		nameExpr, acl := vp.expressionParser.Parse()
		if acl != nil {
			return nil, acl
		}

		// 期望右花括号 }
		vp.nextAndCheck(token.RBRACE)

		from := tracker.EndBefore()
		// 如果后面紧跟 (，则是动态方法调用 $obj->{expr}(...)
		if vp.current().Type() == token.LPAREN {
			stmt, acl := vp.parseFunctionCall()
			if acl != nil {
				return nil, acl
			}
			from = tracker.EndBefore()
			return node.NewCallObjectDynamicMethod(from, object, nameExpr, stmt), nil
		}
		return node.NewCallObjectDynamicProperty(from, object, nameExpr), nil
	}

	// 支持三种情况：
	// 1. 传统形式：$obj->prop / $obj->method()
	// 2. 关键字作为方法/属性名：$obj->class()
	// 3. 动态属性：$obj->$name / $obj->$name()
	if !(vp.checkPositionIs(0, token.IDENTIFIER, token.VARIABLE) || (vp.current().Type() > token.KEYWORD_START && vp.current().Type() < token.VALUE_START)) {
		from := tracker.End()
		return nil, data.NewErrorThrow(from, errors.New("符号'->'后面需要跟随单词"))
	}

	// 动态属性：$obj->$name / $obj->$name()
	// 这里直接将其解析为等价的索引访问：$obj[$name]
	// 这样可以复用 IndexExpression 在运行时对对象/数组的动态属性逻辑
	if vp.current().Type() == token.VARIABLE {
		// 复用变量解析逻辑，确保变量索引、类型信息等保持一致
		nameExpr := vp.parseVariable()
		from := tracker.EndBefore()
		// 如果后面紧跟 (，则是动态方法调用 $obj->$name(...)
		if vp.current().Type() == token.LPAREN {
			stmt, acl := vp.parseFunctionCall()
			if acl != nil {
				return nil, acl
			}
			from = tracker.EndBefore()
			return node.NewCallObjectDynamicMethod(from, object, nameExpr, stmt), nil
		}
		return node.NewCallObjectDynamicProperty(from, object, nameExpr), nil
	}

	method := vp.current().Literal()
	vp.next()

	// 如果后面跟着括号，解析方法调用
	if vp.current().Type() == token.LPAREN {
		stmt, acl := vp.parseFunctionCall()
		if acl != nil {
			return nil, acl
		}
		// 在解析完整个方法调用后设置范围
		from := tracker.EndBefore()
		return node.NewObjectMethod(
			from,
			object,
			method,
			stmt,
		), nil
	}

	// 对于属性访问，在方法名之后设置范围
	from := tracker.EndBefore()
	return node.NewObjectProperty(
		from,
		object,
		method,
	), nil
}
