package runtime

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
)

func TestRequestVMGlobalsOutputAndCallStateIsolated(t *testing.T) {
	base := NewVM(parser.NewParser()).(*VM)
	first := NewRequestVM(base).(*RequestVM)
	second := NewRequestVM(base).(*RequestVM)

	first.EnsureGlobalZVal("request_id").StoreRaw(data.NewStringValue("first"))
	if got := second.EnsureGlobalZVal("request_id").ReadValue().AsString(); got != "" {
		t.Fatalf("request global leaked: %q", got)
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
	if got := second.SnapshotCallStack(); len(got) != 0 {
		t.Fatalf("call stack should be isolated per RequestVM: %#v", got)
	}
	if got := first.SnapshotCallStack(); len(got) != 1 || got[0].Function != "first" {
		t.Fatalf("first RequestVM call stack = %#v", got)
	}
}

func TestExplicitRequestOutputOwnsItsExecutionState(t *testing.T) {
	base := NewVM(parser.NewParser()).(*VM)
	defer BeginRequestOutput()()
	first := NewRequestVM(base).(*RequestVM)
	second := NewRequestVM(base).(*RequestVM)
	var left, right strings.Builder
	first.SetOutputWriter(func(value string) { left.WriteString(value) })
	second.SetOutputWriter(func(value string) { right.WriteString(value) })
	first.WriteOutput("first")
	second.WriteOutput("second")
	if left.String() != "first" || right.String() != "second" {
		t.Fatal("ambient HTTP output overrode explicit request writers")
	}
	first.EnterCall()
	defer first.LeaveCall()
	if first.requestCall().Depth != 1 || second.requestCall().Depth != 0 || currentRequestCallState().Depth != 0 {
		t.Fatal("explicit requests borrowed ambient execution state")
	}
	first.SetErrorHandler(data.NewStringValue("request-error"))
	first.SetExceptionHandler(data.NewStringValue("request-exception"))
	if second.GetErrorHandler() != nil || second.GetExceptionHandler() != nil || base.GetErrorHandler() != nil || base.GetExceptionHandler() != nil {
		t.Fatal("explicit request handlers escaped to the ambient scope")
	}
	if call, ok := first.CreateContext(nil).(*Context); !ok || call.call != &first.call || call.out != first.out {
		t.Fatal("request context did not retain its own execution state")
	}
}

func TestIncludeUsesActiveRequestScopeWithRetainedContext(t *testing.T) {
	base := NewVM(parser.NewParser()).(*VM)
	retained := base.CreateContext(nil).(*Context)
	retained.call = &CallState{}
	file := filepath.Join(t.TempDir(), "view.php")
	if err := os.WriteFile(file, []byte("<?php $viewLocal = ['value']; return $viewLocal;"), 0600); err != nil {
		t.Fatal(err)
	}
	defer BeginRequestOutput()()
	request := NewRequestVM(base).(*RequestVM)
	request.EnterCall()
	defer request.LeaveCall()
	if IncludeBindsToProcessGlobals(retained) {
		t.Fatal("retained zero-depth context overrode active function scope")
	}
	if _, ctl := request.LoadInCallerContext(retained, file); ctl != nil {
		t.Fatal(ctl.AsString())
	}
	if _, found := request.globalVars["viewLocal"]; found {
		t.Fatal("included function local became a request global")
	}
}

func TestRequestIncludeBootstrapSnapshotIsImmutable(t *testing.T) {
	base := NewVM(parser.NewParser()).(*VM)
	base.SetPhpFileCache("boot-one.php")
	first := NewRequestVM(base).(*RequestVM)
	second := NewRequestVM(base).(*RequestVM)
	if first.initialFiles != second.initialFiles {
		t.Fatal("bootstrap snapshot copied for every request")
	}
	base.SetPhpFileCache("boot-two.php")
	third := NewRequestVM(base).(*RequestVM)
	if !first.GetPhpFileCache("boot-one.php") || first.GetPhpFileCache("boot-two.php") || !third.GetPhpFileCache("boot-two.php") {
		t.Fatal("bootstrap generations mixed")
	}
	first.SetPhpFileCache("request-one.php")
	if second.GetPhpFileCache("request-one.php") || base.GetPhpFileCache("request-one.php") {
		t.Fatal("request include state escaped")
	}
}
