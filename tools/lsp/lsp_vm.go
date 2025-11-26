package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/tools/lsp/defines"
	"github.com/php-any/origami/utils"
	"github.com/sirupsen/logrus"
)

// LspVM 只负责为 LSP 服务缓存和管理符号，不依赖运行时 VM。
type LspVM struct {
	mu     sync.RWMutex
	parser *LspParser

	classes        map[string]data.ClassStmt
	classFiles     map[string]string
	interfaces     map[string]data.InterfaceStmt
	interfaceFiles map[string]string
	functions      map[string]data.FuncStmt
	functionFiles  map[string]string

	fileSymbols    map[string]*fileSymbolSet
	classPathCache map[string]string

	throwControl func(data.Control)
}

// NewLspVM 创建一个新的 LSP 虚拟机。
func NewLspVM() *LspVM {
	return NewLspVMWithScanDir("")
}

// NewLspVMWithScanDir 创建一个新的 LSP 虚拟机并扫描指定目录。
func NewLspVMWithScanDir(scanDirectory string) *LspVM {
	vm := &LspVM{
		classes:        make(map[string]data.ClassStmt),
		classFiles:     make(map[string]string),
		interfaces:     make(map[string]data.InterfaceStmt),
		interfaceFiles: make(map[string]string),
		functions:      make(map[string]data.FuncStmt),
		functionFiles:  make(map[string]string),
		fileSymbols:    make(map[string]*fileSymbolSet),
		classPathCache: make(map[string]string),
		throwControl: func(acl data.Control) {
			if acl != nil {
				logrus.Errorf("LspVM 错误：%v", acl.AsString())
				if globalConn != nil {
					sendParseErrorDiagnostic(acl)
				}
			}
		},
	}

	if scanDirectory != "" {
		vm.scanAndParseDirectory(scanDirectory)
	}

	return vm
}

// AddClass 添加或覆盖类定义。
func (vm *LspVM) AddClass(c data.ClassStmt) data.Control {
	if c == nil {
		return utils.NewThrowf("类定义不能为空")
	}

	className := c.GetName()
	if className == "" {
		return utils.NewThrowf("类名不能为空")
	}

	filePath := extractSourcePath(c.GetFrom())

	vm.mu.Lock()
	defer vm.mu.Unlock()

	if prevPath, ok := vm.classFiles[className]; ok && prevPath != "" && prevPath != filePath {
		vm.removeSymbolFromFileLocked(symbolKindClass, prevPath, className)
	}

	vm.classes[className] = c
	vm.classFiles[className] = filePath
	vm.addSymbolToFileLocked(symbolKindClass, filePath, className)

	if filePath != "" {
		vm.classPathCache[className] = filePath
	}
	return nil
}

// GetClass 获取类定义。
func (vm *LspVM) GetClass(className string) (data.ClassStmt, bool) {
	vm.mu.RLock()
	class, exists := vm.classes[className]
	vm.mu.RUnlock()
	return class, exists
}

// GetOrLoadClass 尝试获取或从磁盘加载类定义，若失败返回占位类。
func (vm *LspVM) GetOrLoadClass(pkg string) (data.ClassStmt, data.Control) {
	if class, ok := vm.GetClass(pkg); ok {
		return class, nil
	}

	if err := vm.loadClassFromCache(pkg); err == nil {
		if class, ok := vm.GetClass(pkg); ok {
			return class, nil
		}
	} else {
		logrus.Debugf("从缓存加载类 %s 失败：%v", pkg, err)
	}

	placeholder := newPlaceholderClass(pkg)
	vm.mu.Lock()
	vm.classes[pkg] = placeholder
	vm.mu.Unlock()
	return placeholder, nil
}

// LoadPkg 尝试返回类或接口。
func (vm *LspVM) LoadPkg(pkg string) (data.GetValue, data.Control) {
	if class, ok := vm.GetClass(pkg); ok {
		return class, nil
	}
	if iface, ok := vm.GetInterface(pkg); ok {
		return iface, nil
	}
	return vm.GetOrLoadClass(pkg)
}

