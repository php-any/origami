package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
	"github.com/sirupsen/logrus"
)

// LspVM 是专门为 LSP 服务器设计的虚拟机实现
// 它主要用于存储和管理类、函数、接口的节点信息，以支持代码补全、悬停提示等功能
type LspVM struct {
	mu     sync.RWMutex
	parser *LspParser
	// 存储类定义，key 为类名
	classes map[string]data.ClassStmt
	// 存储接口定义，key 为接口名
	interfaces map[string]data.InterfaceStmt
	// 存储函数定义，key 为函数名
	functions map[string]data.FuncStmt
	// 类解释过程中的缓存, 用于支持循环依赖
	classPathMap map[string]string
	// 错误处理函数
	throwControl func(data.Control)
}

// NewLspVM 创建一个新的 LSP 虚拟机
func NewLspVM() *LspVM {
	return NewLspVMWithScanDir("")
}

// NewLspVMWithScanDir 创建一个新的 LSP 虚拟机并扫描指定目录
func NewLspVMWithScanDir(scanDirectory string) *LspVM {
	vm := &LspVM{
		classes:      make(map[string]data.ClassStmt),
		interfaces:   make(map[string]data.InterfaceStmt),
		functions:    make(map[string]data.FuncStmt),
		classPathMap: make(map[string]string),

		throwControl: func(acl data.Control) {
			if acl != nil {
				logrus.Errorf("LspVM 错误：%v", acl.AsString())

				// 将解析错误转换为诊断信息并发送通知
				if globalConn != nil {
					sendParseErrorDiagnostic(acl)
				}
			}
		},
	}

	// 如果指定了扫描目录，则扫描并解析所有 .zy 文件
	if scanDirectory != "" {
		logrus.Infof("开始扫描目录: %s", scanDirectory)
		vm.scanAndParseDirectory(scanDirectory)
	}

	return vm
}

// AddClass 添加类定义 - 实现 data.VM 接口
func (vm *LspVM) AddClass(c data.ClassStmt) data.Control {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	className := c.GetName()
	if className == "" {
		return utils.NewThrowf("类名不能为空")
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

// GetAllClasses 获取所有类定义
func (vm *LspVM) GetAllClasses() map[string]data.ClassStmt {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	// 创建副本以避免外部修改
	result := make(map[string]data.ClassStmt)
	for className, classStmt := range vm.classes {
		result[className] = classStmt
	}
	return result
}

// AddInterface 添加接口定义 - 实现 data.VM 接口
func (vm *LspVM) AddInterface(i data.InterfaceStmt) data.Control {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	interfaceName := i.GetName()
	if interfaceName == "" {
		return utils.NewThrowf("接口名不能为空")
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
		return utils.NewThrowf("函数名不能为空")
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
	// 解析文件
	_, err := vm.parser.ParseFile(file)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

// scanAndParseDirectory 扫描目录并解析所有 .zy 文件（包括子目录）
func (vm *LspVM) scanAndParseDirectory(directory string) {
	// 添加 panic 恢复机制
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("scanAndParseDirectory 发生 panic：%v", r)
		}
	}()

	logrus.Infof("开始递归扫描目录: %s", directory)

	// 检查目录是否存在
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		logrus.Errorf("目录不存在: %s", directory)
		return
	}

	var fileCount int
	var successCount int
	var dirCount int

	// 使用 filepath.Walk 递归遍历所有子目录
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		// 为每个文件遍历回调添加 panic 恢复
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("遍历文件 %s 时发生 panic：%v", path, r)
			}
		}()

		if err != nil {
			logrus.Warnf("访问文件 %s 时出错: %v", path, err)
			return nil
		}

		// 处理目录
		if info.IsDir() {
			// 跳过隐藏目录和常见的不需要扫描的目录
			if strings.HasPrefix(info.Name(), ".") ||
				info.Name() == "node_modules" ||
				info.Name() == "vendor" ||
				info.Name() == "build" ||
				info.Name() == "dist" {
				logrus.Debugf("跳过目录: %s", path)
				return filepath.SkipDir
			}

			// 记录扫描的目录
			if path != directory {
				dirCount++
				logrus.Debugf("扫描子目录: %s", path)
			}
			return nil
		}

		// 检查文件扩展名
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".zy" || ext == ".php" {
			fileCount++
			logrus.Debugf("发现 .zy 文件: %s", path)

			// 解析文件
			if vm.parseFile(path) {
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

// parseFile 解析单个文件
func (vm *LspVM) parseFile(filePath string) bool {
	// 添加 panic 恢复机制
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("parseFile 发生 panic：%v", r)
		}
	}()

	logrus.Debugf("正在解析文件: %s", filePath)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logrus.Debugf("文件不存在: %s", filePath)
		return false
	}

	// 创建解析器
	parser := NewLspParser()
	parser.SetVM(vm)

	// 解析文件
	_, err := parser.ParseFile(filePath)
	if err != nil {
		logrus.Errorf("解析文件失败 %s: %v", filePath, err)
		return false
	}

	logrus.Debugf("成功解析文件: %s", filePath)
	return true
}

// sendParseErrorDiagnostic 发送解析错误诊断通知
func sendParseErrorDiagnostic(acl data.Control) {
	if globalConn == nil {
		return
	}

	// 尝试从错误控制中提取位置信息
	var uri string
	var range_ Range

	if errorThrow, ok := acl.(*data.ThrowValue); ok {
		if from := errorThrow.GetError().GetFrom(); from != nil {
			// 获取文件路径
			if filePath := from.GetSource(); filePath != "" {
				uri = filePathToURI(filePath)

				// 获取位置范围
				startLine, startCol, endLine, endCol := from.GetRange()
				range_ = Range{
					Start: Position{Line: uint32(startLine - 1), Character: uint32(startCol - 1)},
					End:   Position{Line: uint32(endLine - 1), Character: uint32(endCol - 1)},
				}
			}
		}
	}

	// 如果没有有效的 URI，使用默认值
	if uri == "" {
		uri = "file:///unknown"
	}

	// 创建诊断信息
	diagnostic := Diagnostic{
		Range:    range_,
		Severity: &[]DiagnosticSeverity{DiagnosticSeverityError}[0],
		Message:  acl.AsString(),
	}

	// 发送诊断通知
	params := PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: []Diagnostic{diagnostic},
	}

	globalConn.Notify(context.Background(), "textDocument/publishDiagnostics", params)
	logrus.Infof("已发送解析错误诊断：%#v", params)
}

func (vm *LspVM) SetClassPathCache(name string, path string) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.classPathMap[name] = path
}

func (vm *LspVM) GetClassPathCache(name string) (string, bool) {
	path, ok := vm.classPathMap[name]
	return path, ok
}
