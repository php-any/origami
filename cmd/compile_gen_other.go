package cmd

import "github.com/php-any/origami/node"

// genArray 生成数组字面量
func (g *Generator) genArray(n *node.Array) {
	if len(n.Keys) > 0 {
		// 有键值对：使用 NewArrayWithKeys
		g.printf("node.NewArrayWithKeys(from,\n")
		g.indent++
		g.printf("[]data.GetValue{\n")
		g.indent++
		for _, elem := range n.V {
			g.genGetValue(elem)
			g.printf(",\n")
		}
		g.indent--
		g.printf("},\n")
		g.printf("[]node.KvPair{\n")
		g.indent++
		for _, kv := range n.Keys {
			g.printf("{\n")
			g.indent++
			g.printf("Key: ")
			g.genGetValue(kv.Key)
			g.printf(",\n")
			g.printf("Value: ")
			g.genGetValue(kv.Value)
			g.printf(",\n")
			g.indent--
			g.printf("},\n")
		}
		g.indent--
		g.printf("},\n")
		g.indent--
		g.printf(")")
	} else {
		// 无键值对：使用 NewArray
		g.printf("node.NewArray(from, []data.GetValue{\n")
		g.indent++
		for _, elem := range n.V {
			g.genGetValue(elem)
			g.printf(",\n")
		}
		g.indent--
		g.printf("})")
	}
}

// genIndexExpression 生成数组/对象索引访问表达式
func (g *Generator) genIndexExpression(n *node.IndexExpression) {
	g.printf("node.NewIndexExpression(from,\n")
	g.indent++
	g.genGetValue(n.Array)
	g.printf(",\n")
	if n.Index != nil {
		g.genGetValue(n.Index)
	} else {
		g.printf("nil")
	}
	g.printf(",\n")
	g.indent--
	g.printf(")")
}

// genTernaryExpression 生成三目运算符表达式
func (g *Generator) genTernaryExpression(n *node.TernaryExpression) {
	g.printf("node.NewTernaryExpression(from,\n")
	g.indent++
	g.genGetValue(n.Condition)
	g.printf(",\n")
	g.genGetValue(n.TrueValue)
	g.printf(",\n")
	g.genGetValue(n.FalseValue)
	g.printf(",\n")
	g.indent--
	g.printf(")")
}

// genNullCoalesceExpression 生成空合并运算符表达式 (??)
func (g *Generator) genNullCoalesceExpression(n *node.NullCoalesceExpression) {
	g.printf("node.NewNullCoalesceExpression(from,\n")
	g.indent++
	g.genGetValue(n.Left)
	g.printf(",\n")
	g.genGetValue(n.Right)
	g.printf(",\n")
	g.indent--
	g.printf(")")
}

// genLambdaExpression 生成 Lambda 表达式（匿名函数）
// LambdaExpression 嵌入 FunctionStatement，包含未导出字段 parent/ctx 和闭包捕获语义，
// 无法在编译期完整序列化。
func (g *Generator) genLambdaExpression(n *node.LambdaExpression) {
	g.printf("nil /* TODO: LambdaExpression — 包含闭包捕获语义，暂不支持编译期代码生成 */")
}

// genFunctionStatement 生成函数定义语句
// FunctionStatement 包含未导出字段 vars/defineCtx/staticLocals 和运行时 FuncStmt 接口，
// 无法在编译期完整序列化。
func (g *Generator) genFunctionStatement(n *node.FunctionStatement) {
	g.printf("nil /* TODO: FunctionStatement %q — 太复杂，暂不支持编译期代码生成 */", n.Name)
}

// genIncludeStatement 生成 include/require 语句
func (g *Generator) genIncludeStatement(n *node.IncludeStatement) {
	g.printf("node.NewIncludeStatement(from,\n")
	g.indent++
	g.genGetValue(n.Expr)
	g.printf(",\n")
	g.printf("%v,\n", n.Once)
	g.printf("%v,\n", n.Required)
	g.indent--
	g.printf(")")
}

// genConstStatement 生成常量声明语句
func (g *Generator) genConstStatement(n *node.ConstStatement) {
	g.printf("node.NewConstStatement(from,\n")
	g.indent++
	// Val 是 data.Variable 接口，使用 NewVariable 重建
	if n.Val != nil {
		g.printf("node.NewVariable(from, %q, %d, nil),\n", n.Val.GetName(), n.Val.GetIndex())
	} else {
		g.printf("nil,\n")
	}
	g.genGetValue(n.Initializer)
	g.printf(",\n")
	g.indent--
	g.printf(")")
}

// genSpreadArgument 生成展开参数 (...expr)
func (g *Generator) genSpreadArgument(n *node.SpreadArgument) {
	g.printf("node.NewSpreadArgument(from,\n")
	g.indent++
	g.genGetValue(n.Expr)
	g.printf(",\n")
	g.indent--
	g.printf(")")
}

// genNamedArgument 生成命名参数
func (g *Generator) genNamedArgument(n *node.NamedArgument) {
	g.printf("node.NewNamedArgument(from, %q,\n", n.Name)
	g.indent++
	g.genGetValue(n.Value)
	g.printf(",\n")
	g.indent--
	g.printf(")")
}

// genCompactStatement 生成 compact() 语句
func (g *Generator) genCompactStatement(n *node.CompactStatement) {
	g.printf("node.NewCompactStatement(from, []data.GetValue{\n")
	g.indent++
	for _, expr := range n.VarNames {
		g.genGetValue(expr)
		g.printf(",\n")
	}
	g.indent--
	g.printf("})")
}

// genRange 生成范围表达式 (start..stop 或 $arr[start..stop])
func (g *Generator) genRange(n *node.Range) {
	g.printf("node.NewRange(from,\n")
	g.indent++
	if n.Array != nil {
		g.genGetValue(n.Array)
	} else {
		g.printf("nil")
	}
	g.printf(",\n")
	if n.Start != nil {
		g.genGetValue(n.Start)
	} else {
		g.printf("nil")
	}
	g.printf(",\n")
	if n.Stop != nil {
		g.genGetValue(n.Stop)
	} else {
		g.printf("nil")
	}
	g.printf(",\n")
	g.indent--
	g.printf(")")
}

// genKv 生成键值对字面量
func (g *Generator) genKv(n *node.Kv) {
	g.printf("node.NewKv(from, []node.KvPair{\n")
	g.indent++
	for _, kv := range n.V {
		g.printf("{\n")
		g.indent++
		g.printf("Key: ")
		g.genGetValue(kv.Key)
		g.printf(",\n")
		g.printf("Value: ")
		g.genGetValue(kv.Value)
		g.printf(",\n")
		g.indent--
		g.printf("},\n")
	}
	g.indent--
	g.printf("})")
}
