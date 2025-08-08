package main

import (
	"encoding/json"
	"fmt"
	"strings"

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

	fmt.Printf("[INFO] Document opened: %s\n", uri)

	// 创建解析器
	p := NewLspParser()
	// 设置 LspVM
	if globalLspVM != nil {
		p.SetVM(globalLspVM)
	}

	// 解析 AST
	var ast interface{}
	var err error

	// 如果是文件 URI，直接使用真实文件路径解析
	if strings.HasPrefix(uri, "file://") {
		filePath := strings.TrimPrefix(uri, "file://")
		ast, err = p.ParseFile(filePath)
		if err != nil {
			if *logLevel > 1 {
				fmt.Printf("[WARNING] Failed to parse AST for %s: %v\n", uri, err)
			}
		}
	}

	documents[uri] = &DocumentInfo{
		Content: content,
		Version: int32(version),
		AST:     ast,
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

	if *logLevel > 2 {
		fmt.Printf("[INFO] Document changed: %s\n", uri)
	}

	doc, exists := documents[uri]
	if !exists {
		return nil, fmt.Errorf("document not found: %s", uri)
	}

	// 应用变更
	for _, change := range params.ContentChanges {
		doc.Content = change.Text
	}

	doc.Version = int32(version)

	// 重新解析 AST
	if doc.Parser != nil {
		if p, ok := doc.Parser.(*LspParser); ok {
			var ast interface{}
			var err error

			// 如果是文件 URI，直接使用真实文件路径解析
			if strings.HasPrefix(uri, "file://") {
				filePath := strings.TrimPrefix(uri, "file://")

				// 清除 LspVM 中该文件的旧符号
				if globalLspVM != nil {
					globalLspVM.ClearFile(filePath)
				}

				ast, err = p.ParseFile(filePath)
			}

			if err != nil {
				if *logLevel > 1 {
					fmt.Printf("[WARNING] Failed to re-parse AST for %s: %v\n", uri, err)
				}
			} else {
				doc.AST = ast
			}
		}
	}

	// 验证文档
	validateDocument(conn, uri, doc.Content)

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

	if *logLevel > 2 {
		fmt.Printf("[INFO] Document closed: %s\n", uri)
	}

	// 清除 LspVM 中该文件的符号
	if strings.HasPrefix(uri, "file://") {
		filePath := strings.TrimPrefix(uri, "file://")
		if globalLspVM != nil {
			globalLspVM.ClearFile(filePath)
		}
	}

	delete(documents, uri)

	return nil, nil
}
