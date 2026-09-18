package utils

import (
	"path/filepath"
	"testing"
)

func TestNormalizePhpFilePathCached(t *testing.T) {
	self, err := filepath.Abs("php_file.go")
	if err != nil {
		t.Fatal(err)
	}
	first := NormalizePhpFilePath(self)
	if first == "" {
		t.Fatal("规范化路径为空")
	}
	second := NormalizePhpFilePath(self)
	if first != second {
		t.Fatalf("缓存后路径不一致: %q vs %q", first, second)
	}
	again := NormalizePhpFilePath(first)
	if again != first {
		t.Fatalf("已规范化路径应原样命中缓存: %q vs %q", first, again)
	}
	if !SamePhpFile(self, first) {
		t.Fatal("同一文件 SamePhpFile 应为 true")
	}
	if !SamePhpFile(self, self) {
		t.Fatal("相同字符串 SamePhpFile 应为 true")
	}
	if SamePhpFile("", "") {
		t.Fatal("空路径不应视为同一文件")
	}
}
