package compile

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func (g *Generator) genUseStatement(n *node.UseStatement) {
	g.printf("node.NewUseStatement(from, %q, %q)", n.Namespace, n.Alias)
}

func (g *Generator) genNamespace(n *node.Namespace) {
	g.printf("node.NewNamespace(from, %q, []data.GetValue{\n", n.Name)
	g.indent++
	for _, stmt := range n.Statements {
		g.genGetValue(stmt)
		g.printf(",\n")
	}
	g.indent--
	g.printf("})")
}

func (g *Generator) genBinaryAssignVariable(n *node.BinaryAssignVariable) {
	g.printf("&node.BinaryAssignVariable{Node: node.NewNode(from), Left: ")
	g.genVariable(n.Left)
	g.printf(", Right: ")
	g.genGetValue(n.Right)
	g.printf("}")
}

func (g *Generator) genBinaryAssignVariableList(n *node.BinaryAssignVariableList) {
	g.printf("&node.BinaryAssignVariableList{Node: node.NewNode(from), Left: ")
	g.genVariableList(n.Left)
	g.printf(", Right: ")
	g.genGetValue(n.Right)
	g.printf("}")
}

func (g *Generator) genVariableList(n *node.VariableList) {
	g.printf("node.NewVariableList([]*node.VariableExpression{\n")
	g.indent++
	for _, v := range n.Vars {
		g.printf("&node.VariableExpression{Node: node.NewNode(from), Name: %q, Index: %d},\n", v.Name, v.Index)
	}
	g.indent--
	g.printf("})")
}

func (g *Generator) genVariable(v data.Variable) {
	switch n := v.(type) {
	case *node.VariableExpression:
		g.genVariableExpression(n)
	case *node.VariableReference:
		g.genVariableReference(n)
	default:
		g.printf("data.NewVariable(%q, 0, nil)", v.GetName())
	}
}
