package parser

import (
	"errors"

	"github.com/php-any/origami/data"
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
			next, acl := ep.parseTernary()
			if acl != nil {
				return nil, acl
			}
			valList = append(valList, next)
		}
		if ep.checkPositionIs(0, token.ASSIGN, token.ADD_EQ, token.SUB_EQ, token.MUL_EQ, token.QUO_EQ, token.REM_EQ) {
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

	// 检查各种赋值运算符
	for ep.checkPositionIs(0, token.ASSIGN, token.ADD_EQ, token.SUB_EQ, token.MUL_EQ, token.QUO_EQ, token.REM_EQ) {
		operator := ep.current()
		ep.next()

		right, acl := ep.parseAssignment()
		if acl != nil {
			return nil, acl
		}
		expr = node.NewBinaryExpression(
			ep.FromCurrentToken(),
			expr,
			operator,
			right,
		)
	}

	return expr, nil
}

// parseTernary 解析三目运算符表达式
func (ep *ExpressionParser) parseTernary() (data.GetValue, data.Control) {
	expr, acl := ep.parseConcatenation()
	if acl != nil {
		return nil, acl
	}
	switch ep.current().Type {
	case token.TERNARY:
		// 检查是否是可空类型声明模式：?type $variable
		if isIdentOrTypeToken(ep.peek(1).Type) && ep.checkPositionIs(2, token.VARIABLE) {
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
		if ep.current().Type == token.COLON {
			ep.next() // 跳过 :

			// 解析假值表达式
			falseValue, acl := ep.parseTernary()
			if acl != nil {
				return nil, acl
			}
			// 创建三目运算符表达式
			return node.NewTernaryExpression(
				ep.FromCurrentToken(),
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
			ep.FromCurrentToken(),
			expr,
			right,
		), nil
	default:
		return expr, nil
	}
}

// parseConcatenation 解析字符串连接表达式
func (ep *ExpressionParser) parseConcatenation() (data.GetValue, data.Control) {
	expr, acl := ep.parseLogicalOr()
	if acl != nil {
		return nil, acl
	}
	for ep.current().Type == token.DOT {
		operator := ep.current()
		ep.next()

		right, acl := ep.parseLogicalOr()
		if acl != nil {
			return nil, acl
		}
		expr = node.NewBinaryExpression(
			ep.FromCurrentToken(),
			expr,
			operator,
			right,
		)
	}

	return expr, nil
}

// parseLogicalOr 解析逻辑或表达式
func (ep *ExpressionParser) parseLogicalOr() (data.GetValue, data.Control) {
	expr, acl := ep.parseLogicalAnd()
	if acl != nil {
		return nil, acl
	}
	for ep.current().Type == token.LOR {
		operator := ep.current()
		ep.next()

		right, acl := ep.parseLogicalAnd()
		if acl != nil {
			return nil, acl
		}
		expr = node.NewBinaryExpression(
			ep.FromCurrentToken(),
			expr,
			operator,
			right,
		)
	}

	return expr, nil
}

// parseLogicalAnd 解析逻辑与表达式
func (ep *ExpressionParser) parseLogicalAnd() (data.GetValue, data.Control) {
	expr, acl := ep.parseEquality()
	if acl != nil {
		return nil, acl
	}
	for ep.current().Type == token.LAND {
		operator := ep.current()
		ep.next()

		right, acl := ep.parseEquality()
		if acl != nil {
			return nil, acl
		}
		expr = node.NewBinaryExpression(
			ep.FromCurrentToken(),
			expr,
			operator,
			right,
		)
	}

	return expr, nil
}

// parseEquality 解析相等性表达式
func (ep *ExpressionParser) parseEquality() (data.GetValue, data.Control) {
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
			ep.FromCurrentToken(),
			expr,
			operator,
			right,
		)
	}

	// 处理 instanceof 关键字
	if ep.current().Type == token.INSTANCEOF {
		ep.next() // 跳过 instanceof 关键字

		className, acl := ep.getClassName(true)
		_ = acl
		// 创建 instanceof 表达式
		expr = node.NewInstanceOfExpression(
			ep.FromCurrentToken(),
			expr,
			className,
		)
	}

	// 处理 like 关键字
	if ep.current().Type == token.LIKE {
		ep.next() // 跳过 like 关键字

		className, acl := ep.getClassName(true)
		_ = acl
		// 创建 like 表达式
		expr = node.NewLikeExpression(
			ep.FromCurrentToken(),
			expr,
			className,
		)
	}

	return expr, nil
}

