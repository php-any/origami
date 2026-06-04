package compile

import "github.com/php-any/origami/node"

func (g *Generator) genIntLiteral(n *node.IntLiteral) {
	if n.V != nil {
		s := n.V.AsString()
		g.printf("&node.IntLiteral{Node: node.NewNode(from), V: data.NewIntValue(%s)}", s)
	} else {
		g.printf("&node.IntLiteral{Node: node.NewNode(from), V: data.NewIntValue(0)}")
	}
}

func (g *Generator) genFloatLiteral(n *node.FloatLiteral) {
	if n.V != nil {
		s := n.V.AsString()
		g.printf("&node.FloatLiteral{Node: node.NewNode(from), V: data.NewFloatValue(%s)}", s)
	} else {
		g.printf("&node.FloatLiteral{Node: node.NewNode(from), V: data.NewFloatValue(0)}")
	}
}

func (g *Generator) genStringLiteral(n *node.StringLiteral) {
	g.printf("&node.StringLiteral{Node: node.NewNode(from), Value: %q}", n.Value)
}

func (g *Generator) genBooleanLiteral(n *node.BooleanLiteral) {
	g.printf("&node.BooleanLiteral{Node: node.NewNode(from), Value: %v}", n.Value)
}

func (g *Generator) genNullLiteral(n *node.NullLiteral) {
	g.printf("&node.NullLiteral{Node: node.NewNode(from)}")
}
