package parser

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

type LparenParser struct {
	*Parser
}

// NewLparenParser 类型转换 or lambda or 括号表达式
func NewLparenParser(parser *Parser) StatementParser {
	return &LparenParser{
		parser,
	}
}

// Parse 解析左括号后的内容
func (ep *LparenParser) Parse() (data.GetValue, data.Control) {
	tracking := ep.StartTracking()
	// 检查是否是类型转换: (string) $data
	if ep.isTypeCast() {
		ep.nextAndCheck(token.LPAREN) // 跳过左括号
		return ep.parseTypeCast(tracking)
	}

	// 检查是否是 Lambda 表达式: (a, b) => {}
	if ep.isLambdaExpression() {
		return ep.parseLambdaExpression(tracking)
	}

	// 检查是否是括号表达式: (a + b)
	ep.nextAndCheck(token.LPAREN) // 跳过左括号
	return ep.parseParenthesizedExpression(tracking)
}

// isTypeCast 检查是否是类型转换
func (ep *LparenParser) isTypeCast() bool {
	// 检查模式: (TYPE_NAME) EXPRESSION
	// 其中 TYPE_NAME 可以是标识符、array 关键字或内置类型关键字（如 bool）
	if ep.checkPositionIs(1, token.IDENTIFIER, token.ARRAY, token.BOOL) &&
		ep.checkPositionIs(2, token.RPAREN) &&
		!ep.checkPositionIs(3, token.ARRAY_KEY_VALUE) {
		//  &&
		//		ep.checkPositionIs(3, token.IDENTIFIER, token.VARIABLE, token.LPAREN, token.INT, token.FLOAT, token.STRING, token.NULL, token.TRUE, token.FALSE)
		return true
	}
	return false
}

// parseTypeCast 解析类型转换
func (ep *LparenParser) parseTypeCast(tracking *PositionTracker) (data.GetValue, data.Control) {
	typeName := ep.current().Literal()
	ep.next()                     // 跳过类型名
	ep.nextAndCheck(token.RPAREN) // 跳过右括号

	// PHP 类型转换是一元运算，优先级高于 . / + / == 等；
	// 不能用 parseStatement，否则 (int)$v."'" 会变成 (int)($v."'")。
	exprParser := NewExpressionParser(ep.Parser)
	val, acl := exprParser.parseUnary()
	if acl != nil {
		return nil, acl
	}
	fn, ok := ep.vm.GetFunc(typeName)
	if !ok {
		return nil, data.NewErrorThrow(tracking.EndBefore(), errors.New("未定义的转换函数:"+typeName))
	}
	return node.NewCallExpression(tracking.EndBefore(), typeName, []data.GetValue{val}, fn), nil
}

// isLambdaExpression 检查是否是 Lambda 表达式
func (ep *LparenParser) isLambdaExpression() bool {
	// 如果紧跟 fn 关键字，这是箭头函数 fn() => expr，不是 Lambda (a, b) => expr
	if ep.checkPositionIs(1, token.FN) {
		return false
	}
	// 括号表达式：(new Foo())、(array) 等，绝不是 (params) => body
	if ep.checkPositionIs(1, token.NEW, token.ARRAY, token.ISSET, token.EMPTY, token.ECHO, token.INCLUDE, token.REQUIRE, token.FN) {
		return false
	}
	// 检查是否包含 => 符号，需要正确处理 () 与 [] 嵌套（数组内的 => 不是 lambda）
	pos := 1        // 从 ( 后面开始检查
	parenCount := 1 // 已处于起始 ( 之内
	bracketCount := 0

	for pos < len(ep.tokens)-ep.position {
		tokenType := ep.tokens[ep.position+pos].Type()

		switch tokenType {
		case token.LPAREN:
			parenCount++
		case token.RPAREN:
			parenCount--
			if parenCount == 0 {
				// 找到了与起始 '(' 匹配的 ')'：仅当紧跟 => 才可能是 lambda
				if pos+1 < len(ep.tokens)-ep.position &&
					ep.tokens[ep.position+pos+1].Type() == token.ARRAY_KEY_VALUE {
					// PHP 关联数组可用 (expr) => value 作键；括号内必须像参数列表才是 lambda
					return ep.looksLikeLambdaParameterList(1, pos)
				}
				return false
			}
		case token.LBRACKET:
			bracketCount++
		case token.RBRACKET:
			if bracketCount > 0 {
				bracketCount--
			}
		case token.ARRAY_KEY_VALUE:
			// 仅当不在数组字面量内、且仍在参数列表括号层时，才是 (params) => body
			if parenCount == 0 && bracketCount == 0 {
				return ep.looksLikeLambdaParameterList(1, pos)
			}
		}
		pos++
	}
	return false
}

