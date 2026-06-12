package compile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// generateOutput 生成最终的 Go 源码文件
func generateOutput(parsed []ParsedFile, entryPaths map[string]bool, outputDir, pkgName string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	if err := generateRegisterFile(parsed, entryPaths, outputDir, pkgName); err != nil {
		return err
	}

	if err := generateASTFiles(parsed, outputDir, pkgName); err != nil {
		return err
	}

	if err := generateGoMod(outputDir, pkgName); err != nil {
		return err
	}

	return nil
}

func generateRegisterFile(parsed []ParsedFile, entryPaths map[string]bool, outputDir, pkgName string) error {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("package %s\n\n", pkgName))
	b.WriteString("import (\n")
	b.WriteString("\t\"github.com/php-any/origami/data\"\n")
	b.WriteString(")\n\n")
	b.WriteString("// Register 将预编译的 AST 注册到 VM\n")
	b.WriteString("// - entry 文件：注册到 compiledFiles，供 run_php_file() 显式执行\n")
	b.WriteString("// - 其余文件：直接执行 AST 函数，让类/函数/接口自动挂靠到 VM\n")
	b.WriteString("func Register(vm data.VM) {\n")
	for _, pf := range parsed {
		gen := NewGenerator()
		funcName := gen.funcNameForPath(pf.Path)
		if entryPaths[pf.Path] {
			// entry 文件：注册到 compiledFiles，仅当 run_php_file() 被调用时才执行
			b.WriteString(fmt.Sprintf("\tvm.RegisterCompiledFile(%q, func() (data.GetValue, []data.Variable) {\n", pf.Path))
			b.WriteString(fmt.Sprintf("\t\treturn %s()\n", funcName))
			b.WriteString("\t})\n")
		} else {
			// 声明文件：直接执行，类/函数/接口自动注册进 VM
			b.WriteString(fmt.Sprintf("\tif program, vars := %s(); program != nil {\n", funcName))
			b.WriteString("\t\tctx := vm.CreateContext(vars)\n")
			b.WriteString("\t\tprogram.GetValue(ctx) //nolint:errcheck\n")
			b.WriteString("\t}\n")
		}
	}
	b.WriteString("}\n")

	return writeFormattedGoFile(filepath.Join(outputDir, "register.go"), []byte(b.String()))
}

func generateASTFiles(parsed []ParsedFile, outputDir, pkgName string) error {
	// 移除旧版单文件输出及上次生成的分文件 AST
	_ = os.Remove(filepath.Join(outputDir, "vendor_ast.go"))
	matches, err := filepath.Glob(filepath.Join(outputDir, "ast_*.go"))
	if err != nil {
		return fmt.Errorf("扫描旧 AST 文件失败: %w", err)
	}
	for _, f := range matches {
		if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("删除旧 AST 文件 %s 失败: %w", f, err)
		}
	}

	for _, pf := range parsed {
		gen := NewGenerator()
		code, err := gen.Generate(pf)
		if err != nil {
			return err
		}
		importAliases := map[string]string{
			"github.com/php-any/origami/data": "",
			"github.com/php-any/origami/node": "",
		}
		for imp, alias := range gen.importAliases {
			importAliases[imp] = alias
		}

		var b strings.Builder
		b.WriteString(fmt.Sprintf("package %s\n\n", pkgName))
		b.WriteString("import (\n")
		impList := make([]string, 0, len(importAliases))
		for imp := range importAliases {
			impList = append(impList, imp)
		}
		sort.Strings(impList)
		for _, imp := range impList {
			if alias := importAliases[imp]; alias != "" {
				b.WriteString(fmt.Sprintf("\t%s %q\n", alias, imp))
			} else {
				b.WriteString(fmt.Sprintf("\t%q\n", imp))
			}
		}
		b.WriteString(")\n\n")
		b.WriteString(code)

		outPath := filepath.Join(outputDir, gen.goFileNameForPath(pf.Path))
		if err := writeFormattedGoFile(outPath, []byte(b.String())); err != nil {
			return err
		}
	}
	return nil
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
