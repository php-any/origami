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

	// 获取函数位置信息
	from := fn.GetFrom()
	line := 1
	if from != nil {
		start, _ := from.GetPosition()
		// 简化的行号计算
		line = start/100 + 1
	}

	function := &SimpleFunction{
		Name:     fn.GetName(),
		FilePath: filePath,
		Line:     line,
		Content:  fn.GetName(),
	}

	p.vm.AddFunc(function)
}

// extractClass 提取类信息
func (p *LspParser) extractClass(cls *node.ClassStatement, filePath string) {
	if p.vm == nil {
		return
	}

	// 获取类位置信息
	from := cls.GetFrom()
	line := 1
	if from != nil {
		start, _ := from.GetPosition()
		// 简化的行号计算
		line = start/100 + 1
	}

	class := &SimpleClass{
		Name:     cls.GetName(),
		FilePath: filePath,
		Line:     line,
		Content:  cls.GetName(),
	}

	p.vm.AddClass(class)

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

	// 获取接口位置信息
	from := iface.GetFrom()
	line := 1
	if from != nil {
		start, _ := from.GetPosition()
		// 简化的行号计算
		line = start/100 + 1
	}

	interfaceInfo := &SimpleInterface{
		Name:     iface.GetName(),
		FilePath: filePath,
		Line:     line,
		Content:  iface.GetName(),
	}

	p.vm.AddInterface(interfaceInfo)
}

// extractMethod 提取方法信息
func (p *LspParser) extractMethod(method data.Method, methodName string, className, filePath string) {
	if p.vm == nil {
		return
	}

	// 将方法作为函数添加到 VM（简化处理）
	function := &SimpleFunction{
		Name:     methodName,
		FilePath: filePath,
		Line:     1, // 简化处理
		Content:  methodName,
	}

	p.vm.AddFunc(function)
}

// SimpleAST 是简化的 AST 结构
type SimpleAST struct {
	FilePath   string
	Content    string
	Classes    []*SimpleClass
	Functions  []*SimpleFunction
	Interfaces []*SimpleInterface
}

// SimpleClass 是简化的类结构
type SimpleClass struct {
	Name     string
	FilePath string
	Line     int
	Content  string
}

// GetName 获取类名
func (c *SimpleClass) GetName() string {
	return c.Name
}

// GetFilePath 获取文件路径
func (c *SimpleClass) GetFilePath() string {
	return c.FilePath
}

// GetLine 获取行号
func (c *SimpleClass) GetLine() int {
	return c.Line
}

// GetContent 获取内容
func (c *SimpleClass) GetContent() string {
	return c.Content
}

// SimpleFunction 是简化的函数结构
type SimpleFunction struct {
	Name     string
	FilePath string
	Line     int
	Content  string
}

// GetName 获取函数名
func (f *SimpleFunction) GetName() string {
	return f.Name
}

// GetFilePath 获取文件路径
func (f *SimpleFunction) GetFilePath() string {
	return f.FilePath
}

// GetLine 获取行号
func (f *SimpleFunction) GetLine() int {
	return f.Line
}

// GetContent 获取内容
func (f *SimpleFunction) GetContent() string {
	return f.Content
}

// SimpleInterface 是简化的接口结构
type SimpleInterface struct {
	Name     string
	FilePath string
	Line     int
	Content  string
}

// GetName 获取接口名
func (i *SimpleInterface) GetName() string {
	return i.Name
}

// GetFilePath 获取文件路径
func (i *SimpleInterface) GetFilePath() string {
	return i.FilePath
}

// GetLine 获取行号
func (i *SimpleInterface) GetLine() int {
	return i.Line
}

// GetContent 获取内容
func (i *SimpleInterface) GetContent() string {
	return i.Content
}
