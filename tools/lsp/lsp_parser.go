package main

import (
	"fmt"
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
)

// LspParser 是专门为 LSP 服务器设计的解析器
type LspParser struct {
	vm     *LspVM
	parser *parser.Parser
}

// NewLspParser 创建一个新的 LSP 解析器
func NewLspParser() *LspParser {
	p := parser.NewParser()
	return &LspParser{
		parser: p,
		vm:     globalLspVM,
	}
}

// SetVM 设置虚拟机
func (p *LspParser) SetVM(vm *LspVM) {
	p.vm = vm
	if p.parser != nil && vm != nil {
		// 创建一个适配器，让 LspVM 实现 data.VM 接口
		adapter := &LspVMAdapter{vm: vm}
		p.parser.SetVM(adapter)
	}
}

// LspVMAdapter 让 LspVM 实现 data.VM 接口
type LspVMAdapter struct {
	vm *LspVM
}

func (adapter *LspVMAdapter) AddClass(c data.ClassStmt) data.Control {
	if adapter.vm == nil {
		return nil
	}
	// 将 data.ClassStmt 转换为 interface{} 并添加到 LspVM
	adapter.vm.AddClass(c)
	return nil
}

func (adapter *LspVMAdapter) GetClass(pkg string) (data.ClassStmt, bool) {
	if adapter.vm == nil {
		return nil, false
	}
	// 从 LspVM 获取类
	if class, exists := adapter.vm.GetClass(pkg); exists {
		// 尝试转换为 data.ClassStmt
		if classStmt, ok := class.(data.ClassStmt); ok {
			return classStmt, true
		}
	}
	return nil, false
}

func (adapter *LspVMAdapter) AddInterface(i data.InterfaceStmt) data.Control {
	if adapter.vm == nil {
		return nil
	}
	adapter.vm.AddInterface(i)
	return nil
}

func (adapter *LspVMAdapter) GetInterface(pkg string) (data.InterfaceStmt, bool) {
	if adapter.vm == nil {
		return nil, false
	}
	if iface, exists := adapter.vm.GetInterface(pkg); exists {
		if interfaceStmt, ok := iface.(data.InterfaceStmt); ok {
			return interfaceStmt, true
		}
	}
	return nil, false
}

func (adapter *LspVMAdapter) AddFunc(f data.FuncStmt) data.Control {
	if adapter.vm == nil {
		return nil
	}
	adapter.vm.AddFunc(f)
	return nil
}

func (adapter *LspVMAdapter) GetFunc(pkg string) (data.FuncStmt, bool) {
	if adapter.vm == nil {
		return nil, false
	}
	if function, exists := adapter.vm.GetFunc(pkg); exists {
		if funcStmt, ok := function.(data.FuncStmt); ok {
			return funcStmt, true
		}
	}
	return nil, false
}

func (adapter *LspVMAdapter) RegisterFunction(name string, fn interface{}) data.Control {
	return nil
}

func (adapter *LspVMAdapter) RegisterReflectClass(name string, instance interface{}) data.Control {
	return nil
}

func (adapter *LspVMAdapter) CreateContext(vars []data.Variable) data.Context {
	return nil
}

func (adapter *LspVMAdapter) SetThrowControl(fn func(data.Control)) {
}

func (adapter *LspVMAdapter) ThrowControl(acl data.Control) {
}

func (adapter *LspVMAdapter) LoadAndRun(file string) (data.GetValue, data.Control) {
	return nil, nil
}

// ParseFile 解析文件 - 关键函数
func (p *LspParser) ParseFile(filePath string) (interface{}, error) {
	if _, err := os.Stat(filePath); err != nil {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}

	// 使用真正的解析器解析文件
	program, err := p.parser.ParseFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %v", err)
	}

	// 遍历 AST 并提取符号信息
	p.extractSymbolsFromAST(program, filePath)

	return program, nil
}

// extractSymbolsFromAST 从 AST 中提取符号信息
func (p *LspParser) extractSymbolsFromAST(program *node.Program, filePath string) {
	if p.vm == nil {
		return
	}

	// 遍历所有语句
	for _, stmt := range program.Statements {
		p.extractSymbolFromStatement(stmt, filePath)
	}
}

// extractSymbolFromStatement 从语句中提取符号信息
func (p *LspParser) extractSymbolFromStatement(stmt node.Statement, filePath string) {
	switch s := stmt.(type) {
	case *node.FunctionStatement:
		p.extractFunction(s, filePath)
	case *node.ClassStatement:
		p.extractClass(s, filePath)
	case *node.InterfaceStatement:
		p.extractInterface(s, filePath)
	case *node.Namespace:
		// 遍历命名空间内的语句
		for _, nsStmt := range s.Statements {
			p.extractSymbolFromStatement(nsStmt, filePath)
		}
	}
}

// extractFunction 提取函数信息
func (p *LspParser) extractFunction(fn *node.FunctionStatement, filePath string) {
	if p.vm == nil {
		return
	}

	// 直接使用原始节点，不创建 SimpleFunction
	p.vm.AddFunc(fn)
}

// extractClass 提取类信息
func (p *LspParser) extractClass(cls *node.ClassStatement, filePath string) {
	if p.vm == nil {
		return
	}

	// 直接使用原始节点，不创建 SimpleClass
	p.vm.AddClass(cls)

	// 提取类的方法
	methods := cls.GetMethods()
	for _, method := range methods {
		p.extractMethod(method, method.GetName(), cls.GetName(), filePath)
	}
}

// extractInterface 提取接口信息
func (p *LspParser) extractInterface(iface *node.InterfaceStatement, filePath string) {
	if p.vm == nil {
		return
	}

	// 直接使用原始节点，不创建 SimpleInterface
	p.vm.AddInterface(iface)
}

// extractMethod 提取方法信息
func (p *LspParser) extractMethod(method data.Method, methodName string, className, filePath string) {
	if p.vm == nil {
		return
	}

	// 直接使用方法节点，不创建 SimpleFunction
	p.vm.AddFunc(method)
}
