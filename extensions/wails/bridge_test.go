package wails

import (
	"os"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/php"
)

// captureFunc 是一个测试用 PHP 函数，捕获首个实参，供 Go 端断言。
type captureFunc struct{ captured data.Value }

func (f *captureFunc) Call(ctx data.Context) (data.GetValue, data.Control) {
	if v, ok := ctx.GetIndexValue(0); ok {
		f.captured = v
	}
	return nil, nil
}
func (f *captureFunc) GetName() string { return "__wails_capture" }
func (f *captureFunc) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "v", 0, nil, nil)}
}
func (f *captureFunc) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "v", 0, nil)}
}

// TestEmitTodoListConversion 验证 array_values($todos)（关联数组的列表）
// 经过 valueToGo 后得到 []map，而不是退化成字符串（之前 ObjectValue 缺失分支的 bug）。
func TestEmitTodoListConversion(t *testing.T) {
	vm := runtime.NewVM(parser.NewParser())
	std.Load(vm)
	php.Load(vm)
	Load(vm)

	cap := &captureFunc{}
	vm.AddFunc(cap)

	if _, ctl := vm.LoadAndRun("test_todo_capture.php"); ctl != nil {
		t.Fatalf("script failed: %v", ctl)
	}
	if cap.captured == nil {
		t.Fatal("value was not captured")
	}

	got := valueToGo(cap.captured)
	list, ok := got.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T: %v", got, got)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 todos, got %d", len(list))
	}

	first, ok := list[0].(map[string]any)
	if !ok {
		t.Fatalf("element 0 is not a map (likely degraded to string): %T = %v", list[0], list[0])
	}
	if first["text"] != "hello" {
		t.Fatalf("todo[0].text = %v, want hello", first["text"])
	}
	if first["done"] != false {
		t.Fatalf("todo[0].done = %v, want false", first["done"])
	}

	second, _ := list[1].(map[string]any)
	if second["text"] != "world" || second["done"] != true {
		t.Fatalf("todo[1] mismatch: %v", second)
	}
}

func TestEmitChatMessageConversion(t *testing.T) {
	vm := runtime.NewVM(parser.NewParser())
	std.Load(vm)
	php.Load(vm)
	Load(vm)

	cap := &captureFunc{}
	vm.AddFunc(cap)

	script := `<?php
$msg = ['id' => 'm1', 'channel' => 'general', 'author' => 'Alice', 'text' => '你好', 'type' => 'user', 'time' => '12:00'];
__wails_capture($msg);
`
	path := "test_chat_msg_capture.php"
	if err := os.WriteFile(path, []byte(script), 0644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(path) })

	if _, ctl := vm.LoadAndRun(path); ctl != nil {
		t.Fatalf("script failed: %v", ctl)
	}

	got := valueToGo(cap.captured)
	m, ok := got.(map[string]any)
	if !ok {
		t.Fatalf("expected map, got %T: %v", got, got)
	}
	if m["text"] != "你好" || m["author"] != "Alice" {
		t.Fatalf("message fields wrong: %v", m)
	}
}
