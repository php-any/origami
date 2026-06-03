package cmd

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// Generator 将 AST 节点转换为 Go 构造代码
type Generator struct {
	buf     strings.Builder
	indent  int
	imports map[string]bool
	file    string
}

// NewGenerator 创建新的代码生成器
func NewGenerator() *Generator {
	return &Generator{
		imports: make(map[string]bool),
	}
}

// Generate 为一个解析后的文件生成 Go 代码
func (g *Generator) Generate(pf ParsedFile) string {
	g.file = pf.Path
	g.buf.Reset()
	g.imports = make(map[string]bool)
	g.indent = 0

	funcName := g.funcNameForPath(pf.Path)

	g.printf("func %s() (data.GetValue, []data.Variable) {\n", funcName)
	g.indent++
	g.printf("from := node.NewFrom(%q)\n", pf.Path)
	g.printf("\n")
	g.printf("stmts := []data.GetValue{\n")
	g.indent++
	for _, stmt := range pf.Program.Statements {
		g.genGetValue(stmt)
		g.printf(",\n")
	}
	g.indent--
	g.printf("}\n")
	g.printf("\n")

	// 生成变量列表
	g.printf("vars := []data.Variable{\n")
	g.indent++
	for _, v := range pf.Variables {
		g.printf("data.NewVariable(%q, %d, nil),\n", v.GetName(), v.GetIndex())
	}
	g.indent--
	g.printf("}\n")
	g.printf("\n")
	g.printf("return node.NewProgram(from, stmts), vars\n")
	g.indent--
	g.printf("}\n")

	return g.buf.String()
}

// funcNameForPath 将文件路径转换为合法的 Go 函数名
func (g *Generator) funcNameForPath(path string) string {
	name := strings.TrimSuffix(path, ".php")
	name = strings.TrimSuffix(name, ".zy")
	var b strings.Builder
	for _, c := range name {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			b.WriteRune(c)
		} else {
			b.WriteByte('_')
		}
	}
	result := b.String()
	if len(result) > 0 && result[0] >= 'a' && result[0] <= 'z' {
		result = string(rune(result[0]-32)) + result[1:]
	}
	return "AST_" + result
}

