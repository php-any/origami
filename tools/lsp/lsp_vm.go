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

	// 错误处理函数
	throwControl func(data.Control)
}

// NewLspVM 创建一个新的 LSP 虚拟机
func NewLspVM() *LspVM {
	return &LspVM{
		classes:    make(map[string]data.ClassStmt),
		interfaces: make(map[string]data.InterfaceStmt),
		functions:  make(map[string]data.FuncStmt),

		throwControl: func(ctrl data.Control) {
			if ctrl != nil {
				logger.Error("LspVM 错误：%v", ctrl)
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
