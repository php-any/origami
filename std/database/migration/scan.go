package migration

import (
	"os"
	"path/filepath"
	"strings"
)

// collectPhpFiles 扫描目录下所有 .php 文件
func collectPhpFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(path), ".php") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func normalizePath(path string) string {
	if path == "" {
		return ""
	}
	if !filepath.IsAbs(path) {
		if cwd, err := os.Getwd(); err == nil {
			path = filepath.Join(cwd, path)
		}
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return filepath.Clean(path)
	}
	return filepath.Clean(abs)
}
