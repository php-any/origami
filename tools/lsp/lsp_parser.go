package main

import (
	"fmt"
	"github.com/php-any/origami/data"
	"os"

	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
)

func init() {
	parser.InLSP = true
}

// LspParser 是专门为 LSP 服务器设计的解析器
type LspParser struct {
	vm     *LspVM
	parser *parser.Parser
}

// NewLspParser 创建一个新的 LSP 解析器
func NewLspParser() *LspParser {
	return &LspParser{
		parser: parser.NewParser(),
		vm:     globalLspVM,
	}
}

// SetVM 设置虚拟机
func (p *LspParser) SetVM(vm *LspVM) {
	p.vm = vm
	vm.parser = p
	if p.parser != nil {
		p.parser.SetVM(vm)
	}
}

// ParseFile 解析文件 - 关键函数
func (p *LspParser) ParseFile(filePath string) (*node.Program, error) {
	if _, err := os.Stat(filePath); err != nil {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}

	// 使用真正的解析器解析文件
	program, err := p.parser.Clone().ParseFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %v", err)
	}

	return program, nil
}

// ParseString 从字符串解析程序 - 用于处理编辑器中的最新内容
func (p *LspParser) ParseString(content string, filePath string) (*node.Program, data.Control) {
	// 调用底层解析器的 ParseString 方法
	return p.parser.ParseString(content, filePath)
}
