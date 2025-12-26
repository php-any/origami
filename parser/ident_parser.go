package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

type IdentParser struct {
	*Parser
}

func NewIdentParser(parser *Parser) StatementParser {
	return &IdentParser{
		Parser: parser,
	}
}

// Parse 解析标识符表达式
func (p *IdentParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	name := p.current().Literal()
	startToken := p.current()
	p.next()

	// 标签语句：name: 换行
	// 只在标识符后紧跟一个冒号、且后面不是类型标注（a: string）时，将其视为标签定义
	if p.checkPositionIs(0, token.COLON) && p.peek(1).Line() != p.peek(0).Line() {
		p.next() // 跳过 :
		return node.NewLabelStatement(tracker.EndBefore(), name), nil
	}

	// 函数调用模式 div {} 或者 div []
	if p.checkPositionIs(0, token.LBRACE) {
		if full, ok := p.findFullFunNameByNamespace(name); ok {
			fn, ok := p.vm.GetFunc(full)
			if !ok {
				c, acl := p.vm.LoadPkg(full)
				if acl != nil {
					return nil, acl
				}
				if c != nil {
					switch c.(type) {
					case data.ClassStmt:
						return p.parseClassInit(tracker, full)
					}
				}
				return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("未定义的函数:"+full+" {}。"))
			}
			v, acl := NewLbraceParser(p.Parser).Parse()
			return node.NewCallExpression(tracker.EndBefore(), fn.GetName(), []data.GetValue{v}, fn), acl
		}

		// 检查是否是便捷方式创建 class{}
		if full, ok := p.findFullClassNameByNamespace(name); ok {
			c, acl := p.vm.LoadPkg(full)
			if acl != nil {
				return nil, acl
			}
			if c != nil {
				switch c.(type) {
				case data.ClassStmt:
					return p.parseClassInit(tracker, full)
				}
			}
		}

		v, acl := NewLbraceParser(p.Parser).Parse()
		if p.namespace != nil {
			name = p.namespace.GetName() + "\\" + name
		}
		return node.NewCallExpression(tracker.EndBefore(), name, []data.GetValue{v}, &node.CallFunctionLater{Name: name, Ctx: p.vm.CreateContext(nil)}), acl
	} else if p.checkPositionIs(0, token.LBRACKET) && !p.isTokensAdjacent(startToken, p.current()) {
		if full, ok := p.findFullFunNameByNamespace(name); ok {
			fn, ok := p.vm.GetFunc(full)
			if !ok {
				return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("未定义的函数:"+full+" []。"))
			}
			v, acl := NewLbracketParser(p.Parser).Parse()
			return node.NewCallExpression(tracker.EndBefore(), fn.GetName(), []data.GetValue{v}, fn), acl
		}
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("未定义的函数:"+name+" []。"))
	}

	// 检查是否是变量的类型
	if p.checkPositionIs(0, token.ASSIGN) {
		val := p.scopeManager.CurrentScope().AddVariable(name, nil, tracker.EndBefore())
		return node.NewVariableWithFirst(tracker.EndBefore(), val), nil
	}
	if p.checkPositionIs(0, token.VARIABLE) || p.checkPositionIs(1, token.ASSIGN) {
		// int $num 或者 int i = 0
		ty := name
		name = p.current().Literal()
		p.next()
		val := p.scopeManager.CurrentScope().AddVariable(name, data.NewBaseType(ty), tracker.EndBefore())
		return node.NewVariableWithFirst(tracker.EndBefore(), val), nil
	}

	checkToken := p.current()

	// 检查 startToken 和 checkToken 之间是否连贯
	if p.isTokensAdjacent(startToken, checkToken) {
		// ( 函数调用 div()
		if p.checkPositionIs(0, token.LPAREN) {
			// 创建函数调用表达式
			vp := &VariableParser{p.Parser}
			if full, ok := p.findFullFunNameByNamespace(name); ok {
				stmt, acl := vp.parseFunctionCall()
				if acl != nil {
					return nil, acl
				}
				fn, ok := p.vm.GetFunc(full)
				if !ok {
					return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("函数(%s)先加载后才能使用", name))
				}
				callExpr := node.NewCallExpression(tracker.EndBefore(), full, stmt, fn)
				// 支持函数调用结果继续链式操作：app()->name, app()[0], app()()
				return vp.parseSuffix(callExpr)
			} else if InLSP {
				stmt, acl := vp.parseFunctionCall()
				if acl != nil {
					return nil, acl
				}
				callExpr := node.NewCallExpression(tracker.EndBefore(), name, stmt, nil)
				return vp.parseSuffix(callExpr)
			} else {
				namespace := ""
				if p.namespace != nil {
					namespace = p.namespace.Name
				}
				stmt, acl := vp.parseFunctionCall()
				if acl != nil {
					return nil, acl
				}
				callExpr := node.NewCallTodo(node.NewCallExpression(tracker.EndBefore(), name, stmt, nil), namespace)
				return vp.parseSuffix(callExpr)
			}
		}
		// 变量定义
		if p.checkPositionIs(0, token.COLON) && p.checkPositionIs(1, token.IDENTIFIER) {
			// a: string
			p.next()
			ty := p.current().Literal()
			p.next()
			val := p.scopeManager.CurrentScope().AddVariable(name, data.NewBaseType(ty), tracker.EndBefore())
			expr := node.NewVariableWithFirst(tracker.EndBefore(), val)
			// 解析后续操作（函数调用、数组访问等）
			vp := &VariableParser{p.Parser}
			return vp.parseSuffix(expr)
		}

		// 函数静态调用 Log::info 或 Log::$property
		if p.checkPositionIs(0, token.SCOPE_RESOLUTION) &&
			(p.checkPositionIs(1, token.IDENTIFIER) || p.checkPositionIs(1, token.VARIABLE)) {
			return p.parseStaticCall(tracker, name)
		}

		// 处理 ::class 语法
		if p.checkPositionIs(0, token.SCOPE_RESOLUTION) && p.checkPositionIs(1, token.CLASS) {
			className := name
			if full, ok := p.findFullClassNameByNamespace(className); ok {
				className = full
			}
			p.next() // 跳过 ::
			p.next() // 跳过 class
			// 返回类名字符串
			return data.NewStringValue(className), nil
		}

		if p.checkPositionIs(0, token.OBJECT_OPERATOR, token.DOT) {
			val := p.scopeManager.CurrentScope().AddVariable(name, nil, tracker.EndBefore())
			expr := node.NewVariableWithFirst(tracker.EndBefore(), val)
			vp := &VariableParser{p.Parser}
			return vp.parseSuffix(expr)
		}
	}

	// 检查是否是变量
	varInfo := p.scopeManager.LookupVariable(name)
	if varInfo != nil {
		// 解析后续操作（函数调用、数组访问等）
		vp := &VariableParser{p.Parser}
		return vp.parseSuffix(varInfo)
	}

	if p.checkPositionIs(0, token.LT) && p.checkPositionIs(2, token.GT) && p.checkPositionIs(3, token.LPAREN) {
		// DB<Name>( 才进入分型便捷 new
		className, ok := p.findFullClassNameByNamespace(name)
		if !ok {
			return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("class %s 不存在", name))
		}
		p.next() // <
		generaList := make([]string, 0)
		for !p.checkPositionIs(0, token.GT) {
			generaName, ok := p.tryFindTypes()
			if !ok {
				return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("class %s 不存在", name))
			}
			// p.next()
			generaList = append(generaList, generaName.String())
			if p.checkPositionIs(0, token.COMMA) {
				p.next() // ,
			}
		}
		p.next() // >
		vp := VariableParser{Parser: p.Parser}
		args, acl := vp.parseFunctionCall()
		if acl != nil {
			return nil, acl
		}
		n := &node.NewClassGenerated{
			NewExpression: node.NewNewExpression(
				tracker.EndBefore(),
				className,
				args,
			),
			T: generaList,
		}
		if p.checkPositionIs(0, token.OBJECT_OPERATOR) {
			// 解析链式调用
			return vp.parseSuffix(n)
		}

		return n, nil
	}

	if p.scopeManager.CurrentScope().IsLambda() {
		// 检查是否是变量
		varInfo := p.scopeManager.LookupParentVariable(name)
		if varInfo != nil {
			val := p.scopeManager.CurrentScope().AddVariable(name, varInfo.GetType(), tracker.EndBefore())
			expr := node.NewVariableWithFirst(tracker.EndBefore(), val)
			vp := &VariableParser{p.Parser}
			return vp.parseSuffix(expr)
		}
	}

	// 是否是define后的字符串
	if v, ok := p.vm.GetConstant(name); ok {
		return v, nil
	}

	return node.NewStringLiteral(tracker.EndBefore(), name), nil
}

