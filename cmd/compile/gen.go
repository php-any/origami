package compile

import (
	"fmt"
	"strings"
)

// Generator 将 AST 节点转换为 Go 构造代码
type Generator struct {
	buf           strings.Builder
	indent        int
	importAliases map[string]string // import path -> alias（空字符串表示使用默认包名）
	file          string
	namespace     string // 当前文件的命名空间
}

// NewGenerator 创建新的代码生成器
func NewGenerator() *Generator {
	return &Generator{
		importAliases: make(map[string]string),
	}
}

// Generate 为一个解析后的文件生成 Go 代码
func (g *Generator) Generate(pf ParsedFile) (string, error) {
	g.file = pf.Path
	g.namespace = pf.Namespace
	g.buf.Reset()
	g.importAliases = make(map[string]string)
	g.indent = 0

	funcName := g.funcNameForPath(pf.Path)

	g.printf("func %s() (data.GetValue, []data.Variable) {\n", funcName)
	g.indent++
	g.printf("filePath := %q\n", pf.Path)
	g.printf("from := node.NewTokenFrom(&filePath, 0, 0, 0, 0)\n")
	g.printf("\n")
	g.printf("stmts := []data.GetValue{\n")
	g.indent++
	for _, stmt := range pf.Program.Statements {
		if err := g.Emit(stmt); err != nil {
			return "", wrapEmitError(g.file, stmt, err)
		}
		g.printf(",\n")
	}
	g.indent--
	g.printf("}\n")
	g.printf("\n")

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

	return g.buf.String(), nil
}

// funcNameForPath 将文件路径转换为合法的 Go 函数名
func (g *Generator) funcNameForPath(path string) string {
	return funcNameFromPath(path)
}

// goFileNameForPath 将 PHP 路径映射为同包内的 Go 源文件名
func (g *Generator) goFileNameForPath(path string) string {
	return goFileNameFromPath(path)
}

// printf 带缩进的格式化输出
func (g *Generator) printf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if g.indent <= 0 {
		g.buf.WriteString(msg)
		return
	}
	pad := strings.Repeat("\t", g.indent)
	lines := strings.Split(msg, "\n")
	for i, line := range lines {
		if i > 0 {
			g.buf.WriteString("\n")
		}
		if len(line) > 0 {
			g.buf.WriteString(pad)
		}
		g.buf.WriteString(line)
	}
}
