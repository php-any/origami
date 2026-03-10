package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/lexer"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ExpressionParser 表示表达式解析器
type ExpressionParser struct {
	*Parser
}

// NewExpressionParser 创建一个新的表达式解析器
func NewExpressionParser(parser *Parser) *ExpressionParser {
	return &ExpressionParser{
		parser,
	}
}

// Parse 解析表达式
func (ep *ExpressionParser) Parse() (data.GetValue, data.Control) {
	if ep.checkPositionIs(0, token.START_TAG, token.END_TAG, token.SEMICOLON) {
		ep.next()
		return nil, nil
	}
	return ep.parseAssignment()
}

// parseAssignment 解析赋值表达式
func (ep *ExpressionParser) parseAssignment() (data.GetValue, data.Control) {
	tracker := ep.StartTracking()
	expr, acl := ep.parseTernary()
	if acl != nil {
		return nil, acl
	}

	if v, ok := expr.(*node.VariableExpression); ok && ep.checkPositionIs(0, token.COMMA) {
		resetPosition := ep.position
		assigns := []*node.VariableExpression{v}
		valList := []data.GetValue{v}
		for ep.checkPositionIs(0, token.COMMA) {
			ep.next()
			if ep.checkPositionIs(0, token.VARIABLE) && ep.checkPositionIs(1, token.COMMA, token.ASSIGN) {
				next, acl := ep.parsePrimary()
				if acl != nil {
					return nil, acl
				}
				valList = append(valList, next)
			} else {
				next, acl := ep.parseTernary()
				if acl != nil {
					return nil, acl
				}
				valList = append(valList, next)
			}
		}
		if ep.checkPositionIs(0, token.ASSIGN, token.ADD_EQ, token.SUB_EQ, token.MUL_EQ, token.QUO_EQ, token.REM_EQ, token.CONCAT_EQ, token.NULL_COALESCE_ASSIGN) {
			// 重新构建assigns数组，避免重复
			assigns = []*node.VariableExpression{}
			for _, value := range valList {
				if next, ok := value.(*node.VariableExpression); ok {
					assigns = append(assigns, next)
				} else {
					return nil, data.NewErrorThrow(ep.FromCurrentToken(), errors.New("多赋值表达式只能是变量"))
				}
			}
			expr = node.NewVariableList(assigns)
		} else {
			ep.position = resetPosition
		}
	}

	// 检查各种赋值运算符（含字符串连接赋值 .= 和空合并赋值 ??=，以及位运算赋值）
	for ep.checkPositionIs(0, token.ASSIGN, token.ADD_EQ, token.SUB_EQ, token.MUL_EQ, token.QUO_EQ, token.REM_EQ, token.CONCAT_EQ, token.NULL_COALESCE_ASSIGN, token.BIT_OR_EQ, token.BIT_AND_EQ, token.BIT_XOR_EQ, token.SHL_EQ, token.SHR_EQ, token.POWER_EQ) {
		operator := ep.current()
		ep.next()

		right, acl := ep.parseAssignment()
		if acl != nil {
			return nil, acl
		}
		expr = node.NewBinaryExpression(
			tracker.EndBefore(),
			expr,
			operator,
			right,
		)
	}

	return expr, nil
}

