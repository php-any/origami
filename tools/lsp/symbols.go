package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sourcegraph/jsonrpc2"
)

// 处理文档符号请求
func handleTextDocumentDocumentSymbol(req *jsonrpc2.Request) (interface{}, error) {
	logLSPCommunication("textDocument/documentSymbol", true, req.Params)

	var params DocumentSymbolParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal documentSymbol params: %v", err)
	}

	uri := params.TextDocument.URI

	if *logLevel > 2 {
		fmt.Printf("[INFO] Document symbols requested for %s\n", uri)
	}

	doc, exists := documents[uri]
	if !exists {
		return []DocumentSymbol{}, nil
	}

	// 获取文档符号
	symbols := getDocumentSymbols(doc.Content)

	logLSPResponse("textDocument/documentSymbol", symbols, nil)
	return symbols, nil
}

// 获取文档符号
func getDocumentSymbols(content string) []DocumentSymbol {
	symbols := []DocumentSymbol{}

	// 简化的符号提取
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)

		// 查找函数定义
		if strings.HasPrefix(line, "function ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				funcName := strings.TrimSuffix(parts[1], "(")
				symbol := DocumentSymbol{
					Name: funcName,
					Kind: SymbolKindFunction,
					Range: Range{
						Start: Position{Line: uint32(i), Character: 0},
						End:   Position{Line: uint32(i), Character: uint32(len(line))},
					},
					SelectionRange: Range{
						Start: Position{Line: uint32(i), Character: 0},
						End:   Position{Line: uint32(i), Character: uint32(len(line))},
					},
				}
				symbols = append(symbols, symbol)
			}
		}

		// 查找类定义
		if strings.HasPrefix(line, "class ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				className := parts[1]
				symbol := DocumentSymbol{
					Name: className,
					Kind: SymbolKindClass,
					Range: Range{
						Start: Position{Line: uint32(i), Character: 0},
						End:   Position{Line: uint32(i), Character: uint32(len(line))},
					},
					SelectionRange: Range{
						Start: Position{Line: uint32(i), Character: 0},
						End:   Position{Line: uint32(i), Character: uint32(len(line))},
					},
				}
				symbols = append(symbols, symbol)
			}
		}
	}

	return symbols
}
