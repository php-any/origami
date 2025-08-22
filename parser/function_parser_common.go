package parser

import (
	"errors"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// FunctionParserCommon 提供函数解析的通用功能
type FunctionParserCommon struct {
	*Parser
}

// NewFunctionParserCommon 创建一个新的通用函数解析器
func NewFunctionParserCommon(parser *Parser) *FunctionParserCommon {
	return &FunctionParserCommon{
		parser,
	}
}

// ParseFunctionBody 解析函数体
func (p *FunctionParserCommon) ParseFunctionBody() []node.Statement {
	stmtParser := NewMainStatementParser(p.Parser)
	var body []node.Statement
	if p.current().Type == token.LBRACE {
		p.next()
		for !p.currentIsTypeOrEOF(token.RBRACE) {
			stmt, acl := stmtParser.Parse()
			if acl != nil {
				p.addControl(acl)
			}
			if stmt != nil {
				body = append(body, stmt)
			}
		}
		p.next() // 跳过结束花括号
	}
	return body
}

// ParseParameters 解析参数列表
func (p *FunctionParserCommon) ParseParameters() ([]data.GetValue, data.Control) {
	tracking := p.StartTracking()
	// 检查左括号
	if p.current().Type != token.LPAREN {
		return nil, data.NewErrorThrow(tracking.EndBefore(), errors.New("参数列表前缺少左括号 '('"))
	}
	p.next()

	params := make([]data.GetValue, 0)
	// 如果下一个token是右括号，说明参数列表为空
	if p.current().Type == token.RPAREN {
		p.next()
		return params, nil
	}
	parser := p.Parser

	// 解析参数列表
	for {
		varType := ""
		name := ""
		isParams := false
		isReference := false // 是否引用
		// 解析参数名
		if p.current().Type != token.VARIABLE {
			isVar := false
			// (...args) 省却参数类型。
			if parser.current().Type == token.ELLIPSIS {
				isVar = true
				parser.next()
				name = parser.current().Literal
				isParams = true
				p.next()
			}

			// &$data
			if parser.checkPositionIs(0, token.BIT_AND) {
				isVar = true
				parser.next()
				name = parser.current().Literal
				p.next()
				isReference = true
			}

			// (string $data) 或 (?string $data)
			if !isVar && isIdentOrTypeToken(parser.current().Type) && parser.checkPositionIs(1, token.IDENTIFIER, token.VARIABLE) {
				isVar = true
				varType = parser.current().Literal
				p.next()

				name = parser.current().Literal
				p.next()
			}
			// (?string $data) 可空类型参数
			if !isVar && parser.checkPositionIs(0, token.TERNARY) && parser.checkPositionIs(1, token.IDENTIFIER) && parser.checkPositionIs(2, token.IDENTIFIER, token.VARIABLE) {
				isVar = true
				p.next() // 跳过问号
				varType = "?" + parser.current().Literal
				p.next()

				name = parser.current().Literal
				p.next()
			}
			// (data: string)
			if !isVar && parser.checkPositionIs(0, token.IDENTIFIER, token.VARIABLE) && parser.checkPositionIs(1, token.COLON) && parser.checkPositionIs(2, token.IDENTIFIER) {
				name = parser.current().Literal
				p.next()
				if parser.checkPositionIs(0, token.COLON) {
					p.next()
					varType = parser.current().Literal
					p.next()
				}
				isVar = true
			}
			// fun(data)
			if !isVar && isIdentOrTypeToken(parser.current().Type) && parser.checkPositionIs(1, token.RPAREN, token.COMMA) {
				name = parser.current().Literal
				p.next()
				isVar = true
			}
			if !isVar {
				return nil, data.NewErrorThrow(tracking.EndBefore(), errors.New("参数缺少变量名"))
			}
		} else {
			name = parser.current().Literal
			p.next()
			if parser.checkPositionIs(0, token.COLON) {
				p.next()
				varType = parser.current().Literal
				p.next()
			}
		}

		// 创建参数类型
		var paramType data.Types
		if strings.HasPrefix(varType, "?") {
			// 可空类型
			baseType := data.NewBaseType(varType[1:]) // 去掉问号
			paramType = data.NewNullableType(baseType)
		} else {
			// 普通类型
			paramType = data.NewBaseType(varType)
		}

		// 添加参数到作用域
		val := p.scopeManager.CurrentScope().AddVariable(name, paramType, tracking.EndBefore())

		// 解析默认值
		var defaultValue data.GetValue
		if p.current().Type == token.ASSIGN {
			p.next()
			exprParser := NewExpressionParser(p.Parser)
			var acl data.Control
			defaultValue, acl = exprParser.Parse()
			if acl != nil {
				return nil, acl
			}
		}

		// 创建参数节点
		if isParams {
			param := node.NewParameters(tracking.EndBefore(), val.GetName(), val.GetIndex(), defaultValue, val.GetType())
			params = append(params, param)
		} else if isReference {
			if defaultValue != nil {
				return nil, data.NewErrorThrow(tracking.EndBefore(), errors.New("参数为引用的变量不能有默认值"))
			}
			// 覆盖变量为引用
			p.scopeManager.CurrentScope().variables[val.GetName()] = node.NewVariableReference(tracking.EndBefore(), val.GetName(), val.GetIndex(), val.GetType())
			param := node.NewParameterReference(tracking.EndBefore(), val.GetName(), val.GetIndex(), val.GetType())
			params = append(params, param)
		} else {
			param := node.NewParameter(tracking.EndBefore(), val.GetName(), val.GetIndex(), defaultValue, val.GetType())
			params = append(params, param)
		}

		if p.current().Type == token.COMMA {
			p.next()
			// 检查逗号后是否直接跟着右括号（这是语法错误）
			if p.current().Type == token.RPAREN {
				return nil, data.NewErrorThrow(tracking.EndBefore(), errors.New("逗号后缺少参数"))
			}
		} else if p.current().Type != token.RPAREN {
			return nil, data.NewErrorThrow(tracking.EndBefore(), errors.New("参数后缺少逗号 ',' 或右括号 ')'"))
		} else {
			break
		}
	}
	p.nextAndCheck(token.RPAREN)
	return params, nil
}

// isIdentOrTypeToken 判断 token 是否为标识符或类型关键字（如 string、int、bool、float、array 等）
func isIdentOrTypeToken(t token.TokenType) bool {
	return t == token.IDENTIFIER ||
		t == token.BOOL ||
		t == token.ARRAY
}