// parseTernary 解析三目运算符表达式
func (ep *ExpressionParser) parseTernary() (data.GetValue, data.Control) {
	tracker := ep.StartTracking()
	expr, acl := ep.parseConcatenation()
	if acl != nil {
		return nil, acl
	}
	switch ep.current().Type() {
	case token.ELVIS:
		// 这是 ?: 简写形式（Elvis 运算符）
		// $a ?: $b 等价于 $a ? $a : $b
		ep.next() // 跳过 ?:

		// 解析假值表达式
		falseValue, acl := ep.parseTernary()
		if acl != nil {
			return nil, acl
		}
		// 创建三目运算符表达式，真值使用条件表达式本身
		return node.NewTernaryExpression(
			tracker.EndBefore(),
			expr,
			expr, // 真值就是条件表达式本身
			falseValue,
		), nil
	case token.TERNARY:
		// 检查是否是空安全调用操作符：?-> (PHP 8.0+)
		if ep.checkPositionIs(1, token.OBJECT_OPERATOR) {
			// 解析空安全调用 ?->
			ep.next() // 跳过 ?
			ep.next() // 跳过 ->

			// 使用 VariableParser 解析后续的方法/属性调用
			vp := &VariableParser{ep.Parser}
			callExpr, acl := vp.parseMethodCall(expr)
			if acl != nil {
				return nil, acl
			}
			// 继续解析链式调用（支持 ?->method()?->property）
			callExpr, acl = vp.parseSuffix(callExpr)
			if acl != nil {
				return nil, acl
			}
			// 包装为空安全调用节点
			return node.NewNullsafeCall(tracker.EndBefore(), expr, callExpr), nil
		}

		// 检查是否是可空类型声明模式：?type $variable
		if isIdentOrTypeToken(ep.peek(1).Type()) && ep.checkPositionIs(2, token.VARIABLE) {
			// 这是可空类型声明，交给专门的解析器处理
			if parser, ok := parserRouter[token.TERNARY]; ok {
				return parser(ep.Parser).Parse()
			}
		}

		// 否则按三目运算符处理
		ep.next() // 跳过 ?

		// 解析真值表达式
		trueValue, acl := ep.parseTernary()
		if acl != nil {
			return nil, acl
		}
		// 检查是否有冒号 :
		if ep.current().Type() == token.COLON {
			ep.next() // 跳过 :

			// 解析假值表达式
			falseValue, acl := ep.parseTernary()
			if acl != nil {
				return nil, acl
			}
			// 创建三目运算符表达式
			return node.NewTernaryExpression(
				tracker.EndBefore(),
				expr,
				trueValue,
				falseValue,
			), nil
		} else {
			return nil, data.NewErrorThrow(ep.FromCurrentToken(), errors.New("三目运算符 ?: 缺少冒号"))
		}
	case token.NULL_COALESCE:
		ep.next() // 跳过 ??

		// 解析右操作数
		right, acl := ep.parseTernary()
		if acl != nil {
			return nil, acl
		}
		// 创建空合并运算符表达式
		return node.NewNullCoalesceExpression(
			tracker.EndBefore(),
			expr,
			right,
		), nil
	default:
		return expr, nil
	}
}

// parseConcatenation 解析字符串连接表达式
func (ep *ExpressionParser) parseConcatenation() (data.GetValue, data.Control) {
	tracker := ep.StartTracking()
	expr, acl := ep.parseLogicalOr()
	if acl != nil {
		return nil, acl
	}
	for ep.current().Type() == token.DOT {
		operator := ep.current()
		ep.next()

		right, acl := ep.parseLogicalOr()
		if acl != nil {
			return nil, acl
		}
		expr = node.NewBinaryExpression(
			tracker.EndBefore(),
			expr,
			operator,
			right,
		)
	}

	return expr, nil
}

// parseLogicalOr 解析逻辑或表达式
func (ep *ExpressionParser) parseLogicalOr() (data.GetValue, data.Control) {
	tracker := ep.StartTracking()
	// 逻辑或的优先级低于逻辑与，因此这里从逻辑与开始
	expr, acl := ep.parseLogicalAnd()
	if acl != nil {
		return nil, acl
	}
	for ep.current().Type() == token.LOR {
		operator := ep.current()
		ep.next()

		right, acl := ep.parseLogicalAnd()
		if acl != nil {
			return nil, acl
		}
		expr = node.NewBinaryExpression(
			tracker.EndBefore(),
			expr,
			operator,
			right,
		)
	}

	return expr, nil
}

