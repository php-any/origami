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
		case token.OBJECT_OPERATOR:
			expr, acl = vp.parseMethodCall(expr)
			if acl != nil {
				return nil, acl
			}
		case token.SCOPE_RESOLUTION:
			vp.next() // 跳过 ::
			tracker := vp.StartTracking()
			if vp.checkPositionIs(1, token.CLASS) {
				// 处理 ::class 语法
				vp.next() // 跳过 class
				// 生成 ClassConstant 节点
				return node.NewClassConstant(tracker.EndBefore(), expr), nil
			} else if vp.checkPositionIs(0, token.LPAREN) {
				// 处理 ::test()
				fnName := vp.current().Literal()
				vp.next()
				// 创建函数调用表达式
				vp := &VariableParser{vp.Parser}
				expr := node.NewCallStaticMethod(tracker.EndBefore(), expr, fnName)
				return vp.parseSuffix(expr)
			} else {
				// 处理 ::test;
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
			// 优先检查命名参数
			if vp.checkPositionIs(0, token.IDENTIFIER) && vp.checkPositionIs(1, token.COLON) {
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
			} else {
				expr, acl := vp.parseStatement()
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

// parsePropertyAccess 解析属性访问
func (vp *VariableParser) parsePropertyAccess(object data.GetValue) (data.GetValue, data.Control) {
	tracker := vp.StartTracking()
	vp.next() // 跳过点号

	if vp.checkPositionIs(0, token.IDENTIFIER) || (vp.current().Type() > token.KEYWORD_START && vp.current().Type() < token.VALUE_START) {
		property := vp.current().Literal()
		vp.next()

		if vp.checkPositionIs(0, token.LPAREN) {
			stmt, acl := vp.parseFunctionCall()
			if acl != nil {
				return nil, acl
			}
			from := tracker.EndBefore()
			return node.NewObjectMethod(
				from,
				object,
				property,
				stmt,
			), nil
		} else {
			from := tracker.EndBefore()
			return node.NewObjectProperty(
				from,
				object,
				property,
			), nil
		}
	}

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

	if !(vp.checkPositionIs(0, token.IDENTIFIER) || (vp.current().Type() > token.KEYWORD_START && vp.current().Type() < token.VALUE_START)) {
		from := tracker.End()
		return nil, data.NewErrorThrow(from, errors.New("符号'->'后面需要跟随单词"))
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
