package main

import (
	"context"
	"fmt"
	"github.com/php-any/origami/data"
	"strings"

	"github.com/php-any/origami/node"
	"github.com/sourcegraph/jsonrpc2"
)

// 验证文档 - 针对origami语言(.zy文件)的语法验证
func validateDocument(conn *jsonrpc2.Conn, uri string, content string) {
	diagnostics := []Diagnostic{}

	// 检查文件扩展名
	if !strings.HasSuffix(uri, ".zy") && !strings.HasSuffix(uri, ".php") {
		// 不是origami语言文件，跳过验证
		return
	}

	// 使用专业的AST解析进行诊断，而不是简单的字符串匹配
	astDiagnostics := validateDocumentWithAST(uri, content)
	diagnostics = append(diagnostics, astDiagnostics...)

	// 发布诊断
	params := PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	}

	// 发送通知
	conn.Notify(context.Background(), "textDocument/publishDiagnostics", params)
}

// 使用AST进行专业的文档验证
func validateDocumentWithAST(uri, content string) []Diagnostic {
	var diagnostics []Diagnostic

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
		var err error
		ast, err = parser.ParseFile(filePath)
		if err != nil {
			acl = data.NewErrorThrow(nil, err)
		}
	} else {
		// 对于内存中的内容，使用ParseString
		ast, acl = parser.ParseString(content, "memory_content")
	}

	// 如果解析失败，返回解析错误
	if acl != nil {
		// 解析错误通常意味着语法问题
		diagnostics = append(diagnostics, Diagnostic{
			Range: Range{
				Start: Position{Line: 0, Character: 0},
				End:   Position{Line: 0, Character: 0},
			},
			Severity: &[]DiagnosticSeverity{DiagnosticSeverityError}[0],
			Message:  fmt.Sprintf("解析错误: %v", acl),
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
func validateASTSemantics(ast *node.Program) []Diagnostic {
	var diagnostics []Diagnostic

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
