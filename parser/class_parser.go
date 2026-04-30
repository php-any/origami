package parser

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/lexer"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ClassParser 表示类解析器
type ClassParser struct {
	*Parser
	*FunctionParserCommon
	currentClassName string // 当前正在解析的类名（用于 self 关键字）
}

// NewClassParser 创建一个新的类解析器
func NewClassParser(parser *Parser) StatementParser {
	return &ClassParser{
		Parser:               parser,
		FunctionParserCommon: NewFunctionParserCommon(parser),
	}
}

// Parse 解析类定义
func (p *ClassParser) Parse() (data.GetValue, data.Control) {
	// 解析类前的注解
	var annotations []*node.Annotation
	for p.checkPositionIs(0, token.AT, token.HASH) {
		annotation, acl := p.parseAnnotation()
		if acl != nil {
			return nil, acl
		}
		if annotation != nil {
			annotations = append(annotations, annotation)
		}
	}

	// 跳过class关键字
	p.next()
	tracker := p.StartTracking()

	// 解析类名
	className := p.parseClassName()
	if className == "" {
		return nil, data.NewErrorThrow(p.newFrom(), errors.New("class 后缺少类名"))
	}
	if p.namespace != nil {
		className = p.namespace.GetName() + "\\" + className
	}

	p.currentClass = className
	p.vm.SetClassPathCache(className, *p.source)

	var types []data.Types
	var genericParamNames []string
	if p.checkPositionIs(0, token.LT) {
		// 定义泛型类型解释
		types = p.parseGeneric()
		// 收集泛型参数名
		for _, t := range types {
			if t, ok := t.(data.Generic); ok {
				genericParamNames = append(genericParamNames, t.Name)
			}
		}
		// 预处理：将类体内所有泛型参数名的token类型替换为GENERIC_TYPE
		if len(genericParamNames) > 0 {
			// 找到类体token区间
			braceCount := 0
			startIdx := p.position // 当前token是'{'前
			// 向后找到第一个'{'
			for i := p.position; i < len(p.tokens); i++ {
				if p.tokens[i].Type() == token.LBRACE {
					startIdx = i
					break
				}
			}
			endIdx := startIdx
			braceCount = 1
			for i := startIdx + 1; i < len(p.tokens); i++ {
				if p.tokens[i].Type() == token.LBRACE {
					braceCount++
				} else if p.tokens[i].Type() == token.RBRACE {
					braceCount--
					if braceCount == 0 {
						endIdx = i
						break
					}
				}
			}
			// 替换token类型 - 注意：由于Token现在是接口，我们需要创建新的WorkerToken来替换
			for i := startIdx + 1; i < endIdx; i++ {
				for _, param := range genericParamNames {
					if p.tokens[i].Literal() == param && p.tokens[i].Type() == token.IDENTIFIER {
						// 创建新的WorkerToken替换原来的token
						p.tokens[i] = lexer.NewWorkerToken(
							token.GENERIC_TYPE,
							p.tokens[i].Literal(),
							p.tokens[i].Start(),
							p.tokens[i].End(),
							p.tokens[i].Line(),
							p.tokens[i].Pos(),
						)
					}
				}
			}
		}
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
		return nil, data.NewErrorThrow(p.newFrom(), errors.New("类声明后缺少左花括号 '{'"))
	}
	p.next()

	// 解析类成员
	properties := make([]data.Property, 0)
	staticProperties := make(map[string]data.Property)
	methods := map[string]data.Method{}
	staticMethods := map[string]data.Method{}
	var traits []string                       // trait 列表
	var constructorProperties []data.Property // 构造函数中声明的属性
	for !p.currentIsTypeOrEOF(token.RBRACE) {
		// 先尝试解析注解
		var memberAnnotations []*node.Annotation
		for p.checkPositionIs(0, token.AT, token.HASH) {
			ann, acl := p.parseAnnotation()
			if acl != nil {
				return nil, acl
			}
			if ann != nil {
				memberAnnotations = append(memberAnnotations, ann)
			}
		}

		// 检查是否是 use 语句（用于 trait）
		if p.current().Type() == token.USE {
			traitNames, acl := p.parseTraitUse()
			if acl != nil {
				return nil, acl
			}
			traits = append(traits, traitNames...)
			continue
		}

		// 解析 abstract 关键字（可以在访问修饰符之前）
		isAbstractMethod := false
		if p.current().Type() == token.ABSTRACT {
			isAbstractMethod = true
			p.next()
		}

		// 解析访问修饰符
		modifier := p.parseModifier()
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
			prop, acl := p.parsePropertyWithAnnotations(modifier, isStatic, isReadonly, memberAnnotations)
			if acl != nil {
				return nil, acl
			}
			if prop != nil {
				if isStatic || prop.GetIsStatic() {
					staticProperties[prop.GetName()] = prop
				} else {
					properties = append(properties, prop)
				}
			}
		} else if p.current().Type() == token.FUNC {
			method, props, acl := p.parseMethodWithAnnotations(modifier, isStatic, isAbstractMethod, memberAnnotations, &properties, &staticProperties)
			if acl != nil {
				return nil, acl
			}
			if method != nil {
				if isStatic {
					staticMethods[method.GetName()] = method
				} else {
					methods[method.GetName()] = method
				}
				// 如果是构造函数，添加从参数中声明的属性
				if method.GetName() == token.ConstructName && len(props) > 0 {
					constructorProperties = append(constructorProperties, props...)
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

	// 将构造函数中声明的属性添加到属性列表
	properties = append(properties, constructorProperties...)

	c := node.NewClassStatement(
		tracker.EndBefore(),
		className,
		extends,
		implements,
		properties,
		methods,
	)
	// 两遍评估静态属性/常量：第一遍处理不依赖 self:: 的，第二遍处理依赖 self:: 的
	var deferredStatic []struct {
		name     string
		defValue data.GetValue
	}
	for s, property := range staticProperties {
		defaultValue := property.GetDefaultValue()
		if defaultValue != nil {
			baseCtx := p.vm.CreateContext([]data.Variable{})
			v, acl := defaultValue.GetValue(baseCtx)
			if acl != nil {
				// 可能因为 self:: 引用未解析的常量而失败，延迟处理
				deferredStatic = append(deferredStatic, struct {
					name     string
					defValue data.GetValue
				}{name: s, defValue: defaultValue})
				continue
			}
			c.StaticProperty.Store(s, v)
		} else {
			c.StaticProperty.Store(s, data.NewNullValue())
		}
	}
	// 第二遍：使用 ClassMethodContext 评估延迟的静态属性
	for _, deferred := range deferredStatic {
		baseCtx := p.vm.CreateContext([]data.Variable{})
		classValue := data.NewClassValue(c, baseCtx)
		classCtx := classValue.CreateContext([]data.Variable{})
		v, acl := deferred.defValue.GetValue(classCtx)
		if acl != nil {
			return nil, acl
		}
		c.StaticProperty.Store(deferred.name, v)
	}
	c.StaticMethods = staticMethods

	// 合并 trait 的方法和属性
	if len(traits) > 0 {
		acl := p.mergeTraits(c, traits)
		if acl != nil {
			return nil, acl
		}
	}

	if c.Construct == nil {
		// 寻找父级构造函数
		vm := p.vm
		var last data.ClassStmt = c
		for last != nil && last.GetExtend() != nil {
			ext := *last.GetExtend()

			var acl data.Control
			last, acl = vm.GetOrLoadClass(ext)
			if acl != nil {
				return nil, acl
			}
			if construct, ok := last.GetMethod(token.ConstructName); ok {
				c.Construct = construct
				break
			}
		}
	}

	var acl data.Control

	// 构建类语句：处理泛型
	var classStmt data.ClassStmt = c

	if types != nil {
		// 创建泛型类
		cg := &node.ClassGeneric{
			ClassStatement: c,
			Generic:        types,
		}
		classStmt = cg
		acl = p.vm.AddClass(classStmt)
		if acl != nil {
			return nil, acl
		}
		if addAnn, ok := classStmt.(node.AddAnnotations); ok {
			acl = callClassAnnotation(p.Parser, &annotations, addAnn)
			if acl != nil {
				return nil, acl
			}
		}
		return classStmt, acl
	}

	acl = p.vm.AddClass(classStmt)
	if acl != nil {
		return nil, acl
	}
	if addAnn, ok := classStmt.(node.AddAnnotations); ok {
		acl = callClassAnnotation(p.Parser, &annotations, addAnn)
		if acl != nil {
			return nil, acl
		}
	}
	return classStmt, acl
}

// 调用注解
func callClassAnnotation(p *Parser, ans *[]*node.Annotation, c node.AddAnnotations) data.Control {
	for _, an := range *ans {
		an.Target = c.(data.GetValue)
	}
	for _, an := range *ans {
		obj, acl := an.GetValue(p.vm.CreateContext(nil))
		if acl != nil {
			return acl
		}
		if o, ok := obj.(*data.ClassValue); ok {
			c.AddAnnotations(o)
		}
	}
	return nil
}

// parseClassName 解析类名, 只管定义
func (p *ClassParser) parseClassName() string {
	if p.current().Type() != token.IDENTIFIER {
		return ""
	}
	name := p.current().Literal()
	p.next()
	return name
}

// ParseConstructorParameters 解析构造函数参数，支持在参数中声明类属性
// 语法：public function __construct(private string $name = 'UNKNOWN') {}
// 返回：参数列表、属性列表、错误
func (p *ClassParser) ParseConstructorParameters() ([]data.GetValue, []data.Property, data.Control) {
	tracking := p.StartTracking()
	// 检查左括号
	if p.current().Type() != token.LPAREN {
		return nil, nil, data.NewErrorThrow(tracking.EndBefore(), errors.New("参数列表前缺少左括号 '('"))
	}
	p.next()

	params := make([]data.GetValue, 0)
	properties := make([]data.Property, 0)

	// 如果下一个token是右括号，说明参数列表为空
	if p.current().Type() == token.RPAREN {
		p.next()
		return params, properties, nil
	}

	// 解析参数列表
	for {
		param, property, acl := parseSingleParameter(p.Parser)
		if property != nil {
			properties = append(properties, property)
		}
		if acl != nil {
			return nil, nil, acl
		}
		params = append(params, param)

		if p.current().Type() == token.COMMA {
			p.next()
			// 检查逗号后是否直接跟着右括号（这是语法错误）
			if p.current().Type() == token.RPAREN {
				break
			}
		} else if p.current().Type() != token.RPAREN {
			return nil, nil, data.NewErrorThrow(tracking.EndBefore(), errors.New("参数后缺少逗号 ',' 或右括号 ')'"))
		} else {
			break
		}
	}
	p.nextAndCheck(token.RPAREN)
	return params, properties, nil
}

// parseModifier 解析访问修饰符
func (p *ClassParser) parseModifier() string {
	switch p.current().Type() {
	case token.PUBLIC:
		p.next()
		return "public"
	case token.PROTECTED:
		p.next()
		return "protected"
	case token.PRIVATE:
		p.next()
		return "private"
	default:
		return "public"
	}
}

// parseAnnotation 解析注解
func (p *ClassParser) parseAnnotation() (*node.Annotation, data.Control) {
	tracker := p.StartTracking()

	// 检查是 @ 还是 #[
	if p.checkPositionIs(0, token.HASH) {
		// 处理 #[...] 格式的属性注解 (PHP 8.0+)
		p.next() // 跳过 #
		if p.current().Type() != token.LBRACKET {
			return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("属性注解格式错误，期望 #[...]"))
		}
		p.next() // 跳过 [

		// 解析注解名称
		if p.current().Type() != token.IDENTIFIER && p.current().Type() != token.NAMESPACE_SEPARATOR {
			return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("属性注解缺少名称"))
		}

		annotationName, acl := p.getClassName(true)
		if acl != nil {
			return nil, acl
		}

		// 解析注解参数
		arguments := make([]data.GetValue, 0)
		if p.current().Type() == token.LPAREN {
			vp := VariableParser{Parser: p.Parser}
			arguments, acl = vp.parseFunctionCall()
			if acl != nil {
				return nil, acl
			}
		}

		// 跳过 ]
		if p.current().Type() != token.RBRACKET {
			return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("属性注解缺少右方括号 ']'"))
		}
		p.next()

		// 创建注解节点
		annotation := node.NewAnnotation(
			tracker.EndBefore(),
			annotationName,
			arguments,
		)

		return annotation, nil
	}

	// 处理 @ 格式的注解
	// 跳过 @ 符号
	p.next()

	// 解析注解名称
	if p.current().Type() != token.IDENTIFIER {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("注解缺少名称"))
	}

	annotationName, acl := p.getClassName(true)
	if acl != nil {
		return nil, acl
	}

	// 解析注解参数
	arguments := make([]data.GetValue, 0)
	if p.current().Type() == token.LPAREN {
		vp := VariableParser{Parser: p.Parser}
		arguments, acl = vp.parseFunctionCall()
		if acl != nil {
			return nil, acl
		}
	}
	for p.checkPositionIs(0, token.SEMICOLON) {
		p.next()
	}
	// 创建注解节点
	annotation := node.NewAnnotation(
		tracker.EndBefore(),
		annotationName,
		arguments,
	)

	return annotation, nil
}

