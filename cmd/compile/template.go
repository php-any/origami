package compile

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// TemplateData 传递给模板的变量
type TemplateData struct {
	Pkg       string
	HasEntry  bool
	EntryPath string
	Files     []FileInfo
}

// FileInfo 单个 PHP 文件的信息
type FileInfo struct {
	FuncName  string // AST 函数名
	IsEntry   bool
	Path      string   // 原始 PHP 文件路径
	Namespace string   // 命名空间
	Classes   []string // 文件中定义的类名（全限定名）
}

// defaultRegisterTmpl 内置默认 register.go.tmpl
const defaultRegisterTmpl = `package {{.Pkg}}

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

{{- if .HasEntry}}
const EntryPath = {{printf "%q" .EntryPath}}
{{end}}

func registerClasses(vm data.VM, program data.GetValue) {
	prog, ok := program.(*node.Program)
	if !ok {
		return
	}
	for _, stmt := range prog.Statements {
		ns, ok := stmt.(*node.Namespace)
		if !ok {
			continue
		}
		for _, s := range ns.GetStatements() {
			if cs, ok := s.(data.ClassStmt); ok {
				vm.AddClass(cs)
			}
		}
	}
}

// Register 将预编译的 AST 注册到 VM
func Register(vm data.VM) {
{{- range .Files}}
{{- if .IsEntry}}
	vm.RegisterCompiledFile(EntryPath, func() (data.GetValue, []data.Variable) {
		return {{.FuncName}}()
	})
{{- else}}
	if program, vars := {{.FuncName}}(); program != nil {
		registerClasses(vm, program)
		ctx := vm.CreateContext(vars)
		program.GetValue(ctx) //nolint:errcheck
	}
{{- end}}
{{- end}}
}
`

// defaultMainTmpl 内置默认 main.go.tmpl
const defaultMainTmpl = `package main

import (
	"fmt"
	"os"

	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/net/http"
	"github.com/php-any/origami/std/net/websocket"
	"github.com/php-any/origami/std/net/annotation"
	"github.com/php-any/origami/std/system"
)

func main() {
	p := parser.NewParser()
	vm := runtime.NewVM(p)

	std.Load(vm)
	php.Load(vm)
	http.Load(vm)
	websocket.Load(vm)
	annotation.Load(vm)
	system.Load(vm)

	Register(vm)
{{- if .HasEntry}}
	_, err := vm.RunCompiledFile(EntryPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
{{- end}}
}
`

// defaultModTmpl 内置默认 go.mod.tmpl
const defaultModTmpl = `module {{.Pkg}}

go 1.25.0

require github.com/php-any/origami v0.0.0
`

// builtinTemplates 内置模板名 → 内容映射
var builtinTemplates = map[string]string{
	"register.go.tmpl": defaultRegisterTmpl,
	"main.go.tmpl":     defaultMainTmpl,
	"go.mod.tmpl":      defaultModTmpl,
}

// loadTemplate 加载模板：优先 <.zy/> 目录，否则用内置默认
func loadTemplate(sourceDir, name string) (*template.Template, error) {
	zyPath := filepath.Join(sourceDir, ".zy", name)
	if data, err := os.ReadFile(zyPath); err == nil {
		return template.New(name).Parse(string(data))
	}
	if def, ok := builtinTemplates[name]; ok {
		return template.New(name).Parse(def)
	}
	return nil, fmt.Errorf("模板 %s 不存在，且无内置默认", name)
}

// buildTemplateData 从解析结果构建模板数据
func buildTemplateData(parsed []ParsedFile, entryPaths map[string]bool, pkgName string) *TemplateData {
	td := &TemplateData{
		Pkg:   pkgName,
		Files: make([]FileInfo, 0, len(parsed)),
	}

	for _, pf := range parsed {
		isEntry := entryPaths[filepath.Clean(pf.Path)]
		if isEntry {
			td.HasEntry = true
			absPath, err := filepath.Abs(pf.Path)
			if err != nil {
				absPath = pf.Path
			}
			td.EntryPath = absPath
		}

		gen := NewGenerator()
		fi := FileInfo{
			FuncName:  gen.funcNameForPath(pf.Path),
			IsEntry:   isEntry,
			Path:      pf.Path,
			Namespace: pf.Namespace,
			Classes:   extractClassNames(pf.Program),
		}
		td.Files = append(td.Files, fi)
	}

	return td
}

// extractClassNames 从 Program AST 中提取所有类名（全限定名）
func extractClassNames(program *node.Program) []string {
	var names []string
	for _, stmt := range program.Statements {
		ns, ok := stmt.(*node.Namespace)
		if !ok {
			continue
		}
		nsPrefix := ""
		if ns.Name != "" {
			nsPrefix = ns.Name + "\\"
		}
		for _, s := range ns.GetStatements() {
			if cs, ok := s.(data.ClassStmt); ok {
				names = append(names, nsPrefix+cs.GetName())
			}
		}
	}
	return names
}

// renderTemplate 渲染模板为字节，自动 gofmt（Go 源码）
func renderTemplate(tmpl *template.Template, data *TemplateData) ([]byte, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("模板渲染失败: %w", err)
	}
	raw := buf.Bytes()
	formatted, err := format.Source(raw)
	if err != nil {
		// 格式化失败时返回原始内容
		return raw, nil
	}
	return formatted, nil
}

// renderGoMod 渲染 go.mod（不 gofmt）
func renderGoMod(tmpl *template.Template, data *TemplateData) (string, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("模板渲染失败: %w", err)
	}
	return postProcessGoMod(buf.String()), nil
}

// postProcessGoMod 在 go.mod 后追加 replace 指令
func postProcessGoMod(content string) string {
	origamiPath := findOrigamiPath()
	if origamiPath == "" {
		return content
	}
	// 尾部补充 replace 行
	content = strings.TrimRight(content, "\n")
	content += fmt.Sprintf("\n\nreplace github.com/php-any/origami => %s\n", filepath.ToSlash(origamiPath))
	return content
}
