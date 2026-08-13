package main

import (
	"runtime"
	"testing"
)

func TestUriToFilePath(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		expected string
		os       string
	}{
		{
			name:     "Unix file URI",
			uri:      "file:///home/user/test.txt",
			expected: "/home/user/test.txt",
			os:       "linux",
		},
		{
			name:     "Windows file URI",
			uri:      "file:///C:/Users/user/test.txt",
			expected: "C:\\Users\\user\\test.txt",
			os:       "windows",
		},
		{
			name:     "Non-file URI",
			uri:      "http://example.com",
			expected: "http://example.com",
			os:       "linux",
		},
	}

	originalOS := runtime.GOOS
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 注意：我们无法在测试中真正改变runtime.GOOS，
			// 所以这个测试主要验证当前系统的行为
			if runtime.GOOS != tt.os {
				t.Skipf("Skipping test for %s on %s", tt.os, runtime.GOOS)
			}
			result := uriToFilePath(tt.uri)
			if result != tt.expected {
				t.Errorf("uriToFilePath(%q) = %q, want %q", tt.uri, result, tt.expected)
			}
		})
	}
	_ = originalOS // 避免未使用变量警告
}

func TestFilePathToURI(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected string
		os       string
	}{
		{
			name:     "Unix file path",
			filePath: "/home/user/test.txt",
			expected: "file:///home/user/test.txt",
			os:       "linux",
		},
		{
			name:     "Windows file path",
			filePath: "C:\\Users\\user\\test.txt",
			expected: "file:///C:/Users/user/test.txt",
			os:       "windows",
		},
		{
			name:     "Empty path",
			filePath: "",
			expected: "",
			os:       "linux",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if runtime.GOOS != tt.os {
				t.Skipf("Skipping test for %s on %s", tt.os, runtime.GOOS)
			}
			result := filePathToURI(tt.filePath)
			if result != tt.expected {
				t.Errorf("filePathToURI(%q) = %q, want %q", tt.filePath, result, tt.expected)
			}
		})
	}
}

// TestCurrentSystemBehavior 测试当前系统的实际行为
func TestCurrentSystemBehavior(t *testing.T) {
	// 测试往返转换
	testCases := []string{
		"file:///tmp/test.txt",
		"file:///home/user/document.zy",
	}

	if runtime.GOOS == "windows" {
		testCases = []string{
			"file:///C:/temp/test.txt",
			"file:///D:/Users/user/document.zy",
		}
	}

	for _, uri := range testCases {
		t.Run(uri, func(t *testing.T) {
			// URI -> 文件路径 -> URI
			filePath := uriToFilePath(uri)
			backToURI := filePathToURI(filePath)

			if backToURI != uri {
				t.Errorf("Round trip failed: %s -> %s -> %s", uri, filePath, backToURI)
			}
		})
	}
}
