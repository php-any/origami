package main

import (
	"encoding/json"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/jsonrpc2"
)

// handleTextDocumentDocumentSymbol 处理文档符号请求
func handleTextDocumentDocumentSymbol(req *jsonrpc2.Request) (interface{}, error) {
	logLSPCommunication("textDocument/documentSymbol", true, req.Params)

	var params DocumentSymbolParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal documentSymbol params: %v", err)
	}

	logrus.Infof("处理文档符号请求：%s", params.TextDocument.URI)

	doc, exists := documents[params.TextDocument.URI]
	if !exists {
		logrus.Warnf("文档不存在：%s", params.TextDocument.URI)
		return []DocumentSymbol{}, nil
	}

	symbols := getDocumentSymbolsFromAST(doc)

	logrus.Infof("textDocument/documentSymbol response %#v", symbols)
	return symbols, nil
}

// getDocumentSymbolsFromAST 从 AST 中提取文档符号
func getDocumentSymbolsFromAST(doc *DocumentInfo) []DocumentSymbol {
	if doc.AST == nil {
		logrus.Warn("文档 AST 为空")
		return []DocumentSymbol{}
	}

	var symbols []DocumentSymbol

	// 使用 DocumentInfo.Foreach 遍历 AST，只收集顶级符号
	doc.Foreach(func(ctx *LspContext, parent, child data.GetValue) bool {
		// 只处理顶级符号（parent 为 nil）
		if parent == nil {
			symbol := createSymbolFromNode(child)
			if symbol != nil {
				symbols = append(symbols, *symbol)
			}
		}
		return true // 继续遍历
	})

	return symbols
}

// createSymbolFromNode 从节点创建符号
func createSymbolFromNode(nodeValue data.GetValue) *DocumentSymbol {
	if nodeValue == nil {
		return nil
	}

	switch n := nodeValue.(type) {
	case *node.FunctionStatement:
		detail := fmt.Sprintf("function %s", n.Name)
		return &DocumentSymbol{
			Name:           n.Name,
			Detail:         &detail,
			Kind:           SymbolKindFunction,
			Range:          getNodeRange(n),
			SelectionRange: getNodeRange(n),
		}

	case *node.ClassStatement:
		detail := fmt.Sprintf("class %s", n.Name)
		symbol := &DocumentSymbol{
			Name:           n.Name,
			Detail:         &detail,
			Kind:           SymbolKindClass,
			Range:          getNodeRange(n),
			SelectionRange: getNodeRange(n),
		}

		// 添加类成员
		symbol.Children = extractClassMemberSymbols(n)
		return symbol

	case *node.InterfaceStatement:
		detail := fmt.Sprintf("interface %s", n.Name)
		return &DocumentSymbol{
			Name:           n.Name,
			Detail:         &detail,
			Kind:           SymbolKindInterface,
			Range:          getNodeRange(n),
			SelectionRange: getNodeRange(n),
		}

	case *node.VarStatement:
		detail := fmt.Sprintf("var %s", n.Name)
		return &DocumentSymbol{
			Name:           n.Name,
			Detail:         &detail,
			Kind:           SymbolKindVariable,
			Range:          getNodeRange(n),
			SelectionRange: getNodeRange(n),
		}

	case *node.ConstStatement:
		detail := fmt.Sprintf("const %s", n.Val.GetName())
		return &DocumentSymbol{
			Name:           n.Val.GetName(),
			Detail:         &detail,
			Kind:           SymbolKindConstant,
			Range:          getNodeRange(n),
			SelectionRange: getNodeRange(n),
		}

	case *node.Namespace:
		detail := fmt.Sprintf("namespace %s", n.Name)
		symbol := &DocumentSymbol{
			Name:           n.Name,
			Detail:         &detail,
			Kind:           SymbolKindNamespace,
			Range:          getNodeRange(n),
			SelectionRange: getNodeRange(n),
		}

		// 添加命名空间成员
		var children []DocumentSymbol
		for _, stmt := range n.Statements {
			if childSymbol := createSymbolFromNode(stmt); childSymbol != nil {
				children = append(children, *childSymbol)
			}
		}
		symbol.Children = children
		return symbol
	}

	return nil
}

// extractClassMemberSymbols 提取类成员符号
func extractClassMemberSymbols(class *node.ClassStatement) []DocumentSymbol {
	var symbols []DocumentSymbol

	// 添加属性
	for _, prop := range class.Properties {
		detail := fmt.Sprintf("property %s", prop.GetName())
		symbol := DocumentSymbol{
			Name:           prop.GetName(),
			Detail:         &detail,
			Kind:           SymbolKindProperty,
			Range:          getPropertyRange(prop),
			SelectionRange: getPropertyRange(prop),
		}
		symbols = append(symbols, symbol)
	}

	// 添加方法
	for _, method := range class.Methods {
		detail := fmt.Sprintf("method %s", method.GetName())
		symbol := DocumentSymbol{
			Name:           method.GetName(),
			Detail:         &detail,
			Kind:           SymbolKindMethod,
			Range:          getMethodRange(method),
			SelectionRange: getMethodRange(method),
		}
		symbols = append(symbols, symbol)
	}

	return symbols
}

// getNodeRange 获取节点的范围
func getNodeRange(nodeValue data.GetValue) Range {
	if getFrom, ok := nodeValue.(node.GetFrom); ok {
		if from := getFrom.GetFrom(); from != nil {
			startLine, startCol, endLine, endCol := from.GetRange()
			return Range{
				Start: Position{Line: uint32(startLine), Character: uint32(startCol)},
				End:   Position{Line: uint32(endLine), Character: uint32(endCol)},
			}
		}
	}
	return Range{}
}

// getPropertyRange 获取属性的范围
func getPropertyRange(prop data.Property) Range {
	// 尝试从属性获取位置信息
	if getFrom, ok := prop.(node.GetFrom); ok {
		if from := getFrom.GetFrom(); from != nil {
			startLine, startCol, endLine, endCol := from.GetRange()
			return Range{
				Start: Position{Line: uint32(startLine), Character: uint32(startCol)},
				End:   Position{Line: uint32(endLine), Character: uint32(endCol)},
			}
		}
	}
	return Range{}
}

// getMethodRange 获取方法的范围
func getMethodRange(method data.Method) Range {
	// 尝试从方法获取位置信息
	if getFrom, ok := method.(node.GetFrom); ok {
		if from := getFrom.GetFrom(); from != nil {
			startLine, startCol, endLine, endCol := from.GetRange()
			return Range{
				Start: Position{Line: uint32(startLine), Character: uint32(startCol)},
				End:   Position{Line: uint32(endLine), Character: uint32(endCol)},
			}
		}
	}
	return Range{}
}
