package main

import (
	"os"

	"github.com/php-any/origami/data"

	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/utils"
)

// LspScopeFactory LSP 作用域工厂函数
// 这个函数替换默认的作用域创建逻辑，在 LSP 模式下使用 LspScope
func LspScopeFactory(parent parser.Scope, isLambda bool) parser.Scope {
	// 生成作用域名称和类型
	scopeName := "scope"
	scopeType := "block"
	if isLambda {
		scopeType = "lambda"
	}

	// 创建 LspScope 实例
	scope := NewLspScope(parent, scopeName, scopeType, "")
	scope.SetLambda(isLambda)

	return scope
}

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
	// 设置 LSP 作用域工厂函数
	parser.SetGlobalScopeFactory(LspScopeFactory)

	return &LspParser{
		parser: parser.NewParser(),
		vm:     globalLspVM,
	}
}

// CreateLspScope 创建 LSP 作用域
// 这个方法可以在需要时替换默认的作用域创建逻辑
func (p *LspParser) CreateLspScope(parent parser.Scope, isLambda bool, scopeName, scopeType, filePath string) *LspScope {
	scope := NewLspScope(parent, scopeName, scopeType, filePath)
	scope.SetLambda(isLambda)
	return scope
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
func (p *LspParser) ParseFile(filePath string) (*node.Program, data.Control) {
	if _, err := os.Stat(filePath); err != nil {
		return nil, utils.NewThrowf("file does not exist: %s", filePath)
	}

	// 使用真正的解析器解析文件
	program, acl := p.parser.Clone().ParseFile(filePath)
	if acl != nil {
		return nil, acl
	}

	return program, nil
}

// ParseString 从字符串解析程序 - 用于处理编辑器中的最新内容
func (p *LspParser) ParseString(content string, filePath string) (*node.Program, data.Control) {
	// 调用底层解析器的 ParseString 方法
	return p.parser.ParseString(content, filePath)
}
