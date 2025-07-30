package parser

import (
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

	// 解析函数名
	if fp.current().Type != token.IDENTIFIER {
		fp.addError("缺少函数名")
		return nil, nil
	}
	name := fp.current().Literal

	if fp.namespace != nil {
		name = fp.namespace.GetName() + "\\" + name
	}

	fp.next()

	// 创建新的函数作用域
	fp.scopeManager.NewScope()

	// 解析参数列表
	params, acl := fp.parseParameters()
	if acl != nil {
		return nil, acl
	}
	ret := fp.parserReturnType()

	// 解析函数体
	body := fp.parseBlock()

	vars := fp.scopeManager.CurrentScope().GetVariables()

	// 弹出函数作用域
	fp.scopeManager.PopScope()

	f := node.NewFunctionStatement(
		fp.NewTokenFrom(fp.GetStart()),
		name,
		params,
		body,
		vars,
		ret,
	)

	if acl := fp.vm.AddFunc(f); acl != nil {
		fp.addError(acl.AsString())
	}

	return f, nil
}

// parseParameters 解析参数列表
func (fp *FunctionParser) parseParameters() ([]data.GetValue, data.Control) {
	vp := &FunctionParserCommon{Parser: fp.Parser}
	return vp.ParseParameters()
}

func (fp FunctionParser) parserReturnType() data.Types {
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
			if fp.current().Type == token.IDENTIFIER {
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
				fp.addError("缺少返回类型")
				return nil
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
			return nil
		} else if len(returnTypes) == 1 {
			return returnTypes[0]
		} else {
			// 多个返回类型，创建多返回值类型
			return data.NewMultipleReturnType(returnTypes)
		}
	}

	// 没有返回类型声明，返回 nil
	return nil
}