// parsePropertyWithAnnotations 解析属性（带注解）
func (p *ClassParser) parsePropertyWithAnnotations(modifier string, isStatic bool, isReadonly bool, annotations []*node.Annotation) (data.Property, data.Control) {
	tracker := p.StartTracking()

	// 解析访问修饰符（如果还没有解析）
	if modifier == "" {
		modifier = p.parseModifier()
		if modifier == "" {
			modifier = "public" // 默认为public
		}
	}

	// 解析static关键字（如果还没有解析）
	if !isStatic && p.current().Type() == token.STATIC {
		isStatic = true
		p.next()
	}

	// 解析readonly关键字（如果还没有解析）
	if !isReadonly && p.current().Type() == token.READONLY {
		isReadonly = true
		p.next()
	}

	// 支持 class 常量: public const BLOCKS = 'Mockery_Forward_Blocks';
	if p.current().Type() == token.CONST {
		// 跳过 const 关键字
		p.next()

		// 解析常量名（不带 $）
		if p.current().Type() != token.IDENTIFIER {
			return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("常量声明缺少名称"))
		}
		name := p.current().Literal()
		p.next()

		// 解析默认值
		var defaultValue data.GetValue
		var acl data.Control
		if p.current().Type() == token.ASSIGN {
			p.next()
			exprParser := NewExpressionParser(p.Parser)
			defaultValue, acl = exprParser.Parse()
			if acl != nil {
				return nil, acl
			}
		}

		// 解析分号
		if p.current().Type() == token.SEMICOLON {
			p.next()
		}

		// const 始终视为静态成员
		ret := node.NewPropertyWithReadonly(
			tracker.EndBefore(),
			name,
			modifier,
			true,
			isReadonly,
			defaultValue,
		)
		for _, an := range annotations {
			an.Target = ret
		}
		if len(annotations) != 0 {
			callAnn := make([]*node.CallAnn, 0)
			for _, an := range annotations {
				stmt, acl := p.vm.GetOrLoadClass(an.Name)
				if acl != nil {
					return nil, acl
				}
				object, acl := stmt.GetValue(p.vm.CreateContext(nil))
				if acl != nil {
					return nil, acl
				}
				obj, acl := an.GetValue(p.vm.CreateContext(object.(*data.ClassValue).Class.GetConstruct().GetVariables()))
				if acl != nil {
					if ann, ok := acl.(*node.CallAnn); !ok {
						return nil, acl
					} else {
						callAnn = append(callAnn, ann)
					}
				}
				if o, ok := obj.(*data.ClassValue); ok {
					ret.AddAnnotations(o)
				}
			}
			for i := len(callAnn) - 1; i >= 0; i-- {
				acl := callAnn[i].InitAnnotation()
				if acl != nil {
					return nil, acl
				}
			}
		}

		return ret, acl
	}

	// 解析属性类型（在访问修饰符之后，变量名之前）
	var propertyType data.Types
	if isIdentOrTypeToken(p.current().Type()) || p.checkPositionIs(0, token.NULL, token.FALSE, token.SELF) {
		// 检查是否是联合类型：string|int|null
		var unionTypes []data.Types

		// 解析第一个类型
		var firstType data.Types
		if p.checkPositionIs(0, token.NULL, token.FALSE) {
			firstType = data.NewBaseType(p.current().Literal())
			p.next()
		} else if p.current().Type() == token.SELF {
			// 处理 self 关键字
			p.next()
			if p.currentClassName != "" {
				firstType = data.NewBaseType(p.currentClassName)
			} else {
				firstType = data.NewBaseType("self")
			}
		} else {
			firstType = parseType(p.Parser)
		}

		if firstType != nil {
			unionTypes = append(unionTypes, firstType)

			// 处理后续的 |Type
			for p.current().Type() == token.BIT_OR {
				p.next() // 跳过 |
				var nextType data.Types
				if p.checkPositionIs(0, token.NULL, token.FALSE) {
					nextType = data.NewBaseType(p.current().Literal())
					p.next()
				} else if p.current().Type() == token.SELF {
					// 处理 self 关键字
					p.next()
					if p.currentClassName != "" {
						nextType = data.NewBaseType(p.currentClassName)
					} else {
						nextType = data.NewBaseType("self")
					}
				} else if isIdentOrTypeToken(p.current().Type()) {
					nextType = parseType(p.Parser)
				} else {
					break
				}
				if nextType != nil {
					unionTypes = append(unionTypes, nextType)
				}
			}

			// 创建类型
			if len(unionTypes) == 1 {
				propertyType = unionTypes[0]
			} else {
				propertyType = data.NewUnionType(unionTypes)
			}
		}
	} else if p.checkPositionIs(0, token.TERNARY) && (isIdentOrTypeToken(p.peek(1).Type()) || p.peek(1).Type() == token.SELF) {
		// ?int 或 ?self 方式
		p.next()
		if p.current().Type() == token.SELF {
			p.next()
			var baseType data.Types
			if p.currentClassName != "" {
				baseType = data.NewBaseType(p.currentClassName)
			} else {
				baseType = data.NewBaseType("self")
			}
			propertyType = data.NewNullableType(baseType)
		} else {
			// 处理 ?ClassName 这种可空类类型，需要结合命名空间解析完整类名
			name := p.current().Literal()
			p.next()

			var base data.Types
			// 内置基础类型（int/string/bool 等）保持原样
			if data.ISBaseType(name) {
				base = data.NewBaseType(name)
			} else if full, ok := p.findFullClassNameByNamespace(name); ok {
				// 若当前命名空间下存在对应类，则使用完整类名
				base = data.NewBaseType(full)
			} else {
				// 否则回退为原始名称
				base = data.NewBaseType(name)
			}

			propertyType = data.NewNullableType(base)
		}
	}

	// 解析属性名（普通属性必须是变量）
	if p.current().Type() != token.VARIABLE {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("缺少变量名"))
	}
	name := p.current().Literal()
	p.next()

	// 解析默认值
	var defaultValue data.GetValue
	var acl data.Control
	if p.current().Type() == token.ASSIGN {
		p.next()
		exprParser := NewExpressionParser(p.Parser)
		defaultValue, acl = exprParser.Parse()
		if acl != nil {
			return nil, acl
		}
	}

	// 解析分号
	if p.current().Type() == token.SEMICOLON {
		p.next()
	}

	ret := node.NewPropertyWithReadonly(
		tracker.EndBefore(),
		name,
		modifier,
		isStatic,
		isReadonly,
		defaultValue,
		propertyType,
	)
	for _, an := range annotations {
		an.Target = ret
	}
	if len(annotations) != 0 {
		callAnn := make([]*node.CallAnn, 0)
		for _, an := range annotations {
			stmt, acl := p.vm.GetOrLoadClass(an.Name)
			if acl != nil {
				return nil, acl
			}
			object, acl := stmt.GetValue(p.vm.CreateContext(nil))
			if acl != nil {
				return nil, acl
			}
			obj, acl := an.GetValue(p.vm.CreateContext(object.(*data.ClassValue).Class.GetConstruct().GetVariables()))
			if acl != nil {
				if ann, ok := acl.(*node.CallAnn); !ok {
					return nil, acl
				} else {
					callAnn = append(callAnn, ann)
				}
			}
			if o, ok := obj.(*data.ClassValue); ok {
				ret.AddAnnotations(o)
			}
		}
		for i := len(callAnn) - 1; i >= 0; i-- {
			acl := callAnn[i].InitAnnotation()
			if acl != nil {
				return nil, acl
			}
		}
	}

	return ret, acl
}

