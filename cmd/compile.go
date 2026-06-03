package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var compileCmd = &cobra.Command{
	Use:   "compile [directory]",
	Short: "将 vendor 目录预编译为 Go 源码",
	Long: `将目录下的 PHP 文件解析为 AST，并生成对应的 Go 结构体字面量代码。

默认模式：生成 Go 包，运行时注册到 VM 跳过解析。
--build 模式：生成完整 Go 项目并编译为独立二进制。

示例:
  zy compile vendor/
  zy compile vendor/ -o .zy/build
  zy compile . --build --entry=app.php
  zy compile . --build --entry=app.php -o dist/myapp`,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE:         runCompileCommand,
}

var (
	compileOutput string
	compilePkg    string
	compileBuild  bool
	compileEntry  string
)

func init() {
	compileCmd.Flags().StringVarP(&compileOutput, "output", "o", ".zy/build", "输出目录")
	compileCmd.Flags().StringVar(&compilePkg, "pkg", "build", "生成的 Go 包名")
	compileCmd.Flags().BoolVarP(&compileBuild, "build", "b", false, "编译为独立二进制")
	compileCmd.Flags().StringVar(&compileEntry, "entry", "", "入口 PHP 文件（--build 模式必填）")
}

func runCompileCommand(cmd *cobra.Command, args []string) error {
	if compileBuild {
		if compileEntry == "" {
			return fmt.Errorf("--build 模式需要指定 --entry 参数")
		}
		if compilePkg == "build" {
			compilePkg = "main"
		}
	}

	vendorDir := args[0]

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

	parsed, parseErrs := parseFiles(files)
	if len(parseErrs) > 0 {
		for _, e := range parseErrs {
			fmt.Fprintf(os.Stderr, "警告: %v\n", e)
		}
	}
	if len(parsed) == 0 {
		return fmt.Errorf("没有文件解析成功")
	}

	fmt.Printf("成功解析 %d 个文件\n", len(parsed))

	if err := generateOutput(parsed, compileOutput, compilePkg); err != nil {
		return fmt.Errorf("生成失败: %w", err)
	}

	fmt.Printf("已生成 Go 包到 %s\n", compileOutput)

	if compileBuild {
		if err := generateMainFile(compileEntry, compileOutput, compilePkg); err != nil {
			return fmt.Errorf("生成 main.go 失败: %w", err)
		}
		if err := buildBinary(compileOutput); err != nil {
			return fmt.Errorf("编译失败: %w", err)
		}
		fmt.Println("编译完成！")
	}

	return nil
}
