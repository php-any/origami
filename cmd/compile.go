package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var compileCmd = &cobra.Command{
	Use:   "compile [directory]",
	Short: "将 vendor 目录预编译为 Go 源码",
	Long: `将 vendor/ 目录下的 PHP 文件解析为 AST，
并生成对应的 Go 结构体字面量代码。

生成的 Go 包可以在运行时注册到 VM，
使 LoadAndRun 跳过词法分析和解析阶段，
直接执行预构建的 AST。

示例:
  zy compile vendor/
  zy compile vendor/ -o .zy/build
  zy compile vendor/ --pkg myapp`,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE:         runCompileCommand,
}

var (
	compileOutput string
	compilePkg    string
)

func init() {
	compileCmd.Flags().StringVarP(&compileOutput, "output", "o", ".zy/build", "输出目录")
	compileCmd.Flags().StringVar(&compilePkg, "pkg", "build", "生成的 Go 包名")
}

func runCompileCommand(cmd *cobra.Command, args []string) error {
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

	fmt.Printf("找到 %d 个 PHP 文件\n", len(files))
	return nil
}
