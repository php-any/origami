package main

import (
	"encoding/json"
	"fmt"

	"github.com/sourcegraph/jsonrpc2"
)

// 处理补全请求
func handleTextDocumentCompletion(req *jsonrpc2.Request) (interface{}, error) {
	logLSPCommunication("textDocument/completion", true, req.Params)

	var params CompletionParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal completion params: %v", err)
	}

	uri := params.TextDocument.URI
	position := params.Position

	if *logLevel > 2 {
		fmt.Printf("[INFO] Completion requested for %s at %d:%d\n", uri, position.Line, position.Character)
	}

	doc, exists := documents[uri]
	if !exists {
		return CompletionList{IsIncomplete: false, Items: []CompletionItem{}}, nil
	}

	// 获取补全项
	items := getCompletionItems(doc.Content, position)

	result := CompletionList{
		IsIncomplete: false,
		Items:        items,
	}

	logLSPResponse("textDocument/completion", result, nil)
	return result, nil
}

// 获取补全项
func getCompletionItems(content string, position Position) []CompletionItem {
	items := []CompletionItem{}

	// Origami 关键字
	keywords := []string{
		"fold", "unfold", "crease", "valley", "mountain", "reverse",
		"rotate", "translate", "scale", "reflect", "paper", "point",
		"line", "angle", "distance",
	}

	for _, keyword := range keywords {
		item := CompletionItem{
			Label:  keyword,
			Kind:   &[]CompletionItemKind{CompletionItemKindKeyword}[0],
			Detail: &[]string{"Origami keyword"}[0],
		}
		items = append(items, item)
	}

	// 内置函数
	functions := []string{
		"function", "class", "if", "else", "for", "while", "return",
	}

	for _, function := range functions {
		item := CompletionItem{
			Label:  function,
			Kind:   &[]CompletionItemKind{CompletionItemKindFunction}[0],
			Detail: &[]string{"Built-in function"}[0],
		}
		items = append(items, item)
	}

	return items
}
