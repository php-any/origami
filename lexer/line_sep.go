package lexer

import (
	"runtime"
	"strings"
)

// LineSeparator 当前操作系统的换行符常量。
// Windows 使用 "\r\n"，其他系统使用 "\n"。
var LineSeparator = "\n"

func init() {
	if runtime.GOOS == "windows" {
		LineSeparator = "\r\n"
	}
}

// NormalizeToLF 将字符串中的 "\r\n" 和 "\r" 统一转为 "\n"。
// 用于从源码中提取的内容（如 heredoc 正文），确保跨平台行为一致。
func NormalizeToLF(s string) string {
	if !strings.ContainsRune(s, '\r') {
		return s
	}
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}
