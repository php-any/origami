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
	content := fmt.Sprintf("module %s\n\ngo 1.25.0\n\nrequire github.com/php-any/origami v0.0.0\n", pkgName)

	// 自动检测 origami 仓库路径并添加 replace 指令
	if origamiPath := findOrigamiPath(); origamiPath != "" {
		absOutput, err := filepath.Abs(outputDir)
		if err == nil {
			relPath, err := filepath.Rel(absOutput, origamiPath)
			if err == nil {
				// 统一使用正斜杠，确保跨平台兼容
				relPath = filepath.ToSlash(relPath)
				content += fmt.Sprintf("\nreplace github.com/php-any/origami => %s\n", relPath)
			}
		}
	}

	return os.WriteFile(filepath.Join(outputDir, "go.mod"), []byte(content), 0644)
}

// findOrigamiPath 查找 origami 仓库根目录（包含 go.mod 且 module 名为 github.com/php-any/origami）
func findOrigamiPath() string {
	// 尝试从可执行文件位置向上查找
	execPath, err := os.Executable()
	if err == nil {
		if path := searchGoModUp(filepath.Dir(execPath)); path != "" {
			return path
		}
	}
	// 尝试从当前工作目录向上查找
	if cwd, err := os.Getwd(); err == nil {
		if path := searchGoModUp(cwd); path != "" {
			return path
		}
	}
	return ""
}

// searchGoModUp 从给定目录向上搜索包含 origami module 定义的 go.mod
func searchGoModUp(dir string) string {
	for i := 0; i < 10; i++ { // 最多向上搜索 10 层
		gomod := filepath.Join(dir, "go.mod")
		data, err := os.ReadFile(gomod)
		if err == nil {
			content := string(data)
			if strings.Contains(content, "module github.com/php-any/origami") {
				return dir
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}
