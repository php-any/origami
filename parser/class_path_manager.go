package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/php-any/origami/data"
)

// ClassPathManager 类路径管理器接口
type ClassPathManager interface {
	// AddNamespace 添加命名空间路径
	AddNamespace(namespace string, path string)
	// FindClassFile 查找类文件路径
	FindClassFile(className string) (string, bool)
	// LoadClass 加载类
	LoadClass(className string, parser *Parser) data.Control
}

// NamespaceNode 命名空间节点
type NamespaceNode struct {
	namespace string
	paths     []string // 支持多个路径
	children  map[string]*NamespaceNode
}

// DefaultClassPathManager 默认的类路径管理器实现
type DefaultClassPathManager struct {
	mu   sync.RWMutex   // 互斥锁，保护并发访问
	root *NamespaceNode // 有向无环图的根节点
}

// NewDefaultClassPathManager 创建默认的类路径管理器
func NewDefaultClassPathManager() *DefaultClassPathManager {
	return &DefaultClassPathManager{
		root: &NamespaceNode{
			namespace: "",
			paths:     make([]string, 0),
			children:  make(map[string]*NamespaceNode),
		},
	}
}

// AddNamespace 添加命名空间路径
func (m *DefaultClassPathManager) AddNamespace(namespace string, path string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查路径是否为空
	if path == "" {
		fmt.Printf("警告: 命名空间 %s 的路径为空\n", namespace)
		return
	}

	// 检查路径是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("警告: 命名空间 %s 的路径不存在: %s\n", namespace, path)
		return
	}

	// 转换为绝对路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("警告: 无法获取命名空间 %s 路径的绝对路径: %s, 错误: %v\n", namespace, path, err)
		return
	}

	// 检查转换后的绝对路径是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Printf("警告: 命名空间 %s 的绝对路径不存在: %s\n", namespace, absPath)
		return
	}

	// 将命名空间添加到DAG中
	m.addNamespaceToDAG(namespace, absPath)
}

// addNamespaceToDAG 将命名空间添加到有向无环图中
func (m *DefaultClassPathManager) addNamespaceToDAG(namespace string, path string) {
	// 分割命名空间
	parts := m.splitNamespace(namespace)

	// 从根节点开始构建路径
	current := m.root

	for i, part := range parts {
		if child, exists := current.children[part]; exists {
			current = child
		} else {
			// 创建新节点
			child = &NamespaceNode{
				namespace: part,
				paths:     make([]string, 0),
				children:  make(map[string]*NamespaceNode),
			}
			current.children[part] = child
			current = child
		}

		// 如果是最后一个部分，添加路径（支持多个路径）
		if i == len(parts)-1 {
			// 检查路径是否已经存在，避免重复
			exists := false
			for _, existingPath := range current.paths {
				if existingPath == path {
					exists = true
					break
				}
			}
			if !exists {
				current.paths = append(current.paths, path)
			}
		}
	}
}

// splitNamespace 分割命名空间
func (m *DefaultClassPathManager) splitNamespace(namespace string) []string {
	if namespace == "" {
		return []string{}
	}

	// 使用反斜杠分割命名空间
	parts := make([]string, 0)
	start := 0

	for i, char := range namespace {
		if char == '\\' {
			if i > start {
				parts = append(parts, namespace[start:i])
			}
			start = i + 1
		}
	}

	// 添加最后一部分
	if start < len(namespace) {
		parts = append(parts, namespace[start:])
	}

	return parts
}

// FindClassFile 查找类文件路径
func (m *DefaultClassPathManager) FindClassFile(className string) (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 分割类名以获取命名空间和类名
	namespace, simpleClassName := m.splitClassName(className)

	// 在DAG中查找对应的命名空间节点
	node := m.findNamespaceNode(namespace)
	if node == nil {
		return "", false
	}

	// 在所有路径中搜索类文件
	for _, path := range node.paths {
		// 方法1：直接在路径下查找类文件
		filePath, found := m.searchInPath(path, simpleClassName)
		if found {
			return filePath, true
		}

		// 方法2：在子目录中查找类文件
		// 如果类名包含子目录结构，尝试在子目录中查找
		if strings.Contains(className, "\\") {
			// 将类名中的反斜杠转换为路径分隔符
			relativePath := filepath.FromSlash(className)
			filePath, found := m.searchInPath(path, relativePath)
			if found {
				return filePath, true
			}
		}
	}
	// TODO 忽略大小写方式搜索查找

	return "", false
}