// parseLogicalAnd 解析逻辑与表达式
func (ep *ExpressionParser) parseLogicalAnd() (data.GetValue, data.Control) {
	tracker := ep.StartTracking()
	// 逻辑与的优先级低于按位或/异或/与，因此从按位或开始
	expr, acl := ep.parseBitwiseOr()
	if acl != nil {
		return nil, acl
	}
	for ep.current().Type() == token.LAND {
		operator := ep.current()
		ep.next()

		// PHP 中 && 优先级高于 ?:，故右侧用 parseBitwiseOr，不能再用 parseAssignment，
		// 否则会把 "a && b ? c : d" 解析成 a && (b ? c : d) 而非 (a && b) ? c : d
		right, acl := ep.parseBitwiseOr()
		if acl != nil {
			return nil, acl
		}
		expr = node.NewBinaryExpression(
			tracker.EndBefore(),
			expr,
			operator,
			right,
		)
	}

	return expr, nil
}

// parseEquality 解析相等性表达式
func (ep *ExpressionParser) parseEquality() (data.GetValue, data.Control) {
	tracker := ep.StartTracking()
	expr, acl := ep.parseComparison()
	if acl != nil {
		return nil, acl
	}
	for ep.checkPositionIs(0, token.EQ, token.NE, token.EQ_STRICT, token.NE_STRICT) {
		operator := ep.current()
		ep.next()

		right, acl := ep.parseComparison()
		if acl != nil {
			return nil, acl
		}
		expr = node.NewBinaryExpression(
			tracker.EndBefore(),
			expr,
			operator,
			right,
		)
	}

	// 处理 instanceof 关键字
	if ep.current().Type() == token.INSTANCEOF {
		ep.next() // 跳过 instanceof 关键字
		var right data.GetValue
		switch ep.current().Type() {
		case token.VARIABLE:
			vp := &VariableParser{ep.Parser}
			right = vp.parseVariable()
		case token.IDENTIFIER, token.PARENT, token.SELF, token.STATIC:
			// 支持 `instanceof Foo` / `instanceof parent` / `instanceof self` / `instanceof static`
			className, acl := ep.getClassName(true)
			if acl != nil {
				return nil, acl
			}
			right = node.NewStringLiteral(tracker.EndBefore(), className)
		default:
			return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("expected variable, string or identifier; str(%s)", ep.current().Literal()))
		}

		// 创建 instanceof 表达式
		expr = node.NewInstanceOfExpression(
			tracker.EndBefore(),
			expr,
			right,
		)
	}

	// 处理 like 关键字
	if ep.current().Type() == token.LIKE {
		ep.next() // 跳过 like 关键字

		className, acl := ep.getClassName(true)
		_ = acl
		// 创建 like 表达式
		expr = node.NewLikeExpression(
			tracker.EndBefore(),
			expr,
			className,
		)
	}

	return expr, nil
}

// 解析按位与/异或/或： & ^ |
// 优先级：比较 < 按位与 < 按位异或 < 按位或 < 逻辑与 < 逻辑或
func (ep *ExpressionParser) parseBitwiseAnd() (data.GetValue, data.Control) {
	tracker := ep.StartTracking()
	expr, acl := ep.parseEquality()
	if acl != nil {
		return nil, acl
	}
	for ep.current().Type() == token.BIT_AND {
		operator := ep.current()
		ep.next()

		right, acl := ep.parseEquality()
		if acl != nil {
			return nil, acl
		}
		expr = node.NewBinaryExpression(
			tracker.EndBefore(),
			expr,
			operator,
			right,
		)
	}
	return expr, nil
}

func (ep *ExpressionParser) parseBitwiseXor() (data.GetValue, data.Control) {
	tracker := ep.StartTracking()
	expr, acl := ep.parseBitwiseAnd()
	if acl != nil {
		return nil, acl
	}
	for ep.current().Type() == token.BIT_XOR {
		operator := ep.current()
		ep.next()

		right, acl := ep.parseBitwiseAnd()
		if acl != nil {
			return nil, acl
		}
		expr = node.NewBinaryExpression(
			tracker.EndBefore(),
			expr,
			operator,
			right,
		)
	}
	return expr, nil
}

