package cmd

import "github.com/php-any/origami/node"

func (g *Generator) genVariableExpression(n *node.VariableExpression) {
	g.printf("node.NewVariable(from, %q, %d, nil)", n.Name, n.Index)
}

func (g *Generator) genVariableReference(n *node.VariableReference) {
	g.printf("node.NewVariableReference(from, %q, %d, nil)", n.Name, n.Index)
}
