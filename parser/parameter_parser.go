package parser

import (
	"errors"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

func parseSingleParameter(parser *Parser) (data.GetValue, data.Property, data.Control) {
	tracking := parser.StartTracking()

	// 检查是否有访问修饰符（private、public、protected、readonly）
	var paramModifier string
	var isReadonly bool

	// 检查 readonly 关键字
	if parser.checkPositionIs(0, token.READONLY) {
		isReadonly = true
		parser.next()
	}

	// 检查访问修饰符
	if parser.checkPositionIs(0, token.PUBLIC, token.PRIVATE, token.PROTECTED) {
		switch parser.current().Type() {
		case token.PUBLIC:
			paramModifier = "public"
		case token.PRIVATE:
			paramModifier = "private"
		case token.PROTECTED:
			paramModifier = "protected"
		}
		parser.next()
	} else {
		// 没有访问修饰符，使用默认值（但不会创建属性）
		paramModifier = ""
	}

	// 检查是否有属性注解（如 #[\SensitiveParameter]）
	var paramAnnotations []*node.Annotation
	for parser.checkPositionIs(0, token.HASH) {
		cp := &ClassParser{Parser: parser, FunctionParserCommon: NewFunctionParserCommon(parser)}
		ann, acl := cp.parseAnnotation()
		if acl != nil {
			return nil, nil, acl
		}
		if ann != nil {
			paramAnnotations = append(paramAnnotations, ann)
		}
	}

	varType := ""
	name := ""
	isParams := false
	isReference := false // 是否引用
	// 解析参数名
	if parser.current().Type() != token.VARIABLE {
		isVar := false
		// (...args) 省却参数类型。
		if parser.current().Type() == token.ELLIPSIS {
			isVar = true
			parser.next()
			name = parser.current().Literal()
			isParams = true
			parser.next()
		}

		// &$data 或 &...$vars
		if parser.checkPositionIs(0, token.BIT_AND) {
			isVar = true
			parser.next() // 跳过 &

			// 支持 &...$vars 语法：引用可变参数
			if parser.current().Type() == token.ELLIPSIS {
				isParams = true
				parser.next() // 跳过 ...
			}

			name = parser.current().Literal()
			parser.next()
			isReference = true
		}

		// (string|int|null $data) 联合类型参数，兼容引用 & 和可变参数 ...
		if !isVar && isIdentOrTypeToken(parser.current().Type()) &&
			parser.checkPositionIs(1, token.IDENTIFIER, token.VARIABLE, token.BIT_OR, token.ELLIPSIS, token.BIT_AND) {

			// 先解析第一个类型
			varType = parserType(parser, parser.current().Literal())
			parser.next()

			// 处理后续的 |Type
			for parser.checkPositionIs(0, token.BIT_OR) {
				parser.next() // 跳过 |
				varType = varType + "|" + parserType(parser, parser.current().Literal())
				parser.next()
			}

			isVar = true

			// 引用参数: type &$var 或 type &...$vars
			if parser.checkPositionIs(0, token.BIT_AND) {
				isReference = true
				parser.next() // 跳过 &
			}

			// 可变参数: type ...$vars 或 type &...$vars
			if parser.checkPositionIs(0, token.ELLIPSIS) {
				isParams = true
				parser.next()
			}

			name = parser.current().Literal()
			parser.next()
		}
		// (?string $data) 可空类型参数
		if !isVar && parser.checkPositionIs(0, token.TERNARY) && isIdentOrTypeToken(parser.peek(1).Type()) && parser.checkPositionIs(2, token.IDENTIFIER, token.VARIABLE, token.ELLIPSIS) {
			isVar = true
			parser.next() // 跳过问号
			varType = "?" + parserType(parser, parser.current().Literal())
			parser.next()
			if parser.checkPositionIs(0, token.ELLIPSIS) {
				isParams = true
				parser.next()
			}
			name = parser.current().Literal()
			parser.next()
		}
		// (data: string)
		if !isVar && parser.checkPositionIs(0, token.IDENTIFIER, token.VARIABLE) && parser.checkPositionIs(1, token.COLON) && parser.checkPositionIs(2, token.IDENTIFIER) {
			name = parser.current().Literal()
			parser.next()
			if parser.checkPositionIs(0, token.COLON) {
				parser.next()
				varType = parserType(parser, parser.current().Literal())
				parser.next()
			}
			isVar = true
		}
		// fun(data)
		if !isVar && isIdentOrTypeToken(parser.current().Type()) && parser.checkPositionIs(1, token.RPAREN, token.COMMA) {
			name = parser.current().Literal()
			parser.next()
			isVar = true
		}
		if !isVar {
			return nil, nil, data.NewErrorThrow(tracking.EndBefore(), errors.New("参数缺少变量名"))
		}
	} else {
		name = parser.current().Literal()
		parser.next()
		if parser.checkPositionIs(0, token.COLON) {
			parser.next()
			varType = parser.current().Literal()
			parser.next()
		}
	}

	// 解析参数类型（支持联合类型：string|int|null）
	paramType := parseConstructorParameterType(parser)
	// 对于普通带类型声明但未被 parseConstructorParameterType 捕获的情况（如简单的 string $a, int $b），
	// 回退使用前面解析到的 varType 文本来构造基础类型，确保类型信息不会丢失，
	// 以便 ReflectionParameter::getType() 等功能可以拿到非空的类型。
	if paramType == nil && varType != "" {
		paramType = data.NewBaseType(varType)
	}

	// 添加参数到作用域
	val := parser.scopeManager.CurrentScope().AddVariable(name, paramType, tracking.EndBefore())

	// 解析默认值
	var defaultValue data.GetValue
	if parser.current().Type() == token.ASSIGN {
		parser.next()
		exprParser := NewExpressionParser(parser)
		var acl data.Control
		defaultValue, acl = exprParser.Parse()
		if acl != nil {
			return nil, nil, acl
		}
	}
	// 如果有访问修饰符，创建属性（属性提升）
	if paramModifier != "" {
		// 属性类型直接使用 paramType（已经支持联合类型）
		propertyType := paramType
		return node.NewPromotedParameter(tracking.EndBefore(), val.GetName(), val.GetIndex(), defaultValue, val.GetType()), node.NewPropertyWithPromoted(
			tracking.EndBefore(),
			name,
			paramModifier,
			false, // 构造函数参数不能是静态的
			isReadonly,
			true, // 标记为属性提升
			defaultValue,
			propertyType,
		), nil
	}

	// 创建参数节点
	if isParams {
		return node.NewParameters(tracking.EndBefore(), val.GetName(), val.GetIndex(), defaultValue, val.GetType()), nil, nil
	} else if isReference {
		if defaultValue != nil {
			// return nil, data.NewErrorThrow(tracking.EndBefore(), errors.New("参数为引用的变量不能有默认值"))
		}
		// 覆盖变量为引用
		parser.scopeManager.CurrentScope().SetVariable(val.GetName(), node.NewVariableReference(tracking.EndBefore(), val.GetName(), val.GetIndex(), val.GetType()))
		return node.NewParameterReference(tracking.EndBefore(), val.GetName(), val.GetIndex(), val.GetType()), nil, nil
	} else {
		return node.NewParameter(tracking.EndBefore(), val.GetName(), val.GetIndex(), defaultValue, val.GetType()), nil, nil
	}
}

func parserType(parser *Parser, name string) string {
	if data.ISBaseType(name) {
		return name
	}
	if strings.Contains(name, "\\") {
		return name[1:]
	}

	class, ok := parser.findFullClassNameByNamespace(name)
	if ok {
		return class
	}

	return name
}

// parseConstructorParameterType 解析构造函数参数类型（支持联合类型：string|int|null）
func parseConstructorParameterType(p *Parser) data.Types {
	if isIdentOrTypeToken(p.current().Type()) || p.checkPositionIs(0, token.NULL, token.FALSE) {
		// 检查是否是联合类型：string|int|null
		var unionTypes []data.Types

		// 解析第一个类型
		var firstType data.Types
		if p.checkPositionIs(0, token.NULL, token.FALSE) {
			firstType = data.NewBaseType(p.current().Literal())
			p.next()
		} else {
			firstType = parseType(p)
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
				} else if isIdentOrTypeToken(p.current().Type()) {
					nextType = parseType(p)
				} else {
					break
				}
				if nextType != nil {
					unionTypes = append(unionTypes, nextType)
				}
			}

			// 创建类型
			if len(unionTypes) == 1 {
				return unionTypes[0]
			} else {
				return data.NewUnionType(unionTypes)
			}
		}
	} else if p.checkPositionIs(0, token.TERNARY) && (isIdentOrTypeToken(p.peek(1).Type()) || p.peek(1).Type() == token.SELF) {
		// ?int 或 ?self 方式
		p.next()
		if p.current().Type() == token.SELF {
			p.next()
			var baseType data.Types
			if p.currentClass != "" {
				baseType = data.NewBaseType(p.currentClass)
			} else {
				baseType = data.NewBaseType("self")
			}
			return data.NewNullableType(baseType)
		} else {
			base := data.NewBaseType(p.current().Literal())
			p.next()
			return data.NewNullableType(base)
		}
	}
	// 没有类型声明，使用 mixed
	return data.NewBaseType("mixed")
}

func parseType(p *Parser) data.Types {
	// 支持类型关键字（bool, int, string, float, array 等）
	if !isIdentOrTypeToken(p.current().Type()) {
		if p.current().Type() == token.GENERIC_TYPE {
			T := p.current().Literal()
			p.next()
			return data.NewGenericType(T, nil)
		}
		return nil
	}

	typeName := p.current().Literal()

	// 处理 self 关键字
	if p.current().Type() == token.SELF {
		p.next()
		if p.currentClass != "" {
			// 使用当前类名作为类型
			return data.NewBaseType(p.currentClass)
		}
		// 如果没有当前类名，返回 self 作为类型名（运行时解析）
		return data.NewBaseType("self")
	}

	p.next()

	subTypes := make([]data.Types, 0)
	// 检查是否为泛型类型
	if p.current().Type() == token.LT {
		p.next() // 跳过 <
		for {
			typ := parseType(p)
			if typ == nil {
				break
			}
			subTypes = append(subTypes, typ)
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