// parseComparison 解析比较表达式
func (ep *ExpressionParser) parseComparison() (data.GetValue, data.Control) {
	expr, acl := ep.parseTerm()
	if acl != nil {
		return nil, acl
	}
	if expr == nil && ep.checkPositionIs(0, token.LT) && ep.checkPositionIs(1, token.IDENTIFIER) {
		// <html
		return NewHtmlParser(ep.Parser).Parse()
	}
	for ep.current().Type == token.LT || ep.current().Type == token.LE ||
		ep.current().Type == token.GT || ep.current().Type == token.GE {
		operator := ep.current()
		ep.next()

		right, acl := ep.parseTerm()
		if acl != nil {
			return nil, acl
		}
		expr = node.NewBinaryExpression(
			ep.FromCurrentToken(),
			expr,
			operator,
			right,
		)
	}

	return expr, nil
}

// parseTerm 解析加减表达式
func (ep *ExpressionParser) parseTerm() (data.GetValue, data.Control) {
	expr, acl := ep.parseFactor()
	if acl != nil {
		return nil, acl
	}
	for ep.current().Type == token.ADD || ep.current().Type == token.SUB {
		operator := ep.current()
		ep.next()

		right, acl := ep.parseFactor()
		if acl != nil {
			return nil, acl
		}
		expr = node.NewBinaryExpression(
			ep.FromCurrentToken(),
			expr,
			operator,
			right,
		)
	}

	return expr, nil
}

// parseFactor 解析乘除表达式
func (ep *ExpressionParser) parseFactor() (data.GetValue, data.Control) {
	expr, acl := ep.parseUnary()
	if acl != nil {
		return expr, acl
	}
	for ep.current().Type == token.MUL || ep.current().Type == token.QUO ||
		ep.current().Type == token.REM {
		operator := ep.current()
		ep.next()

		right, acl := ep.parseUnary()
		if acl != nil {
			return nil, acl
		}
		expr = node.NewBinaryExpression(
			ep.FromCurrentToken(),
			expr,
			operator,
			right,
		)
	}

	return expr, nil
}

// parseUnary 解析一元表达式
func (ep *ExpressionParser) parseUnary() (data.GetValue, data.Control) {
	if ep.current().Type == token.SUB || ep.current().Type == token.NOT {
		operator := ep.current().Literal
		ep.next()

		right, acl := ep.parseUnary()
		if acl != nil {
			return nil, acl
		}
		return node.NewUnaryExpression(
			ep.FromCurrentToken(),
			operator,
			right,
		), nil
	}

	// 处理前缀自增自减
	if ep.current().Type == token.INCR || ep.current().Type == token.DECR {
		operator := ep.current()
		ep.next()

		right, acl := ep.parseUnary()
		if acl != nil {
			return nil, acl
		}
		for ep.current().Type == token.SEMICOLON {
			// 跳过没意义的分号
			ep.next()
		}
		if operator.Type == token.INCR {
			return node.NewUnaryIncr(
				ep.FromCurrentToken(),
				right,
			), nil
		} else {
			return node.NewUnaryDecr(
				ep.FromCurrentToken(),
				right,
			), nil
		}
	}

	return ep.parsePrimary()
}

// parsePrimary 解析基本表达式
func (ep *ExpressionParser) parsePrimary() (data.GetValue, data.Control) {
	switch ep.current().Type {
	case token.INT:
		value := ep.current().Literal
		ep.next()
		return node.NewIntLiteral(ep.FromCurrentToken(), value), nil
	case token.FLOAT:
		value := ep.current().Literal
		ep.next()
		return node.NewFloatLiteral(ep.FromCurrentToken(), value), nil
	case token.STRING:
		value := ep.current().Literal
		ep.next()
		return node.NewStringLiteral(ep.FromCurrentToken(), value), nil

	case token.TRUE:
		ep.next()
		return node.NewBooleanLiteral(ep.FromCurrentToken(), true), nil

	case token.FALSE:
		ep.next()
		return node.NewBooleanLiteral(ep.FromCurrentToken(), false), nil

	case token.NULL:
		ep.next()
		return node.NewNullLiteral(ep.FromCurrentToken()), nil
	case token.START_TAG, token.END_TAG, token.SEMICOLON:
		ep.next()
		return nil, nil
	default:
		if parser, ok := parserRouter[ep.current().Type]; ok {
			expr, acl := parser(ep.Parser).Parse()
			if acl != nil {
				return nil, acl
			}

			// 检查是否有后缀自增自减
			if ep.current().Type == token.INCR || ep.current().Type == token.DECR {
				operator := ep.current()
				ep.next()
				for ep.current().Type == token.SEMICOLON {
					// 跳过没意义的分号
					ep.next()
				}
				if operator.Type == token.INCR {
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

			return expr, nil
		}
		return nil, nil
	}
}
