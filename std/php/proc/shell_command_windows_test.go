//go:build windows

package proc

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestShellCommandQuotedRedirection(t *testing.T) {
	target := filepath.Join(t.TempDir(), "quoted output.txt")
	cmd := shellCommand(context.Background(), `"`+os.Getenv("ComSpec")+`" /D /C echo composer>"`+target+`"`)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	content, err := os.ReadFile(target)
	if err != nil || strings.TrimSpace(string(content)) != "composer" {
		t.Fatalf("content=%q err=%v", content, err)
	}
}
