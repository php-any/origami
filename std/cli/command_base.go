package cli

import (
	"fmt"
	"os"
)

// CommandBase 命令基类，提供通用功能
type CommandBase struct {
	name        string
	description string
	parser      *OptionParser
}

// NewCommandBase 创建命令基类
func NewCommandBase(name, description string) *CommandBase {
	return &CommandBase{
		name:        name,
		description: description,
		parser:      NewOptionParser(),
	}
}

// GetName 获取命令名称
func (c *CommandBase) GetName() string {
	return c.name
}

// GetDescription 获取命令描述
func (c *CommandBase) GetDescription() string {
	return c.description
}

// AddOption 添加选项
func (c *CommandBase) AddOption(opt Option) {
	c.parser.AddOption(opt)
}

// ParseArgs 解析参数
func (c *CommandBase) ParseArgs() (map[string]string, []string, error) {
	return c.parser.Parse(os.Args[1:])
}

// Output 输出信息
func (c *CommandBase) Output(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// OutputLine 输出一行信息
func (c *CommandBase) OutputLine(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Error 输出错误信息
func (c *CommandBase) Error(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
}

// ErrorExit 输出错误信息并退出
func (c *CommandBase) ErrorExit(code int, format string, args ...interface{}) {
	c.Error(format, args...)
	os.Exit(code)
}

// ShowHelp 显示帮助信息
func (c *CommandBase) ShowHelp() {
	fmt.Printf("Usage: %s [options] [arguments]\n\n", c.name)
	fmt.Printf("Description:\n  %s\n\n", c.description)
	fmt.Println("Options:")

	for _, opt := range c.parser.options {
		if opt.ShortName != "" {
			fmt.Printf("  -%s, --%-20s %s\n", opt.ShortName, opt.Name, opt.Description)
		} else {
			fmt.Printf("      --%-20s %s\n", opt.Name, opt.Description)
		}
	}
}
