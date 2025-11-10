package parser

import (
	"errors"
	"fmt"

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

	// 解析类名
	if !p.checkPositionIs(0, token.IDENTIFIER, token.GENERIC_TYPE) {
		return nil, data.NewErrorThrow(p.newFrom(), errors.New("new关键字后面必须跟类名"))
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
