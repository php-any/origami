package compile

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
			// 跳过 .zy 内部运行时桩文件目录（非用户业务代码）
			if d.Name() == ".zy" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(strings.ToLower(path), ".php") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
