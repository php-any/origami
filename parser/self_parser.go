package parser

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

type SelfParser struct {
	*VariableParser
}

func NewSelfParser(parser *Parser) StatementParser {
	return &SelfParser{
		VariableParser: &VariableParser{
			parser,
		},
	}
}

func (sp *SelfParser) Parse() (data.GetValue, data.Control) {
	// 获取变量名
	tracker := sp.StartTracking()
	sp.next()

	// 检查是否是 self:: 语法（支持 IDENTIFIER、VARIABLE 或 CLASS）
	if sp.checkPositionIs(0, token.SCOPE_RESOLUTION) &&
		(sp.checkPositionIs(1, token.IDENTIFIER, token.VARIABLE, token.CLASS)) {
		sp.next() // 跳过 ::

		// 检查是否是 self::class
		if sp.current().Type() == token.CLASS {
			sp.next() // 跳过 class
			return node.NewSelfClass(tracker.EndBefore()), nil
		}

		isVariable := sp.current().Type() == token.VARIABLE
		memberName := sp.current().Literal()
		sp.next()
		tokenFrom := tracker.EndBefore()

		// 如果是 VARIABLE，去掉 $ 前缀
		if isVariable && len(memberName) > 0 && memberName[0] == '$' {
			memberName = memberName[1:]
		}

		if sp.checkPositionIs(0, token.LPAREN) {
			vp := &VariableParser{sp.Parser}
			// 如果已知当前类名（解析时），直接使用 CallStaticMethodLater
			// 这样在闭包等非 ClassMethodContext 中也能正确工作
			if sp.currentClass != "" {
				expr := node.NewCallStaticMethodLater(tokenFrom, sp.currentClass, memberName, sp.currentClass)
				return vp.parseSuffix(expr)
			}
			// 否则回退到运行时解析
			expr := node.NewCallSelfMethod(tokenFrom, memberName)
			return vp.parseSuffix(expr)
		} else {
			// 创建静态属性访问表达式
			vp := &VariableParser{sp.Parser}
			expr := node.NewCallSelfProperty(tokenFrom, memberName)
			return vp.parseSuffix(expr)
		}
	}

	// self 关键字单独使用时，应该报错（因为 self 必须配合 :: 使用）
	return nil, data.NewErrorThrow(tracker.EndBefore(), fmt.Errorf("self 关键字必须配合 :: 使用，如 self::$property 或 self::method()"))
}
