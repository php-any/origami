package runtime

import (
	"testing"

	"github.com/php-any/origami/parser"
)

// TestVMFullOutputBufferMethods 验证新增的 ob_* 能力在 VM 上的完整行为。
func TestVMFullOutputBufferMethods(t *testing.T) {
	vm := NewVM(parser.NewParser()).(*VM)

	// 初始无缓冲。
	if n := vm.OutputBufferLevel(); n != 0 {
		t.Fatalf("expected level 0, got %d", n)
	}
	if _, ok := vm.OutputBufferLength(); ok {
		t.Fatalf("expected no buffer, got length ok")
	}

	vm.StartOutputBuffer()
	vm.WriteOutput("hello")

	// ob_get_length
	if n, ok := vm.OutputBufferLength(); !ok || n != 5 {
		t.Fatalf("ob_get_length = %d, ok=%v", n, ok)
	}

	// ob_clean（清空但保留缓冲）
	if !vm.CleanCurrentBuffer() {
		t.Fatalf("CleanCurrentBuffer returned false")
	}
	if c, _ := vm.OutputBufferContents(); c != "" {
		t.Fatalf("after clean contents = %q", c)
	}
	if n := vm.OutputBufferLevel(); n != 1 {
		t.Fatalf("after clean level = %d, want 1", n)
	}

	// ob_flush 等价：弹出写出到上层，再压入空缓冲（层级保持不变）。
	vm.WriteOutput("data")
	vm.StartOutputBuffer() // level 2
	vm.WriteOutput("inner")
	if _, ok := vm.FlushOutputBuffer(); !ok {
		t.Fatalf("FlushOutputBuffer returned not ok")
	}
	vm.StartOutputBuffer() // 还原 level 2
	if n := vm.OutputBufferLevel(); n != 2 {
		t.Fatalf("after flush level = %d, want 2", n)
	}

	// ob_get_status
	status := vm.OutputBufferStatus(true)
	if len(status) != 2 {
		t.Fatalf("OutputBufferStatus(true) len = %d, want 2", len(status))
	}
	if status[0].Level != 1 || status[1].Level != 2 {
		t.Fatalf("status levels = %d,%d want 1,2", status[0].Level, status[1].Level)
	}
	single := vm.OutputBufferStatus(false)
	if len(single) != 1 || single[0].Level != 2 {
		t.Fatalf("OutputBufferStatus(false) = %#v", single)
	}

	// ob_list_handlers
	handlers := vm.ListOutputHandlers()
	if len(handlers) != 2 {
		t.Fatalf("ListOutputHandlers len = %d, want 2", len(handlers))
	}

	// ob_implicit_flush
	if vm.IsImplicitFlush() {
		t.Fatalf("implicit flush should start false")
	}
	vm.SetImplicitFlush(true)
	if !vm.IsImplicitFlush() {
		t.Fatalf("implicit flush should be true after set")
	}
	vm.SetImplicitFlush(false)

	// 清理栈。
	vm.CleanOutputBuffer()
	vm.CleanOutputBuffer()
	if n := vm.OutputBufferLevel(); n != 0 {
		t.Fatalf("final level = %d, want 0", n)
	}
	if vm.hasOutputBuffer.Load() {
		t.Fatalf("hasOutputBuffer flag should be false after all popped")
	}
}

// TestVMFastPathNoBuffer 验证无缓冲时 WriteOutput 走零锁快速路径且标记正确。
func TestVMFastPathNoBuffer(t *testing.T) {
	vm := NewVM(parser.NewParser()).(*VM)
	if vm.hasOutputBuffer.Load() {
		t.Fatalf("initial hasOutputBuffer should be false")
	}
	// 直接输出（无缓冲）不应 panic。
	vm.WriteOutput("no-buffer")
	vm.StartOutputBuffer()
	if !vm.hasOutputBuffer.Load() {
		t.Fatalf("hasOutputBuffer should be true after start")
	}
	vm.CleanOutputBuffer()
	if vm.hasOutputBuffer.Load() {
		t.Fatalf("hasOutputBuffer should be false after clean")
	}
}
