package parser

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// AbstractClassParser 表示 abstract 关键字解析器
type AbstractClassParser struct {
	*Parser
}

// NewAbstractClassParser 创建一个新的 abstract 解析器
func NewAbstractClassParser(parser *Parser) StatementParser {
	return &AbstractClassParser{
		parser,
	}
}

// Parse 解析 abstract class 定义
func (p *AbstractClassParser) Parse() (data.GetValue, data.Control) {
	// 跳过 abstract 关键字
	p.next()

	// 确保下一个是 class 关键字
	if p.current().Type() != token.CLASS {
		return nil, data.NewErrorThrow(p.newFrom(), errors.New("abstract 关键字后必须是 class 关键字"))
	}

	// 使用 ClassParser 解析类定义
	classParser := NewClassParser(p.Parser)
	classStmt, acl := classParser.Parse()
	if acl != nil {
		return nil, acl
	}

	// 将类定义包装为抽象类
	var abstractClassStmt data.ClassStmt
	if c, ok := classStmt.(*node.ClassStatement); ok {
		abstractClassStmt = node.NewAbstractClassStatement(c)
	} else if cg, ok := classStmt.(*node.ClassGeneric); ok {
		// 如果已经是泛型类，需要特殊处理
		abstractClass := node.NewAbstractClassStatement(cg.ClassStatement)
		abstractClassStmt = &node.ClassGeneric{
			ClassStatement: abstractClass.ClassStatement,
			Generic:        cg.Generic,
		}
	} else {
		// 如果无法识别类型，直接返回
		return classStmt, nil
	}

	// ClassParser 已经注册了普通类到 VM
	// 由于 VM.AddClass 不允许重复注册，我们无法再次注册抽象类
	// 暂时直接返回抽象类，VM 中存储的仍然是普通类
	// TODO: 需要提供 VM.ReplaceClass 方法或修改 AddClass 支持替换

	return abstractClassStmt, nil
}
