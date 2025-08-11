package main

import (
	"fmt"
	"sync"

	"github.com/php-any/origami/data"
)

// LspVM 是专门为 LSP 服务器设计的虚拟机实现
// 它主要用于存储和管理类、函数、接口的节点信息，以支持代码补全、悬停提示等功能
type LspVM struct {
	mu sync.RWMutex

	// 存储类定义，key 为类名
	classes map[string]data.ClassStmt
	// 存储接口定义，key 为接口名
	interfaces map[string]data.InterfaceStmt
	// 存储函数定义，key 为函数名
	functions map[string]data.FuncStmt

	// 文件到符号的映射，用于快速查找文件中定义的符号
	fileToClasses    map[string][]string // 文件路径 -> 类名列表
	fileToInterfaces map[string][]string // 文件路径 -> 接口名列表
	fileToFunctions  map[string][]string // 文件路径 -> 函数名列表

	// 符号到文件的映射，用于查找符号定义的位置
	classToFile     map[string]string // 类名 -> 文件路径
	interfaceToFile map[string]string // 接口名 -> 文件路径
	functionToFile  map[string]string // 函数名 -> 文件路径

	// 错误处理函数
	throwControl func(data.Control)
}

// NewLspVM 创建一个新的 LSP 虚拟机
func NewLspVM() *LspVM {
	return &LspVM{
		classes:          make(map[string]data.ClassStmt),
		interfaces:       make(map[string]data.InterfaceStmt),
		functions:        make(map[string]data.FuncStmt),
		fileToClasses:    make(map[string][]string),
		fileToInterfaces: make(map[string][]string),
		fileToFunctions:  make(map[string][]string),
		classToFile:      make(map[string]string),
		interfaceToFile:  make(map[string]string),
		functionToFile:   make(map[string]string),
		throwControl: func(ctrl data.Control) {
			if ctrl != nil {
				fmt.Printf("LspVM Error: %v\n", ctrl)
			}
		},
	}
}

// AddClass 添加类定义 - 实现 data.VM 接口
func (vm *LspVM) AddClass(c data.ClassStmt) data.Control {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	className := c.GetName()
	if className == "" {
		return data.NewErrorThrow(nil, fmt.Errorf("类名不能为空"))
	}

	vm.classes[className] = c

	// 更新文件映射
	filePath := vm.getSourcePath(c)
	if filePath != "" {
		vm.addClassToFile(filePath, className)
	}

	return nil
}

// GetClass 获取类定义 - 实现 data.VM 接口
func (vm *LspVM) GetClass(className string) (data.ClassStmt, bool) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	class, exists := vm.classes[className]
	return class, exists
}

// AddInterface 添加接口定义 - 实现 data.VM 接口
func (vm *LspVM) AddInterface(i data.InterfaceStmt) data.Control {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	interfaceName := i.GetName()
	if interfaceName == "" {
		return data.NewErrorThrow(nil, fmt.Errorf("接口名不能为空"))
	}

	vm.interfaces[interfaceName] = i

	// 更新文件映射
	filePath := vm.getSourcePath(i)
	if filePath != "" {
		vm.addInterfaceToFile(filePath, interfaceName)
	}

	return nil
}

// GetInterface 获取接口定义 - 实现 data.VM 接口
func (vm *LspVM) GetInterface(interfaceName string) (data.InterfaceStmt, bool) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	iface, exists := vm.interfaces[interfaceName]
	return iface, exists
}

// AddFunc 添加函数定义 - 实现 data.VM 接口
func (vm *LspVM) AddFunc(f data.FuncStmt) data.Control {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	funcName := f.GetName()
	if funcName == "" {
		return data.NewErrorThrow(nil, fmt.Errorf("函数名不能为空"))
	}

	vm.functions[funcName] = f

	// 更新文件映射
	filePath := vm.getSourcePath(f)
	if filePath != "" {
		vm.addFunctionToFile(filePath, funcName)
	}

	return nil
}

// GetFunc 获取函数定义 - 实现 data.VM 接口
func (vm *LspVM) GetFunc(funcName string) (data.FuncStmt, bool) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	function, exists := vm.functions[funcName]
	return function, exists
}

