package node

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/token"
)

// ClassRegisterStmt 在程序执行时完成类的延迟初始化：
// 解析父类构造函数、应用注解。类本身在 ClassParser 解析期已 AddClass。
type ClassRegisterStmt struct {
	*Node
	Class       *ClassStatement
	Annotations []*Annotation
	Generic     []data.Types
}

func NewClassRegisterStmt(from data.From, class *ClassStatement, annotations []*Annotation, generic []data.Types) *ClassRegisterStmt {
	return &ClassRegisterStmt{
		Node:        NewNode(from),
		Class:       class,
		Annotations: annotations,
		Generic:     generic,
	}
}

func (s *ClassRegisterStmt) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	vm := ctx.GetVM()
	c := s.Class

	if acl := resolveParentConstruct(vm, c); acl != nil {
		return nil, acl
	}

	classStmt := s.classStmtForVM()
	existing, ok := vm.GetClass(c.GetName())
	if !ok {
		if acl := vm.AddClass(classStmt); acl != nil {
			return nil, acl
		}
		existing = classStmt
	} else {
		// 同文件再次 require 时 VM 仍持有首次注册的 ClassStatement；注解目标必须用它
		classStmt = existing
	}

	// 注解只应用一次（InitAnnotation 有全局副作用：路由/命令注册）
	if registered := classStmtFromAny(classStmt); registered != nil && registered.AnnotationsApplied {
		return data.NewNullValue(), nil
	}

	if addAnn, ok := classStmt.(AddAnnotations); ok {
		if acl := applyClassAnnotations(vm, s.Annotations, addAnn); acl != nil {
			return nil, acl
		}
	} else if addAnn, ok := any(c).(AddAnnotations); ok {
		if acl := applyClassAnnotations(vm, s.Annotations, addAnn); acl != nil {
			return nil, acl
		}
	}
	if registered := classStmtFromAny(classStmt); registered != nil {
		registered.AnnotationsApplied = true
	}
	return data.NewNullValue(), nil
}

func (s *ClassRegisterStmt) classStmtForVM() data.ClassStmt {
	// VM 中存 *ClassStatement，用 IsAbstract 标记；避免 *AbstractClassStatement 导致 static::$prop 类型断言失败
	if len(s.Generic) > 0 {
		return &ClassGeneric{
			ClassStatement: s.Class,
			Generic:        s.Generic,
		}
	}
	return s.Class
}

func classStmtFromAny(v data.ClassStmt) *ClassStatement {
	switch t := v.(type) {
	case *ClassStatement:
		return t
	case *ClassGeneric:
		return t.ClassStatement
	case *AbstractClassStatement:
		return t.ClassStatement
	default:
		return nil
	}
}

// AddAnnotations 转发到内部 ClassStatement（AnnotationParser 以 ClassRegisterStmt 为 next）
func (s *ClassRegisterStmt) AddAnnotations(a *data.ClassValue) {
	if s.Class != nil {
		s.Class.AddAnnotations(a)
	}
}

func resolveParentConstruct(vm data.VM, c *ClassStatement) data.Control {
	if c.Construct != nil || c.Extends == nil {
		return nil
	}
	var last data.ClassStmt = c
	for last != nil && last.GetExtend() != nil {
		ext := *last.GetExtend()
		var acl data.Control
		last, acl = vm.GetOrLoadClass(ext)
		if acl != nil {
			return acl
		}
		if construct, ok := last.GetMethod(token.ConstructName); ok {
			c.Construct = construct
			break
		}
	}
	return nil
}

func applyClassAnnotations(vm data.VM, annotations []*Annotation, target AddAnnotations) data.Control {
	if len(annotations) == 0 {
		return nil
	}
	if gv, ok := target.(data.GetValue); ok {
		for _, an := range annotations {
			an.Target = gv
		}
	}
	callAnn := make([]*CallAnn, 0, len(annotations))
	for _, an := range annotations {
		stmt, acl := vm.GetOrLoadClass(an.Name)
		if acl != nil {
			if data.CompileMode {
				continue
			}
			return acl
		}
		object, acl := stmt.GetValue(vm.CreateContext(nil))
		if acl != nil {
			return acl
		}
		cv, ok := object.(*data.ClassValue)
		if !ok || cv.Class.GetConstruct() == nil {
			continue
		}
		obj, acl := an.GetValue(vm.CreateContext(cv.Class.GetConstruct().GetVariables()))
		if acl != nil {
			if ann, ok := acl.(*CallAnn); ok {
				callAnn = append(callAnn, ann)
			} else if data.CompileMode {
				continue
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