// searchInPath 在指定路径下搜索类文件
func (m *DefaultClassPathManager) searchInPath(basePath, className string) (string, bool) {
	// 构造可能的文件名
	possibleFiles := []string{
		className + ".zy",
		className + ".php",
	}

	// 首先尝试精确匹配
	for _, fileName := range possibleFiles {
		filePath := filepath.Join(basePath, fileName)
		if _, err := os.Stat(filePath); err == nil {
			return filePath, true
		}
	}

	// 如果精确匹配失败，尝试大小写不敏感的匹配（用于跨平台兼容）
	for _, fileName := range possibleFiles {
		if foundPath := m.findFileCaseInsensitive(basePath, fileName); foundPath != "" {
			return foundPath, true
		}
	}

	return "", false
}

// splitClassName 分割类名，返回命名空间和简单类名
func (m *DefaultClassPathManager) splitClassName(className string) (string, string) {
	lastSeparator := -1
	for i, char := range className {
		if char == '\\' {
			lastSeparator = i
		}
	}

	if lastSeparator == -1 {
		// 没有命名空间
		return "", className
	}

	namespace := className[:lastSeparator]
	simpleClassName := className[lastSeparator+1:]
	return namespace, simpleClassName
}

// findNamespaceNode 在DAG中查找命名空间节点
func (m *DefaultClassPathManager) findNamespaceNode(namespace string) *NamespaceNode {
	if namespace == "" {
		return m.root
	}

	parts := m.splitNamespace(namespace)
	current := m.root

	for _, part := range parts {
		if child, exists := current.children[part]; exists {
			current = child
		} else {
			// 根据current值，直接拼接完整路径，识别目录是否存在
			// 检查当前节点的路径配置，看是否可以动态创建子目录
			found := false
			for _, path := range current.paths {
				// 构建可能的目录路径
				dirPath := filepath.Join(path, part)
				if info, err := os.Stat(dirPath); err == nil && info.IsDir() {
					// 目录存在，创建一个新的节点
					child := &NamespaceNode{
						namespace: part,
						paths:     []string{dirPath},
						children:  make(map[string]*NamespaceNode),
					}
					current.children[part] = child
					current = child
					found = true
					continue
				}
				// 如果精确匹配失败，尝试大小写不敏感的匹配（用于跨平台兼容）
				if foundDir := m.findDirectoryCaseInsensitive(path, part); foundDir != "" {
					child := &NamespaceNode{
						namespace: part,
						paths:     []string{foundDir},
						children:  make(map[string]*NamespaceNode),
					}
					current.children[part] = child
					current = child
					found = true
					continue
				}
			}
			if !found {
				return nil
			}
		}
	}

	return current
}

// findDirectoryCaseInsensitive 在指定路径下大小写不敏感地查找目录
func (m *DefaultClassPathManager) findDirectoryCaseInsensitive(basePath, dirName string) string {
	// 读取目录内容
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return ""
	}

	// 转换为小写进行比较
	lowerDirName := strings.ToLower(dirName)
	for _, entry := range entries {
		if entry.IsDir() && strings.ToLower(entry.Name()) == lowerDirName {
			foundPath := filepath.Join(basePath, entry.Name())
			if info, err := os.Stat(foundPath); err == nil && info.IsDir() {
				return foundPath
			}
		}
	}

	return ""
}

// findFileCaseInsensitive 在指定路径下大小写不敏感地查找文件
func (m *DefaultClassPathManager) findFileCaseInsensitive(basePath, fileName string) string {
	// 读取目录内容
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return ""
	}

	// 转换为小写进行比较
	lowerFileName := strings.ToLower(fileName)
	for _, entry := range entries {
		if !entry.IsDir() && strings.ToLower(entry.Name()) == lowerFileName {
			foundPath := filepath.Join(basePath, entry.Name())
			if info, err := os.Stat(foundPath); err == nil && !info.IsDir() {
				return foundPath
			}
		}
	}

	return ""
}

// LoadClass 加载类
func (m *DefaultClassPathManager) LoadClass(className string, parser *Parser) data.Control {
	filePath, found := m.FindClassFile(className)
	if !found {
		return data.TryErrorThrow(parser.newFrom(), fmt.Errorf("类 %s 不存在或无法加载", className))
	}

	// 加载文件
	_, acl := parser.vm.LoadAndRun(filePath)
	if acl != nil {
		return acl
	}

	// 检查类是否成功加载
	if _, ok := parser.vm.GetClass(className); ok {
		return nil
	}
	if _, ok := parser.vm.GetInterface(className); ok {
		return nil
	}

	return data.TryErrorThrow(parser.newFrom(), fmt.Errorf("文件 file://%s 中未找到类 %s", filePath, className))
}
