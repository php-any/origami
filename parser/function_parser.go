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
	// 跳过 & 引用返回标记（function &name() {}）
	returnsReference := false
	if fp.checkPositionIs(0, token.BIT_AND) {
		returnsReference = true
		fp.next()
	}
	outerReference := fp.currentReturnsReference
	fp.currentReturnsReference = returnsReference
	defer func() { fp.currentReturnsReference = outerReference }()
	// 解析函数名
	if !fp.checkPositionIs(0, token.IDENTIFIER) {
		if fp.checkPositionIs(0, token.LPAREN) {
			// 直接解析闭包值: function() {}
			fp.enterStaticScope()
			defer fp.leaveStaticScope()
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
			// use 变量必须在解析函数体前注册到闭包作用域，否则若函数体未直接引用
			//（例如只再传给内层 use），parent 映射会丢失。
			fp.registerClosureUseCaptures(captures)
			ret, acl := fp.parserReturnType()
			if acl != nil {
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
				fp.strictTypes,
			)

			fn.Ret = data.DeclaredTypeRef(ret)
			fn.ReturnsReference = returnsReference
			return fn, nil
		}

		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("缺少函数名"))
	}
	name := fp.current().Literal()

	if fp.namespace != nil {
		name = fp.namespace.GetName() + "\\" + name
	}

	// 设置当前函数名（用于 __FUNCTION__/__METHOD__ 魔术常量）
	fp.currentFunction = name

	fp.next()

	fp.enterStaticScope()
	defer fp.leaveStaticScope()

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
		returnsReference,
		fp.strictTypes,
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

// registerClosureUseCaptures 将 use ($a, &$b) 中的变量立即注册到闭包作用域。
// PHP：即便函数体未直接出现该变量（例如只再传给内层 use），也应存在于闭包作用域。
func (fp *FunctionParser) registerClosureUseCaptures(captures []UseCapture) {
	for _, c := range captures {
		fp.scopeManager.CurrentScope().AddVariable(c.Name, nil, fp.FromCurrentToken())
	}
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
	return parseDeclaredReturn(fp.Parser)
}
