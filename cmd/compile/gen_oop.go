package compile

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func (g *Generator) genClassStatement(n *node.ClassStatement) {
	g.printf("func() data.GetValue {\n")
	g.indent++
	g.genClassStatementInit(n)
	g.printf("return cls\n")
	g.indent--
	g.printf("}()")
}

func (g *Generator) genClassStatementInit(n *node.ClassStatement) {
	extends := ""
	if n.Extends != nil {
		extends = *n.Extends
	}
	g.printf("cls := node.NewClassStatement(from, %q, %q, ", n.Name, extends)
	g.genStringSlice(n.Implements)
	g.printf(", []data.Property{\n")
	g.indent++
	for _, name := range n.PropertiesIndex {
		g.genClassProperty(n.Properties[name])
		g.printf(",\n")
	}
	g.indent--
	g.printf("}, map[string]data.Method{\n")
	g.indent++
	for name, method := range n.Methods {
		g.printf("%q: ", name)
		g.genClassMethod(method)
		g.printf(",\n")
	}
	g.indent--
	g.printf("})\n")
	if n.IsAbstract {
		g.printf("cls.IsAbstract = true\n")
	}
	if len(n.StaticMethods) > 0 {
		g.printf("cls.StaticMethods = map[string]data.Method{\n")
		g.indent++
		for name, method := range n.StaticMethods {
			g.printf("%q: ", name)
			g.genClassMethod(method)
			g.printf(",\n")
		}
		g.indent--
		g.printf("}\n")
	}
	for _, ann := range n.Annotations {
		g.printf("cls.AddAnnotations(")
		g.genClassAnnotation(ann)
		g.printf(")\n")
	}
}

