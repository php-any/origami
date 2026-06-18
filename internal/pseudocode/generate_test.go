package pseudocode

import (
	"os"
	"strings"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func TestFormatPHPTypeUnionClassNames(t *testing.T) {
	namespace := "Net\\Http"
	ty := data.NewUnionType([]data.Types{
		data.Class{Name: "Net\\Http\\Header"},
		data.Class{Name: "Net\\Http\\Response"},
	})
	got := formatPHPType(ty, namespace)
	want := "Header|Response"
	if got != want {
		t.Fatalf("formatPHPType() = %q, want %q", got, want)
	}
}

func TestFormatPHPTypeUnionScalars(t *testing.T) {
	namespace := "Net\\Http"
	ty := data.NewUnionType([]data.Types{
		data.NewBaseType("array"),
		data.NewBaseType("string"),
		data.NewBaseType("null"),
	})
	got := formatPHPType(ty, namespace)
	want := "array|string|null"
	if got != want {
		t.Fatalf("formatPHPType() = %q, want %q", got, want)
	}
}

func TestFormatPHPTypeMultipleReturnAsUnion(t *testing.T) {
	ty := data.NewMultipleReturnType([]data.Types{
		data.NewBaseType("object"),
		data.NewBaseType("array"),
	})
	got := formatPHPType(ty, "")
	want := "object|array"
	if got != want {
		t.Fatalf("formatPHPType() = %q, want %q", got, want)
	}
}

func TestGenerateHTTPReturnTypes(t *testing.T) {
	if err := Generate(dir); err != nil {
		t.Fatal(err)
	}

	requestPHP, err := os.ReadFile(dir + "/Net/Http/request.php")
	if err != nil {
		t.Fatal(err)
	}
	content := string(requestPHP)
	checks := []string{
		"function pathValue($param0) : string",
		"function all() : array",
		"function input($key) : array|string|null",
		"function header($key) : array|string",
		"function clone() : Request",
	}
	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Fatalf("request.php missing %q\n%s", check, content)
		}
	}

	responsePHP, err := os.ReadFile(dir + "/Net/Http/response.php")
	if err != nil {
		t.Fatal(err)
	}
	resp := string(responsePHP)
	if !strings.Contains(resp, "function header($key, $value) : Header|Response") {
		t.Fatalf("response.php header return type invalid:\n%s", resp)
	}
	if !strings.Contains(resp, "function write($param0) : Response") {
		t.Fatalf("response.php write return type invalid:\n%s", resp)
	}
	if !strings.Contains(resp, "function json(object|array $data)") {
		t.Fatalf("response.php json param type invalid:\n%s", resp)
	}
}

func TestAnalyzeAnnotationTargetParameter(t *testing.T) {
	param := analyzeParam(node.NewAnnotationTargetParameter(nil, 1), true, 1, "Container")
	if param == nil {
		t.Fatal("annotation target param should not be skipped")
	}
	if param.Name != "target" {
		t.Fatalf("param name = %q, want target", param.Name)
	}
	if param.Type != "?AstNode" {
		t.Fatalf("param type = %q, want ?AstNode", param.Type)
	}
	if param.Default != " = null" {
		t.Fatalf("param default = %q, want ' = null'", param.Default)
	}
}
