package compile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteFormattedGoFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.go")
	src := []byte("package main\n\nimport(\"fmt\")\n\nfunc main(){fmt.Println(1)}\n")
	if err := writeFormattedGoFile(path, src); err != nil {
		t.Fatal(err)
	}
	out, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	formatted := string(out)
	if !strings.Contains(formatted, "import (\n") {
		t.Fatalf("expected formatted imports, got:\n%s", formatted)
	}
	if strings.Contains(formatted, "import(\"") {
		t.Fatalf("unformatted import block remains:\n%s", formatted)
	}
}
