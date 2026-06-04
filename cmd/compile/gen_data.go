package compile

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/net/annotation"
)

func (g *Generator) genTypes(ty data.Types) {
	if ty == nil {
		g.printf("nil")
		return
	}
	switch t := ty.(type) {
	case data.NullableType:
		g.printf("data.NewNullableType(")
		g.genTypes(t.BaseType)
		g.printf(")")
	case data.UnionType:
		g.printf("data.NewUnionType([]data.Types{")
		for i, ut := range t.Types {
			if i > 0 {
				g.printf(", ")
			}
			g.genTypes(ut)
		}
		g.printf("})")
	default:
		g.printf("data.NewBaseType(%q)", ty.String())
	}
}

func (g *Generator) genClassAnnotation(cv *data.ClassValue) {
	if cv == nil {
		g.printf("nil")
		return
	}
	g.needAnnotationImport()
	switch c := cv.Class.(type) {
	case *annotation.RouteClass:
		g.printf("annotation.CompiledRouteValue(%q)", c.Prefix())
	case *annotation.ControllerClass:
		name := ""
		if c.GetConstruct() != nil {
			// Controller 通常无 name 参数，保持空字符串
		}
		_ = name
		g.printf("annotation.CompiledControllerValue(%q)", name)
	case *annotation.GetMappingClass:
		g.printf("annotation.CompiledGetMappingValue(%q)", c.Path())
	case *annotation.PostMappingClass:
		g.printf("annotation.CompiledPostMappingValue(%q)", c.Path())
	case *annotation.PutMappingClass:
		g.printf("annotation.CompiledPutMappingValue(%q)", c.Path())
	case *annotation.DeleteMappingClass:
		g.printf("annotation.CompiledDeleteMappingValue(%q)", c.Path())
	default:
		g.printf("nil /* annotation %T */", cv.Class)
	}
}

func (g *Generator) needAnnotationImport() {
	g.imports["github.com/php-any/origami/std/net/annotation"] = true
}

func modifierName(m data.Modifier) string {
	switch m {
	case data.ModifierPrivate:
		return "private"
	case data.ModifierProtected:
		return "protected"
	default:
		return "public"
	}
}

func (g *Generator) genParentMap(parent map[int]int) {
	if len(parent) == 0 {
		g.printf("nil")
		return
	}
	g.printf("map[int]int{")
	first := true
	for cID, pID := range parent {
		if !first {
			g.printf(", ")
		}
		first = false
		g.printf("%d: %d", cID, pID)
	}
	g.printf("}")
}

func staticCallClassName(stmt data.GetValue) string {
	if stmt == nil {
		return ""
	}
	switch s := stmt.(type) {
	case data.ClassStmt:
		return s.GetName()
	case *node.VariableExpression:
		return s.Name
	case *node.StaticClass:
		return "static"
	case *node.SelfClass:
		return "self"
	case *node.Parent:
		return "parent"
	default:
		return fmt.Sprintf("%T", s)
	}
}