// parseMethodWithAnnotations 解析方法（带注解）
func (p *ClassParser) parseMethodWithAnnotations(modifier string, isStatic bool, isAbstract bool, annotations []*node.Annotation, properties *[]data.Property, staticProperties *map[string]data.Property) (data.Method, []data.Property, data.Control) {
	// 跳过function关键字
	p.next()
	tracker := p.StartTracking()
	p.scopeManager.NewScope(false)
	// 解析方法名
	// PHP 允许使用大部分关键字作为方法名（例如 unset、clone 等），
	// 这里仅禁止明显的符号/非法 token，放开关键字作为方法名的场景。
	if !(p.checkPositionIs(0, token.IDENTIFIER, token.SELF) ||
		(p.current().Type() > token.KEYWORD_START && p.current().Type() < token.VALUE_START)) {
		return nil, nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("方法名不符合规范, 不能使用符号(%s)", p.current().Literal()))
	}
	name := p.current().Literal()
	p.next()

	// 检查是否是构造函数
	isConstructor := name == token.ConstructName

	// 使用通用函数解析器解析参数和方法体
	var params []data.GetValue
	var constructorProps []data.Property
	var acl data.Control
	if isConstructor {
		params, constructorProps, acl = p.ParseConstructorParameters()
	} else {
		params, acl = p.ParseParameters()
	}
	if acl != nil {
		return nil, nil, acl
	}

	// 解析返回类型
	var retType data.Types
	if p.current().Type() == token.COLON {
		p.next() // 跳过冒号

		// 解析返回类型列表
		var returnTypes []data.Types

		for {
			// 检查是否是可空类型语法 ?type
			isNullable := false
			if p.current().Type() == token.TERNARY {
				isNullable = true
				p.next() // 跳过问号
			}

			// 解析一个"返回类型表达式"，支持联合类型：string|int|false
			// 其中每个原子类型可以是标识符、内置类型、null、false 等
			var unionTypes []data.Types

			parseOneTypeAtom := func() (data.Types, data.Control) {
				if !p.checkPositionIs(0,
					token.IDENTIFIER,
					token.STRING,
					token.INT,
					token.FLOAT,
					token.BOOL,
					token.ARRAY,
					token.NULL,
					token.FALSE,
					token.STATIC,
					token.SELF,
					token.PARENT,
				) {
					return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("无法识别返回类型的定义符号"+p.current().Literal()))
				}

				// 处理 self 关键字
				if p.current().Type() == token.SELF {
					p.next()
					if p.currentClassName != "" {
						return data.NewBaseType(p.currentClassName), nil
					}
					return data.NewBaseType("self"), nil
				}

				// 处理 parent 关键字：返回父类名（若有），否则字符串 "parent"
				if p.current().Type() == token.PARENT {
					p.next()
					if p.currentClassName != "" {
						// 当前类名可能是完整名，通过 VM 查找父类
						if cls, ok := p.vm.GetClass(p.currentClassName); ok && cls.GetExtend() != nil {
							return data.NewBaseType(*cls.GetExtend()), nil
						}
					}
					return data.NewBaseType("parent"), nil
				}

				name := p.current().Literal()
				p.next()

				// 如果是基础类型，直接返回
				if data.ISBaseType(name) {
					return data.NewBaseType(name), nil
				}

				// 尝试解析完整的类名（包括命名空间）
				if full, ok := p.findFullClassNameByNamespace(name); ok {
					return data.NewBaseType(full), nil
				}

				// 如果无法解析，返回原始名称
				return data.NewBaseType(name), nil
			}

			// 第一个类型原子
			firstType, acl := parseOneTypeAtom()
			if acl != nil {
				return nil, nil, acl
			}
			unionTypes = append(unionTypes, firstType)

			// 后续的 |Type 原子
			for p.current().Type() == token.BIT_OR {
				p.next() // 跳过 |
				nextType, acl := parseOneTypeAtom()
				if acl != nil {
					return nil, nil, acl
				}
				unionTypes = append(unionTypes, nextType)
			}

			// 将本次解析出的类型（可能是单一，也可能是联合）加入返回类型列表
			var thisType data.Types
			if len(unionTypes) == 1 {
				thisType = unionTypes[0]
			} else {
				// 联合类型：string|int|false 之类
				thisType = data.NewUnionType(unionTypes)
			}
			if isNullable {
				thisType = data.NewNullableType(thisType)
			}
			returnTypes = append(returnTypes, thisType)

			// 检查是否有更多类型（逗号分隔）
			if p.current().Type() == token.COMMA {
				p.next() // 跳过逗号
				continue
			}

			// 没有更多类型，结束解析
			break
		}

		// 根据返回类型数量决定返回类型
		if len(returnTypes) == 0 {
			retType = nil
		} else if len(returnTypes) == 1 {
			retType = returnTypes[0]
		} else {
			// 多个返回类型，创建多返回值类型
			retType = data.NewMultipleReturnType(returnTypes)
		}
	}

	var body []data.GetValue
	var vars []data.Variable

	// 抽象方法没有方法体，只有分号
	if isAbstract {
		if p.current().Type() != token.SEMICOLON {
			return nil, nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("抽象方法必须以分号结尾"))
		}
		p.next() // 跳过分号
		body = []data.GetValue{}
		vars = []data.Variable{}
	} else {
		body, acl = p.ParseFunctionBody()
		if acl != nil {
			return nil, nil, acl
		}
		vars = p.GetVariables()
	}

	p.scopeManager.PopScope()

	method := node.NewMethod(
		tracker.EndBefore(),
		name,
		modifier,
		isStatic,
		params,
		body,
		vars,
		retType,
	)

	// 如果是抽象方法，包装为 AbstractMethod
	var ret data.Method = method
	if isAbstract {
		if classMethod, ok := method.(*node.ClassMethod); ok {
			ret = node.NewAbstractMethod(classMethod)
		}
	}
	for _, an := range annotations {
		an.Target = ret.(data.GetValue)
	}
	if len(annotations) != 0 {
		callAnn := make([]*node.CallAnn, 0)

		for _, an := range annotations {
			stmt, acl := p.vm.GetOrLoadClass(an.Name)
			if acl != nil {
				return nil, nil, acl
			}
			object, acl := stmt.GetValue(p.vm.CreateContext(nil))
			if acl != nil {
				return nil, nil, acl
			}
			if o, ok := object.(*data.ClassValue); ok {
				if o.Class.GetConstruct() != nil {
					obj, acl := an.GetValue(p.vm.CreateContext(o.Class.GetConstruct().GetVariables()))
					if acl != nil {
						if ann, ok := acl.(*node.CallAnn); !ok {
							return nil, nil, acl
						} else {
							callAnn = append(callAnn, ann)
						}
					}
					if c, ok := ret.(node.AddAnnotations); ok {
						if o, ok := obj.(*data.ClassValue); ok {
							c.AddAnnotations(o)
						}
					}
				}
			}
		}

		for i := len(callAnn) - 1; i >= 0; i-- {
			acl := callAnn[i].InitAnnotation()
			if acl != nil {
				return nil, nil, acl
			}
		}
	}

	return ret, constructorProps, nil
}

