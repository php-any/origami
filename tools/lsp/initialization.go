package main

import (
	"encoding/json"
	"fmt"
	"os"

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
	return result, nil
}

// 处理初始化完成通知
func handleInitialized(req *jsonrpc2.Request) (interface{}, error) {
	logLSPCommunication("initialized", false, req.Params)

	if *logLevel > 0 {
		fmt.Fprintf(os.Stderr, "[INFO] Origami LSP server initialized successfully\n")
	}

	return nil, nil
}

// 处理关闭请求
func handleShutdown(req *jsonrpc2.Request) (interface{}, error) {
	logLSPCommunication("shutdown", true, req.Params)

	if *logLevel > 0 {
		fmt.Fprintf(os.Stderr, "[INFO] Shutting down Origami LSP server...\n")
	}

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
	if *logLevel > 0 {
		fmt.Fprintf(os.Stderr, "[INFO] Setting trace value: %s\n", params.Value)
	}

	return nil, nil
}
