package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/tools/lsp/defines"

	"github.com/php-any/origami/node"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/jsonrpc2"
)

// 处理文档打开通知
func handleTextDocumentDidOpen(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	logLSPCommunication("textDocument/didOpen", false, req.Params)

	var params defines.DidOpenTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal didOpen params: %v", err)
	}

	uri := params.TextDocument.URI
	content := params.TextDocument.Text
	version := params.TextDocument.Version

	logrus.Infof("文档已打开：%s", uri)

	// 创建解析器
	p := NewLspParser()
	// 设置 LspVM
	if globalLspVM != nil {
		p.SetVM(globalLspVM)
	}

	// 解析 AST
	var ast *node.Program
	var acl data.Control

	// 如果是文件 URI，直接使用真实文件路径解析
	if strings.HasPrefix(uri, "file://") {
		filePath := uriToFilePath(uri)
		ast, acl = p.ParseFile(filePath)
		if acl != nil {
			logrus.Warnf("解析 AST 失败 %s：%v", uri, acl)
			// 解析失败：不在此处发送诊断，统一由 validateDocument 处理
		}
	}

	// 无论解析是否成功都创建 DocumentInfo
	if old, ok := documents[uri]; ok {
		if ast != nil {
			old.AST = ast
			old.Version = int32(version)
		}
		old.Content = content
	} else {
		documents[uri] = &DocumentInfo{
			Content: content,
			Version: int32(version),
			AST:     ast, // 可能为 nil
			Parser:  p,
		}
	}

	// 验证文档
	validateDocument(conn, uri, content)

	return nil, nil
}

// 解析错误的诊断由 validateDocument 统一发送

// 处理文档变更通知
func handleTextDocumentDidChange(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	logLSPCommunication("textDocument/didChange", false, req.Params)

	var params defines.DidChangeTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal didChange params: %v", err)
	}

	uri := params.TextDocument.URI
	version := params.TextDocument.Version

	logrus.Infof("文档已变更：%s", uri)

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
	var acl data.Control

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
	ast, acl = p.ParseString(content, filePath)
	if acl != nil {
		logrus.Warnf("重新解析 AST 失败 %s：%v", uri, acl)
		// 解析失败时，只更新内容和版本，保留原有的 AST 和解析器
		if existingDoc, exists := documents[uri]; exists {
			existingDoc.Content = content
			existingDoc.Version = int32(version)
		}
		// 统一由 validateDocument 发送诊断
		validateDocument(conn, uri, content)
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
	logrus.Infof("重新解析 AST 成功 %s; %v", uri, ast)
	// 验证文档
	validateDocument(conn, uri, content)

	return nil, nil
}

// 处理文档关闭通知
func handleTextDocumentDidClose(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	logLSPCommunication("textDocument/didClose", false, req.Params)

	var params defines.DidCloseTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal didClose params: %v", err)
	}

	uri := params.TextDocument.URI

	logrus.Infof("文档已关闭：%s", uri)

	// 注意：不再在此处清除 LspVM 中该文件的符号。
	// 原先这里调用 ClearFile，会导致类/函数符号从全局索引中删除，
	// 关闭文件后再次执行「跳转到定义」就找不到目标了（只能跳转一次）。
	// 为了保持符号索引稳定，仅移除文档缓存，让 LspVM 继续保留已解析的符号。
	// 修改：不删除 DocumentInfo，仅清空 Content，以便保留 AST 供后续分析使用，同时减少内存占用。
	if doc, ok := documents[uri]; ok {
		doc.Content = ""
	}
	// delete(documents, uri)

	return nil, nil
}
