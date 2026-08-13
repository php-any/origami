package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/cli/annotation"
)

// CliRuntime CLI 运行时，负责命令分发和执行
type CliRuntime struct {
	appName    string
	appVersion string
}

// NewCliRuntime 创建 CLI 运行时
func NewCliRuntime(appName, appVersion string) *CliRuntime {
	return &CliRuntime{
		appName:    appName,
		appVersion: appVersion,
	}
}

// Run 运行 CLI 应用
func (r *CliRuntime) Run(ctx data.Context) data.Control {
	args := os.Args[1:]

	if len(args) == 0 {
		r.showHelp()
		return nil
	}

	commandName := args[0]

	// 处理特殊命令
	switch commandName {
	case "--help", "-h":
		r.showHelp()
		return nil
	case "--version", "-v":
		r.showVersion()
		return nil
	}

	// 执行命令
	return annotation.ExecuteCommand(ctx, commandName)
}

// showHelp 显示帮助信息
func (r *CliRuntime) showHelp() {
	fmt.Printf("%s %s\n\n", r.appName, r.appVersion)
	fmt.Println("Usage:")
	fmt.Printf("  %s <command> [options]\n\n", os.Args[0])
	fmt.Println("Available Commands:")

	commands := annotation.GetRegisteredCommands()
	maxNameLength := 0

	// 计算最大命令名长度
	for name := range commands {
		if len(name) > maxNameLength {
			maxNameLength = len(name)
		}
	}

	// 显示命令列表
	for name, cmd := range commands {
		padding := strings.Repeat(" ", maxNameLength-len(name)+2)
		fmt.Printf("  %s%s%s\n", name, padding, cmd.GetDescription())
	}

	fmt.Printf("\nUse \"%s <command> --help\" for more information about a command.\n", os.Args[0])
}

// showVersion 显示版本信息
func (r *CliRuntime) showVersion() {
	fmt.Printf("%s version %s\n", r.appName, r.appVersion)
}