// 解释泛型定义 class<T>、class<T, Y>、class<string, int>、class<T<int>>
func (p *ClassParser) parseGeneric() []data.Types {
	if p.current().Type() != token.LT {
		return nil
	}
	p.next() // 跳过 <

	var types []data.Types
	for {
		typ := parseType(p.Parser)
		if typ == nil {
			break
		}
		types = append(types, typ)
		if p.current().Type() == token.GT {
			p.next() // 跳过 >
			break
		}
		if p.current().Type() == token.COMMA {
			p.next() // 跳过 ,
		} else {
			break
		}
	}
	return types
}

// 解析类型（支持嵌套泛型）
// parseTraitUse 解析 use 语句，用于嵌入 trait
// 语法：use Trait1, Trait2;
func (p *ClassParser) parseTraitUse() ([]string, data.Control) {
	p.next() // 跳过 use 关键字

	var traitNames []string
	for {
		// 解析 trait 名称
		traitName, acl := p.getClassName(true)
		if acl != nil {
			return nil, acl
		}
		traitNames = append(traitNames, traitName)

		// 检查是否有逗号，继续解析下一个 trait
		if p.current().Type() == token.COMMA {
			p.next() // 跳过逗号
		} else {
			break
		}
	}

	// 跳过分号
	if p.current().Type() == token.SEMICOLON {
		p.next()
	}

	return traitNames, nil
}

