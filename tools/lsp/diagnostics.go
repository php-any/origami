package main

import (
	"context"
	"strings"

	"github.com/sourcegraph/jsonrpc2"
)

// 验证文档
func validateDocument(conn *jsonrpc2.Conn, uri string, content string) {
	diagnostics := []Diagnostic{}

	// 基本的标记化和验证
	lines := strings.Split(content, "\n")
	for lineNum, line := range lines {
		tokens := strings.Fields(line)
		for colNum, token := range tokens {
			// 检查未知标记（简化）
			if !isKnownToken(token) {
				diagnostic := Diagnostic{
					Range: Range{
						Start: Position{
							Line:      uint32(lineNum),    // 从0开始，与lexer保持一致
							Character: uint32(colNum * 5), // 简化的位置计算
						},
						End: Position{
							Line:      uint32(lineNum), // 从0开始，与lexer保持一致
							Character: uint32(colNum*5 + len(token)),
						},
					},
					Severity: &[]DiagnosticSeverity{DiagnosticSeverityWarning}[0],
					Message:  "Unknown token: " + token,
				}
				diagnostics = append(diagnostics, diagnostic)
			}
		}
	}

	// 发布诊断
	params := PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	}

	// 发送通知
	conn.Notify(context.Background(), "textDocument/publishDiagnostics", params)
}

// 判断是否为已知标记
func isKnownToken(token string) bool {
	// 简化的标记验证
	keywords := []string{
		"fold", "unfold", "crease", "valley", "mountain", "reverse",
		"rotate", "translate", "scale", "reflect", "paper", "point",
		"line", "angle", "distance", "function", "class", "if", "else",
		"for", "while", "return", "var", "let", "const",
	}

	for _, keyword := range keywords {
		if token == keyword {
			return true
		}
	}

	// 允许数字、字符串和标识符
	if len(token) > 0 {
		first := token[0]
		if (first >= 'a' && first <= 'z') || (first >= 'A' && first <= 'Z') || first == '_' {
			return true
		}
		if first >= '0' && first <= '9' {
			return true
		}
		if first == '"' || first == '\'' {
			return true
		}
	}

	return false
}
