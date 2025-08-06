package server

import (
	"log"
)

// handleInitialize 处理 LSP 初始化请求
// 这是 LSP 协议的第一个请求，客户端发送此请求来初始化服务器
// 服务器需要返回其支持的功能能力（capabilities）
func (s *Server) handleInitialize(request map[string]interface{}) error {
	// 获取请求 ID，用于响应
	id := request["id"]

	// 定义服务器支持的功能能力
	capabilities := map[string]interface{}{
		// 文档同步配置
		"textDocumentSync": map[string]interface{}{
			"openClose": true, // 支持文档打开/关闭事件
			"change":    1,    // 完整文档同步模式（1 = Full, 2 = Incremental）
		},
		// 代码补全提供器配置
		"completionProvider": map[string]interface{}{
			// 触发代码补全的字符
			"triggerCharacters": []string{"$", "->", "::", "\\"},
		},
		// 其他语言服务功能
		"hoverProvider":           true, // 悬停提示支持
		"definitionProvider":      true, // 定义跳转支持
		"referencesProvider":      true, // 引用查找支持
		"documentSymbolProvider":  true, // 文档符号支持
		"workspaceSymbolProvider": true, // 工作区符号支持
	}

	// 构建初始化响应
	result := map[string]interface{}{
		"capabilities": capabilities,
		"serverInfo": map[string]interface{}{
			"name":    "Origami Language Server", // 服务器名称
			"version": "1.0.0",                   // 服务器版本
		},
	}

	// 发送响应给客户端
	return s.sendResponse(id, result)
}

// handleInitialized 处理初始化完成通知
// 客户端在收到 initialize 响应后发送此通知，表示初始化过程完成
// 这是一个通知（notification），不需要响应
func (s *Server) handleInitialized(request map[string]interface{}) error {
	log.Println("LSP 客户端初始化完成")
	// 在这里可以执行一些初始化后的操作，比如：
	// - 注册文件监听器
	// - 初始化工作区配置
	// - 发送初始化完成的日志
	return nil
}
