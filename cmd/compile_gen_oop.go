package cmd

import "github.com/php-any/origami/node"

// genClassStatement 生成类定义语句
// ClassStatement 字段非常多（Name, Extends, Implements, Properties, Methods, StaticMethods, Construct, IsAbstract 等），
// 且 Properties/Methods 包含不可序列化的运行时类型（sync.Map, data.Property, data.Method），无法在编译期直接生成。
func (g *Generator) genClassStatement(n *node.ClassStatement) {
	g.printf("nil /* TODO: ClassStatement %q — 太复杂，暂不支持编译期代码生成 */", n.Name)
}

// genInterfaceStatement 生成接口定义语句
// InterfaceStatement 的 Methods 列表包含不可序列化的 data.Method 接口。
func (g *Generator) genInterfaceStatement(n *node.InterfaceStatement) {
	g.printf("nil /* TODO: InterfaceStatement %q — 太复杂，暂不支持编译期代码生成 */", n.Name)
}

// genNewExpression 生成 new 表达式（静态类名，如 new Foo()）
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

// genNewVariableExpression 生成 new $var() 表达式（变量类名）
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

// genNewExpressionDynamic 生成 new $expr() 表达式（动态类名表达式）
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

// genNewSelfExpression 生成 new self() 表达式
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

// genNewStaticExpression 生成 new static() 表达式
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

// genInstanceOfExpression 生成 instanceof 表达式
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

// genCloneExpression 生成 clone 表达式
func (g *Generator) genCloneExpression(n *node.CloneExpression) {
	g.printf("node.NewCloneExpression(from,\n")
	g.indent++
	g.genGetValue(n.Target)
	g.printf(",\n")
	g.indent--
	g.printf(")")
}

// genInitClass 生成 ClassName {} 初始化表达式
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

// genClassConstant 生成 ::class 语法节点
func (g *Generator) genClassConstant(n *node.ClassConstant) {
	g.printf("node.NewClassConstant(from,\n")
	g.indent++
	g.genGetValue(n.Expr)
	g.printf(",\n")
	g.indent--
	g.printf(")")
}

// genStaticClass 生成 static::class 表达式
func (g *Generator) genStaticClass(n *node.StaticClass) {
	g.printf("node.NewStaticClass(from)")
}

// genSelfClass 生成 self::class 表达式
func (g *Generator) genSelfClass(n *node.SelfClass) {
	g.printf("node.NewSelfClass(from)")
}

// genParent 生成 parent 表达式
func (g *Generator) genParent(n *node.Parent) {
	g.printf("node.NewParent(from)")
}

// genClassProperty 生成类属性节点
func (g *Generator) genClassProperty(n *node.ClassProperty) {
	g.printf("nil /* TODO: ClassProperty %q — 需要运行时类型 data.Property */", n.Name)
}

// genClassMethod 生成类方法节点
func (g *Generator) genClassMethod(n *node.ClassMethod) {
	g.printf("nil /* TODO: ClassMethod %q — 需要运行时类型 data.Method */", n.Name)
}
