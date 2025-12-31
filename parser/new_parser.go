package parser

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// NewStructParser 表示 new 表达式解析器
type NewStructParser struct {
	*Parser
}

// NewNewParser 创建一个新的 new 表达式解析器
func NewNewParser(parser *Parser) StatementParser {
	return &NewStructParser{
		parser,
	}
}

// Parse 解析 new 表达式
func (p *NewStructParser) Parse() (data.GetValue, data.Control) {
	// 允许 new() 作为函数调用, 方便兼容 go 语言的库
	if p.checkPositionIs(1, token.LPAREN) {
		tracker := p.StartTracking()
		name := "new"
		p.next()
		if full, ok := p.uses[name]; ok {
			// 创建函数调用表达式
			vp := &VariableParser{p.Parser}
			stmt, acl := vp.parseFunctionCall()
			fn, ok := p.vm.GetFunc(full)
			if !ok {
				return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("new关键字作为函数调用, 函数(%s)先加载后才能使用", name))
			}
			return node.NewCallExpression(tracker.EndBefore(), full, stmt, fn), acl
		}
	}

	tracker := p.StartTracking()
	// 跳过 new 关键字
	p.next()

	// 检查是否是匿名类 new class { ... }
	if p.checkPositionIs(0, token.CLASS) {
		return p.parseAnonymousClass(tracker)
	}

	// 检查是否是变量类名 new $variable()
	var classNameExpr data.GetValue
	if p.checkPositionIs(0, token.VARIABLE) {
		// 解析变量表达式作为类名
		vp := &VariableParser{p.Parser}
		classNameExpr = vp.parseVariable()

		// 解析参数列表
		var args []data.GetValue
		if p.checkPositionIs(0, token.LPAREN) {
			var acl data.Control
			args, acl = vp.parseFunctionCall()
			if acl != nil {
				return nil, acl
			}
		}

		// 创建使用变量类名的 new 表达式
		n := node.NewNewVariableExpression(
			tracker.EndBefore(),
			classNameExpr,
			args,
		)

		if p.checkPositionIs(0, token.OBJECT_OPERATOR) {
			// 解析链式调用
			return vp.parseSuffix(n)
		}

		return n, nil
	}

	// 检查是否是 new self
	if p.checkPositionIs(0, token.SELF) {
		p.next() // 跳过 self

		// 检查是否有参数列表（括号）
		var args []data.GetValue
		var acl data.Control
		vp := VariableParser{Parser: p.Parser}

		if p.checkPositionIs(0, token.LPAREN) {
			// 有括号，解析参数列表
			args, acl = vp.parseFunctionCall()
			if acl != nil {
				return nil, acl
			}
		} else {
			// 没有括号，使用空参数列表
			args = []data.GetValue{}
		}

		// 创建 new self 表达式
		n := node.NewNewSelfExpression(
			tracker.EndBefore(),
			args,
		)

		if p.checkPositionIs(0, token.OBJECT_OPERATOR) {
			// 解析链式调用
			return vp.parseSuffix(n)
		}

		return n, nil
	}

	// 检查是否是 new static
	if p.checkPositionIs(0, token.STATIC) {
		p.next() // 跳过 static

		// 检查是否有参数列表（括号）
		var args []data.GetValue
		var acl data.Control
		vp := VariableParser{Parser: p.Parser}

		if p.checkPositionIs(0, token.LPAREN) {
			// 有括号，解析参数列表
			args, acl = vp.parseFunctionCall()
			if acl != nil {
				return nil, acl
			}
		} else {
			// 没有括号，使用空参数列表
			args = []data.GetValue{}
		}

		// 创建 new static 表达式
		n := node.NewNewStaticExpression(
			tracker.EndBefore(),
			args,
		)

		if p.checkPositionIs(0, token.OBJECT_OPERATOR) {
			// 解析链式调用
			return vp.parseSuffix(n)
		}

		return n, nil
	}

	// 解析类名
	if !p.checkPositionIs(0, token.IDENTIFIER, token.GENERIC_TYPE) {
		return nil, data.NewErrorThrow(p.newFrom(), errors.New("new关键字后面必须跟类名或变量; 当前="+p.current().Literal()))
	}

	isGenerated := false
	generateType := ""
	if p.checkPositionIs(0, token.GENERIC_TYPE) {
		isGenerated = true
		generateType = p.current().Literal()
	}

	// 获取完整的类名路径
	className, acl := p.getClassName(true)
	if acl != nil {
		return nil, acl
	}

	if p.checkPositionIs(0, token.COMMA, token.RPAREN, token.SEMICOLON) {
		return node.NewNewExpression(
			tracker.EndBefore(),
			className,
			[]data.GetValue{},
		), nil
	}

	// 处理泛型参数
	var genericTypes []data.Types
	if p.checkPositionIs(0, token.LT) {
		p.next() // 跳过 <
		genericTypes = make([]data.Types, 0)

		for !p.checkPositionIs(0, token.GT) {
			genericType, ok := p.tryFindTypes()
			if !ok {
				return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("泛型类型解析失败"))
			}
			// p.next()
			genericTypes = append(genericTypes, genericType)

			if p.checkPositionIs(0, token.COMMA) {
				p.next() // 跳过 ,
			}
		}
		p.next() // 跳过 >
	}
	// 解析参数列表
	vp := VariableParser{Parser: p.Parser}
	args, acl := vp.parseFunctionCall()
	if acl != nil {
		return nil, acl
	}

	// 如果有泛型参数，创建泛型 new 表达式
	if len(genericTypes) > 0 {
		// 将泛型类型转换为字符串列表
		genericStrings := make([]string, len(genericTypes))
		for i, genericType := range genericTypes {
			genericStrings[i] = genericType.String()
		}

		n := &node.NewClassGenerated{
			NewExpression: node.NewNewExpression(
				tracker.EndBefore(),
				className,
				args,
			),
			T: genericStrings,
		}

		if p.checkPositionIs(0, token.OBJECT_OPERATOR) {
			// 解析链式调用
			return vp.parseSuffix(n)
		}

		return n, nil
	}

	// 普通 new 表达式
	n := node.NewNewExpression(
		tracker.EndBefore(),
		className,
		args,
	)

	if p.checkPositionIs(0, token.OBJECT_OPERATOR) {
		// 解析链式调用
		return vp.parseSuffix(n)
	}

	if isGenerated {
		return &node.NewGenerated{
			NewExpression: n,
			T:             generateType,
		}, nil
	}

	// 创建 new 表达式节点
	return n, nil
}

