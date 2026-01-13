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
	if !fp.checkPositionIs(0, token.IDENTIFIER) {
		if fp.checkPositionIs(0, token.LPAREN) {
			// 直接解析闭包值: function() {}
			// 创建新的函数作用域
			fp.scopeManager.NewScope(false)

			// 解析参数列表
			params, acl := fp.parseParameters()
			if acl != nil {
				return nil, acl
			}
			// 解析 use 捕获列表（可选）：function () use ($a, $b) {}
			captures, acl := fp.parseClosureUse()
			if acl != nil {
				return nil, acl
			}
			if _, acl := fp.parserReturnType(); acl != nil {
				return nil, acl
			}
			// 解析函数体
			body, acl := fp.parseBlock()
			if acl != nil {
				return nil, acl
			}

			// 当前作用域中的局部变量（闭包内部）
			vars := fp.scopeManager.CurrentScope().GetVariables()

			// 先根据 use(&$var) 把需要按引用捕获的变量替换成 VariableReference
			if len(captures) > 0 {
				for _, c := range captures {
					if !c.IsReference {
						continue
					}
					if childVar, ok := fp.scopeManager.CurrentScope().GetVariable(c.Name); ok {
						fp.scopeManager.CurrentScope().SetVariable(
							c.Name,
							node.NewVariableReference(
								fp.FromCurrentToken(),
								childVar.GetName(),
								childVar.GetIndex(),
								childVar.GetType(),
							),
						)
					}
				}
				// 更新 vars，确保其中的引用变量已经变成 VariableReference
				vars = fp.scopeManager.CurrentScope().GetVariables()
			}

			// 弹出函数作用域，返回到外部作用域
			fp.scopeManager.PopScope()

			// 构建 parent 映射，仅捕获 use 声明的变量
			parent := make(map[int]int)
			if len(captures) > 0 {
				for _, outer := range fp.scopeManager.CurrentScope().GetVariables() {
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
			}

			fn := node.NewLambdaExpression(
				tracker.EndBefore(),
				params,
				body,
				vars,
				parent,
			)

			return fn, nil
		}

		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("缺少函数名"))
	}
	name := fp.current().Literal()

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

	//if acl := fp.vm.AddFunc(f); acl != nil {
	//	return nil, acl
	//}

	return f, nil
}

// parseParameters 解析参数列表
func (fp *FunctionParser) parseParameters() ([]data.GetValue, data.Control) {
	vp := &FunctionParserCommon{Parser: fp.Parser}
	return vp.ParseParameters()
}

// UseCapture 表示 use 子句中的捕获变量信息
// Name 为变量名（不含 $），IsReference 表示是否使用 & 按引用捕获
type UseCapture struct {
	Name        string
	IsReference bool
}

// parseClosureUse 解析闭包的 use 捕获列表
// 支持:
//   - use ($a, $b)
//   - use (&$a, $b)
func (fp *FunctionParser) parseClosureUse() ([]UseCapture, data.Control) {
	if !fp.checkPositionIs(0, token.USE) {
		return nil, nil
	}
	fp.next() // 跳过 use
	if acl := fp.nextAndCheck(token.LPAREN); acl != nil {
		return nil, acl
	}
	var captures []UseCapture
	for {
		isRef := false
		// 可选的引用符号 &（use (&$var)）
		if fp.checkPositionIs(0, token.BIT_AND) {
			isRef = true
			fp.next()
		}

		if !fp.checkPositionIs(0, token.VARIABLE) {
			return nil, data.NewErrorThrow(fp.FromCurrentToken(), errors.New("use 语法错误，期望变量"))
		}
		name := fp.current().Literal()
		if len(name) > 0 && name[0] == '$' {
			name = name[1:]
		}
		captures = append(captures, UseCapture{
			Name:        name,
			IsReference: isRef,
		})
		fp.next() // 跳过变量
		if fp.current().Type() == token.COMMA {
			fp.next()
			continue
		}
		break
	}
	if acl := fp.nextAndCheck(token.RPAREN); acl != nil {
		return nil, acl
	}
	return captures, nil
}

func (fp FunctionParser) parserReturnType() (data.Types, data.Control) {
	// 检查是否有返回类型声明
	// 语法: function name(): returnType 或 function name(): ?returnType
	// 或者: function name(): type1, type2, type3 (多返回值)
	if fp.current().Type() == token.COLON {
		fp.next() // 跳过冒号

		// 解析返回类型列表
		var returnTypes []data.Types

		for {
			// 检查是否是可空类型语法 ?type
			isNullable := false
			if fp.current().Type() == token.TERNARY {
				isNullable = true
				fp.next() // 跳过问号
			}

			// 解析一个“返回类型表达式”，支持联合类型：array|string|false
			// 其中每个原子类型可以是标识符、内置类型、null、false 等
			var unionTypes []data.Types

			parseOneTypeAtom := func() (data.Types, data.Control) {
				if !fp.checkPositionIs(0,
					token.IDENTIFIER,
					token.STRING,
					token.INT,
					token.FLOAT,
					token.BOOL,
					token.ARRAY,
					token.NULL,
					token.FALSE,
					token.STATIC,
				) {
					return nil, data.NewErrorThrow(fp.newFrom(), errors.New("无法识别返回类型的定义符号"))
				}
				name := fp.current().Literal()
				fp.next()

				// 如果是基础类型，直接返回
				if data.ISBaseType(name) {
					return data.NewBaseType(name), nil
				}

				// 尝试解析完整的类名（包括命名空间）
				if full, ok := fp.findFullClassNameByNamespace(name); ok {
					return data.NewBaseType(full), nil
				}

				// 如果无法解析，返回原始名称
				return data.NewBaseType(name), nil
			}

			// 第一个类型原子
			firstType, acl := parseOneTypeAtom()
			if acl != nil {
				return nil, acl
			}
			unionTypes = append(unionTypes, firstType)

			// 后续的 |Type 原子
			for fp.current().Type() == token.BIT_OR {
				fp.next() // 跳过 |
				nextType, acl := parseOneTypeAtom()
				if acl != nil {
					return nil, acl
				}
				unionTypes = append(unionTypes, nextType)
			}

			// 将本次解析出的类型（可能是单一，也可能是联合）加入返回类型列表
			var thisType data.Types
			if len(unionTypes) == 1 {
				thisType = unionTypes[0]
			} else {
				// 联合类型：array|string|false 之类
				thisType = data.NewUnionType(unionTypes)
			}
			if isNullable {
				thisType = data.NewNullableType(thisType)
			}
			returnTypes = append(returnTypes, thisType)

			// 检查是否有更多类型（逗号分隔）
			if fp.current().Type() == token.COMMA {
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
