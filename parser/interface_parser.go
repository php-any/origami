package parser

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/lexer"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// InterfaceParser 表示接口解析器
type InterfaceParser struct {
	parser *Parser
	*FunctionParserCommon
}

// NewInterfaceParser 创建一个新的接口解析器
func NewInterfaceParser(parser *Parser) StatementParser {
	return &InterfaceParser{
		parser:               parser,
		FunctionParserCommon: NewFunctionParserCommon(parser),
	}
}

// Parse 解析接口定义
func (p *InterfaceParser) Parse() (data.GetValue, data.Control) {
	tracker := p.StartTracking()
	// 跳过interface关键字
	p.next()

	// 解析接口名
	shortName := p.parseInterfaceName()
	if shortName == "" {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("缺少接口名"))
	}

	interfaceName := shortName

	if p.namespace != nil {
		interfaceName = p.namespace.GetName() + "\\" + interfaceName
	}

	// 解析继承
	var extends []string
	if p.current().Type() == token.EXTENDS {
		p.next()

		for {
			extendsName, acl := p.getClassName(true)
			if acl != nil {
				return nil, acl
			}
			if extendsName == "" {
				return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("缺少继承的接口名"))
			}
			extends = append(extends, extendsName)

			// 支持多个继承接口：interface A extends B, C, D
			if p.current().Type() == token.COMMA {
				p.next()
				continue
			}
			break
		}
	}

	// 在接口体内，将 self::X 重写为 ShortInterfaceName::X，避免在接口常量表达式中依赖 self:: 的运行时语义。
	// 只修改当前接口 body 范围内的 token，不影响其它位置。
	{
		startIdx := p.parser.position
		// 找到当前接口的左花括号
		for startIdx < len(p.parser.tokens) && p.parser.tokens[startIdx].Type() != token.LBRACE {
			startIdx++
		}
		if startIdx < len(p.parser.tokens) && p.parser.tokens[startIdx].Type() == token.LBRACE {
			braceCount := 1
			endIdx := startIdx
			for i := startIdx + 1; i < len(p.parser.tokens); i++ {
				switch p.parser.tokens[i].Type() {
				case token.LBRACE:
					braceCount++
				case token.RBRACE:
					braceCount--
					if braceCount == 0 {
						endIdx = i
						i = len(p.parser.tokens)
					}
				}
			}

			for i := startIdx + 1; i < endIdx; i++ {
				if p.parser.tokens[i].Type() == token.SELF {
					// 仅处理 self::X 形式
					if i+1 < len(p.parser.tokens) && p.parser.tokens[i+1].Type() == token.SCOPE_RESOLUTION {
						tok := p.parser.tokens[i]
						p.parser.tokens[i] = lexer.NewWorkerToken(
							token.IDENTIFIER,
							shortName,
							tok.Start(),
							tok.End(),
							tok.Line(),
							tok.Pos(),
						)
					}
				}
			}
		}
	}

	// 解析接口体
	if p.current().Type() != token.LBRACE {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("接口声明后缺少左花括号 '{'"))
	}
	p.next()

	// 解析接口成员（方法和静态属性）
	var methods []data.Method
	staticProperties := make(map[string]data.Property)
	staticPropertyOrder := make([]string, 0)

	for !p.isEOF() && p.current().Type() != token.RBRACE {
		tracker := p.StartTracking()

		// 跳过 PHP 8 属性 #[...]（目前接口成员上的属性仅用于元信息，这里先忽略其语义）
		for p.checkPositionIs(0, token.HASH) && p.checkPositionIs(1, token.LBRACKET) {
			p.next() // 跳过 #
			if p.current().Type() != token.LBRACKET {
				return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("属性注解格式错误，期望 #[...]"))
			}
			p.next() // 跳过 [

			depth := 1
			for depth > 0 && !p.isEOF() {
				if p.current().Type() == token.LBRACKET {
					depth++
				} else if p.current().Type() == token.RBRACKET {
					depth--
				}
				p.next()
			}
		}

		// 解析访问修饰符（接口成员默认 public）
		modifier := p.parseModifier()

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
			// 解析属性（使用 ClassParser 的 parsePropertyWithAnnotations 方法）
			cp := &ClassParser{
				Parser:               p.parser,
				FunctionParserCommon: NewFunctionParserCommon(p.parser),
			}
			prop, acl := cp.parsePropertyWithAnnotations(modifier, isStatic, false, nil)
			if acl != nil {
				return nil, acl
			}
			if prop != nil {
				if isStatic || prop.GetIsStatic() {
					name := prop.GetName()
					staticProperties[name] = prop
					staticPropertyOrder = append(staticPropertyOrder, name)
				}
			}
		} else if p.current().Type() == token.FUNC {
			// 解析方法
			method, acl := p.parseInterfaceMethod(modifier)
			if acl != nil {
				return nil, acl
			}
			if method != nil {
				methods = append(methods, method)
			}
		} else if p.current().Type() == token.SEMICOLON {
			p.next()
			continue
		} else {
			return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("缺少属性或方法声明"))
		}

		// 跳过分号
		if p.current().Type() == token.SEMICOLON {
			p.next()
		}
	}

	// 解析右花括号
	if p.current().Type() != token.RBRACE {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("接口定义后缺少右花括号 '}'"))
	}
	p.next()

	i := node.NewInterfaceStatement(
		tracker.EndBefore(),
		interfaceName,
		extends,
		methods,
	)

	// 先把接口本身注册到 VM，避免在后续静态属性求值或类型检查中通过 VM 再次触发加载，造成自身依赖死循环
	if acl := p.vm.AddInterface(i); acl != nil {
		return nil, acl
	}

	// 处理静态属性（按解析顺序求值，确保依赖已有常量的表达式能正确工作）
	for _, s := range staticPropertyOrder {
		property := staticProperties[s]
		defaultValue := property.GetDefaultValue()
		if defaultValue != nil {
			baseCtx := p.vm.CreateContext([]data.Variable{})
			v, acl := defaultValue.GetValue(baseCtx)
			if acl != nil {
				return nil, acl
			}
			i.StaticProperty.Store(s, v)
		} else {
			i.StaticProperty.Store(s, data.NewNullValue())
		}
	}

	// 规范化父接口名（如果已加载则替换为其完整名称）
	if len(i.Extends) > 0 {
		for idx, name := range i.Extends {
			if ext, ok := p.vm.GetInterface(name); ok {
				i.Extends[idx] = ext.GetName()
			}
		}
	}

	// 接口定义本身不返回值，这里仅返回语法树节点；注册已在前面完成
	return i, nil
}