// looksLikeLambdaParameterList 判断 tokens[start, end) 是否像 lambda 参数列表。
// 用于区分 (a, b) => expr 与 PHP 数组键 (is_int($k) ? $v : $k) => $v。
func (ep *LparenParser) looksLikeLambdaParameterList(start, end int) bool {
	if start >= end {
		return true // () =>
	}

	i := start
	for i < end {
		if !ep.segmentLooksLikeLambdaParameter(&i, end) {
			return false
		}
		if i >= end {
			return true
		}
		if ep.tokens[ep.position+i].Type() != token.COMMA {
			return false
		}
		i++ // 跳过 ,
	}
	return true
}

// segmentLooksLikeLambdaParameter 判断从 *i 起的一段是否像单个参数，成功后 *i 指向段末（逗号或 end）。
func (ep *LparenParser) segmentLooksLikeLambdaParameter(i *int, end int) bool {
	tok := func(off int) token.TokenType {
		return ep.tokens[ep.position+off].Type()
	}

	// 跳过参数属性 #[...]
	for *i < end && tok(*i) == token.HASH {
		*i++
		if *i < end && tok(*i) == token.LBRACKET {
			depth := 1
			*i++
			for *i < end && depth > 0 {
				switch tok(*i) {
				case token.LBRACKET:
					depth++
				case token.RBRACKET:
					depth--
				}
				*i++
			}
		}
	}

	// 跳过可见性 / readonly
	for *i < end {
		switch tok(*i) {
		case token.PUBLIC, token.PRIVATE, token.PROTECTED, token.READONLY:
			*i++
			continue
		}
		break
	}

	// 可选类型：?Type、Type|Type、namespace\Type
	if *i < end && tok(*i) == token.TERNARY && *i+1 < end && isIdentOrTypeToken(tok(*i+1)) {
		*i += 2
		for *i < end && (tok(*i) == token.BIT_OR || tok(*i) == token.NAMESPACE_SEPARATOR || isIdentOrTypeToken(tok(*i))) {
			*i++
		}
	} else if *i < end && isIdentOrTypeToken(tok(*i)) {
		// ident( 是函数调用，不是类型声明（类型后应为 $var / | / & / ...）
		if *i+1 < end && tok(*i+1) == token.LPAREN {
			return false
		}
		*i++
		for *i < end {
			t := tok(*i)
			if t == token.NAMESPACE_SEPARATOR || t == token.BIT_OR || isIdentOrTypeToken(t) {
				*i++
				continue
			}
			// 泛型 Type<...>
			if t == token.LT {
				depth := 1
				*i++
				for *i < end && depth > 0 {
					switch tok(*i) {
					case token.LT:
						depth++
					case token.GT:
						depth--
					}
					*i++
				}
				continue
			}
			break
		}
	}

	if *i < end && tok(*i) == token.BIT_AND {
		*i++
	}
	if *i < end && tok(*i) == token.ELLIPSIS {
		*i++
	}

	if *i >= end {
		return false
	}

	switch tok(*i) {
	case token.VARIABLE:
		*i++
	case token.IDENTIFIER, token.UNUSED:
		// Origami 简写：(a, b) =>；PHP 参数名必须是 $var
		next := *i + 1
		if next < end {
			switch tok(next) {
			case token.COMMA, token.ASSIGN:
				*i++
			default:
				return false
			}
		} else {
			*i++
		}
	default:
		return false
	}

	// 默认值 = expr：跳过到本段结束（顶层逗号）
	if *i < end && tok(*i) == token.ASSIGN {
		*i++
		parenDepth := 0
		bracketDepth := 0
		braceDepth := 0
		for *i < end {
			switch tok(*i) {
			case token.LPAREN:
				parenDepth++
			case token.RPAREN:
				if parenDepth > 0 {
					parenDepth--
				}
			case token.LBRACKET:
				bracketDepth++
			case token.RBRACKET:
				if bracketDepth > 0 {
					bracketDepth--
				}
			case token.LBRACE:
				braceDepth++
			case token.RBRACE:
				if braceDepth > 0 {
					braceDepth--
				}
			case token.COMMA:
				if parenDepth == 0 && bracketDepth == 0 && braceDepth == 0 {
					return true
				}
			}
			*i++
		}
		return true
	}

	return true
}

