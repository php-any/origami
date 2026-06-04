package genstd

import (
	"github.com/php-any/origami/internal/pseudocode"
	"github.com/spf13/cobra"
)

// NewCommand 创建 gen-std 子命令。
func NewCommand() *cobra.Command {
	var output string
	cmd := &cobra.Command{
		Use:   "gen-std",
		Short: "从 Go 标准库实现反射生成 PHP 伪代码（IDE 提示）",
		Long: `扫描 Go 实现的标准库（函数与类），通过反射提取签名并生成 PHP 伪代码。

生成的 .php 伪代码文件可用于 IDE 自动补全与类型提示，默认输出到当前目录的 .zy/std/ 目录。

示例:
  zy gen-std
  zy gen-std -o docs/std`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return pseudocode.Generate(output)
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", ".zy/std", "伪代码输出目录")
	return cmd
}
