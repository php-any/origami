package compile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/php-any/origami/data"
	"github.com/spf13/cobra"
)

var runtimeLoader func(vm data.VM)

var (
	compileOutput string
	compilePkg    string
	compileBuild  bool
	compileEntry  string
)

// NewCommand 创建 compile 子命令。
func NewCommand(loader func(vm data.VM)) *cobra.Command {
	runtimeLoader = loader
	cmd := &cobra.Command{
		Use:   "compile [directory]",
		Short: "将 vendor 目录预编译为 Go 源码",
		Long: `将目录下的 PHP 文件解析为 AST，并生成对应的 Go 结构体字面量代码。

默认模式：生成 Go 包，运行时注册到 VM 跳过解析。
--build 模式：生成完整 Go 项目并编译为独立二进制。

示例:
  zy compile vendor/
  zy compile vendor/ -o dist
  zy compile . --build --entry=app.php
  zy compile . --build --entry=app.php -o dist/myapp`,
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE:         runCompileCommand,
	}
	cmd.Flags().StringVarP(&compileOutput, "output", "o", "", "输出目录（默认 <dir>/build/dist）")
	cmd.Flags().StringVar(&compilePkg, "pkg", "build", "生成的 Go 包名")
	cmd.Flags().BoolVarP(&compileBuild, "build", "b", false, "编译为独立二进制")
	cmd.Flags().StringVar(&compileEntry, "entry", "", "入口 PHP 文件（--build 模式必填）")
	return cmd
}

func runCompileCommand(cmd *cobra.Command, args []string) error {
	vendorDir := args[0]

	// 默认输出到 <dir>/build/dist
	if compileOutput == "" {
		compileOutput = filepath.Join(vendorDir, "build", "dist")
	}

	if compileBuild {
		if compileEntry == "" {
			compileEntry = filepath.Join(vendorDir, "index.php")
		}
		if compilePkg == "build" {
			compilePkg = "main"
		}
	}

	info, err := os.Stat(vendorDir)
	if err != nil {
		return fmt.Errorf("目录不存在: %s", vendorDir)
	}
	if !info.IsDir() {
		return fmt.Errorf("不是目录: %s", vendorDir)
	}

	files, err := collectPhpFiles(vendorDir)
	if err != nil {
		return fmt.Errorf("扫描失败: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("未找到 .php 文件: %s", vendorDir)
	}

	fmt.Printf("找到 %d 个 PHP 文件，开始解析...\n", len(files))

	// 解析 --entry 为规范路径集合（支持文件或目录）
	entryPaths, err := resolveEntryPaths(compileEntry)
	if err != nil {
		return fmt.Errorf("解析 --entry 失败: %w", err)
	}

	parsed, parseErrs := parseFiles(files)
	if len(parseErrs) > 0 {
		for _, e := range parseErrs {
			fmt.Fprintf(os.Stderr, "解析失败: %v\n", e)
		}
		return fmt.Errorf("有 %d 个文件解析失败，编译终止", len(parseErrs))
	}
	if len(parsed) == 0 {
		return fmt.Errorf("没有文件解析成功")
	}

	fmt.Printf("成功解析 %d 个文件（其中 entry %d 个）\n", len(parsed), len(entryPaths))

	if err := generateOutput(parsed, entryPaths, compileOutput, compilePkg); err != nil {
		return fmt.Errorf("生成失败: %w", err)
	}

	fmt.Printf("已生成 Go 包到 %s\n", compileOutput)

	if compileBuild {
		// --build 时 entry 必须是单个文件
		entryFile, err := resolveSingleEntry(compileEntry)
		if err != nil {
			return fmt.Errorf("--build 模式 --entry 需要指定单个 PHP 文件: %w", err)
		}
		if err := generateMainFile(entryFile, compileOutput, compilePkg); err != nil {
			return fmt.Errorf("生成 main.go 失败: %w", err)
		}
		if err := buildBinary(compileOutput, vendorDir); err != nil {
			return fmt.Errorf("编译失败: %w", err)
		}
		fmt.Println("编译完成！")
	}

	return nil
}

// resolveEntryPaths 解析 --entry 参数，支持单文件或目录，返回规范路径集合。
// 若 entry 为空，返回空 map（无 entry 文件）。
func resolveEntryPaths(entry string) (map[string]bool, error) {
	result := make(map[string]bool)
	if entry == "" {
		return result, nil
	}
	info, err := os.Stat(entry)
	if err != nil {
		return nil, fmt.Errorf("路径不存在: %s", entry)
	}
	if info.IsDir() {
		files, err := collectPhpFiles(entry)
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			result[filepath.Clean(f)] = true
		}
	} else {
		result[filepath.Clean(entry)] = true
	}
	return result, nil
}

// resolveSingleEntry 确保 --build entry 是单个文件，返回其路径。
func resolveSingleEntry(entry string) (string, error) {
	info, err := os.Stat(entry)
	if err != nil {
		return "", fmt.Errorf("路径不存在: %s", entry)
	}
	if info.IsDir() {
		return "", fmt.Errorf("%s 是目录，--build 模式需要指定单个入口文件", entry)
	}
	return entry, nil
}