// mergeTraitsIntoMaps 将 trait 的方法和属性合并到 trait 解析过程中的 maps 中（用于 trait 内的 use 语句）
func (p *ClassParser) mergeTraitsIntoMaps(traitNames []string, properties *[]data.Property, methods map[string]data.Method, staticProperties map[string]data.Property, staticMethods map[string]data.Method) data.Control {
	vm := p.vm

	for _, traitName := range traitNames {
		trait, acl := vm.GetOrLoadClass(traitName)
		if acl != nil {
			return acl
		}
		if trait == nil {
			return data.NewErrorThrow(p.newFrom(), fmt.Errorf("trait %s 不存在", traitName))
		}

		// 合并 trait 的实例方法
		for _, method := range trait.GetMethods() {
			methodName := method.GetName()
			if _, exists := methods[methodName]; !exists {
				methods[methodName] = method
			}
		}
		// 合并 trait 的静态方法（静态方法只存放在 ClassStatement.StaticMethods 中）
		if cs, ok := trait.(*node.ClassStatement); ok {
			for methodName, method := range cs.StaticMethods {
				if _, exists := staticMethods[methodName]; !exists {
					staticMethods[methodName] = method
				}
			}
		}

		// 合并 trait 的属性
		for _, property := range trait.GetPropertyList() {
			propertyName := property.GetName()
			hasProperty := false
			for _, prop := range *properties {
				if prop.GetName() == propertyName {
					hasProperty = true
					break
				}
			}
			if hasProperty || staticProperties[propertyName] != nil {
				continue
			}
			if property.GetIsStatic() {
				staticProperties[propertyName] = property
			} else {
				*properties = append(*properties, property)
			}
		}
	}
	return nil
}

