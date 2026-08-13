package utils

import (
	"path/filepath"
	"runtime"
	"strings"
)

// NormalizePhpFilePath 将 PHP 文件路径规范化为统一的绝对路径（用于 php 文件加载缓存）。
func NormalizePhpFilePath(file string) string {
	if file == "" {
		return ""
	}
	cleaned := filepath.Clean(file)
	if abs, err := filepath.Abs(cleaned); err == nil {
		cleaned = abs
	}
	if resolved, err := filepath.EvalSymlinks(cleaned); err == nil {
		if filepath.IsAbs(resolved) {
			cleaned = resolved
		} else if abs, err := filepath.Abs(resolved); err == nil {
			cleaned = abs
		} else {
			cleaned = resolved
		}
	} else if dir, err := filepath.EvalSymlinks(filepath.Dir(cleaned)); err == nil {
		// 文件可能已被删除，但仍需与存在时解析出的路径保持一致（如 macOS /var -> /private/var）。
		cleaned = filepath.Join(dir, filepath.Base(cleaned))
	}
	// macOS/Windows 默认大小写不敏感，PSR-4 命名空间大小写可能与目录名不一致，统一为小写
	if isCaseInsensitiveFS() {
		cleaned = strings.ToLower(cleaned)
	}
	return cleaned
}

func isCaseInsensitiveFS() bool {
	switch runtime.GOOS {
	case "darwin", "windows":
		return true
	default:
		return false
	}
}

// SamePhpFile 判断两个路径是否指向同一 PHP 文件。
func SamePhpFile(a, b string) bool {
	na := NormalizePhpFilePath(a)
	nb := NormalizePhpFilePath(b)
	return na != "" && na == nb
}
