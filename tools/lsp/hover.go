package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/php-any/origami/tools/lsp/defines"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/jsonrpc2"
)

// 处理悬停请求
func handleTextDocumentHover(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	logLSPCommunication("textDocument/hover", true, req.Params)

	var params defines.HoverParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal hover params: %v", err)
	}

	uri := params.TextDocument.URI
	position := params.Position

	logrus.Infof("请求悬停提示：%s 位置 %d:%d", uri, position.Line, position.Character)

	doc, exists := documents[uri]
	if !exists {
		return nil, nil
	}

	// 获取悬停信息
	hoverInfo := getHoverInfo(doc.Content, position)
	if hoverInfo == "" {
		return nil, nil
	}

	result := &defines.Hover{
		Contents: defines.MarkupContent{
			Kind:  defines.MarkupKindMarkdown,
			Value: hoverInfo,
		},
	}

	logLSPResponse("textDocument/hover", result, nil)
	return result, nil
}

// 获取悬停信息
func getHoverInfo(content string, position defines.Position) string {
	// 获取光标位置的单词
	word := getWordAtPosition(content, position)
	if word == "" {
		return ""
	}

	// 简化的悬停信息
	switch word {
	case "fold":
		return "**fold** - Creates a fold in the paper"
	case "unfold":
		return "**unfold** - Unfolds a previously made fold"
	case "crease":
		return "**crease** - Creates a crease line"
	case "valley":
		return "**valley** - Creates a valley fold"
	case "mountain":
		return "**mountain** - Creates a mountain fold"
	case "function":
		return "**function** - Defines a function"
	case "class":
		return "**class** - Defines a class"
	default:
		return ""
	}
}