func (ep *ExpressionParser) parseBitwiseOr() (data.GetValue, data.Control) {
	tracker := ep.StartTracking()
	expr, acl := ep.parseBitwiseXor()
	if acl != nil {
		return nil, acl
	}
	for ep.current().Type() == token.BIT_OR {
		operator := ep.current()
		ep.next()

		right, acl := ep.parseBitwiseXor()
		if acl != nil {
			return nil, acl
		}
		expr = node.NewBinaryExpression(
			tracker.EndBefore(),
			expr,
			operator,
			right,
		)
	}
	return expr, nil
}

// parseComparison 解析比较表达式
func (ep *ExpressionParser) parseComparison() (data.GetValue, data.Control) {
	tracker := ep.StartTracking()
	expr, acl := ep.parseShift()
	if acl != nil {
		return nil, acl
	}
	if expr == nil {
		if ep.checkPositionIs(0, token.LT) && ep.checkPositionIs(1, token.IDENTIFIER) {
			// <html
			return NewHtmlParser(ep.Parser).Parse()
		}
	}
	for ep.checkPositionIs(0, token.LT, token.LE, token.GT, token.GE) {
		operator := ep.current()
		ep.next()

		right, acl := ep.parseShift()
		if acl != nil {
			return nil, acl
		}
		expr = node.NewBinaryExpression(
			tracker.EndBefore(),
			expr,
			operator,
			right,
		)
	}

	return expr, nil
}

// parseShift 解析位移表达式 (<< >>)
func (ep *ExpressionParser) parseShift() (data.GetValue, data.Control) {
	tracker := ep.StartTracking()
	expr, acl := ep.parseTerm()
	if acl != nil {
		return nil, acl
	}
	for ep.checkPositionIs(0, token.SHL, token.SHR) {
		operator := ep.current()
		ep.next()

		right, acl := ep.parseTerm()
		if acl != nil {
			return nil, acl
		}
		expr = node.NewBinaryExpression(
			tracker.EndBefore(),
			expr,
			operator,
			right,
		)
	}

	return expr, nil
}

// parseTerm 解析加减表达式
func (ep *ExpressionParser) parseTerm() (data.GetValue, data.Control) {
	tracker := ep.StartTracking()
	expr, acl := ep.parseFactor()
	if acl != nil {
		return nil, acl
	}
	for ep.current().Type() == token.ADD || ep.current().Type() == token.SUB {
		operator := ep.current()
		ep.next()

		right, acl := ep.parseFactor()
		if acl != nil {
			return nil, acl
		}
		expr = node.NewBinaryExpression(
			tracker.EndBefore(),
			expr,
			operator,
			right,
		)
	}

	return expr, nil
}

// parseFactor 解析乘除表达式
func (ep *ExpressionParser) parseFactor() (data.GetValue, data.Control) {
	tracker := ep.StartTracking()
	expr, acl := ep.parseUnary()
	if acl != nil {
		return expr, acl
	}
	for ep.current().Type() == token.MUL || ep.current().Type() == token.QUO ||
		ep.current().Type() == token.REM {
		operator := ep.current()
		ep.next()

		right, acl := ep.parseUnary()
		if acl != nil {
			return nil, acl
		}
		expr = node.NewBinaryExpression(
			tracker.EndBefore(),
			expr,
			operator,
			right,
		)
	}

	return expr, nil
}

