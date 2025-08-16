package parser

import (
	"errors"
	"fmt"

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
	name := p.current().Literal
	startToken := p.current()
	p.next()

	// 函数调用模式 div {} 或者 div []
	if p.checkPositionIs(0, token.LBRACE) {
		if full, ok := p.findFullFunNameByNamespace(name); ok {
			fn, ok := p.vm.GetFunc(full)
			if !ok {
				return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("未定义的函数:"+full))
			}
			v, acl := NewLbraceParser(p.Parser).Parse()
			return node.NewCallExpression(tracker.EndBefore(), fn.GetName(), []data.GetValue{v}, fn), acl
		}
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("未定义的函数:"+name))
	} else if p.checkPositionIs(0, token.LBRACKET) {
		if full, ok := p.findFullFunNameByNamespace(name); ok {
			fn, ok := p.vm.GetFunc(full)
			if !ok {
				return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("未定义的函数:"+full))
			}
			v, acl := NewLbracketParser(p.Parser).Parse()
			return node.NewCallExpression(tracker.EndBefore(), fn.GetName(), []data.GetValue{v}, fn), acl
		}
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("未定义的函数:"+name))
	}

	// 检查是否是变量的类型
	if p.checkPositionIs(0, token.ASSIGN) {
		index := p.scopeManager.CurrentScope().AddVariable(name, nil, tracker.EndBefore())
		return node.NewVariable(tracker.EndBefore(), name, index, nil), nil
	}
	if p.checkPositionIs(0, token.VARIABLE) || p.checkPositionIs(1, token.ASSIGN) {
		// int $num 或者 int i = 0
		ty := name
		name = p.current().Literal
		p.next()
		index := p.scopeManager.CurrentScope().AddVariable(name, data.NewBaseType(ty), tracker.EndBefore())
		return node.NewVariable(tracker.EndBefore(), name, index, data.NewBaseType(ty)), nil
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
				fn, ok := p.vm.GetFunc(full)
				if !ok {
					return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("未定义的函数:"+name))
				}
				return node.NewCallExpression(tracker.EndBefore(), full, stmt, fn), acl
			} else if InLSP {
				stmt, acl := vp.parseFunctionCall()
				return node.NewCallExpression(tracker.EndBefore(), full, stmt, nil), acl
			} else {
				return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("未定义的函数:"+name))
			}
		}
		// 变量定义
		if p.checkPositionIs(0, token.COLON) && p.checkPositionIs(1, token.IDENTIFIER) {
			// a: string
			p.next()
			ty := p.current().Literal
			p.next()
			index := p.scopeManager.CurrentScope().AddVariable(name, data.NewBaseType(ty), tracker.EndBefore())
			expr := node.NewVariable(tracker.EndBefore(), name, index, data.NewBaseType(ty))
			// 解析后续操作（函数调用、数组访问等）
			vp := &VariableParser{p.Parser}
			return vp.parseSuffix(expr)
		}

		// 函数静态调用 Log::info
		if p.checkPositionIs(0, token.SCOPE_RESOLUTION) && p.checkPositionIs(1, token.IDENTIFIER) {
			className := name
			if full, ok := p.findFullClassNameByNamespace(className); ok {
				className = full
			}
			p.next()
			fnName := p.current().Literal
			p.next()

			if p.checkPositionIs(0, token.LPAREN) {
				// 创建函数调用表达式
				vp := &VariableParser{p.Parser}
				expr := node.NewCallStaticMethod(tracker.EndBefore(), className, fnName)
				return vp.parseSuffix(expr)
			} else {
				vp := &VariableParser{p.Parser}
				expr := node.NewCallStaticProperty(tracker.EndBefore(), className, fnName)
				return vp.parseSuffix(expr)
			}
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
			index := p.scopeManager.CurrentScope().AddVariable(name, nil, tracker.EndBefore())
			expr := node.NewVariable(tracker.EndBefore(), name, index, nil)
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
			generaName, ok := p.findFullClassNameByNamespace(p.current().Literal)
			if !ok {
				return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("class %s 不存在", name))
			}
			p.next()
			generaList = append(generaList, generaName)
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

	if p.scopeManager.current.isLambda {
		// 检查是否是变量
		varInfo := p.scopeManager.LookupParentVariable(name)
		if varInfo != nil {
			index := p.scopeManager.CurrentScope().AddVariable(name, varInfo.GetType(), tracker.EndBefore())
			expr := node.NewVariable(tracker.EndBefore(), name, index, varInfo.GetType())
			vp := &VariableParser{p.Parser}
			return vp.parseSuffix(expr)
		}
	}

	return node.NewStringLiteral(tracker.EndBefore(), name), nil
}
