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

	// PHP 允许函数名与 '(' 之间存在空白：andi (1, 2)
	if p.checkPositionIs(0, token.LPAREN) {
		// 词法器会把 \func_get_args 收成单个 IDENTIFIER "\func_get_args"
		// （Carbon 等大量使用 \func_get_args()），这里还原为语言结构节点
		if builtin, ok := globalLanguageConstruct(name); ok {
			return p.parseGlobalLanguageConstruct(tracker, builtin)
		}
		if call, acl := p.parseIdentCall(tracker, name); acl != nil || call != nil {
			return call, acl
		}
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
		// ( 函数调用 div() 或可调用变量 describe()
		if p.checkPositionIs(0, token.LPAREN) {
			if builtin, ok := globalLanguageConstruct(name); ok {
				return p.parseGlobalLanguageConstruct(tracker, builtin)
			}
			return p.parseIdentCall(tracker, name)
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

		// 处理 ::class 语法（须优先于「关键字作方法名」，避免 class 被当成静态成员）
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

		// 函数静态调用 Log::info / Factory::new() 等（PHP 允许关键字作方法名，但 class 除外）
		if p.checkPositionIs(0, token.SCOPE_RESOLUTION) {
			nextType := p.tokens[p.position+1].Type()
			isStaticMember := nextType == token.IDENTIFIER || nextType == token.VARIABLE ||
				nextType == token.COMPACT || nextType == token.UNSET || nextType == token.ISSET ||
				(nextType > token.KEYWORD_START && nextType < token.VALUE_START && nextType != token.CLASS)
			if isStaticMember {
				return p.parseStaticCall(tracker, name)
			}
		}

		if p.checkPositionIs(0, token.DOT) {
			// 字符串连接左侧的裸标识符按常量运行时解析（APP_ROOT . '/x'）
			expr := node.NewConstantName(tracker.EndBefore(), name)
			vp := &VariableParser{p.Parser}
			return vp.parseSuffix(expr)
		}
		if p.checkPositionIs(0, token.OBJECT_OPERATOR) {
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

	// 常量一律运行时查找：define() 与请求级覆盖（如 PHP_SAPI）在解析期尚不可见或可能变化。
	// 赋值目标：describe = ... 将裸标识符视为变量名
	if p.checkPositionIs(0, token.ASSIGN, token.ADD_EQ, token.SUB_EQ, token.MUL_EQ, token.QUO_EQ, token.REM_EQ, token.CONCAT_EQ, token.NULL_COALESCE_ASSIGN) {
		val := p.scopeManager.CurrentScope().AddVariable(name, nil, tracker.EndBefore())
		return node.NewVariableWithFirst(tracker.EndBefore(), val), nil
	}

	return node.NewConstantName(tracker.EndBefore(), name), nil
}

// parseIdentCall 解析标识符函数调用 name(...)
// PHP：函数名与变量名空间分离，bare name(...) 始终是函数调用；
// 变量作为可调用必须写 $name(...)。同名形参不得抢占函数调用。
func (p *IdentParser) parseIdentCall(tracker *PositionTracker, name string) (data.GetValue, data.Control) {
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
		return vp.parseSuffix(callExpr)
	}
	if InLSP {
		stmt, acl := vp.parseFunctionCall()
		if acl != nil {
			return nil, acl
		}
		callExpr := node.NewCallExpression(tracker.EndBefore(), name, stmt, nil)
		return vp.parseSuffix(callExpr)
	}
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

// parseStaticCall 解析静态调用（如 Log::info 或 Log::property）
func (p *IdentParser) parseStaticCall(tracker *PositionTracker, className string) (data.GetValue, data.Control) {
	// 尝试获取完整的类名。
	// findFullClassNameByNamespace 在 ok=false 时仍可能返回命名空间候选名（供 autoload），
	// 必须采用，否则未加载类的静态访问会丢前缀。
	fullClassName := className
	if full, _ := p.findFullClassNameByNamespace(className); full != "" {
		fullClassName = full
	}

	// 尝试获取类
	stmt, has := p.vm.GetClass(fullClassName)

	// 跳过 ::
	p.next()
	// 获取方法名或属性名
	isVariable := p.current().Type() == token.VARIABLE
	fnName := p.current().Literal()

	// Class::$method()：变量方法名需运行时求值（DateFactory / Str::$method）
	if isVariable && p.checkPositionIs(1, token.LPAREN) {
		varInfo := p.scopeManager.LookupVariable(fnName)
		from := tracker.EndBefore()
		if varInfo == nil {
			val := p.scopeManager.CurrentScope().AddVariable(fnName, nil, from)
			varInfo = node.NewVariableWithFirst(from, val)
		}
		methodExpr := node.NewVariableWithFirst(from, varInfo)
		p.next() // 跳过变量
		vp := &VariableParser{p.Parser}
		var classRef data.GetValue
		if has {
			classRef = stmt
		} else {
			classRef = node.NewStringLiteral(from, fullClassName)
		}
		expr := node.NewCallStaticDynamicMethod(from, classRef, methodExpr)
		return vp.parseSuffix(expr)
	}

	p.next()

	// 如果是 VARIABLE，去掉 $ 前缀（静态属性 Class::$prop）
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
			// 类未加载，创建延迟调用。
			// 不要再拼当前命名空间：findFullClassNameByNamespace 已处理
			// use / 当前命名空间 / 全局类；盲目前缀会破坏 `use SortDirection`
			// 这类全局导入（变成 Illuminate\...\SortDirection）。
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

// globalLanguageConstruct 识别被词法器收成 "\func_get_args" 形式的全局语言结构
func globalLanguageConstruct(name string) (string, bool) {
	n := name
	for strings.HasPrefix(n, "\\") {
		n = n[1:]
	}
	switch n {
	case "func_get_args", "func_num_args":
		return n, true
	default:
		return "", false
	}
}

func (p *IdentParser) parseGlobalLanguageConstruct(tracker *PositionTracker, builtin string) (data.GetValue, data.Control) {
	if p.checkPositionIs(0, token.LPAREN) {
		p.next()
		if p.current().Type() != token.RPAREN {
			return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("%s() 不接受参数", builtin))
		}
		p.next()
	}
	from := tracker.EndBefore()
	switch builtin {
	case "func_get_args":
		return node.NewFuncGetArgs(from), nil
	case "func_num_args":
		return node.NewFuncNumArgs(from), nil
	default:
		return nil, data.NewErrorThrow(from, fmt.Errorf("未知语言结构: %s", builtin))
	}
}
