package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sourcegraph/jsonrpc2"
)

// 处理初始化请求
func handleInitialize(req *jsonrpc2.Request) (interface{}, error) {
	logLSPCommunication("initialize", true, req.Params)

	var params InitializeParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal initialize params: %v", err)
	}

	// 设置服务器能力
	capabilities := ServerCapabilities{
		TextDocumentSync: &TextDocumentSyncOptions{
			OpenClose: &[]bool{true}[0],
			Change:    &[]int{1}[0], // Full sync
		},
		CompletionProvider: &CompletionOptions{
			TriggerCharacters: []string{".", "$", ":", "\\"},
		},
		HoverProvider:          &HoverOptions{},
		DefinitionProvider:     &DefinitionOptions{},
		DocumentSymbolProvider: &DocumentSymbolOptions{},
	}

	version := lsVersion
	result := InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &ServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}

	logLSPResponse("initialize", result, nil)

	// 异步加载所有脚本文件，不影响 initialize 响应
	go loadAllScriptFiles(params)

	return result, nil
}

// 处理初始化完成通知
func handleInitialized(req *jsonrpc2.Request) (interface{}, error) {
	logLSPCommunication("initialized", false, req.Params)

	logger.Info("Origami LSP 服务器初始化成功")

	return nil, nil
}

// 处理关闭请求
func handleShutdown(req *jsonrpc2.Request) (interface{}, error) {
	logLSPCommunication("shutdown", true, req.Params)

	logger.Info("正在关闭 Origami LSP 服务器...")

	return nil, nil
}

// 处理设置跟踪请求
func handleSetTrace(req *jsonrpc2.Request) (interface{}, error) {
	logLSPCommunication("$/setTrace", false, req.Params)

	var params SetTraceParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal setTrace params: %v", err)
	}

	// 设置跟踪级别
	logger.Info("设置跟踪值：%s", params.Value)

	return nil, nil
}

// loadAllScriptFiles 异步加载所有脚本文件
func loadAllScriptFiles(params InitializeParams) {
	// 添加 panic 恢复机制
	defer func() {
		if r := recover(); r != nil {
			logger.Error("loadAllScriptFiles 发生 panic：%v", r)
		}
	}()

	logger.Info("开始异步加载工作区中的所有脚本文件...")

	// 获取 LSP 工作区根目录
	workspaceRoot := getWorkspaceRoot(params)
	if workspaceRoot == "" {
		logger.Error("无法获取工作区根目录")
		return
	}

	logger.Info("工作区根目录：%s", workspaceRoot)

	// 直接遍历并立即加载文件，避免收集所有文件
	loadScriptFilesInDirectory(workspaceRoot)
}

// getWorkspaceRoot 从 LSP 参数获取工作区根目录
func getWorkspaceRoot(params InitializeParams) string {
	// 添加 panic 恢复机制
	defer func() {
		if r := recover(); r != nil {
			logger.Error("getWorkspaceRoot 发生 panic：%v", r)
		}
	}()

	// 优先使用 RootURI（更现代的 LSP 客户端）
	if params.RootURI != nil && *params.RootURI != "" {
		uri := *params.RootURI
		// 将 URI 转换为文件路径
		if strings.HasPrefix(uri, "file://") {
			filePath := strings.TrimPrefix(uri, "file://")
			// 在 Windows 上，file:///C:/path 需要特殊处理
			if strings.HasPrefix(filePath, "/") && len(filePath) > 3 && filePath[2] == ':' {
				filePath = filePath[1:] // 移除开头的 /
			}
			return filePath
		}
	}

	// 备选使用 RootPath
	if params.RootPath != nil && *params.RootPath != "" {
		return *params.RootPath
	}

	// 如果都没有，尝试使用工作区文件夹
	if len(params.WorkspaceFolders) > 0 {
		uri := params.WorkspaceFolders[0].URI
		if strings.HasPrefix(uri, "file://") {
			filePath := strings.TrimPrefix(uri, "file://")
			if strings.HasPrefix(filePath, "/") && len(filePath) > 3 && filePath[2] == ':' {
				filePath = filePath[1:]
			}
			return filePath
		}
	}

	logger.Error("LSP 参数中未找到有效的工作区根目录")
	return ""
}

// findScriptFiles 查找所有脚本文件
func findScriptFiles(workspaceRoot string) []string {
	// 添加 panic 恢复机制
	defer func() {
		if r := recover(); r != nil {
			logger.Error("findScriptFiles 发生 panic：%v", r)
		}
	}()

	var scriptFiles []string

	// 遍历目录查找脚本文件
	err := filepath.Walk(workspaceRoot, func(path string, info os.FileInfo, err error) error {
		// 为每个文件遍历回调添加 panic 恢复
		defer func() {
			if r := recover(); r != nil {
				logger.Error("遍历文件 %s 时发生 panic：%v", path, r)
			}
		}()

		if err != nil {
			return err
		}

		// 跳过隐藏目录和 .git 目录
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") || info.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		// 检查文件扩展名
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".cjp" || ext == ".ori" {
			scriptFiles = append(scriptFiles, path)
		}

		return nil
	})

	if err != nil {
		logger.Error("遍历工作区目录失败：%v", err)
	}

	return scriptFiles
}

// loadScriptFilesInDirectory 在目录中查找并立即加载脚本文件
func loadScriptFilesInDirectory(workspaceRoot string) {
	// 添加 panic 恢复机制
	defer func() {
		if r := recover(); r != nil {
			logger.Error("loadScriptFilesInDirectory 发生 panic：%v", r)
		}
	}()

	var fileCount int

	// 遍历目录查找脚本文件，找到后立即加载
	err := filepath.Walk(workspaceRoot, func(path string, info os.FileInfo, err error) error {
		// 为每个文件遍历回调添加 panic 恢复
		defer func() {
			if r := recover(); r != nil {
				logger.Error("遍历文件 %s 时发生 panic：%v", path, r)
			}
		}()

		if err != nil {
			return err
		}

		// 跳过隐藏目录和 .git 目录
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") || info.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		// 检查文件扩展名
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".cjp" || ext == ".ori" {
			fileCount++
			logger.Debug("发现脚本文件：%s", path)

			// 立即异步加载文件
			go func(filePath string) {
				// 为每个文件加载 goroutine 添加 panic 恢复
				defer func() {
					if r := recover(); r != nil {
						logger.Error("加载文件 %s 时发生 panic：%v", filePath, r)
					}
				}()
				// 创建共享的 LspParser 实例
				parser := NewLspParser()
				if globalLspVM != nil {
					parser.SetVM(globalLspVM)
				}
				loadScriptFile(filePath, parser)
			}(path)
		}

		return nil
	})

	if err != nil {
		logger.Error("遍历工作区目录失败：%v", err)
	}

	logger.Info("发现并开始加载 %d 个脚本文件", fileCount)
}

// loadScriptFile 加载单个脚本文件
func loadScriptFile(filePath string, parser *LspParser) {
	// 添加 panic 恢复机制
	defer func() {
		if r := recover(); r != nil {
			logger.Error("loadScriptFile 发生 panic：%v", r)
		}
	}()

	logger.Debug("正在加载脚本文件：%s", filePath)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logger.Debug("文件不存在：%s", filePath)
		return
	}

	// 使用传入的共享解析器解析文件
	if parser != nil {
		parser.ParseFile(filePath)
		logger.Debug("成功加载脚本文件：%s", filePath)
	} else {
		logger.Error("解析器未初始化，无法加载文件：%s", filePath)
	}
}