// parseLambdaExpression 解析 Lambda 表达式
func (ep *LparenParser) parseLambdaExpression(tracking *PositionTracker) (data.GetValue, data.Control) {
	fp := &FunctionParser{
		ep.Parser,
	}

	// 创建新的函数作用域
	fp.scopeManager.NewScope(true)

	// 解析参数列表
	params, acl := fp.parseParameters()
	if acl != nil {
		return nil, acl
	}
	fp.nextAndCheck(token.ARRAY_KEY_VALUE)

	// 解析函数体
	body, acl := fp.parseBlock()
	if acl != nil {
		return nil, acl
	}
	vars := fp.scopeManager.CurrentScope().GetVariables()

	// 弹出函数作用域
	fp.scopeManager.PopScope()

	parent := make(map[int]int)
	for _, parentVariable := range fp.scopeManager.CurrentScope().GetVariables() {
		for _, childVariable := range vars {
			if childVariable.GetName() == parentVariable.GetName() {
				// 形参由调用方传入，不应从父作用域捕获，否则体内会误读外层同名变量
				if isParameterName(params, childVariable.GetName()) {
					continue
				}
				parent[childVariable.GetIndex()] = parentVariable.GetIndex()
			}
		}
	}

	return node.NewLambdaExpression(
		tracking.EndBefore(),
		params,
		body,
		vars,
		parent,
	), nil
}

// parseParenthesizedExpression 解析括号表达式
func (ep *LparenParser) parseParenthesizedExpression(tracking *PositionTracker) (data.GetValue, data.Control) {
	// 括号内须解析完整表达式；仅用 parseTernary 会在 (a ? b : c) + d 处提前结束
	exprParser := NewExpressionParser(ep.Parser)
	expr, acl := exprParser.parseTernary()
	if acl != nil {
		return nil, acl
	}
	termTracker := ep.StartTracking()
	for ep.current().Type() == token.ADD || ep.current().Type() == token.SUB {
		op := ep.current()
		ep.next()
		right, acl := exprParser.parseRangeOperand()
		if acl != nil {
			return nil, acl
		}
		expr = node.NewBinaryExpression(termTracker.EndBefore(), expr, op, right)
	}
	// 检查是否有右括号
	if ep.current().Type() != token.RPAREN {
		// 某些表达式路径会在 parseStatement 阶段提前消费掉 ')'（如 ($i-1)），
		// 这里兼容该情况，避免误报“缺少右括号”。
		if ep.position > 0 && ep.tokens[ep.position-1].Type() == token.RPAREN {
			vp := &VariableParser{ep.Parser}
			return vp.parseSuffix(expr)
		}
		return nil, data.NewErrorThrow(tracking.EndBefore(), fmt.Errorf("缺少右括号 ')'"))
	}
	ep.next() // 跳过右括号

	// 括号表达式后应支持后缀操作（-> 属性/方法、[] 索引、() 调用等）
	vp := &VariableParser{ep.Parser}
	return vp.parseSuffix(expr)
}
