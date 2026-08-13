package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/sirupsen/logrus"

	"github.com/php-any/origami/node"
	"github.com/php-any/origami/tools/lsp/defines"
	"github.com/sourcegraph/jsonrpc2"
)

// 验证文档 - 针对origami语言(.zy文件)的语法验证
func validateDocument(conn *jsonrpc2.Conn, uri string, content string) {
	diagnostics := []defines.Diagnostic{}

	// 检查文件扩展名
	if !strings.HasSuffix(uri, ".zy") && !strings.HasSuffix(uri, ".php") {
		// 不是origami语言文件，跳过验证
		return
	}

	// 使用专业的AST解析进行诊断，而不是简单的字符串匹配
	astDiagnostics := validateDocumentWithAST(uri, content)
	diagnostics = append(diagnostics, astDiagnostics...)

	// 发布诊断
	params := defines.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	}

	// 发送通知
	if b, err := json.Marshal(params); err == nil {
		logrus.Infof("textDocument/publishDiagnostics response %s", string(b))
	}
	conn.Notify(context.Background(), "textDocument/publishDiagnostics", params)
}

// 使用AST进行专业的文档验证
func validateDocumentWithAST(uri, content string) []defines.Diagnostic {
	var diagnostics []defines.Diagnostic

	// 创建解析器
	parser := NewLspParser()
	if globalLspVM != nil {
		parser.SetVM(globalLspVM)
	}

	// 解析AST
	var ast *node.Program
	var acl data.Control

	// 根据URI类型选择解析方法
	if strings.HasPrefix(uri, "file://") {
		filePath := uriToFilePath(uri)
		ast, acl = parser.ParseFile(filePath)
	} else {
		// 对于内存中的内容，使用ParseString
		ast, acl = parser.ParseString(content, "memory_content")
	}

	// 如果解析失败，返回解析错误
	if acl != nil {
		// 使用 acl 自带的 from 位置定位
		startLine, startChar, endLine, endChar := uint32(0), uint32(0), uint32(0), uint32(0)
		if gf, ok := acl.(node.GetFrom); ok && gf.GetFrom() != nil {
			sl, sc, el, ec := gf.GetFrom().ToLSPPosition()
			startLine, startChar, endLine, endChar = uint32(sl), uint32(sc), uint32(el), uint32(ec)
		}
		diagnostics = append(diagnostics, defines.Diagnostic{
			Range:    defines.Range{Start: defines.Position{Line: startLine, Character: startChar}, End: defines.Position{Line: endLine, Character: endChar}},
			Severity: &[]defines.DiagnosticSeverity{defines.DiagnosticSeverityError}[0],
			Message:  fmt.Sprintf("解析错误: %v", acl.AsString()),
			Source:   &[]string{"origami-lsp"}[0],
		})
		return diagnostics
	}

	// 如果AST解析成功，进行更深层的语义检查
	if ast != nil {
		semanticDiagnostics := validateASTSemantics(ast)
		diagnostics = append(diagnostics, semanticDiagnostics...)
	}

	return diagnostics
}

// 验证AST的语义
func validateASTSemantics(ast *node.Program) []defines.Diagnostic {
	var diagnostics []defines.Diagnostic

	// 这里可以添加更专业的语义检查
	// 例如：
	// - 类型检查
	// - 未定义变量检查
	// - 未使用变量检查
	// - 函数签名检查
	// - 等等

	// 暂时返回空列表，后续可以扩展
	return diagnostics
}
