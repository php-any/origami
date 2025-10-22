package parser

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ClassParser 表示类解析器
type ClassParser struct {
	*Parser
	*FunctionParserCommon
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
	for p.current().Type == token.AT {
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
				if p.tokens[i].Type == token.LBRACE {
					startIdx = i
					break
				}
			}
			endIdx := startIdx
			braceCount = 1
			for i := startIdx + 1; i < len(p.tokens); i++ {
				if p.tokens[i].Type == token.LBRACE {
					braceCount++
				} else if p.tokens[i].Type == token.RBRACE {
					braceCount--
					if braceCount == 0 {
						endIdx = i
						break
					}
				}
			}
			// 替换token类型
			for i := startIdx + 1; i < endIdx; i++ {
				for _, param := range genericParamNames {
					if p.tokens[i].Literal == param && p.tokens[i].Type == token.IDENTIFIER {
						p.tokens[i].Type = token.GENERIC_TYPE
					}
				}
			}
		}
	}

	// 解析继承
	var extends string
	if p.current().Type == token.EXTENDS {
		p.next()
		var acl data.Control
		extends, acl = p.getClassName(true)
		if acl != nil {
			return nil, acl
		}
		acl = p.tryLoadClass(extends)
		if acl != nil {
			return nil, acl
		}
	}

	// 解析实现的接口
	var implements []string
	if p.current().Type == token.IMPLEMENTS {
		p.next()
		for {
			interfaceName, acl := p.getClassName(true)
			if acl != nil {
				return nil, acl
			}
			acl = p.tryLoadClass(interfaceName)
			if acl != nil {
				return nil, acl
			}
			implements = append(implements, interfaceName)

			if p.current().Type != token.COMMA {
				break
			}
			p.next()
		}
	}

	// 解析类体
	if p.current().Type != token.LBRACE {
		return nil, data.NewErrorThrow(p.newFrom(), errors.New("类声明后缺少左花括号 '{'"))
	}
	p.next()

	// 解析类成员
	properties := make([]data.Property, 0)
	staticProperties := make(map[string]data.Property)
	methods := map[string]data.Method{}
	staticMethods := map[string]data.Method{}
	for !p.currentIsTypeOrEOF(token.RBRACE) {
		// 先尝试解析注解
		var memberAnnotations []*node.Annotation
		for p.current().Type == token.AT {
			ann, acl := p.parseAnnotation()
			if acl != nil {
				return nil, acl
			}
			if ann != nil {
				memberAnnotations = append(memberAnnotations, ann)
			}
		}

		// 解析访问修饰符
		modifier := p.parseModifier()
		if modifier == "" {
			return nil, data.NewErrorThrow(p.newFrom(), errors.New("缺少访问修饰符"))
		}

		// 解析static关键字
		isStatic := false
		if p.current().Type == token.STATIC {
			isStatic = true
			p.next()
		}

		// 解析属性或方法
		if p.current().Type == token.VAR ||
			p.current().Type == token.CONST ||
			p.current().Type == token.VARIABLE ||
			p.checkPositionIs(0, token.IDENTIFIER) {
			prop, acl := p.parsePropertyWithAnnotations(modifier, isStatic, memberAnnotations)
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
		} else if p.current().Type == token.FUNC {
			method, acl := p.parseMethodWithAnnotations(modifier, isStatic, memberAnnotations)
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
		} else if p.current().Type == token.SEMICOLON {
			p.next()
			continue
		} else {
			return nil, data.NewErrorThrow(p.newFrom(), errors.New("缺少属性或方法声明"))
		}
	}
	p.next() // 跳过结束花括号

	c := node.NewClassStatement(
		tracker.EndBefore(),
		className,
		extends,
		implements,
		properties,
		methods,
	)
	for s, property := range staticProperties {
		c.StaticProperty.Store(s, property.GetDefaultValue())
	}
	c.StaticMethods = staticMethods

	if c.Construct == nil {
		// 寻找父级构造函数
		vm := p.vm
		var ok = false
		var last data.ClassStmt = c
		for last != nil && last.GetExtend() != nil {
			ext := last.GetExtend()
			last, ok = vm.GetClass(*ext)
			if ok {
				if construct, ok := last.GetMethod(token.ConstructName); ok {
					c.Construct = construct
					break
				}
			} else {
				break
			}
		}
	}

	var acl data.Control

	if types != nil {
		cg := &node.ClassGeneric{
			ClassStatement: c,
			Generic:        types,
		}
		acl = p.vm.AddClass(cg)
		if acl != nil {
			return nil, acl
		}
		acl = callClassAnnotation(p.Parser, &annotations, cg)
		if acl != nil {
			return nil, acl
		}
		return cg, acl
	}

	acl = p.vm.AddClass(c)
	if acl != nil {
		return nil, acl
	}
	acl = callClassAnnotation(p.Parser, &annotations, c)
	if acl != nil {
		return nil, acl
	}
	return c, acl
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
	if p.current().Type != token.IDENTIFIER {
		return ""
	}
	name := p.current().Literal
	p.next()
	return name
}