// parseInterfaceName 解析接口名
func (p *InterfaceParser) parseInterfaceName() string {
	if p.current().Type() != token.IDENTIFIER {
		return ""
	}
	name := p.current().Literal()
	p.next()
	return name
}

// parseModifier 解析访问修饰符
func (p *InterfaceParser) parseModifier() string {
	switch p.current().Type() {
	case token.PUBLIC:
		p.next()
		return "public"
	case token.PRIVATE:
		p.next()
		return "private"
	case token.PROTECTED:
		p.next()
		return "protected"
	default:
		return ""
	}
}

// parseInterfaceMethod 解析接口方法
func (p *InterfaceParser) parseInterfaceMethod(modifier string) (data.Method, data.Control) {
	tracker := p.StartTracking()

	// 跳过function关键字
	if p.current().Type() == token.FUNC {
		p.next()
	}

	// 解析方法名
	if !(p.checkPositionIs(0, token.IDENTIFIER, token.SELF) ||
		(p.current().Type() > token.KEYWORD_START && p.current().Type() < token.VALUE_START)) {
		return nil, data.NewErrorThrow(p.FromCurrentToken(), errors.New("缺少方法名"))
	}
	name := p.current().Literal()
	p.next()

	// 解析参数列表
	params, acl := p.ParseParameters()
	if acl != nil {
		return nil, acl
	}
	// 解析返回类型（可选，支持可空与联合类型：?Type、A|B）
	var returnType data.Types
	if p.current().Type() == token.COLON {
		p.next() // 跳过冒号

		// 解析返回类型列表（与 ClassParser 中的方法返回类型逻辑保持一致）
		var returnTypes []data.Types

		for {
			// 检查是否是可空类型语法 ?type
			isNullable := false
			if p.current().Type() == token.TERNARY {
				isNullable = true
				p.next() // 跳过问号
			}

			// 解析一个“返回类型原子”，支持联合类型：string|int|false
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
				) {
					return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("无法识别返回类型的定义符号"))
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
				return nil, acl
			}
			unionTypes = append(unionTypes, firstType)

			// 后续的 |Type 原子
			for p.current().Type() == token.BIT_OR {
				p.next() // 跳过 |
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
				thisType = data.NewUnionType(unionTypes)
			}

			if isNullable {
				thisType = data.NewNullableType(thisType)
			}

			returnTypes = append(returnTypes, thisType)

			// 支持多返回值列表时可以继续解析；接口方法当前只需要一个返回类型，遇到逗号/其他符号时停止
			break
		}

		if len(returnTypes) == 1 {
			returnType = returnTypes[0]
		} else if len(returnTypes) > 1 {
			returnType = data.NewMultipleReturnType(returnTypes)
		}
	}

	// 解析分号
	if p.checkPositionIs(0, token.SEMICOLON) {
		p.next()
	}
	// 创建接口方法（接口方法没有方法体）
	return node.NewInterfaceMethod(
		tracker.EndBefore(),
		name,
		modifier,
		params,
		returnType,
	), nil
}
