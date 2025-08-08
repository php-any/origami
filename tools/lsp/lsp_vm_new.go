package main

import (
	"fmt"
	"sync"
	"time"
)

// LspVM 是专门为 LSP 服务器设计的虚拟机实现
// 它主要用于存储和管理类、函数、接口的节点信息，以支持代码补全、悬停提示等功能
type LspVM struct {
	mu sync.RWMutex

	// 存储类定义，key 为类名
	classes map[string]interface{}
	// 存储接口定义，key 为接口名
	interfaces map[string]interface{}
	// 存储函数定义，key 为函数名
	functions map[string]interface{}

	// 文件到符号的映射，用于快速查找文件中定义的符号
	fileToClasses    map[string][]string // 文件路径 -> 类名列表
	fileToInterfaces map[string][]string // 文件路径 -> 接口名列表
	fileToFunctions  map[string][]string // 文件路径 -> 函数名列表

	// 符号到文件的映射，用于查找符号定义的位置
	classToFile     map[string]string // 类名 -> 文件路径
	interfaceToFile map[string]string // 接口名 -> 文件路径
	functionToFile  map[string]string // 函数名 -> 文件路径

	// 变量跟踪 - 新增
	variables map[string]*LspVariableInfo // 变量名 -> 变量信息

	// 错误处理函数
	throwControl func(interface{})
}

// LspVariableInfo 扩展的变量信息结构
type LspVariableInfo struct {
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
	FilePath  string      `json:"filePath"`
	Line      int         `json:"line"`
	Column    int         `json:"column"`
	Scope     string      `json:"scope"` // global, function, class
	UpdatedAt int64       `json:"updatedAt"`
}

// NewLspVM 创建一个新的 LSP 虚拟机
func NewLspVM() *LspVM {
	return &LspVM{
		classes:          make(map[string]interface{}),
		interfaces:       make(map[string]interface{}),
		functions:        make(map[string]interface{}),
		fileToClasses:    make(map[string][]string),
		fileToInterfaces: make(map[string][]string),
		fileToFunctions:  make(map[string][]string),
		classToFile:      make(map[string]string),
		interfaceToFile:  make(map[string]string),
		functionToFile:   make(map[string]string),
		variables:        make(map[string]*LspVariableInfo), // 新增
		throwControl: func(ctrl interface{}) {
			if ctrl != nil {
				fmt.Printf("LspVM Error: %v\n", ctrl)
			}
		},
	}
}

// AddClass 添加类定义 - 关键函数
func (vm *LspVM) AddClass(c interface{}) interface{} {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	className := vm.getClassName(c)
	if className == "" {
		return fmt.Errorf("类名不能为空")
	}

	vm.classes[className] = c

	// 更新文件映射
	filePath := vm.getSourcePath(c)
	if filePath != "" {
		vm.addClassToFile(filePath, className)
	}

	return nil
}

// GetClass 获取类定义 - 关键函数
func (vm *LspVM) GetClass(className string) (interface{}, bool) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	class, exists := vm.classes[className]
	return class, exists
}

// AddInterface 添加接口定义 - 关键函数
func (vm *LspVM) AddInterface(i interface{}) interface{} {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	interfaceName := vm.getInterfaceName(i)
	if interfaceName == "" {
		return fmt.Errorf("接口名不能为空")
	}

	vm.interfaces[interfaceName] = i

	// 更新文件映射
	filePath := vm.getSourcePath(i)
	if filePath != "" {
		vm.addInterfaceToFile(filePath, interfaceName)
	}

	return nil
}

// GetInterface 获取接口定义 - 关键函数
func (vm *LspVM) GetInterface(interfaceName string) (interface{}, bool) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	iface, exists := vm.interfaces[interfaceName]
	return iface, exists
}

// AddFunc 添加函数定义 - 关键函数
func (vm *LspVM) AddFunc(f interface{}) interface{} {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	funcName := vm.getFunctionName(f)
	if funcName == "" {
		return fmt.Errorf("函数名不能为空")
	}

	vm.functions[funcName] = f

	// 更新文件映射
	filePath := vm.getSourcePath(f)
	if filePath != "" {
		vm.addFunctionToFile(filePath, funcName)
	}

	return nil
}

// GetFunc 获取函数定义 - 关键函数
func (vm *LspVM) GetFunc(funcName string) (interface{}, bool) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	function, exists := vm.functions[funcName]
	return function, exists
}

// AddVariable 添加变量跟踪 - 新增关键函数
func (vm *LspVM) AddVariable(name string, varType string, value interface{}, filePath string, line, column int, scope string) {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	vm.variables[name] = &LspVariableInfo{
		Name:      name,
		Type:      varType,
		Value:     value,
		FilePath:  filePath,
		Line:      line,
		Column:    column,
		Scope:     scope,
		UpdatedAt: getCurrentTimestamp(),
	}
}

// GetVariable 获取变量信息 - 新增关键函数
func (vm *LspVM) GetVariable(name string) (*LspVariableInfo, bool) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	info, exists := vm.variables[name]
	return info, exists
}

// GetAllVariables 获取所有变量 - 新增关键函数
func (vm *LspVM) GetAllVariables() map[string]*LspVariableInfo {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	result := make(map[string]*LspVariableInfo)
	for k, v := range vm.variables {
		result[k] = v
	}
	return result
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

	// 清除变量 - 新增
	for name, info := range vm.variables {
		if info.FilePath == filePath {
			delete(vm.variables, name)
		}
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

	// 查找变量 - 新增
	if info, exists := vm.variables[symbolName]; exists {
		return info.FilePath, "variable", true
	}

	return "", "", false
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

// 辅助方法：获取类名
func (vm *LspVM) getClassName(c interface{}) string {
	if class, ok := c.(*SimpleClass); ok {
		return class.GetName()
	}
	if classWithName, ok := c.(interface{ GetName() string }); ok {
		return classWithName.GetName()
	}
	return ""
}

// 辅助方法：获取接口名
func (vm *LspVM) getInterfaceName(i interface{}) string {
	if iface, ok := i.(*SimpleInterface); ok {
		return iface.GetName()
	}
	if ifaceWithName, ok := i.(interface{ GetName() string }); ok {
		return ifaceWithName.GetName()
	}
	return ""
}

// 辅助方法：获取函数名
func (vm *LspVM) getFunctionName(f interface{}) string {
	if function, ok := f.(*SimpleFunction); ok {
		return function.GetName()
	}
	if funcWithName, ok := f.(interface{ GetName() string }); ok {
		return funcWithName.GetName()
	}
	return ""
}

// 辅助方法：获取源文件路径
func (vm *LspVM) getSourcePath(obj interface{}) string {
	if class, ok := obj.(*SimpleClass); ok {
		return class.GetFilePath()
	}
	if function, ok := obj.(*SimpleFunction); ok {
		return function.GetFilePath()
	}
	if iface, ok := obj.(*SimpleInterface); ok {
		return iface.GetFilePath()
	}
	if objWithPath, ok := obj.(interface{ GetFilePath() string }); ok {
		return objWithPath.GetFilePath()
	}
	return ""
}

// 获取当前时间戳
func getCurrentTimestamp() int64 {
	return int64(time.Now().Unix())
}