// parseUnary 解析一元表达式
func (ep *ExpressionParser) parseUnary() (data.GetValue, data.Control) {
	tracker := ep.StartTracking()
	if ep.current().Type() == token.SUB || ep.current().Type() == token.NOT || ep.current().Type() == token.BIT_NOT {
		operator := ep.current().Literal()
		ep.next()

		right, acl := ep.parseUnary()
		if acl != nil {
			return nil, acl
		}
		return node.NewUnaryExpression(
			tracker.EndBefore(),
			operator,
			right,
		), nil
	}

	// 处理引用取值 &$var
	if ep.current().Type() == token.BIT_AND {
		ep.next()
		right, acl := ep.parseUnary()
		if acl != nil {
			return nil, acl
		}
		return node.NewValueReference(tracker.EndBefore(), right), nil
	}

	// 处理前缀自增自减
	if ep.current().Type() == token.INCR || ep.current().Type() == token.DECR {
		operator := ep.current()
		ep.next()

		right, acl := ep.parseUnary()
		if acl != nil {
			return nil, acl
		}
		for ep.current().Type() == token.SEMICOLON {
			// 跳过没意义的分号
			ep.next()
		}
		if operator.Type() == token.INCR {
			return node.NewUnaryIncr(
				tracker.EndBefore(),
				right,
			), nil
		} else {
			return node.NewUnaryDecr(
				tracker.EndBefore(),
				right,
			), nil
		}
	}
	expr, acl := ep.parsePrimary()
	if acl == nil {
		// 检查各种赋值运算符（含字符串连接赋值 .=、位移赋值 >>= <<= 和空合并赋值 ??=，以及位运算赋值）
		for ep.checkPositionIs(0, token.ASSIGN, token.ADD_EQ, token.SUB_EQ, token.MUL_EQ, token.QUO_EQ, token.REM_EQ, token.CONCAT_EQ, token.SHL_EQ, token.SHR_EQ, token.NULL_COALESCE_ASSIGN, token.BIT_OR_EQ, token.BIT_AND_EQ, token.BIT_XOR_EQ, token.POWER_EQ) {
			operator := ep.current()
			ep.next()

			right, acl := ep.parseAssignment()
			if acl != nil {
				return nil, acl
			}
			expr = node.NewBinaryExpression(
				tracker.EndBefore(),
				expr,
				operator,
				right,
			)
		}
	}

	return expr, acl
}

