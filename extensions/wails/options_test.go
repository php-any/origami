package wails

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/php"
)

func TestOptionsAppStoresHTMLFromObjectValue(t *testing.T) {
	vm := runtime.NewVM(parser.NewParser())
	php.Load(vm)
	Load(vm)

	classStmt, ok := vm.GetClass("Wails\\Options\\App")
	if !ok {
		t.Fatal("App class not found")
	}

	baseCtx := vm.CreateContext(nil)
	gv, ctl := classStmt.GetValue(baseCtx)
	if ctl != nil {
		t.Fatalf("GetValue failed: %v", ctl)
	}
	appCV, ok := gv.(*data.ClassValue)
	if !ok {
		t.Fatalf("expected ClassValue, got %T", gv)
	}
	appCV.SetVM(vm)

	htmlStr := "<title>UNIQUE_MARKER_OBJECT</title>"
	opts := data.NewObjectValue()
	opts.SetProperty("HTML", data.NewStringValue(htmlStr))
	opts.SetProperty("Title", data.NewStringValue("Probe"))

	construct := classStmt.GetConstruct()
	methodCtx := appCV.CreateContext(construct.GetVariables())
	methodCtx.SetIndexZVal(0, data.NewZVal(opts))

	if _, ctl := construct.Call(methodCtx); ctl != nil {
		t.Fatalf("construct failed: %v", ctl)
	}

	html := getPropString(appCV, "HTML", "")
	if html != htmlStr {
		t.Fatalf("HTML mismatch: got %q want %q", html, htmlStr)
	}
}

func TestOptionsAppStoresHTMLFromArray(t *testing.T) {
	vm := runtime.NewVM(parser.NewParser())
	php.Load(vm)
	Load(vm)

	classStmt, ok := vm.GetClass("Wails\\Options\\App")
	if !ok {
		t.Fatal("App class not found")
	}

	baseCtx := vm.CreateContext(nil)
	gv, ctl := classStmt.GetValue(baseCtx)
	if ctl != nil {
		t.Fatalf("GetValue failed: %v", ctl)
	}
	appCV, ok := gv.(*data.ClassValue)
	if !ok {
		t.Fatalf("expected ClassValue, got %T", gv)
	}
	appCV.SetVM(vm)

	htmlStr := "<title>UNIQUE_MARKER_TEST</title>"
	opts := &data.ArrayValue{
		List: []*data.ZVal{
			data.NewNamedZVal("HTML", data.NewStringValue(htmlStr)),
			data.NewNamedZVal("Title", data.NewStringValue("Probe")),
		},
	}

	construct := classStmt.GetConstruct()
	methodCtx := appCV.CreateContext(construct.GetVariables())
	methodCtx.SetIndexZVal(0, data.NewZVal(opts))

	if _, ctl := construct.Call(methodCtx); ctl != nil {
		t.Fatalf("construct failed: %v", ctl)
	}

	html := getPropString(appCV, "HTML", "")
	if html != htmlStr {
		t.Fatalf("HTML mismatch: got %q want %q", html, htmlStr)
	}
}

func TestOptionsAppHTMLFromPHPHeredoc(t *testing.T) {
	vm := runtime.NewVM(parser.NewParser())
	std.Load(vm)
	php.Load(vm)
	Load(vm)

	probe := filepath.Join(".html_probe.txt")
	_ = os.Remove(probe)

	script := filepath.Join("test_heredoc_html.php")
	if _, ctl := vm.LoadAndRun(script); ctl != nil {
		t.Fatalf("script failed: %v", ctl)
	}

	raw, err := os.ReadFile(probe)
	if err != nil {
		t.Fatalf("probe file missing: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(probe) })

	got := strings.TrimSpace(string(raw))
	t.Logf("html_len=%d", len(got))
	if !strings.Contains(got, "FROM_HEREDOC_PROBE") {
		t.Fatalf("App HTML not stored from PHP constructor: %q", got[:min(80, len(got))])
	}
}