// GetAllClasses 返回当前类的副本。
func (vm *LspVM) GetAllClasses() map[string]data.ClassStmt {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	result := make(map[string]data.ClassStmt, len(vm.classes))
	for k, v := range vm.classes {
		result[k] = v
	}
	return result
}

// AddInterface 添加或覆盖接口。
func (vm *LspVM) AddInterface(i data.InterfaceStmt) data.Control {
	if i == nil {
		return utils.NewThrowf("接口定义不能为空")
	}

	name := i.GetName()
	if name == "" {
		return utils.NewThrowf("接口名不能为空")
	}

	filePath := extractSourcePath(i.GetFrom())

	vm.mu.Lock()
	defer vm.mu.Unlock()

	if prevPath, ok := vm.interfaceFiles[name]; ok && prevPath != "" && prevPath != filePath {
		vm.removeSymbolFromFileLocked(symbolKindInterface, prevPath, name)
	}

	vm.interfaces[name] = i
	vm.interfaceFiles[name] = filePath
	vm.addSymbolToFileLocked(symbolKindInterface, filePath, name)
	return nil
}

// GetInterface 获取接口。
func (vm *LspVM) GetInterface(interfaceName string) (data.InterfaceStmt, bool) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	iface, exists := vm.interfaces[interfaceName]
	return iface, exists
}

// AddFunc 添加或覆盖函数。
func (vm *LspVM) AddFunc(f data.FuncStmt) data.Control {
	if f == nil {
		return utils.NewThrowf("函数定义不能为空")
	}

	funcName := f.GetName()
	if funcName == "" {
		return utils.NewThrowf("函数名不能为空")
	}

	filePath := extractFuncSource(f)

	vm.mu.Lock()
	defer vm.mu.Unlock()

	if prevPath, ok := vm.functionFiles[funcName]; ok && prevPath != "" && prevPath != filePath {
		vm.removeSymbolFromFileLocked(symbolKindFunction, prevPath, funcName)
	}

	vm.functions[funcName] = f
	vm.functionFiles[funcName] = filePath
	vm.addSymbolToFileLocked(symbolKindFunction, filePath, funcName)
	return nil
}

// GetFunc 获取函数。
func (vm *LspVM) GetFunc(funcName string) (data.FuncStmt, bool) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	fn, exists := vm.functions[funcName]
	return fn, exists
}

// ClearFile 清除文件中的符号。
func (vm *LspVM) ClearFile(filePath string) {
	normalized := normalizeFilePath(filePath)
	if normalized == "" {
		return
	}

	vm.mu.Lock()
	defer vm.mu.Unlock()

	symbols, exists := vm.fileSymbols[normalized]
	if !exists {
		return
	}

	for className := range symbols.classes {
		delete(vm.classes, className)
		delete(vm.classFiles, className)
	}
	for iface := range symbols.interfaces {
		delete(vm.interfaces, iface)
		delete(vm.interfaceFiles, iface)
	}
	for funcName := range symbols.functions {
		delete(vm.functions, funcName)
		delete(vm.functionFiles, funcName)
	}

	delete(vm.fileSymbols, normalized)
}

// RegisterFunction 仅为满足接口，直接返回。
func (vm *LspVM) RegisterFunction(name string, fn interface{}) data.Control {
	return nil
}

// RegisterReflectClass 仅为满足接口。
func (vm *LspVM) RegisterReflectClass(name string, instance interface{}) data.Control {
	return nil
}

// CreateContext LSP 模式无需上下文。
func (vm *LspVM) CreateContext(vars []data.Variable) data.Context {
	return runtime.NewContext(vm).CreateContext(vars)
}

// SetThrowControl 设置异常回调。
func (vm *LspVM) SetThrowControl(fn func(data.Control)) {
	vm.throwControl = fn
}