func (g *Generator) genAbstractClassStatement(n *node.AbstractClassStatement) {
	g.printf("node.NewAbstractClassStatement(func() *node.ClassStatement {\n")
	g.indent++
	g.genClassStatementInit(n.ClassStatement)
	g.printf("return cls\n")
	g.indent--
	g.printf("}())")
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

func (g *Generator) genClassProperty(p data.Property) {
	cp, ok := p.(*node.ClassProperty)
	if !ok {
		g.printf("nil")
		return
	}
	mod := modifierName(cp.GetModifier())
	if cp.IsPromoted {
		g.printf("node.NewPropertyWithPromoted(from, %q, %q, %v, %v, %v, ", cp.Name, mod, cp.IsStatic, cp.IsReadonly, cp.IsPromoted)
	} else if cp.IsReadonly {
		g.printf("node.NewPropertyWithReadonly(from, %q, %q, %v, %v, ", cp.Name, mod, cp.IsStatic, cp.IsReadonly)
	} else {
		g.printf("node.NewProperty(from, %q, %q, %v, ", cp.Name, mod, cp.IsStatic)
	}
	if cp.DefaultValue != nil {
		g.genGetValue(cp.DefaultValue)
	} else {
		g.printf("nil")
	}
	if cp.Type != nil {
		g.printf(", ")
		g.genTypes(cp.Type)
	}
	g.printf(")")
}

func (g *Generator) genClassMethod(method data.Method) {
	if am, ok := method.(*node.AbstractMethod); ok {
		g.printf("node.NewAbstractMethod(")
		g.genClassMethod(am.ClassMethod)
		g.printf(")")
		return
	}
	cm, ok := method.(*node.ClassMethod)
	if !ok {
		g.printf("nil")
		return
	}
	mod := modifierName(cm.GetModifier())
	if len(cm.Annotations) == 0 {
		g.printf("node.NewMethod(from, %q, %q, %v, ", cm.Name, mod, cm.IsStatic)
		g.genParamList(cm.Params)
		g.printf(", []data.GetValue{\n")
		g.indent++
		for _, stmt := range cm.Body {
			g.genGetValue(stmt)
			g.printf(",\n")
		}
		g.indent--
		g.printf("}, ")
		g.genMethodVars(cm.GetVariables())
		g.printf(", ")
		g.genTypes(cm.Ret)
		g.printf(")")
		return
	}
	g.printf("func() data.Method {\n")
	g.indent++
	g.printf("m := node.NewMethod(from, %q, %q, %v, ", cm.Name, mod, cm.IsStatic)
	g.genParamList(cm.Params)
	g.printf(", []data.GetValue{\n")
	g.indent++
	for _, stmt := range cm.Body {
		g.genGetValue(stmt)
		g.printf(",\n")
	}
	g.indent--
	g.printf("}, ")
	g.genMethodVars(cm.GetVariables())
	g.printf(", ")
	g.genTypes(cm.Ret)
	g.printf(").(*node.ClassMethod)\n")
	for _, ann := range cm.Annotations {
		g.printf("m.AddAnnotations(")
		g.genClassAnnotation(ann)
		g.printf(")\n")
	}
	g.printf("return m\n")
	g.indent--
	g.printf("}()")
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

func (g *Generator) genParamList(params []data.GetValue) {
	g.printf("[]data.GetValue{\n")
	g.indent++
	for _, p := range params {
		g.genGetValue(p)
		g.printf(",\n")
	}
	g.indent--
	g.printf("}")
}

func (g *Generator) genParameter(n *node.Parameter) {
	g.printf("node.NewParameter(from, %q, %d, ", n.Name, n.Index)
	if n.DefaultValue != nil {
		g.genGetValue(n.DefaultValue)
	} else {
		g.printf("nil")
	}
	g.printf(", ")
	g.genTypes(n.Type)
	g.printf(")")
}

func (g *Generator) genPromotedParameter(n *node.PromotedParameter) {
	g.genParameter(n.Parameter)
}

func (g *Generator) genParameters(n *node.Parameters) {
	g.genParameter(n.Parameter)
}

func (g *Generator) genParameterReference(n *node.ParameterReference) {
	g.printf("node.NewParameterReference(from, %q, %d, ", n.Name, n.Index)
	if n.DefaultValue != nil {
		g.genGetValue(n.DefaultValue)
	} else {
		g.printf("nil")
	}
	g.printf(", ")
	g.genTypes(n.Type)
	g.printf(")")
}

func (g *Generator) genThis(n *node.This) {
	g.printf("node.NewThis(from)")
}

func (g *Generator) genCallObjectProperty(n *node.CallObjectProperty) {
	g.printf("&node.CallObjectProperty{Node: node.NewNode(from), Object: ")
	g.genGetValue(n.Object)
	g.printf(", Property: %q}", n.Property)
}

func (g *Generator) genInterfaceStatement(n *node.InterfaceStatement) {
	g.printf("nil /* TODO: InterfaceStatement %q */", n.Name)
}

func (g *Generator) genNewExpression(n *node.NewExpression) {
	g.printf("node.NewNewExpression(from,\n")
	g.indent++
	g.printf("%q,\n", n.ClassName)
	g.printf("[]data.GetValue{\n")
	g.indent++
	for _, arg := range n.Arguments {
		g.genGetValue(arg)
		g.printf(",\n")
	}
	g.indent--
	g.printf("},\n")
	g.indent--
	g.printf(")")
}

func (g *Generator) genNewVariableExpression(n *node.NewVariableExpression) {
	g.printf("node.NewNewVariableExpression(from,\n")
	g.indent++
	g.genGetValue(n.ClassNameExpr)
	g.printf(",\n")
	g.printf("[]data.GetValue{\n")
	g.indent++
	for _, arg := range n.Arguments {
		g.genGetValue(arg)
		g.printf(",\n")
	}
	g.indent--
	g.printf("},\n")
	g.indent--
	g.printf(")")
}

func (g *Generator) genNewExpressionDynamic(n *node.NewExpressionDynamic) {
	g.printf("node.NewNewExpressionDynamic(from,\n")
	g.indent++
	g.genGetValue(n.ClassExpr)
	g.printf(",\n")
	g.printf("[]data.GetValue{\n")
	g.indent++
	for _, arg := range n.Arguments {
		g.genGetValue(arg)
		g.printf(",\n")
	}
	g.indent--
	g.printf("},\n")
	g.indent--
	g.printf(")")
}

func (g *Generator) genNewSelfExpression(n *node.NewSelfExpression) {
	g.printf("node.NewNewSelfExpression(from,\n")
	g.indent++
	g.printf("[]data.GetValue{\n")
	g.indent++
	for _, arg := range n.Arguments {
		g.genGetValue(arg)
		g.printf(",\n")
	}
	g.indent--
	g.printf("},\n")
	g.indent--
	g.printf(")")
}

func (g *Generator) genNewStaticExpression(n *node.NewStaticExpression) {
	g.printf("node.NewNewStaticExpression(from,\n")
	g.indent++
	g.printf("[]data.GetValue{\n")
	g.indent++
	for _, arg := range n.Arguments {
		g.genGetValue(arg)
		g.printf(",\n")
	}
	g.indent--
	g.printf("},\n")
	g.indent--
	g.printf(")")
}

func (g *Generator) genInstanceOfExpression(n *node.InstanceOfExpression) {
	g.printf("node.NewInstanceOfExpression(from,\n")
	g.indent++
	g.genGetValue(n.Object)
	g.printf(",\n")
	g.genGetValue(n.ClassName)
	g.printf(",\n")
	g.indent--
	g.printf(")")
}

func (g *Generator) genCloneExpression(n *node.CloneExpression) {
	g.printf("node.NewCloneExpression(from,\n")
	g.indent++
	g.genGetValue(n.Target)
	g.printf(",\n")
	g.indent--
	g.printf(")")
}

func (g *Generator) genInitClass(n *node.InitClass) {
	g.printf("node.NewInitClass(from,\n")
	g.indent++
	g.printf("%q,\n", n.ClassName)
	g.printf("map[string]data.GetValue{\n")
	g.indent++
	for k, v := range n.KV {
		g.printf("%q: ", k)
		g.genGetValue(v)
		g.printf(",\n")
	}
	g.indent--
	g.printf("},\n")
	g.indent--
	g.printf(")")
}

func (g *Generator) genClassConstant(n *node.ClassConstant) {
	g.printf("node.NewClassConstant(from,\n")
	g.indent++
	g.genGetValue(n.Expr)
	g.printf(",\n")
	g.indent--
	g.printf(")")
}

func (g *Generator) genStaticClass(n *node.StaticClass) {
	g.printf("node.NewStaticClass(from)")
}

func (g *Generator) genSelfClass(n *node.SelfClass) {
	g.printf("node.NewSelfClass(from)")
}

func (g *Generator) genParent(n *node.Parent) {
	g.printf("node.NewParent(from)")
}
