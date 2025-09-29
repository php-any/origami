package parser

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// FunctionParser 表示函数解析器
type FunctionParser struct {
	*Parser
}

// NewFunctionParser 创建一个新的函数解析器
func NewFunctionParser(parser *Parser) StatementParser {
	return &FunctionParser{
		parser,
	}
}

// Parse 解析函数声明
func (fp *FunctionParser) Parse() (data.GetValue, data.Control) {
	// 跳过 function 关键字
	fp.next()
	tracker := fp.StartTracking()
	// 解析函数名
	if fp.current().Type != token.IDENTIFIER {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("缺少函数名"))
	}
	name := fp.current().Literal

	if fp.namespace != nil {
		name = fp.namespace.GetName() + "\\" + name
	}

	fp.next()

	// 创建新的函数作用域
	fp.scopeManager.NewScope(false)

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
	vars := fp.scopeManager.CurrentScope().GetVariables()

	// 弹出函数作用域
	fp.scopeManager.PopScope()

	f := node.NewFunctionStatement(
		tracker.EndBefore(),
		name,
		params,
		body,
		vars,
		ret,
	)

	if acl := fp.vm.AddFunc(f); acl != nil {
		return nil, acl
	}

	return f, nil
}

// parseParameters 解析参数列表
func (fp *FunctionParser) parseParameters() ([]data.GetValue, data.Control) {
	vp := &FunctionParserCommon{Parser: fp.Parser}
	return vp.ParseParameters()
}

func (fp FunctionParser) parserReturnType() (data.Types, data.Control) {
	// 检查是否有返回类型声明
	// 语法: function name(): returnType 或 function name(): ?returnType
	// 或者: function name(): type1, type2, type3 (多返回值)
	if fp.current().Type == token.COLON {
		fp.next() // 跳过冒号

		// 解析返回类型列表
		var returnTypes []data.Types

		for {
			// 检查是否是可空类型语法 ?type
			isNullable := false
			if fp.current().Type == token.TERNARY {
				isNullable = true
				fp.next() // 跳过问号
			}

			// 解析返回类型
			if fp.checkPositionIs(0, token.IDENTIFIER, token.STRING, token.INT, token.FLOAT, token.BOOL, token.ARRAY) {
				returnType := fp.current().Literal
				fp.next()

				// 创建基础类型
				baseType := data.NewBaseType(returnType)

				// 如果是可空类型，包装为基础类型的可空版本
				if isNullable {
					baseType = data.NewNullableType(baseType)
				}

				returnTypes = append(returnTypes, baseType)
			} else {
				return nil, data.NewErrorThrow(fp.newFrom(), errors.New("无法识别返回类型的定义符号"))
			}

			// 检查是否有更多类型（逗号分隔）
			if fp.current().Type == token.COMMA {
				fp.next() // 跳过逗号
				continue
			}

			// 没有更多类型，结束解析
			break
		}

		// 根据返回类型数量决定返回类型
		if len(returnTypes) == 0 {
			return nil, nil
		} else if len(returnTypes) == 1 {
			return returnTypes[0], nil
		} else {
			// 多个返回类型，创建多返回值类型
			return data.NewMultipleReturnType(returnTypes), nil
		}
	}

	// 没有返回类型声明，返回 nil
	return nil, nil
}
