package runtime

import (
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
)

func TestTempVMOnlyIsolatesOutputState(t *testing.T) {
	base := NewVM(parser.NewParser()).(*VM)
	first := NewTempVM(base).(*TempVM)
	second := NewTempVM(base).(*TempVM)

	first.EnsureGlobalZVal("request_id").Value = data.NewStringValue("first")
	if got := second.EnsureGlobalZVal("request_id").Value.AsString(); got != "first" {
		t.Fatalf("PHP global should be shared through base VM: %q", got)
	}

	first.StartOutputBuffer()
	first.WriteOutput("first")
	second.StartOutputBuffer()
	second.WriteOutput("second")
	if got, _ := first.CleanOutputBuffer(); got != "first" {
		t.Fatalf("first output buffer = %q", got)
	}
	if got, _ := second.CleanOutputBuffer(); got != "second" {
		t.Fatalf("second output buffer = %q", got)
	}

	first.PushCallFrame(data.CallFrame{Function: "first"})
	defer first.PopCallFrame()
	if got := second.SnapshotCallStack(); len(got) != 1 || got[0].Function != "first" {
		t.Fatalf("call stack should be delegated to base VM: %#v", got)
	}
}
