package compile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

// ParsedFile 表示一个已解析的 PHP 文件
type ParsedFile struct {
	Path      string
	Program   *node.Program
	Variables []data.Variable
	Namespace string // 文件的命名空间
}

// parseFiles 批量解析 PHP 文件为 AST。
// 使用统一的全局 VM，编译模式下只解析不执行代码。
func parseFiles(files []string) ([]ParsedFile, []error) {
	p := parser.NewParser()
	baseVM := runtime.NewVM(p)
	if runtimeLoader != nil {
		runtimeLoader(baseVM)
	}

	// 编译模式下跳过 #[Application] 注解的扫描/启动逻辑
	data.CompileMode = true
	defer func() { data.CompileMode = false }()

	// 优先加载 composer autoload 配置，注册 namespace->path 映射
	if err := loadComposerAutoload(p.GetClassPathManager(), files); err != nil {
		fmt.Printf("加载 composer autoload 配置失败: %v\n", err)
	}

	var parsed []ParsedFile
	var errs []error

	for _, file := range files {
		clone := p.Clone()
		program, acl := clone.ParseFile(file)
		if acl != nil {
			errs = append(errs, fmt.Errorf("解析 %s 失败: %v", file, acl))
			continue
		}
		augmentProgramASTFromBase(program, baseVM.(*runtime.VM), file)
		vars := clone.GetVariables()
		parsed = append(parsed, ParsedFile{
			Path:      file,
			Program:   program,
			Variables: vars,
			Namespace: clone.GetNamespace(),
		})
	}
	return parsed, errs
}

// loadComposerAutoload 查找并加载 composer autoload 配置，注册 namespace->path 映射。
func loadComposerAutoload(cpm parser.ClassPathManager, files []string) error {
	// 从文件列表中推断项目根目录
	var projectDir string
	for _, f := range files {
		dir := filepath.Dir(f)
		for d := dir; d != "/" && d != "."; d = filepath.Dir(d) {
			jsonPath := filepath.Join(d, "composer.json")
			if info, err := os.Stat(jsonPath); err == nil && !info.IsDir() {
				projectDir = d
				break
			}
		}
		if projectDir != "" {
			break
		}
	}
	if projectDir == "" {
		return nil
	}

	vendorDir := filepath.Join(projectDir, "vendor")

	// 1. 加载项目根 composer.json 的 autoload 配置
	if err := loadComposerJSONAutoload(cpm, filepath.Join(projectDir, "composer.json"), projectDir); err != nil {
		fmt.Printf("加载项目 composer.json 失败: %v\n", err)
	}

	// 2. 加载 vendor/composer/installed.json 中所有包的 autoload 配置
	installedPath := filepath.Join(vendorDir, "composer", "installed.json")
	if info, err := os.Stat(installedPath); err == nil && !info.IsDir() {
		if err := loadInstalledJSONAutoload(cpm, installedPath, vendorDir); err != nil {
			fmt.Printf("加载 installed.json 失败: %v\n", err)
		}
	}

	return nil
}

// loadComposerJSONAutoload 从单个 composer.json 加载 autoload 映射。
func loadComposerJSONAutoload(cpm parser.ClassPathManager, jsonPath, baseDir string) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return err
	}

	var composer struct {
		Autoload struct {
			PSR4 map[string]string `json:"psr-4"`
			PSR0 map[string]string `json:"psr-0"`
		} `json:"autoload"`
	}
	if err := json.Unmarshal(data, &composer); err != nil {
		return fmt.Errorf("解析 composer.json 失败: %w", err)
	}

	registerAutoloadMappings(cpm, composer.Autoload.PSR4, composer.Autoload.PSR0, baseDir)
	return nil
}

// loadInstalledJSONAutoload 从 vendor/composer/installed.json 加载所有包的 autoload 映射。
func loadInstalledJSONAutoload(cpm parser.ClassPathManager, installedPath, vendorDir string) error {
	data, err := os.ReadFile(installedPath)
	if err != nil {
		return err
	}

	var installed struct {
		Packages []struct {
			Name     string `json:"name"`
			Autoload struct {
				PSR4 map[string]string `json:"psr-4"`
				PSR0 map[string]string `json:"psr-0"`
			} `json:"autoload"`
		} `json:"packages"`
	}
	if err := json.Unmarshal(data, &installed); err != nil {
		return fmt.Errorf("解析 installed.json 失败: %w", err)
	}

	for _, pkg := range installed.Packages {
		pkgDir := filepath.Join(vendorDir, pkg.Name)
		registerAutoloadMappings(cpm, pkg.Autoload.PSR4, pkg.Autoload.PSR0, pkgDir)
	}

	return nil
}

// registerAutoloadMappings 注册 PSR-4 和 PSR-0 映射到 ClassPathManager。
func registerAutoloadMappings(cpm parser.ClassPathManager, psr4, psr0 map[string]string, baseDir string) {
	for ns, relPath := range psr4 {
		absPath := filepath.Join(baseDir, relPath)
		cpm.AddNamespace(ns, absPath)
	}
	for ns, relPath := range psr0 {
		absPath := filepath.Join(baseDir, relPath)
		cpm.AddNamespace(ns, absPath)
	}
}

// augmentProgramASTFromBase 从全局 baseVM 中找出当前文件注册的类，补充到 AST 中。
func augmentProgramASTFromBase(program *node.Program, baseVM *runtime.VM, file string) {
	var classes []data.GetValue
	for _, c := range baseVM.AllClasses() {
		if classSourceFile(c) != file {
			continue
		}
		if gv, ok := c.(data.GetValue); ok {
			classes = append(classes, gv)
		}
	}
	if len(classes) == 0 {
		return
	}
	for _, stmt := range program.Statements {
		ns, ok := stmt.(*node.Namespace)
		if !ok {
			continue
		}
		ns.Statements = append(filterNamespaceStmts(ns.Statements, ns.Name), classes...)
	}
}

func filterNamespaceStmts(stmts []data.GetValue, nsName string) []data.GetValue {
	out := make([]data.GetValue, 0, len(stmts))
	for _, s := range stmts {
		if lit, ok := s.(*node.StringLiteral); ok && lit.Value == nsName {
			continue
		}
		out = append(out, s)
	}
	return out
}

func classSourceFile(c data.ClassStmt) string {
	switch s := c.(type) {
	case *node.ClassStatement:
		if f := s.GetFrom(); f != nil {
			return f.GetSource()
		}
	case *node.AbstractClassStatement:
		if f := s.ClassStatement.GetFrom(); f != nil {
			return f.GetSource()
		}
	}
	return ""
}