// parsePrimary 解析基本表达式
func (ep *ExpressionParser) parsePrimary() (data.GetValue, data.Control) {
	switch ep.current().Type() {
	case token.DOLLAR:
		// 处理 PHP 变量变量的基础形式：$$field
		tracker := ep.StartTracking()
		ep.next() // 跳过第一个 $
		if ep.current().Type() != token.VARIABLE {
			return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("当前仅支持 $$var 形式的变量变量"))
		}
		// 解析名称变量（例如 $field）
		vp := &VariableParser{ep.Parser}
		nameExpr := vp.parseVariable()
		// 捕获当前作用域中的所有变量，作为运行时名称解析的候选集合
		vars := ep.GetVariables()
		return node.NewVarVar(tracker.EndBefore(), nameExpr, vars), nil
	case token.INT:
		value := ep.current().Literal()
		ep.next()
		return node.NewIntLiteral(ep.FromCurrentToken(), value), nil
	case token.FLOAT:
		value := ep.current().Literal()
		ep.next()
		return node.NewFloatLiteral(ep.FromCurrentToken(), value), nil
	case token.STRING:
		// 普通字符串
		value := ep.current().Literal()
		ep.next()
		return node.NewStringLiteral(ep.FromCurrentToken(), value), nil
	case token.INTERPOLATION_TOKEN:
		// 检查是否是 LingToken（插值字符串）
		if lingToken, ok := ep.current().(*lexer.LingToken); ok {
			ep.next()
			return ep.parseLingToken(lingToken), nil
		}
	case token.INTERPOLATION_VALUE:
		if lingToken, ok := ep.current().(*lexer.LingToken); ok {
			ep.next()
			return ep.parseTokensAsExpression(lingToken.Children())
		}
	case token.DOCTYPE:
		// 解析 <!DOCTYPE ...>
		tracker := ep.StartTracking()
		// 跳过 <!DOCTYPE
		ep.next()
		// 收集直到 '>' 的所有字面量，作为 DocType 内容
		var parts []string
		for !ep.isEOF() && ep.current().Type() != token.GT {
			lit := ep.current().Literal()
			if lit != "" {
				parts = append(parts, lit)
			}
			ep.next()
		}
		// 必须有 '>' 结束
		if !ep.checkPositionIs(0, token.GT) {
			return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("DOCTYPE 缺少 > 结束符"))
		}
		ep.next() // 消费 '>'
		docType := strings.TrimSpace(strings.Join(parts, " "))

		// 使用 HtmlParser 解析后续的 HTML 内容
		if hp, ok := NewHtmlParser(ep.Parser).(*HtmlParser); ok {
			// 解析后续的 HTML 子节点，直至 EOF
			children, acl := hp.parseHtmlChildren()
			if acl != nil {
				return nil, acl
			}
			return node.NewHtmlDocTypeNode(tracker.EndBefore(), docType, children), nil
		}

		return node.NewHtmlDocTypeNode(tracker.EndBefore(), docType, nil), nil
	case token.TRUE:
		ep.next()
		return node.NewBooleanLiteral(ep.FromCurrentToken(), true), nil

	case token.FALSE:
		ep.next()
		return node.NewBooleanLiteral(ep.FromCurrentToken(), false), nil

	case token.NULL:
		ep.next()
		return node.NewNullLiteral(ep.FromCurrentToken()), nil
	case token.JS_SERVER:
		// 处理 $.SERVER 表达式
		if parser, ok := parserRouter[ep.current().Type()]; ok {
			expr, acl := parser(ep.Parser).Parse()
			if acl != nil {
				return nil, acl
			}
			return expr, nil
		}
	case token.START_TAG, token.END_TAG, token.SEMICOLON:
		ep.next()
		return nil, nil
	case token.NUMBER:
		currentToken := ep.current()
		value := currentToken.Literal()
		tokenFrom := node.NewTokenFrom(ep.source, currentToken.Start(), currentToken.End(), currentToken.Line(), currentToken.Pos())
		ep.next()
		return node.NewNumberLiteral(tokenFrom, value), nil
	default:
		startType := ep.current().Type()
		if parser, ok := parserRouter[startType]; ok {
			expr, acl := parser(ep.Parser).Parse()
			if acl != nil {
				return nil, acl
			}

			// 检查是否有后缀自增自减：
			// 仅当起始 token 不是语句关键字时，才将后缀 ++ / -- 绑定到该表达式，
			// 避免诸如 "if (...) ++$i;" 被错误解析成 "if(...)++"。
			if ep.current().Type() == token.INCR || ep.current().Type() == token.DECR {
				switch startType {
				// 这些是语句级关键字，不应该在表达式里直接绑定后缀 ++ / --
				case token.IF, token.ELSE, token.FOR, token.FOREACH, token.WHILE,
					token.SWITCH, token.TRY, token.CATCH, token.FINALLY:
					// 跳过，让外层语句级解析去处理
				default:
					operator := ep.current()
					ep.next()
					for ep.current().Type() == token.SEMICOLON {
						// 跳过没意义的分号
						ep.next()
					}
					// 对于后缀自增自减，使用当前 token 的位置信息即可
					if operator.Type() == token.INCR {
						return node.NewPostfixIncr(
							ep.FromCurrentToken(),
							expr,
						), nil
					} else {
						return node.NewPostfixDecr(
							ep.FromCurrentToken(),
							expr,
						), nil
					}
				}
			}

			return expr, nil
		}
		return nil, nil
	}
	return nil, nil
}

// parseLingToken 解析 LingToken（插值字符串），创建链接节点
func (ep *ExpressionParser) parseLingToken(lingToken *lexer.LingToken) data.GetValue {
	return ep.Parser.parseLingToken(lingToken)
}
