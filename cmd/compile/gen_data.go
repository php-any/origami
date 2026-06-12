package compile

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
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

func (g *Generator) needImport(path, alias string) {
	g.importAliases[path] = alias
}

func (g *Generator) needAnnotationImport() {
	g.needImport("github.com/php-any/origami/std/net/annotation", "annotation")
}

func (g *Generator) needDatabaseAnnotationImport() {
	g.needImport("github.com/php-any/origami/std/database/annotation", "dbannotation")
}

func (g *Generator) needContainerImport() {
	g.needImport("github.com/php-any/origami/std/container", "container")
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

func (g *Generator) genMethodVars(vars []data.Variable) {
	g.printf("[]data.Variable{\n")
	g.indent++
	for _, v := range vars {
		g.printf("data.NewVariable(%q, %d, nil),\n", v.GetName(), v.GetIndex())
	}
	g.indent--
	g.printf("}")
}

func (g *Generator) genStringSlice(ss []string) {
	if len(ss) == 0 {
		g.printf("nil")
		return
	}
	g.printf("[]string{")
	for i, s := range ss {
		if i > 0 {
			g.printf(", ")
		}
		g.printf("%q", s)
	}
	g.printf("}")
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
