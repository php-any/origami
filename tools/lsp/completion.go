package main

import (
	"encoding/json"
	"fmt"

	"github.com/php-any/origami/tools/lsp/completion"
	"github.com/php-any/origami/tools/lsp/defines"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/jsonrpc2"
)

// 处理补全请求
func handleTextDocumentCompletion(req *jsonrpc2.Request) (interface{}, error) {
	logLSPCommunication("textDocument/completion", true, req.Params)

	var params defines.CompletionParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal completion params: %v", err)
	}

	uri := params.TextDocument.URI
	position := params.Position

	logrus.Infof("请求代码补全：%s 位置 %d:%d", uri, position.Line, position.Character)

	doc, exists := documents[uri]
	if !exists {
		return defines.CompletionList{IsIncomplete: false, Items: []defines.CompletionItem{}}, nil
	}

	// 创建 SymbolProvider
	provider := &LSPSymbolProvider{
		doc: doc,
		vm:  globalLspVM,
	}

	// 获取补全项
	items := completion.GetCompletionItems(doc.Content, position, provider)

	result := defines.CompletionList{
		IsIncomplete: false,
		Items:        items,
	}

	logLSPResponse("textDocument/completion", result, nil)
	return result, nil
}
