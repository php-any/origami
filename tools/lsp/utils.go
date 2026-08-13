package main

import (
	"net/url"
	"runtime"
	"strings"

	"github.com/php-any/origami/tools/lsp/defines"
)

// 获取光标位置的单词
func getWordAtPosition(content string, position defines.Position) string {
	lines := strings.Split(content, "\n")
	if int(position.Line) >= len(lines) {
		return ""
	}

	line := lines[position.Line]
	if int(position.Character) >= len(line) {
		return ""
	}

	// 查找单词边界
	start := int(position.Character)
	end := int(position.Character)

	// 向前查找单词开始
	for start > 0 && isWordChar(line[start-1]) {
		start--
	}

	// 向后查找单词结束
	for end < len(line) && isWordChar(line[end]) {
		end++
	}

	if start == end {
		return ""
	}

	return line[start:end]
}

// 判断是否为单词字符
func isWordChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

// uriToFilePath 将 file:// URI 转换为本地文件路径，兼容不同操作系统
func uriToFilePath(uri string) string {
	if !strings.HasPrefix(uri, "file://") {
		return uri
	}

	// 使用 url.Parse 来正确解析 URI
	parsedURL, err := url.Parse(uri)
	if err != nil {
		// 如果解析失败，回退到简单的字符串处理
		return strings.TrimPrefix(uri, "file://")
	}

	filePath := parsedURL.Path

	// Windows 系统特殊处理
	if runtime.GOOS == "windows" {
		// Windows 路径格式：file:///C:/path/to/file
		// 需要去掉开头的 / 并将 / 替换为 \
		if len(filePath) > 0 && filePath[0] == '/' {
			filePath = filePath[1:] // 去掉开头的 /
		}
		filePath = strings.ReplaceAll(filePath, "/", "\\")
	}

	return filePath
}

// filePathToURI 将本地文件路径转换为 file:// URI，兼容不同操作系统
func filePathToURI(filePath string) string {
	if filePath == "" {
		return ""
	}

	// Windows 系统特殊处理
	if runtime.GOOS == "windows" {
		// 将 \ 替换为 /
		filePath = strings.ReplaceAll(filePath, "\\", "/")
		// 确保以 / 开头
		if !strings.HasPrefix(filePath, "/") {
			filePath = "/" + filePath
		}
	}

	// 使用 url.URL 来正确构建 URI
	u := &url.URL{
		Scheme: "file",
		Path:   filePath,
	}

	return u.String()
}
