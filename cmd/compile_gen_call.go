package cmd

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// genCallExpression 生成函数调用表达式
// CallExpression.Fun (FuncStmt) 在编译期无法直接序列化，使用 CallTodo 延迟解析
func (g *Generator) genCallExpression(n *node.CallExpression) {
	g.printf("node.NewCallTodo(&node.CallExpression{Node: node.NewNode(from), FunName: %q, Args: []data.GetValue{\n", n.FunName)
	g.indent++
	for _, arg := range n.Args {
		g.genGetValue(arg)
		g.printf(",\n")
	}
	g.indent--
	g.printf("}}, %q)", g.namespace)
}

// genCallMethod 生成方法调用节点（通过方法引用调用）
func (g *Generator) genCallMethod(n *node.CallMethod) {
	g.printf("node.NewCallMethod(from,\n")
	g.indent++
	g.genGetValue(n.Method)
	g.printf(",\n")
	g.printf("[]data.GetValue{\n")
	g.indent++
	for _, arg := range n.Args {
		g.genGetValue(arg)
		g.printf(",\n")
	}
	g.indent--
	g.printf("},\n")
	g.indent--
	g.printf(")")
}

// genCallStaticMethod 生成静态方法调用节点（Class::method()）
// stmt 是未导出字段且编译期无法直接序列化，使用 CallStaticMethodLater 延迟解析
func (g *Generator) genCallStaticMethod(n *node.CallStaticMethod) {
	// 从 stmt 中提取类名（如果 stmt 是 ClassStmt）
	className := ""
	if cls, ok := n.GetStmt().(data.ClassStmt); ok {
		className = cls.GetName()
	}
	g.printf("node.NewCallStaticMethodLater(from, %q, %q, %q)", className, n.Method, g.namespace)
}

// genCallObjectMethod 生成对象方法调用节点（$obj->method()）
func (g *Generator) genCallObjectMethod(n *node.CallObjectMethod) {
	g.printf("node.NewObjectMethod(from,\n")
	g.indent++
	g.genGetValue(n.Object)
	g.printf(",\n")
	g.printf("%q,\n", n.Method)
	g.printf("[]data.GetValue{\n")
	g.indent++
	for _, arg := range n.Args {
		g.genGetValue(arg)
		g.printf(",\n")
	}
	g.indent--
	g.printf("},\n")
	g.indent--
	g.printf(")")
}

// genCallParentMethod 生成父类方法调用节点（parent::method()）
func (g *Generator) genCallParentMethod(n *node.CallParentMethod) {
	g.printf("node.NewCallParentMethod(from, %q, %q, []data.GetValue{\n", n.CurrentClass, n.Method)
	g.indent++
	for _, arg := range n.Arguments {
		g.genGetValue(arg)
		g.printf(",\n")
	}
	g.indent--
	g.printf("})")
}

// genCallSelfMethod 生成 self 方法调用节点（self::method()）
func (g *Generator) genCallSelfMethod(n *node.CallSelfMethod) {
	g.printf("node.NewCallSelfMethod(from, %q)", n.Method)
}

// genNullsafeCall 生成空安全调用节点（?->）
func (g *Generator) genNullsafeCall(n *node.NullsafeCall) {
	g.printf("node.NewNullsafeCall(from,\n")
	g.indent++
	g.genGetValue(n.Object)
	g.printf(",\n")
	g.genGetValue(n.CallExpr)
	g.printf(",\n")
	g.indent--
	g.printf(")")
}
