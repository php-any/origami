package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// generateOutput 生成最终的 Go 源码文件
func generateOutput(parsed []ParsedFile, outputDir, pkgName string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	if err := generateRegisterFile(parsed, outputDir, pkgName); err != nil {
		return err
	}

	if err := generateASTFile(parsed, outputDir, pkgName); err != nil {
		return err
	}

	if err := generateGoMod(outputDir, pkgName); err != nil {
		return err
	}

	return nil
}

func generateRegisterFile(parsed []ParsedFile, outputDir, pkgName string) error {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("package %s\n\n", pkgName))
	b.WriteString("import (\n")
	b.WriteString("\t\"github.com/php-any/origami/data\"\n")
	b.WriteString(")\n\n")
	b.WriteString("// Register 将预编译的 vendor AST 注册到 VM\n")
	b.WriteString("func Register(vm data.VM) {\n")
	for _, pf := range parsed {
		gen := NewGenerator()
		funcName := gen.funcNameForPath(pf.Path)
		b.WriteString(fmt.Sprintf("\tvm.RegisterCompiledFile(%q, func() (data.GetValue, []data.Variable) {\n", pf.Path))
		b.WriteString(fmt.Sprintf("\t\treturn %s()\n", funcName))
		b.WriteString("\t})\n")
	}
	b.WriteString("}\n")

	return os.WriteFile(filepath.Join(outputDir, "register.go"), []byte(b.String()), 0644)
}

func generateASTFile(parsed []ParsedFile, outputDir, pkgName string) error {
	gen := NewGenerator()
	var b strings.Builder
	b.WriteString(fmt.Sprintf("package %s\n\n", pkgName))
	b.WriteString("import (\n")
	b.WriteString("\t\"github.com/php-any/origami/data\"\n")
	b.WriteString("\t\"github.com/php-any/origami/node\"\n")
	b.WriteString(")\n\n")

	for _, pf := range parsed {
		code := gen.Generate(pf)
		b.WriteString(code)
		b.WriteString("\n\n")
	}

	return os.WriteFile(filepath.Join(outputDir, "vendor_ast.go"), []byte(b.String()), 0644)
}

func generateGoMod(outputDir, pkgName string) error {
	content := fmt.Sprintf("module %s\n\ngo 1.21\n\nrequire github.com/php-any/origami v0.0.0\n", pkgName)
	return os.WriteFile(filepath.Join(outputDir, "go.mod"), []byte(content), 0644)
}