// ThrowControl 调用异常回调。
func (vm *LspVM) ThrowControl(acl data.Control) {
	if vm.throwControl != nil {
		vm.throwControl(acl)
	}
}

// LoadAndRun LSP 不执行代码。
func (vm *LspVM) LoadAndRun(file string) (data.GetValue, data.Control) {
	return nil, nil
}

// ParseFile LSP 不执行解释。
func (vm *LspVM) ParseFile(file string, value data.Value) (data.Value, data.Control) {
	return nil, nil
}

// scanAndParseDirectory 扫描目录并解析 .zy/.php 文件。
func (vm *LspVM) scanAndParseDirectory(directory string) {
	root := normalizeFilePath(directory)
	if root == "" {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("scanAndParseDirectory 发生 panic：%v", r)
		}
	}()

	if _, err := os.Stat(root); os.IsNotExist(err) {
		logrus.Warnf("目录不存在: %s", root)
		return
	}

	var fileCount, successCount, dirCount int
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Warnf("访问文件 %s 时出错: %v", path, err)
			return nil
		}

		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") ||
				info.Name() == "node_modules" ||
				info.Name() == "vendor" ||
				info.Name() == "build" ||
				info.Name() == "dist" {
				return filepath.SkipDir
			}
			if path != root {
				dirCount++
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".zy" || ext == ".php" {
			fileCount++
			if vm.parseFile(path) == nil {
				successCount++
			}
		}
		return nil
	})

	if err != nil {
		logrus.Errorf("遍历目录失败: %v", err)
		return
	}

	logrus.Infof("递归扫描完成: 扫描了 %d 个目录，发现 %d 个 .zy 文件，成功解析 %d 个", dirCount, fileCount, successCount)
}

// parseFile 解析单个文件。
func (vm *LspVM) parseFile(filePath string) data.Control {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("parseFile 发生 panic：%v", r)
		}
	}()

	normalized := normalizeFilePath(filePath)
	if normalized == "" {
		return utils.NewThrowf("无法解析空文件路径")
	}

	if _, err := os.Stat(normalized); os.IsNotExist(err) {
		logrus.Debugf("文件不存在: %s", normalized)
		return nil
	}

	vm.ClearFile(normalized)

	parser := NewLspParser()
	parser.SetVM(vm)
	if _, acl := parser.ParseFile(normalized); acl != nil {
		logrus.Errorf("解析文件失败 %s: %v", normalized, acl)
		return acl
	}

	return nil
}

func sendParseErrorDiagnostic(acl data.Control) {
	if globalConn == nil {
		return
	}

	var uri string
	var range_ defines.Range

	if errorThrow, ok := acl.(*data.ThrowValue); ok {
		if from := errorThrow.GetError().GetFrom(); from != nil {
			if filePath := from.GetSource(); filePath != "" {
				uri = filePathToURI(filePath)
				startLine, startCol, endLine, endCol := from.GetRange()
				range_ = defines.Range{
					Start: defines.Position{Line: uint32(startLine - 1), Character: uint32(startCol - 1)},
					End:   defines.Position{Line: uint32(endLine - 1), Character: uint32(endCol - 1)},
				}
			}
		}
	}

	if uri == "" {
		uri = "file:///unknown"
	}

	diagnostic := defines.Diagnostic{
		Range:    range_,
		Severity: &[]defines.DiagnosticSeverity{defines.DiagnosticSeverityError}[0],
		Message:  acl.AsString(),
	}

	params := defines.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: []defines.Diagnostic{diagnostic},
	}

	globalConn.Notify(context.Background(), "textDocument/publishDiagnostics", params)
	logrus.Infof("已发送解析错误诊断：%#v", params)
}

func (vm *LspVM) SetClassPathCache(name string, path string) {
	normalized := normalizeFilePath(path)
	if normalized == "" {
		return
	}

	vm.mu.Lock()
	vm.classPathCache[name] = normalized
	vm.mu.Unlock()
}

