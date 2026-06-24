package runtime

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
)

func TestParseFileAcceptsEmptyArray(t *testing.T) {
	dir := t.TempDir()
	htmlPath := filepath.Join(dir, "page.html")
	const want = "<p>hello</p>"
	if err := os.WriteFile(htmlPath, []byte("<!DOCTYPE html><html><body><p>hello</p></body></html>"), 0o644); err != nil {
		t.Fatal(err)
	}

	p := parser.NewParser()
	vm := NewVM(p).(*VM)

	rendered, acl := vm.ParseFile(htmlPath, data.NewArrayValue(nil))
	if acl != nil {
		t.Fatalf("ParseFile with [] failed: %v", acl)
	}
	sv, ok := rendered.(*data.StringValue)
	if !ok {
		t.Fatalf("expected StringValue, got %T", rendered)
	}
	if !strings.Contains(sv.AsString(), want) {
		t.Fatalf("rendered %q does not contain %q", sv.AsString(), want)
	}
}

func TestParseFileBindsAssociativeArray(t *testing.T) {
	dir := t.TempDir()
	htmlPath := filepath.Join(dir, "page.html")
	if err := os.WriteFile(htmlPath, []byte("<!DOCTYPE html><html><body><p>{$greeting}</p></body></html>"), 0o644); err != nil {
		t.Fatal(err)
	}

	p := parser.NewParser()
	vm := NewVM(p).(*VM)

	arr := &data.ArrayValue{
		List: []*data.ZVal{
			data.NewNamedZVal("greeting", data.NewStringValue("from-array")),
		},
	}
	rendered, acl := vm.ParseFile(htmlPath, arr)
	if acl != nil {
		t.Fatalf("ParseFile with associative array failed: %v", acl)
	}
	sv, ok := rendered.(*data.StringValue)
	if !ok {
		t.Fatalf("expected StringValue, got %T", rendered)
	}
	if !strings.Contains(sv.AsString(), "from-array") {
		t.Fatalf("rendered %q does not contain injected greeting", sv.AsString())
	}
}
