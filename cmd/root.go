package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "zy [脚本路径]",
	Short: "折言(origami-lang) - 融合型脚本语言",
	Long: `折言(origami-lang) - 融合型脚本语言

直接运行脚本:
  zy script.php
  zy tests/run_tests.zy

支持的脚本格式:
  .zy  - 折言脚本文件
  .php - PHP 兼容脚本文件`,
}

func Execute() {
	localizeCompletionCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// IsDirectScriptArg 判断参数是否应作为脚本路径直接执行（而非子命令）。
func IsDirectScriptArg(arg string) bool {
	if strings.HasPrefix(arg, "-") {
		return false
	}
	switch arg {
	case "gen-std", "help", "completion", "phpt":
		return false
	}
	return true
}

// RunScriptFile 直接运行指定脚本，等价于 zy <脚本路径>。
func RunScriptFile(scriptPath string) error {
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		_, _ = fmt.Fprintf(os.Stderr, "错误: 文件 '%s' 不存在\n\n", scriptPath)
		return rootCmd.Help()
	}

	vm, p := getRuntimeVM()
	_, err := vm.LoadAndRun(scriptPath)
	if err != nil {
		p.ShowControl(err)
	}
	vm.RunShutdownCallbacks()
	return nil
}

func init() {
	rootCmd.AddCommand(genStdCmd)
	rootCmd.AddCommand(phptCmd)
}

// RootHelp 显示根命令帮助信息。
func RootHelp() error {
	localizeCompletionCmd()
	return rootCmd.Help()
}
