package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/php-any/origami/node"

	"github.com/sourcegraph/jsonrpc2"
)

// 处理文档打开通知
func handleTextDocumentDidOpen(conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	logLSPCommunication("textDocument/didOpen", false, req.Params)

	var params DidOpenTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal didOpen params: %v", err)
	}

	uri := params.TextDocument.URI
	content := params.TextDocument.Text
	version := params.TextDocument.Version

	logger.Info("文档已打开：%s", uri)

	// 创建解析器
	p := NewLspParser()
	// 设置 LspVM
	if globalLspVM != nil {
		p.SetVM(globalLspVM)
	}

	// 解析 AST
	var ast *node.Program
	var err error

	// 如果是文件 URI，直接使用真实文件路径解析
	if strings.HasPrefix(uri, "file://") {
		filePath := uriToFilePath(uri)
		ast, err = p.ParseFile(filePath)
		if err != nil {
			logger.Warn("解析 AST 失败 %s：%v", uri, err)
			// 解析失败时，AST 为 nil，但仍然创建文档信息
		}
	}

	// 无论解析是否成功都创建 DocumentInfo
	documents[uri] = &DocumentInfo{
		Content: content,
		Version: int32(version),
		AST:     ast, // 可能为 nil
		Parser:  p,
	}

	// 验证文档
	validateDocument(conn, uri, content)

	return nil, nil
}

// 处理文档变更通知
func handleTextDocumentDidChange(conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	logLSPCommunication("textDocument/didChange", false, req.Params)

	var params DidChangeTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal didChange params: %v", err)
	}

	uri := params.TextDocument.URI
	version := params.TextDocument.Version

	logger.Info("文档已变更：%s", uri)

	_, exists := documents[uri]
	if !exists {
		return nil, fmt.Errorf("document not found: %s", uri)
	}

	// 获取更新后的内容
	var content string
	for _, change := range params.ContentChanges {
		content = change.Text
	}

	// 创建解析器
	p := NewLspParser()
	// 设置 LspVM
	if globalLspVM != nil {
		p.SetVM(globalLspVM)
	}

	// 解析 AST
	var ast *node.Program
	var err error

	// 使用编辑器提供的最新内容来解析，而不是从磁盘读取

	// 清除 LspVM 中该文件的旧符号（如果是文件 URI）
	if strings.HasPrefix(uri, "file://") {
		filePath := uriToFilePath(uri)
		if globalLspVM != nil {
			globalLspVM.ClearFile(filePath)
		}
	}

	// 使用最新内容解析 AST
	var filePath string
	if strings.HasPrefix(uri, "file://") {
		filePath = uriToFilePath(uri)
	} else {
		filePath = "memory_content" // 非文件 URI 使用虚拟路径
	}
	ast, err = p.ParseString(content, filePath)
	if err != nil {
		logger.Warn("重新解析 AST 失败 %s：%v", uri, err)
		// 解析失败时，只更新内容和版本，保留原有的 AST 和解析器
		if existingDoc, exists := documents[uri]; exists {
			existingDoc.Content = content
			existingDoc.Version = int32(version)
		}
		return nil, nil
	}
	delete(documents, uri)
	// 只有解析成功时才重新创建 DocumentInfo
	documents[uri] = &DocumentInfo{
		Content: content,
		Version: int32(version),
		AST:     ast,
		Parser:  p,
	}
	logger.Info("重新解析 AST 成功 %s; %v", uri, ast)
	// 验证文档
	validateDocument(conn, uri, content)

	return nil, nil
}

// 处理文档关闭通知
func handleTextDocumentDidClose(conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	logLSPCommunication("textDocument/didClose", false, req.Params)

	var params DidCloseTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal didClose params: %v", err)
	}

	uri := params.TextDocument.URI

	logger.Info("文档已关闭：%s", uri)

	// 清除 LspVM 中该文件的符号
	if strings.HasPrefix(uri, "file://") {
		filePath := uriToFilePath(uri)
		if globalLspVM != nil {
			globalLspVM.ClearFile(filePath)
		}
	}

	delete(documents, uri)

	return nil, nil
}
