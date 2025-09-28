package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/sourcegraph/jsonrpc2"
)

func Test_handleTextDocumentDefinition_hello_function_jump(t *testing.T) {
	globalLspVM = NewLspVMWithScanDir("/Users/lvluo/Desktop/github.com/php-any/origami/docs/std")

	testURI := "file:///Users/lvluo/Desktop/github.com/php-any/tests/app/http.zy"
	filePath := "/Users/lvluo/Desktop/github.com/php-any/tests/app/http.zy"

	// 构造测试参数：点击位于 "echo hello();" 这一行（期望跳转到 hello() 函数定义）
	params := DefinitionParams{
		TextDocumentPositionParams: TextDocumentPositionParams{
			TextDocument: TextDocumentIdentifier{
				URI: testURI,
			},
			Position: Position{
				Line:      4, // 第12行（LSP 0-based），包含 "echo hello();"
				Character: 9, // 第10列（LSP 0-based），落在 "hello" 字符上
			},
		},
	}
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal params: %v", err)
	}

	req := &jsonrpc2.Request{
		ID:     jsonrpc2.ID{Num: 1},
		Method: "textDocument/definition",
		Params: (*json.RawMessage)(&paramsJSON),
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Skipf("Test file %s does not exist, skipping test", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read test file %s: %v", filePath, err)
	}

	parser := NewLspParser()
	parser.SetVM(globalLspVM)

	ast, err := parser.ParseFile(filePath)
	if err != nil {
		t.Fatalf("Failed to parse test file %s: %v", filePath, err)
	}

	documents[testURI] = &DocumentInfo{
		Content: string(content),
		Version: 1,
		AST:     ast,
		Parser:  parser,
	}

	got, err := handleTextDocumentDefinition(req)
	if err != nil {
		t.Fatalf("Function returned error: %v", err)
	}

	if got == nil {
		t.Fatalf("Expected definition result, got nil")
	}

	locations, ok := got.([]Location)
	if !ok {
		t.Fatalf("Expected []Location, got %T", got)
	}

	if len(locations) != 1 {
		t.Fatalf("Expected exactly 1 Location, got %d", len(locations))
	}

	loc := locations[0]
	if loc.URI != testURI {
		t.Fatalf("Expected URI %s, got %s", testURI, loc.URI)
	}

	// 期望指向 hello() 函数的定义范围：第 3 行到第 5 行（0-based 行号）
	// function hello() {
	//     return "hello world";
	// }
	if loc.Range.Start.Line != 2 || loc.Range.End.Line != 4 {
		t.Fatalf("Expected Range lines [2,4] for hello() function, got [%d,%d]", loc.Range.Start.Line, loc.Range.End.Line)
	}

	t.Logf("Successfully jumped to hello() function definition at lines [%d,%d]", loc.Range.Start.Line, loc.Range.End.Line)
}