// parseStaticCall 解析静态调用（如 Log::info 或 Log::property）
func (p *IdentParser) parseStaticCall(tracker *PositionTracker, className string) (data.GetValue, data.Control) {
	// 尝试获取完整的类名
	fullClassName := className
	if full, ok := p.findFullClassNameByNamespace(className); ok {
		fullClassName = full
	}

	// 尝试获取类
	stmt, has := p.vm.GetClass(fullClassName)

	// 跳过 ::
	p.next()
	// 获取方法名或属性名
	isVariable := p.current().Type() == token.VARIABLE
	fnName := p.current().Literal()
	p.next()

	// 如果是 VARIABLE，去掉 $ 前缀
	if isVariable && len(fnName) > 0 && fnName[0] == '$' {
		fnName = fnName[1:]
	}

	// 获取命名空间
	namespace := ""
	if p.namespace != nil {
		namespace = p.namespace.Name
	}

	// 判断是方法调用还是属性访问
	if p.checkPositionIs(0, token.LPAREN) {
		// 静态方法调用 Class::method()
		if has {
			// 类已加载，创建静态方法调用
			vp := &VariableParser{p.Parser}
			expr := node.NewCallStaticMethod(tracker.EndBefore(), stmt, fnName)
			return vp.parseSuffix(expr)
		} else {
			// 类未加载，创建延迟调用
			vp := &VariableParser{p.Parser}
			expr := node.NewCallStaticMethodLater(tracker.EndBefore(), fullClassName, fnName, namespace)
			return vp.parseSuffix(expr)
		}
	} else {
		// 静态属性访问 Class::property
		if has {
			// 类已加载，创建静态属性访问
			vp := &VariableParser{p.Parser}
			expr := node.NewCallStaticProperty(tracker.EndBefore(), stmt, fnName)
			return vp.parseSuffix(expr)
		} else {
			if strings.Index(fullClassName, "\\") == -1 && p.namespace != nil {
				fullClassName = p.namespace.Name + "\\" + fullClassName
			}

			// 类未加载，创建延迟调用
			vp := &VariableParser{p.Parser}
			expr := node.NewCallStaticPropertyLater(tracker.EndBefore(), fullClassName, fnName, namespace)
			return vp.parseSuffix(expr)
		}
	}
}

func (p *IdentParser) parseClassInit(tracker *PositionTracker, className string) (data.GetValue, data.Control) {
	p.nextAndCheck(token.LBRACE)

	kv := map[string]data.GetValue{}
	// 解释 key: stmt
	for !p.checkPositionIs(0, token.RBRACE, token.EOF) {
		if !p.checkPositionIs(0, token.IDENTIFIER) {
			return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("初始类 %s 的属性名必须是标识符", className))
		}
		key := p.current().Literal()
		p.next()
		acl := p.nextAndCheck(token.COLON)
		if acl != nil {
			return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("初始类 %s 的属性名后面必须是(:)符号", className))
		}
		value, acl := p.parseStatement()
		if acl != nil {
			return nil, acl
		}
		if p.checkPositionIs(0, token.COMMA) {
			p.next()
		}

		kv[key] = value
	}
	p.nextAndCheck(token.RBRACE)

	return node.NewInitClass(tracker.EndBefore(), className, kv), nil
}
