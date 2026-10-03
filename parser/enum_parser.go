package parser

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
	"github.com/php-any/origami/token"
)

// EnumParser retains nominal enum flags and lazy case declarations.
type EnumParser struct {
	*Parser
}

func NewEnumParser(p *Parser) StatementParser {
	return &EnumParser{Parser: p}
}

func (p *EnumParser) Parse() (data.GetValue, data.Control) {
	// 跳过 enum 关键字
	p.next()
	tracker := p.StartTracking()

	// 解析枚举名
	if p.current().Type() != token.IDENTIFIER {
		return nil, data.NewErrorThrow(p.newFrom(), errors.New("enum 后缺少名称"))
	}
	enumName := p.current().Literal()
	p.next()

	// 加上命名空间前缀
	if p.namespace != nil {
		enumName = p.namespace.GetName() + "\\" + enumName
	}
	// 解析可选的底层类型: enum Status: string
	noParent := ""
	defer p.enterClassDeclaration(enumName, &noParent)()
	backingType := ""
	// PHP only permits int and string backing types.
	if p.current().Type() == token.COLON {
		p.next()
		if p.current().Type() != token.IDENTIFIER && p.current().Type() != token.STRING && p.current().Type() != token.INT {
			return nil, data.NewErrorThrow(p.newFrom(), errors.New("enum 底层类型缺失或非法"))
		}
		// Retain the type for case initialization.
		backingType = p.current().Literal()
		if backingType != "string" && backingType != "int" {
			return nil, data.NewCompileFatal(p.newFrom(), "Enum backing type must be int or string")
		}
		p.next()
	}

	// 可选 implements：enum Heroicon: string implements ScalableIcon
	implements := []string{"UnitEnum"}
	if backingType != "" {
		implements = append(implements, "BackedEnum")
	}
	if p.current().Type() == token.IMPLEMENTS {
		p.next()
		for {
			if p.current().Type() != token.IDENTIFIER && p.current().Type() != token.NAMESPACE_SEPARATOR {
				return nil, data.NewErrorThrow(p.newFrom(), errors.New("enum implements 后缺少接口名"))
			}
			ifaceName := p.current().Literal()
			p.next()
			for p.current().Type() == token.NAMESPACE_SEPARATOR {
				p.next()
				if p.current().Type() != token.IDENTIFIER {
					return nil, data.NewErrorThrow(p.newFrom(), errors.New("enum implements 接口名非法"))
				}
				ifaceName += "\\" + p.current().Literal()
				p.next()
			}
			if full, ok := p.findFullClassNameByNamespace(ifaceName); ok {
				implements = append(implements, full)
			} else {
				implements = append(implements, ifaceName)
			}
			if p.current().Type() != token.COMMA {
				break
			}
			p.next()
		}
	}

	// 解析枚举体
	if p.current().Type() != token.LBRACE {
		return nil, data.NewErrorThrow(p.newFrom(), errors.New("enum 声明后缺少 '{'"))
	}
	p.next()

	type enumCase struct {
		name  string
		value data.GetValue // 底层值表达式
	}
	var cases []enumCase
	methods := map[string]data.Method{}
	staticMethods := map[string]data.Method{}
	// enum public const 只进入 StaticProperties，不能放入实例 Properties：
	// GetValue 实例化时会求值 Properties 默认值，早于 case 静态属性注入。
	var staticConstProps []data.Property

	for !p.currentIsTypeOrEOF(token.RBRACE) {
		if p.current().Type() == token.SEMICOLON {
			p.next()
			continue
		}

		// 先尝试解析注解
		var memberAnnotations []*node.Annotation
		cp := &ClassParser{
			Parser:               p.Parser,
			FunctionParserCommon: NewFunctionParserCommon(p.Parser),
		}
		for p.checkPositionIs(0, token.AT, token.HASH) {
			ann, acl := cp.parseAnnotation()
			if acl != nil {
				return nil, acl
			}
			if ann != nil {
				memberAnnotations = append(memberAnnotations, ann)
			}
		}

		// 解析访问修饰符（方法需要）
		modifier := cp.parseModifier()

		// 解析 static 关键字
		isStatic := false
		if p.current().Type() == token.STATIC {
			isStatic = true
			p.next()
		}

		// 检查是否是 case 声明
		if p.current().Type() == token.CASE {
			p.next()

			if p.current().Type() != token.IDENTIFIER {
				return nil, data.NewErrorThrow(p.newFrom(), errors.New("enum case 缺少名称"))
			}
			caseName := p.current().Literal()
			p.next()

			var val data.GetValue

			// 支持: case OPEN = 'open';
			if p.current().Type() == token.ASSIGN {
				if backingType == "" {
					return nil, data.NewCompileFatal(p.newFrom(), "Unit enum cannot have a value")
				}
				p.next()
				exprParser := NewExpressionParser(p.Parser)
				var acl data.Control
				val, acl = exprParser.Parse()
				if acl != nil {
					return nil, acl
				}
			} else {
				if backingType != "" {
					return nil, data.NewCompileFatal(p.newFrom(), "Backed enum case must have a value")
				}
			}

			// 跳过可选分号
			if p.current().Type() == token.SEMICOLON {
				p.next()
			}

			cases = append(cases, enumCase{name: caseName, value: val})
		} else if p.current().Type() == token.CONST {
			// parse enum constant: public const VALUES = [...];
			p.next()

			if p.current().Type() != token.IDENTIFIER {
				return nil, data.NewErrorThrow(p.newFrom(), errors.New("enum constant missing name"))
			}
			constName := p.current().Literal()
			p.next()

			var defaultValue data.GetValue
			if p.current().Type() == token.ASSIGN {
				p.next()
				exprParser := NewExpressionParser(p.Parser)
				var acl data.Control
				defaultValue, acl = exprParser.Parse()
				if acl != nil {
					return nil, acl
				}
			}

			if p.current().Type() == token.SEMICOLON {
				p.next()
			}

			if defaultValue != nil {
				prop := node.NewProperty(p.newFrom(), constName, modifier, true, defaultValue)
				staticConstProps = append(staticConstProps, prop)
			}
		} else if p.current().Type() == token.FUNC {
			// 解析方法
			method, _, acl := cp.parseMethodWithAnnotations(modifier, isStatic, false, memberAnnotations, nil, nil)
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
		} else {
			return nil, data.NewErrorThrow(p.newFrom(), errors.New("enum 体内只支持 case 声明、常量声明和方法声明"))
		}
	}

	// 跳过枚举体结束右括号
	p.next()

	// properties 中包含 enum 常量（已通过 const 声明解析添加）

	// Enums implement UnitEnum and, when backed, BackedEnum.
	extends := ""
	properties := []data.Property{node.NewPropertyWithReadonly(tracker.EndBefore(), "name", "public", false, true, nil, data.String{})}
	if backingType != "" {
		properties = append(properties, node.NewPropertyWithReadonly(tracker.EndBefore(), "value", "public", false, true, nil, data.NewBaseType(backingType)))
		staticMethods["tryFrom"] = &core.BackedEnumTryFromMethod{}
		staticMethods["from"] = &core.BackedEnumFromMethod{}
	}
	staticMethods["cases"] = &core.BackedEnumCasesMethod{}
	classStmt := node.NewClassStatement(
		tracker.EndBefore(),
		enumName,
		extends,
		implements,
		properties,
		methods,
	)
	classStmt.StaticMethods = staticMethods
	classStmt.Flags = data.ClassEnum | data.ClassFinal

	// 将枚举作为类注册到 VM
	if acl := p.vm.AddClass(classStmt); acl != nil {
		return nil, acl
	}

	// enum 内 public const 延迟到首次访问时求值（对齐原生 PHP 惰性常量语义，
	// 允许前向引用）。case 静态属性仍在下方立即求值（需要实例化枚举对象）。
	staticProps := make(map[string]data.Property)
	var staticIdx []string
	for _, prop := range staticConstProps {
		cp, ok := prop.(*node.ClassProperty)
		if !ok || !cp.GetIsStatic() {
			continue
		}
		staticProps[cp.GetName()] = cp
		staticIdx = append(staticIdx, cp.GetName())
	}
	classStmt.StaticProperties = staticProps
	classStmt.StaticPropertiesIndex = staticIdx
	classStmt.SetStaticPropertyContext(data.NewClassValue(classStmt, p.vm.CreateContext([]data.Variable{})))

	// Cases are lazy constants, preserving execution-VM identity.
	for _, ccase := range cases {
		initializer := &node.EnumCaseInitializer{Node: node.NewNode(tracker.EndBefore()), ClassName: enumName, CaseName: ccase.name, BackingValue: ccase.value}
		if backingType != "" {
			initializer.BackingType = data.DeclaredTypeRef(data.NewBaseType(backingType))
		}
		prop := node.NewProperty(tracker.EndBefore(), ccase.name, "public", true, initializer)
		prop.IsConstant = true
		classStmt.EnumCases = append(classStmt.EnumCases, ccase.name)
		classStmt.StaticProperties[ccase.name] = prop
		classStmt.StaticPropertiesIndex = append(classStmt.StaticPropertiesIndex, ccase.name)
	}

	return node.NewClassRegisterStmt(tracker.EndBefore(), classStmt, nil, nil), nil
}