// parseModifier 解析访问修饰符
func (p *ClassParser) parseModifier() string {
	switch p.current().Type {
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

	// 跳过 @ 符号
	p.next()

	// 解析注解名称
	if p.current().Type != token.IDENTIFIER {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("注解缺少名称"))
	}

	annotationName, acl := p.getClassName(true)
	if acl != nil {
		return nil, acl
	}

	// 解析注解参数
	arguments := make([]data.GetValue, 0)
	if p.current().Type == token.LPAREN {
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
func (p *ClassParser) parsePropertyWithAnnotations(modifier string, isStatic bool, annotations []*node.Annotation) (data.Property, data.Control) {
	tracker := p.StartTracking()
	if p.current().Type != token.VARIABLE {
		// 跳过var或const关键字
		p.next()
	}

	// 解析属性名
	if p.current().Type != token.VARIABLE {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("缺少变量名"))
	}
	name := p.current().Literal
	p.next()

	// 解析默认值
	var defaultValue data.GetValue
	var acl data.Control
	if p.current().Type == token.ASSIGN {
		p.next()
		exprParser := NewExpressionParser(p.Parser)
		defaultValue, acl = exprParser.Parse()
		if acl != nil {
			return nil, acl
		}
	}

	// 解析分号
	if p.current().Type == token.SEMICOLON {
		p.next()
	}

	ret := node.NewProperty(
		tracker.EndBefore(),
		name,
		modifier,
		isStatic,
		defaultValue,
	)
	for _, an := range annotations {
		an.Target = ret
	}
	for _, an := range annotations {
		obj, acl := an.GetValue(p.vm.CreateContext(nil))
		if acl != nil {
			return nil, acl
		}
		if o, ok := obj.(*data.ClassValue); ok {
			ret.AddAnnotations(o)
		}
	}
	return ret, acl
}

// parseMethodWithAnnotations 解析方法（带注解）
func (p *ClassParser) parseMethodWithAnnotations(modifier string, isStatic bool, annotations []*node.Annotation) (data.Method, data.Control) {
	// 跳过function关键字
	p.next()
	tracker := p.StartTracking()
	p.scopeManager.NewScope(false)
	// 解析方法名
	if p.current().Type != token.IDENTIFIER {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("缺少方法名"))
	}
	name := p.current().Literal
	p.next()

	// 使用通用函数解析器解析参数和方法体
	params, acl := p.ParseParameters()
	if acl != nil {
		return nil, acl
	}

	// 解析返回类型
	var retType data.Types
	if p.current().Type == token.COLON {
		p.next() // 跳过冒号

		// 解析返回类型列表
		var returnTypes []data.Types

		for {
			// 检查是否是可空类型语法 ?type
			isNullable := false
			if p.current().Type == token.TERNARY {
				isNullable = true
				p.next() // 跳过问号
			}

			// 解析返回类型
			if isIdentOrTypeToken(p.current().Type) {
				returnType := p.current().Literal
				p.next()

				// 创建基础类型
				baseType := data.NewBaseType(returnType)

				// 如果是可空类型，包装为基础类型的可空版本
				if isNullable {
					baseType = data.NewNullableType(baseType)
				}

				returnTypes = append(returnTypes, baseType)
			} else {
				return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("缺少返回类型"))
			}

			// 检查是否有更多类型（逗号分隔）
			if p.current().Type == token.COMMA {
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

	body, acl := p.ParseFunctionBody()
	if acl != nil {
		return nil, acl
	}
	vars := p.GetVariables()

	p.scopeManager.PopScope()

	ret := node.NewMethod(
		tracker.EndBefore(),
		name,
		modifier,
		isStatic,
		params,
		body,
		vars,
		retType,
	)
	for _, an := range annotations {
		an.Target = ret.(data.GetValue)
	}
	for _, an := range annotations {
		obj, acl := an.GetValue(p.vm.CreateContext(nil))
		if acl != nil {
			return nil, acl
		}
		if c, ok := ret.(node.AddAnnotations); ok {
			if o, ok := obj.(*data.ClassValue); ok {
				c.AddAnnotations(o)
			}
		}
	}

	return ret, nil
}

// 解释泛型定义 class<T>、class<T, Y>、class<string, int>、class<T<int>>
func (p *ClassParser) parseGeneric() []data.Types {
	if p.current().Type != token.LT {
		return nil
	}
	p.next() // 跳过 <

	var types []data.Types
	for {
		typ := p.parseType()
		if typ == nil {
			break
		}
		types = append(types, typ)
		if p.current().Type == token.GT {
			p.next() // 跳过 >
			break
		}
		if p.current().Type == token.COMMA {
			p.next() // 跳过 ,
		} else {
			break
		}
	}
	return types
}

// 解析类型（支持嵌套泛型）
func (p *ClassParser) parseType() data.Types {
	if p.current().Type != token.IDENTIFIER {
		return nil
	}

	typeName := p.current().Literal
	p.next()

	subTypes := make([]data.Types, 0)
	// 检查是否为泛型类型
	if p.current().Type == token.LT {
		p.next() // 跳过 <
		for {
			typ := p.parseType()
			if typ == nil {
				break
			}
			subTypes = append(subTypes, typ)
			if p.current().Type == token.GT {
				p.next() // 跳过 >
				break
			}
			if p.current().Type == token.COMMA {
				p.next() // 跳过 ,
			} else {
				break
			}
		}
	}

	if !data.ISBaseType(typeName) {
		if full, ok := p.findFullClassNameByNamespace(typeName); ok {
			typeName = full
			return data.NewBaseType(typeName)
		}
		return data.NewGenericType(typeName, subTypes)
	} else {
		return data.NewBaseType(typeName)
	}
}
