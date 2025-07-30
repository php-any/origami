package parser

import (
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// AnnotationParser 表示注解解析器
type AnnotationParser struct {
	*Parser
}

// NewAnnotationParser 创建一个新的注解解析器
func NewAnnotationParser(parser *Parser) StatementParser {
	return &AnnotationParser{
		Parser: parser,
	}
}

// Parse 解析注解
func (p *AnnotationParser) Parse() (data.GetValue, data.Control) {
	var annotations []*node.Annotation

	for p.current().Type == token.AT {
		start := p.GetStart()

		// 跳过 @ 符号
		p.next()

		// 解析注解名称
		if p.current().Type != token.IDENTIFIER {
			return nil, data.NewErrorThrow(p.newFrom(), errors.New("注解缺少名称"))
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
			p.NewTokenFrom(start),
			annotationName,
			arguments,
		)

		annotations = append(annotations, annotation)
	}

	next, acl := p.parseStatement()
	if acl != nil {
		return nil, acl
	}
	for _, an := range annotations {
		an.Target = next
	}
	for _, an := range annotations {
		obj, acl := an.GetValue(p.vm.CreateContext(nil))
		if acl != nil {
			return nil, acl
		}
		if c, ok := next.(node.AddAnnotations); ok {
			if o, ok := obj.(*data.ClassValue); ok {
				c.AddAnnotations(o)
			}
		}
	}

	return next, nil
}
