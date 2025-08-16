package parser

import (
	"errors"
	"github.com/php-any/origami/data"
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
	interfaceName := p.parseInterfaceName()
	if interfaceName == "" {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("缺少接口名"))
	}
	interfaceName = p.namespace.GetName() + "\\" + interfaceName

	// 解析继承
	var extends *string
	if p.current().Type == token.EXTENDS {
		p.next()
		extendsName, acl := p.getClassName(true)
		if acl != nil {
			return nil, acl
		}
		if extendsName == "" {
			return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("缺少继承的接口名"))
		}
		extends = &extendsName
	}

	// 解析接口体
	if p.current().Type != token.LBRACE {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("接口声明后缺少左花括号 '{'"))
	}
	p.next()

	// 解析接口方法
	var methods []data.Method
	for !p.isEOF() && p.current().Type != token.RBRACE {
		// 解析方法修饰符
		modifier := p.parseModifier()
		if modifier == "" {
			return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("缺少方法修饰符"))
		}

		// 解析方法
		method, acl := p.parseInterfaceMethod(modifier)
		if acl != nil {
			return nil, acl
		}
		if method != nil {
			methods = append(methods, method)
		}

		// 跳过分号
		if p.current().Type == token.SEMICOLON {
			p.next()
		}
	}

	// 解析右花括号
	if p.current().Type != token.RBRACE {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("接口定义后缺少右花括号 '}'"))
	}
	p.next()

	i := node.NewInterfaceStatement(
		tracker.EndBefore(),
		interfaceName,
		extends,
		methods,
	)

	if i.Extends != nil {
		extends, ok := p.vm.GetInterface(*i.Extends)
		if ok {
			temp := extends.GetName()
			i.Extends = &temp
		}
	}

	// 接口定义本身不返回值，但需要注册到虚拟机中
	// 创建接口定义
	return i, p.vm.AddInterface(i)
}

// parseInterfaceName 解析接口名
func (p *InterfaceParser) parseInterfaceName() string {
	if p.current().Type != token.IDENTIFIER {
		return ""
	}
	name := p.current().Literal
	p.next()
	return name
}

// parseModifier 解析访问修饰符
func (p *InterfaceParser) parseModifier() string {
	switch p.current().Type {
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
	if p.current().Type == token.FUNC {
		p.next()
	}

	// 解析方法名
	if p.current().Type != token.IDENTIFIER {
		return nil, data.NewErrorThrow(p.FromCurrentToken(), errors.New("缺少方法名"))
	}
	name := p.current().Literal
	p.next()

	// 解析参数列表
	params, acl := p.ParseParameters()
	if acl != nil {
		return nil, acl
	}
	// 解析返回类型（可选）
	var returnType data.Types
	if p.current().Type == token.COLON {
		p.next()
		if p.current().Type == token.IDENTIFIER {
			returnType = data.NewBaseType(p.current().Literal)
			p.next()
		}
	}

	// 解析分号
	p.nextAndCheck(token.SEMICOLON)
	// 创建接口方法（接口方法没有方法体）
	return node.NewInterfaceMethod(
		tracker.EndBefore(),
		name,
		modifier,
		params,
		returnType,
	), nil
}
