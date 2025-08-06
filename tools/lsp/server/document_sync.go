package server

import (
	"log"
)

// handleTextDocumentDidOpen 处理文档打开事件
// 当客户端打开一个文档时会发送此通知
// 服务器需要开始跟踪该文档的状态并进行语法验证
func (s *Server) handleTextDocumentDidOpen(request map[string]interface{}) error {
	// 解析请求参数
	params, ok := request["params"].(map[string]interface{})
	if !ok {
		log.Println("文档打开事件：无效的参数格式")
		return nil
	}

	// 获取文档信息
	textDocument, ok := params["textDocument"].(map[string]interface{})
	if !ok {
		log.Println("文档打开事件：无效的文档信息")
		return nil
	}

	// 提取文档 URI 和内容
	uri, _ := textDocument["uri"].(string)
	text, _ := textDocument["text"].(string)

	log.Printf("文档打开: %s", uri)

	// 将文档内容保存到服务器缓存中
	// 这样可以在后续的代码补全、悬停等功能中快速访问文档内容
	s.documents[uri] = text

	// 对新打开的文档进行语法验证
	// 如果发现语法错误，会通过诊断消息发送给客户端
	s.validateDocument(uri, text)
	return nil
}

// handleTextDocumentDidChange 处理文档变更事件
// 当文档内容发生变化时，客户端会发送此通知
// 服务器需要更新文档缓存并重新进行语法验证
func (s *Server) handleTextDocumentDidChange(request map[string]interface{}) error {
	// 解析请求参数
	params, ok := request["params"].(map[string]interface{})
	if !ok {
		log.Println("文档变更事件：无效的参数格式")
		return nil
	}

	// 获取文档信息
	textDocument, ok := params["textDocument"].(map[string]interface{})
	if !ok {
		log.Println("文档变更事件：无效的文档信息")
		return nil
	}

	// 提取文档 URI
	uri, _ := textDocument["uri"].(string)

	// 获取内容变更信息
	contentChanges, ok := params["contentChanges"].([]interface{})
	if !ok || len(contentChanges) == 0 {
		log.Println("文档变更事件：无效的变更内容")
		return nil
	}

	// 获取最新的文档内容
	// 在完整文档同步模式下，变更数组只包含一个元素，即完整的新文档内容
	change, ok := contentChanges[0].(map[string]interface{})
	if !ok {
		log.Println("文档变更事件：无效的变更格式")
		return nil
	}

	// 提取新的文档内容
	text, _ := change["text"].(string)

	log.Printf("文档变更: %s", uri)

	// 更新服务器中的文档缓存
	s.documents[uri] = text

	// 重新验证更新后的文档
	// 这会检查新的语法错误并发送诊断信息给客户端
	s.validateDocument(uri, text)
	return nil
}

// validateDocument 函数已移至 validation.go 文件中
// 该函数负责验证文档语法并发送诊断信息