// parseAnonymousClass 解析匿名类 new class { ... }
func (p *NewStructParser) parseAnonymousClass(tracker *PositionTracker) (data.GetValue, data.Control) {
	// 获取当前 token 的行号信息用于生成匿名类名
	currentLine := p.current().Line()
	currentPos := p.current().Pos()

	// 获取文件名
	var fileName string
	if p.source != nil && *p.source != "" {
		fileName = filepath.Base(*p.source)
		// 移除文件扩展名
		if ext := filepath.Ext(fileName); ext != "" {
			fileName = strings.TrimSuffix(fileName, ext)
		}
		// 将文件名中的特殊字符替换为下划线，确保类名合法
		fileName = strings.ReplaceAll(fileName, ".", "_")
		fileName = strings.ReplaceAll(fileName, "-", "_")
		fileName = strings.ReplaceAll(fileName, " ", "_")
	} else {
		fileName = "unknown"
	}

	// 跳过 class 关键字
	p.next()

	// 根据文件名、行号和列号信息生成匿名类名
	anonymousClassName := fmt.Sprintf("class@anonymous@%s@%d:%d", fileName, currentLine, currentPos)
	if p.namespace != nil {
		anonymousClassName = p.namespace.GetName() + "\\" + anonymousClassName
	}

	// 解析泛型参数（如果存在）
	var types []data.Types
	if p.checkPositionIs(0, token.LT) {
		cp := &ClassParser{
			Parser:               p.Parser,
			FunctionParserCommon: NewFunctionParserCommon(p.Parser),
		}
		types = cp.parseGeneric()
	}

	// 解析构造函数参数（在 extends 之前）
	var constructorArgs []data.GetValue
	if p.checkPositionIs(0, token.LPAREN) {
		vp := VariableParser{Parser: p.Parser}
		args, acl := vp.parseFunctionCall()
		if acl != nil {
			return nil, acl
		}
		constructorArgs = args
	}

	// 解析继承
	var extends string
	if p.current().Type() == token.EXTENDS {
		p.next()
		var acl data.Control
		extends, acl = p.getClassName(true)
		if acl != nil {
			return nil, acl
		}
	}

	// 解析实现的接口
	var implements []string
	if p.current().Type() == token.IMPLEMENTS {
		p.next()
		for {
			interfaceName, acl := p.getClassName(true)
			if acl != nil {
				return nil, acl
			}

			implements = append(implements, interfaceName)

			if p.current().Type() != token.COMMA {
				break
			}
			p.next()
		}
	}

	// 解析类体
	if p.current().Type() != token.LBRACE {
		currentToken := p.current()
		return nil, data.NewErrorThrow(p.newFrom(), fmt.Errorf("匿名类声明后缺少左花括号 '{'，当前token: %s (类型: %v)", currentToken.Literal(), currentToken.Type()))
	}
	p.next()

	// 使用 ClassParser 的解析逻辑来解析类成员
	cp := &ClassParser{
		Parser:               p.Parser,
		FunctionParserCommon: NewFunctionParserCommon(p.Parser),
	}

	// 解析类成员
	properties := make([]data.Property, 0)
	staticProperties := make(map[string]data.Property)
	methods := map[string]data.Method{}
	staticMethods := map[string]data.Method{}
	for !p.currentIsTypeOrEOF(token.RBRACE) {
		// 先尝试解析注解
		var memberAnnotations []*node.Annotation
		for p.current().Type() == token.AT {
			ann, acl := cp.parseAnnotation()
			if acl != nil {
				return nil, acl
			}
			if ann != nil {
				memberAnnotations = append(memberAnnotations, ann)
			}
		}

		// 解析 abstract 关键字（可以在访问修饰符之前）
		isAbstractMethod := false
		if p.current().Type() == token.ABSTRACT {
			isAbstractMethod = true
			p.next()
		}

		// 解析访问修饰符
		modifier := cp.parseModifier()
		if modifier == "" {
			return nil, data.NewErrorThrow(p.newFrom(), errors.New("缺少访问修饰符"))
		}

		// 解析 abstract 关键字（也可以在访问修饰符之后）
		if p.current().Type() == token.ABSTRACT {
			isAbstractMethod = true
			p.next()
		}

		// 解析readonly关键字（在访问修饰符之后）
		isReadonly := false
		if p.current().Type() == token.READONLY {
			isReadonly = true
			p.next()
		}

		// 解析static关键字
		isStatic := false
		if p.current().Type() == token.STATIC {
			isStatic = true
			p.next()
		}

		// 解析属性或方法
		if p.current().Type() == token.VAR ||
			p.current().Type() == token.CONST ||
			p.current().Type() == token.VARIABLE ||
			isIdentOrTypeToken(p.current().Type()) ||
			(p.checkPositionIs(0, token.TERNARY) && isIdentOrTypeToken(p.peek(1).Type())) {
			prop, acl := cp.parsePropertyWithAnnotations(modifier, isStatic, isReadonly, memberAnnotations)
			if acl != nil {
				return nil, acl
			}
			if prop != nil {
				if isStatic {
					staticProperties[prop.GetName()] = prop
				} else {
					properties = append(properties, prop)
				}
			}
		} else if p.current().Type() == token.FUNC {
			method, acl := cp.parseMethodWithAnnotations(modifier, isStatic, isAbstractMethod, memberAnnotations)
			if acl != nil {
				return nil, acl
			}
			if method != nil {
				if isStatic {
					staticMethods[method.GetName()] = method
				} else {
					methods[method.GetName()] = method
				}
			}
		} else if p.current().Type() == token.SEMICOLON {
			p.next()
			continue
		} else {
			return nil, data.NewErrorThrow(p.newFrom(), errors.New("缺少属性或方法声明"))
		}
	}
	p.next() // 跳过结束花括号

	// 创建匿名类
	c := node.NewClassStatement(
		tracker.EndBefore(),
		anonymousClassName,
		extends,
		implements,
		properties,
		methods,
	)
	for s, property := range staticProperties {
		defaultValue := property.GetDefaultValue()
		if defaultValue != nil {
			baseCtx := p.vm.CreateContext([]data.Variable{})
			v, acl := defaultValue.GetValue(baseCtx)
			if acl != nil {
				return nil, acl
			}
			c.StaticProperty.Store(s, v)
		} else {
			c.StaticProperty.Store(s, data.NewNullValue())
		}
	}
	c.StaticMethods = staticMethods

	// 处理父类构造函数
	if c.Construct == nil {
		vm := p.vm
		var last data.ClassStmt = c
		for last != nil && last.GetExtend() != nil {
			ext := last.GetExtend()
			var acl data.Control
			last, acl = vm.GetOrLoadClass(*ext)
			if acl != nil {
				return nil, acl
			}
			if construct, ok := last.GetMethod(token.ConstructName); ok {
				c.Construct = construct
				break
			}
		}
	}

	// 处理泛型
	var classStmt data.ClassStmt = c
	if types != nil {
		cg := &node.ClassGeneric{
			ClassStatement: c,
			Generic:        types,
		}
		classStmt = cg
	}

	// 创建匿名类 new 表达式节点（延迟执行）
	// 在执行阶段才注册类、实例化对象并调用构造函数
	anonymousClassExpr := node.NewNewAnonymousClassExpression(
		tracker.EndBefore(),
		classStmt,
		constructorArgs,
		types,
	)

	// 检查是否有链式调用
	if p.checkPositionIs(0, token.OBJECT_OPERATOR) {
		vp := VariableParser{Parser: p.Parser}
		return vp.parseSuffix(anonymousClassExpr)
	}

	return anonymousClassExpr, nil
}

// findVariable 在变量列表中查找指定名称的变量
func findVariable(varies []data.Variable, name string) (data.Variable, error) {
	for _, vary := range varies {
		if vary.GetName() == name {
			return vary, nil
		}
	}
	return nil, fmt.Errorf("无法找到变量: %s", name)
}