// ClearFile 清除文件中的符号 - 关键函数
func (vm *LspVM) ClearFile(filePath string) {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	// 清除类
	if classNames, exists := vm.fileToClasses[filePath]; exists {
		for _, className := range classNames {
			delete(vm.classes, className)
			delete(vm.classToFile, className)
		}
		delete(vm.fileToClasses, filePath)
	}

	// 清除接口
	if interfaceNames, exists := vm.fileToInterfaces[filePath]; exists {
		for _, interfaceName := range interfaceNames {
			delete(vm.interfaces, interfaceName)
			delete(vm.interfaceToFile, interfaceName)
		}
		delete(vm.fileToInterfaces, filePath)
	}

	// 清除函数
	if functionNames, exists := vm.fileToFunctions[filePath]; exists {
		for _, functionName := range functionNames {
			delete(vm.functions, functionName)
			delete(vm.functionToFile, functionName)
		}
		delete(vm.fileToFunctions, filePath)
	}

}

// FindSymbolFile 查找符号所在文件 - 关键函数
func (vm *LspVM) FindSymbolFile(symbolName string) (filePath string, symbolType string, found bool) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	// 查找类
	if filePath, exists := vm.classToFile[symbolName]; exists {
		return filePath, "class", true
	}

	// 查找接口
	if filePath, exists := vm.interfaceToFile[symbolName]; exists {
		return filePath, "interface", true
	}

	// 查找函数
	if filePath, exists := vm.functionToFile[symbolName]; exists {
		return filePath, "function", true
	}

	return "", "", false
}

// RegisterFunction 注册函数 - 实现 data.VM 接口
func (vm *LspVM) RegisterFunction(name string, fn interface{}) data.Control {
	// LSP VM 不需要实现这个功能，返回 nil
	return nil
}

// RegisterReflectClass 注册反射类 - 实现 data.VM 接口
func (vm *LspVM) RegisterReflectClass(name string, instance interface{}) data.Control {
	// LSP VM 不需要实现这个功能，返回 nil
	return nil
}

// CreateContext 创建上下文 - 实现 data.VM 接口
func (vm *LspVM) CreateContext(vars []data.Variable) data.Context {
	// LSP VM 不需要实现这个功能，返回 nil
	return nil
}

// SetThrowControl 设置异常控制函数 - 实现 data.VM 接口
func (vm *LspVM) SetThrowControl(fn func(data.Control)) {
	vm.throwControl = fn
}

// ThrowControl 抛出异常控制 - 实现 data.VM 接口
func (vm *LspVM) ThrowControl(acl data.Control) {
	if vm.throwControl != nil {
		vm.throwControl(acl)
	}
}

// LoadAndRun 加载并运行文件 - 实现 data.VM 接口
func (vm *LspVM) LoadAndRun(file string) (data.GetValue, data.Control) {
	// LSP VM 不需要实现这个功能，返回 nil
	return nil, nil
}

// 辅助方法：添加类到文件映射
func (vm *LspVM) addClassToFile(filePath, className string) {
	if !vm.containsString(vm.fileToClasses[filePath], className) {
		vm.fileToClasses[filePath] = append(vm.fileToClasses[filePath], className)
	}
	vm.classToFile[className] = filePath
}

// 辅助方法：添加接口到文件映射
func (vm *LspVM) addInterfaceToFile(filePath, interfaceName string) {
	if !vm.containsString(vm.fileToInterfaces[filePath], interfaceName) {
		vm.fileToInterfaces[filePath] = append(vm.fileToInterfaces[filePath], interfaceName)
	}
	vm.interfaceToFile[interfaceName] = filePath
}

// 辅助方法：添加函数到文件映射
func (vm *LspVM) addFunctionToFile(filePath, funcName string) {
	if !vm.containsString(vm.fileToFunctions[filePath], funcName) {
		vm.fileToFunctions[filePath] = append(vm.fileToFunctions[filePath], funcName)
	}
	vm.functionToFile[funcName] = filePath
}

// 辅助方法：检查字符串是否在切片中
func (vm *LspVM) containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// 辅助方法：获取源文件路径
func (vm *LspVM) getSourcePath(obj interface{}) string {
	// 尝试从原始节点获取源文件路径
	if nodeWithFrom, ok := obj.(interface{ GetFrom() data.From }); ok {
		from := nodeWithFrom.GetFrom()
		if from != nil {
			return from.GetSource()
		}
	}

	// 尝试从其他类型获取文件路径
	if objWithPath, ok := obj.(interface{ GetFilePath() string }); ok {
		return objWithPath.GetFilePath()
	}
	return ""
}