func (vm *LspVM) GetClassPathCache(name string) (string, bool) {
	vm.mu.RLock()
	path, ok := vm.classPathCache[name]
	vm.mu.RUnlock()
	return path, ok
}

func (vm *LspVM) loadClassFromCache(name string) data.Control {
	vm.mu.RLock()
	path, ok := vm.classPathCache[name]
	vm.mu.RUnlock()
	if !ok || path == "" {
		return utils.NewThrowf("class %s not found in cache", name)
	}
	return vm.parseFile(path)
}

type symbolKind int

const (
	symbolKindClass symbolKind = iota
	symbolKindInterface
	symbolKindFunction
)

type fileSymbolSet struct {
	classes    map[string]struct{}
	interfaces map[string]struct{}
	functions  map[string]struct{}
}

func newFileSymbolSet() *fileSymbolSet {
	return &fileSymbolSet{
		classes:    make(map[string]struct{}),
		interfaces: make(map[string]struct{}),
		functions:  make(map[string]struct{}),
	}
}

func (vm *LspVM) addSymbolToFileLocked(kind symbolKind, filePath, name string) {
	if filePath == "" || name == "" {
		return
	}
	set := vm.ensureFileSymbolSetLocked(filePath)
	switch kind {
	case symbolKindClass:
		set.classes[name] = struct{}{}
	case symbolKindInterface:
		set.interfaces[name] = struct{}{}
	case symbolKindFunction:
		set.functions[name] = struct{}{}
	}
}

func (vm *LspVM) removeSymbolFromFileLocked(kind symbolKind, filePath, name string) {
	if filePath == "" || name == "" {
		return
	}

	set, ok := vm.fileSymbols[filePath]
	if !ok {
		return
	}

	switch kind {
	case symbolKindClass:
		delete(set.classes, name)
	case symbolKindInterface:
		delete(set.interfaces, name)
	case symbolKindFunction:
		delete(set.functions, name)
	}

	if len(set.classes) == 0 && len(set.interfaces) == 0 && len(set.functions) == 0 {
		delete(vm.fileSymbols, filePath)
	}
}

func (vm *LspVM) ensureFileSymbolSetLocked(filePath string) *fileSymbolSet {
	if set, ok := vm.fileSymbols[filePath]; ok {
		return set
	}
	set := newFileSymbolSet()
	vm.fileSymbols[filePath] = set
	return set
}

func extractSourcePath(from data.From) string {
	if from == nil {
		return ""
	}
	return normalizeFilePath(from.GetSource())
}

func extractFuncSource(f data.FuncStmt) string {
	type fromProvider interface {
		GetFrom() data.From
	}
	if fp, ok := f.(fromProvider); ok {
		return extractSourcePath(fp.GetFrom())
	}
	return ""
}

func normalizeFilePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	cleaned := filepath.Clean(path)
	if abs, err := filepath.Abs(cleaned); err == nil {
		return abs
	}
	return cleaned
}

type placeholderClass struct {
	name string
	from data.From
}

func newPlaceholderClass(name string) data.ClassStmt {
	if name == "" {
		name = "UnknownClass"
	}
	return &placeholderClass{
		name: name,
		from: data.NewBaseFrom("placeholder://"+name, 0, 0),
	}
}

func (p *placeholderClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(p, ctx), nil
}
func (p *placeholderClass) GetFrom() data.From                            { return p.from }
func (p *placeholderClass) GetName() string                               { return p.name }
func (p *placeholderClass) GetExtend() *string                            { return nil }
func (p *placeholderClass) GetImplements() []string                       { return nil }
func (p *placeholderClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (p *placeholderClass) GetPropertyList() []data.Property              { return nil }
func (p *placeholderClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (p *placeholderClass) GetMethods() []data.Method                     { return nil }
func (p *placeholderClass) GetConstruct() data.Method                     { return nil }
