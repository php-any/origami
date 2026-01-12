package parser

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// TraitParser 表示 trait 解析器
// trait 和 class 结构类似，但不能实例化，只能被 use
type TraitParser struct {
	*ClassParser
}

// NewTraitParser 创建一个新的 trait 解析器
func NewTraitParser(parser *Parser) StatementParser {
	return &TraitParser{
		ClassParser: &ClassParser{
			Parser:               parser,
			FunctionParserCommon: NewFunctionParserCommon(parser),
		},
	}
}

// Parse 解析 trait 定义
func (p *TraitParser) Parse() (data.GetValue, data.Control) {
	// 解析 trait 前的注解
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

	// 跳过 trait 关键字
	p.next()
	tracker := p.StartTracking()

	// 解析 trait 名
	traitName := p.parseClassName()
	if traitName == "" {
		return nil, data.NewErrorThrow(p.newFrom(), errors.New("trait 后缺少 trait 名"))
	}
	if p.namespace != nil {
		traitName = p.namespace.GetName() + "\\" + traitName
	}

	p.vm.SetClassPathCache(traitName, *p.source)

	// trait 不支持泛型、继承和接口实现
	// 解析 trait 体
	if p.current().Type() != token.LBRACE {
		return nil, data.NewErrorThrow(p.newFrom(), errors.New("trait 声明后缺少左花括号 '{'"))
	}
	p.next()

	// 解析 trait 成员（方法和属性）
	properties := make([]data.Property, 0)
	staticProperties := make(map[string]data.Property)
	methods := map[string]data.Method{}
	staticMethods := map[string]data.Method{}
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
			method, _, acl := p.parseMethodWithAnnotations(modifier, isStatic, isAbstractMethod, memberAnnotations, nil, nil)
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

	// 创建 trait（使用 ClassStatement，因为结构相同）
	trait := node.NewClassStatement(
		tracker.EndBefore(),
		traitName,
		"",         // trait 不支持继承
		[]string{}, // trait 不支持实现接口
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
			trait.StaticProperty.Store(s, v)
		} else {
			trait.StaticProperty.Store(s, data.NewNullValue())
		}
	}
	trait.StaticMethods = staticMethods

	// trait 不支持构造函数
	trait.Construct = nil

	// 注册 trait 到 VM（trait 和 class 一样存储在 classMap 中）
	acl := p.vm.AddClass(trait)
	if acl != nil {
		return nil, acl
	}
	// 处理注解（ClassStatement 实现了 AddAnnotations 接口）
	acl = callClassAnnotation(p.Parser, &annotations, trait)
	if acl != nil {
		return nil, acl
	}
	return trait, nil
}
