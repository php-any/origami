package parser

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ApplyAnnotations 将已解析的注解绑定到目标并延迟执行构造。
func ApplyAnnotations(p *Parser, target node.AddAnnotations, annotations []*node.Annotation) data.Control {
	if len(annotations) == 0 {
		return nil
	}
	if gv, ok := target.(data.GetValue); ok {
		for _, an := range annotations {
			an.Target = gv
		}
	}
	callAnn := make([]*node.CallAnn, 0, len(annotations))
	for _, an := range annotations {
		stmt, acl := p.vm.GetOrLoadClass(an.Name)
		if acl != nil {
			return acl
		}
		object, acl := stmt.GetValue(p.vm.CreateContext(nil))
		if acl != nil {
			return acl
		}
		cv, ok := object.(*data.ClassValue)
		if !ok {
			continue
		}
		var vars []data.Variable
		if construct := cv.Class.GetConstruct(); construct != nil {
			vars = construct.GetVariables()
		}
		obj, acl := an.GetValue(p.vm.CreateContext(vars))
		if acl != nil {
			if ann, ok := acl.(*node.CallAnn); ok {
				callAnn = append(callAnn, ann)
			} else {
				return acl
			}
		}
		if o, ok := obj.(*data.ClassValue); ok {
			target.AddAnnotations(o)
		}
	}
	for i := len(callAnn) - 1; i >= 0; i-- {
		if acl := callAnn[i].InitAnnotation(); acl != nil {
			return acl
		}
	}
	return nil
}
