package cmd

import "github.com/php-any/origami/node"

// genUnaryExpression 生成一元表达式（-, !, ~）
func (g *Generator) genUnaryExpression(n *node.UnaryExpression) {
	g.printf("node.NewUnaryExpression(from, %q,\n", n.Operator)
	g.indent++
	g.genGetValue(n.Right)
	g.printf(",\n")
	g.indent--
	g.printf(")")
}

// genUnaryIncr 生成前缀自增（++$var）
func (g *Generator) genUnaryIncr(n *node.UnaryIncr) {
	g.printf("node.NewUnaryIncr(from,\n")
	g.indent++
	g.genGetValue(n.Right)
	g.printf(",\n")
	g.indent--
	g.printf(")")
}

// genUnaryDecr 生成前缀自减（--$var）
func (g *Generator) genUnaryDecr(n *node.UnaryDecr) {
	g.printf("node.NewUnaryDecr(from,\n")
	g.indent++
	g.genGetValue(n.Right)
	g.printf(",\n")
	g.indent--
	g.printf(")")
}

// genPostfixIncr 生成后缀自增（$var++）
func (g *Generator) genPostfixIncr(n *node.PostfixIncr) {
	g.printf("node.NewPostfixIncr(from,\n")
	g.indent++
	g.genGetValue(n.Left)
	g.printf(",\n")
	g.indent--
	g.printf(")")
}

// genPostfixDecr 生成后缀自减（$var--）
func (g *Generator) genPostfixDecr(n *node.PostfixDecr) {
	g.printf("node.NewPostfixDecr(from,\n")
	g.indent++
	g.genGetValue(n.Left)
	g.printf(",\n")
	g.indent--
	g.printf(")")
}

// genErrorSuppress 生成错误抑制（@expr）
func (g *Generator) genErrorSuppress(n *node.ErrorSuppress) {
	g.printf("node.NewErrorSuppress(from,\n")
	g.indent++
	g.genGetValue(n.Inner)
	g.printf(",\n")
	g.indent--
	g.printf(")")
}
