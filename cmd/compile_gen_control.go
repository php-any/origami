package cmd

import "github.com/php-any/origami/node"

func (g *Generator) genIfStatement(n *node.IfStatement) {
	g.printf("node.NewIfStatement(from,\n")
	g.indent++
	g.genGetValue(n.Condition)
	g.printf(",\n")
	g.printf("[]data.GetValue{\n")
	g.indent++
	for _, stmt := range n.ThenBranch {
		g.genGetValue(stmt)
		g.printf(",\n")
	}
	g.indent--
	g.printf("},\n")
	g.printf("nil,\n") // ElseIf
	g.printf("nil,\n") // ElseBranch
	g.indent--
	g.printf(")")
}

func (g *Generator) genReturnStatement(n *node.ReturnStatement) {
	g.printf("&node.ReturnStatement{Node: node.NewNode(from), Value: ")
	if n.Value != nil {
		g.genGetValue(n.Value)
	} else {
		g.printf("nil")
	}
	g.printf("}")
}

func (g *Generator) genEchoStatement(n *node.EchoStatement) {
	g.printf("node.NewEchoStatement(from, []data.GetValue{\n")
	g.indent++
	for _, expr := range n.Expressions {
		g.genGetValue(expr)
		g.printf(",\n")
	}
	g.indent--
	g.printf("})")
}
