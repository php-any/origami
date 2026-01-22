package parser

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
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
func (p *FunctionParserCommon) ParseFunctionBody() ([]data.GetValue, data.Control) {
	stmtParser := NewMainStatementParser(p.Parser)
	var body []data.GetValue
	if p.current().Type() == token.LBRACE {
		p.next()
		last := p.position
		for !p.currentIsTypeOrEOF(token.RBRACE) {
			stmt, acl := stmtParser.Parse()
			if acl != nil {
				return nil, acl
			}
			if stmt != nil {
				body = append(body, stmt)
			} else if last == p.position {
				t := p.current()
				fmt.Println(t)
				panic("出现死循环")
			}
		}
		p.next() // 跳过结束花括号
	}
	return body, nil
}

// ParseParameters 解析参数列表
func (p *FunctionParserCommon) ParseParameters() ([]data.GetValue, data.Control) {
	tracking := p.StartTracking()
	// 检查左括号
	if p.current().Type() != token.LPAREN {
		return nil, data.NewErrorThrow(tracking.EndBefore(), errors.New("参数列表前缺少左括号 '('"))
	}
	p.next()

	params := make([]data.GetValue, 0)
	// 如果下一个token是右括号，说明参数列表为空
	if p.current().Type() == token.RPAREN {
		p.next()
		return params, nil
	}

	// 解析参数列表
	for {
		param, _, acl := parseSingleParameter(p.Parser)
		if acl != nil {
			return nil, acl
		}
		params = append(params, param)

		if p.current().Type() == token.COMMA {
			p.next()
			// 检查逗号后是否直接跟着右括号（这是语法错误）
			if p.current().Type() == token.RPAREN {
				break
			}
		} else if p.current().Type() != token.RPAREN {
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
		t == token.ARRAY ||
		t == token.NULL || // 支持 null 作为类型声明的一部分
		t == token.FALSE || // 支持 false 作为类型声明的一部分
		t == token.INT ||
		t == token.STRING ||
		t == token.FLOAT ||
		t == token.GENERIC_TYPE ||
		t == token.SELF // 支持 self 作为类型关键字
}
