package main

import (
	"strings"
)

// 获取光标位置的单词
func getWordAtPosition(content string, position Position) string {
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
