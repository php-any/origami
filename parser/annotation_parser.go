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

	for p.current().Type() == token.AT {
		tracker := p.StartTracking()

		// 跳过 @ 符号
		p.next()

		// 解析注解名称
		if p.current().Type() != token.IDENTIFIER {
			return nil, data.NewErrorThrow(p.FromCurrentToken(), errors.New("注解缺少名称"))
		}

		annotationName, acl := p.getClassName(true)
		if acl != nil {
			return nil, acl
		}

		// 解析注解参数
		arguments := make([]data.GetValue, 0)
		if p.current().Type() == token.LPAREN {
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
			tracker.EndBefore(),
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

	// 注解的构造处理是需要延后执行的
	if len(annotations) != 0 {
		callAnn := make([]*node.CallAnn, 0)

		for _, an := range annotations {
			stmt, acl := p.vm.GetOrLoadClass(an.Name)
			if acl != nil {
				return nil, acl
			}
			object, acl := stmt.GetValue(p.vm.CreateContext(nil))
			if acl != nil {
				return nil, acl
			}
			if o, ok := object.(*data.ClassValue); ok {
				if o.Class.GetConstruct() != nil {
					obj, acl := an.GetValue(p.vm.CreateContext(o.Class.GetConstruct().GetVariables()))
					if acl != nil {
						if ann, ok := acl.(*node.CallAnn); !ok {
							return nil, acl
						} else {
							callAnn = append(callAnn, ann)
						}
					}
					if c, ok := next.(node.AddAnnotations); ok {
						if o, ok := obj.(*data.ClassValue); ok {
							c.AddAnnotations(o)
						}
					}
				}
			}
		}

		for i := len(callAnn) - 1; i >= 0; i-- {
			acl := callAnn[i].InitAnnotation()
			if acl != nil {
				return nil, acl
			}
		}
	}

	return next, nil
}
