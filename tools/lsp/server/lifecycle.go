package server

// handleShutdown 处理关闭请求
// 当客户端准备关闭 LSP 服务器时发送此请求
// 服务器应该清理资源并准备退出，但不应该立即退出
func (s *Server) handleShutdown(request map[string]interface{}) error {
	// 获取请求 ID
	id := request["id"]

	// 执行清理操作
	// 这里可以添加清理文档缓存、关闭文件句柄等操作
	s.cleanup()

	// 发送成功响应
	// 根据 LSP 规范，shutdown 请求应该返回 null
	return s.sendResponse(id, nil)
}

// handleExit 处理退出通知
// 这是一个通知（notification），不需要响应
// 收到此通知后，服务器应该立即退出
func (s *Server) handleExit(request map[string]interface{}) error {
	// 执行最终清理
	s.cleanup()

	// 退出程序
	// 注意：在实际实现中，这里可能需要更优雅的退出机制
	// 比如设置一个标志位，让主循环检测到后退出
	return nil
}

// cleanup 执行清理操作
// 清理服务器使用的资源，如文档缓存等
func (s *Server) cleanup() {
	// 清理文档缓存
	if s.documents != nil {
		// 清空文档缓存映射
		for uri := range s.documents {
			delete(s.documents, uri)
		}
	}

	// 这里可以添加其他清理操作
	// 例如：关闭数据库连接、清理临时文件等
}