// mergeTraits 合并 trait 的方法和属性到类中
func (p *ClassParser) mergeTraits(class *node.ClassStatement, traitNames []string) data.Control {
	vm := p.vm

	for _, traitName := range traitNames {
		// 从 VM 加载 trait（trait 和 class 一样存储在 classMap 中）
		trait, acl := vm.GetOrLoadClass(traitName)
		if acl != nil {
			return acl
		}

		if trait == nil {
			return data.NewErrorThrow(class.GetFrom(), fmt.Errorf("trait %s 不存在", traitName))
		}

		// 合并 trait 的实例方法
		traitMethods := trait.GetMethods()
		for _, method := range traitMethods {
			methodName := method.GetName()
			// 如果类中已经有同名方法，跳过（类的方法优先级更高）
			if _, exists := class.Methods[methodName]; !exists {
				class.Methods[methodName] = method
			}
		}
		// 合并 trait 的静态方法（只存放在 ClassStatement.StaticMethods 中）
		if cs, ok := trait.(*node.ClassStatement); ok {
			for methodName, method := range cs.StaticMethods {
				if _, exists := class.StaticMethods[methodName]; !exists {
					class.StaticMethods[methodName] = method
				}
			}
		}

		// 合并 trait 的属性
		traitProperties := trait.GetPropertyList()
		for _, property := range traitProperties {
			propertyName := property.GetName()
			// 如果类中已经有同名属性，跳过（类的属性优先级更高）
			if _, exists := class.Properties[propertyName]; !exists {
				// 检查是否是静态属性
				if property.GetIsStatic() {
					// 静态属性需要设置默认值
					defaultValue := property.GetDefaultValue()
					if defaultValue != nil {
						baseCtx := vm.CreateContext([]data.Variable{})
						v, acl := defaultValue.GetValue(baseCtx)
						if acl != nil {
							return acl
						}
						class.StaticProperty.Store(propertyName, v)
					} else {
						class.StaticProperty.Store(propertyName, data.NewNullValue())
					}
				} else {
					// 添加到属性列表
					class.Properties[propertyName] = property
					class.PropertiesIndex = append(class.PropertiesIndex, propertyName)
				}
			}
		}
	}

	return nil
}