// genGetValue 根据节点类型分派到对应的生成函数
func (g *Generator) genGetValue(v data.GetValue) {
	if v == nil {
		g.printf("nil")
		return
	}
	switch n := v.(type) {
	// 字面量
	case *node.IntLiteral:
		g.genIntLiteral(n)
	case *node.FloatLiteral:
		g.genFloatLiteral(n)
	case *node.StringLiteral:
		g.genStringLiteral(n)
	case *node.BooleanLiteral:
		g.genBooleanLiteral(n)
	case *node.NullLiteral:
		g.genNullLiteral(n)
	// 变量
	case *node.VariableExpression:
		g.genVariableExpression(n)
	case *node.VariableReference:
		g.genVariableReference(n)
	// 二元运算
	case *node.BinaryAdd:
		g.genBinaryOp("BinaryAdd", n.Left, n.Right)
	case *node.BinarySub:
		g.genBinaryOp("BinarySub", n.Left, n.Right)
	case *node.BinaryMul:
		g.genBinaryOp("BinaryMul", n.Left, n.Right)
	case *node.BinaryQuo:
		g.genBinaryOp("BinaryQuo", n.Left, n.Right)
	case *node.BinaryRem:
		g.genBinaryOp("BinaryRem", n.Left, n.Right)
	case *node.BinaryPow:
		g.genBinaryOp("BinaryPow", n.Left, n.Right)
	case *node.BinaryDot:
		g.genBinaryOp("BinaryDot", n.Left, n.Right)
	case *node.BinaryEq:
		g.genBinaryOp("BinaryEq", n.Left, n.Right)
	case *node.BinaryNe:
		g.genBinaryOp("BinaryNe", n.Left, n.Right)
	case *node.BinaryEqStrict:
		g.genBinaryOp("BinaryEqStrict", n.Left, n.Right)
	case *node.BinaryNeStrict:
		g.genBinaryOp("BinaryNeStrict", n.Left, n.Right)
	case *node.BinaryLt:
		g.genBinaryOp("BinaryLt", n.Left, n.Right)
	case *node.BinaryLe:
		g.genBinaryOp("BinaryLe", n.Left, n.Right)
	case *node.BinaryGt:
		g.genBinaryOp("BinaryGt", n.Left, n.Right)
	case *node.BinaryGe:
		g.genBinaryOp("BinaryGe", n.Left, n.Right)
	case *node.BinaryLand:
		g.genBinaryOp("BinaryLand", n.Left, n.Right)
	case *node.BinaryLor:
		g.genBinaryOp("BinaryLor", n.Left, n.Right)
	case *node.BinarySpaceship:
		g.genBinaryOp("BinarySpaceship", n.Left, n.Right)
	// 一元/后缀运算
	case *node.UnaryExpression:
		g.genUnaryExpression(n)
	case *node.UnaryIncr:
		g.genUnaryIncr(n)
	case *node.UnaryDecr:
		g.genUnaryDecr(n)
	case *node.PostfixIncr:
		g.genPostfixIncr(n)
	case *node.PostfixDecr:
		g.genPostfixDecr(n)
	case *node.ErrorSuppress:
		g.genErrorSuppress(n)
	// 控制流
	case *node.IfStatement:
		g.genIfStatement(n)
	case *node.ReturnStatement:
		g.genReturnStatement(n)
	case *node.EchoStatement:
		g.genEchoStatement(n)
	// 函数/方法调用
	case *node.CallExpression:
		g.genCallExpression(n)
	case *node.CallMethod:
		g.genCallMethod(n)
	case *node.CallStaticMethod:
		g.genCallStaticMethod(n)
	case *node.CallObjectMethod:
		g.genCallObjectMethod(n)
	case *node.CallParentMethod:
		g.genCallParentMethod(n)
	case *node.CallSelfMethod:
		g.genCallSelfMethod(n)
	case *node.NullsafeCall:
		g.genNullsafeCall(n)
	// OOP 节点
	case *node.ClassStatement:
		g.genClassStatement(n)
	case *node.InterfaceStatement:
		g.genInterfaceStatement(n)
	case *node.NewExpression:
		g.genNewExpression(n)
	case *node.NewVariableExpression:
		g.genNewVariableExpression(n)
	case *node.NewExpressionDynamic:
		g.genNewExpressionDynamic(n)
	case *node.NewSelfExpression:
		g.genNewSelfExpression(n)
	case *node.NewStaticExpression:
		g.genNewStaticExpression(n)
	case *node.InstanceOfExpression:
		g.genInstanceOfExpression(n)
	case *node.CloneExpression:
		g.genCloneExpression(n)
	case *node.InitClass:
		g.genInitClass(n)
	case *node.ClassConstant:
		g.genClassConstant(n)
	case *node.StaticClass:
		g.genStaticClass(n)
	case *node.SelfClass:
		g.genSelfClass(n)
	case *node.Parent:
		g.genParent(n)
	// 未支持的类型
	default:
		g.printf("nil /* TODO: unsupported %T */", v)
	}
}

// genBinaryOp 生成二元运算的通用方法
func (g *Generator) genBinaryOp(typeName string, left, right data.GetValue) {
	g.printf("node.New%s(from,\n", typeName)
	g.indent++
	g.genGetValue(left)
	g.printf(",\n")
	g.genGetValue(right)
	g.printf(",\n")
	g.indent--
	g.printf(")")
}

// printf 带缩进的格式化输出
func (g *Generator) printf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	lines := strings.Split(msg, "\n")
	for i, line := range lines {
		if i > 0 {
			g.buf.WriteString("\n")
			if len(line) > 0 {
				for j := 0; j < g.indent; j++ {
					g.buf.WriteString("\t")
				}
			}
		}
		g.buf.WriteString(line)
	}
}
