package wails

import (
	"os"
	"strings"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std/php"
)

func TestKeysCmdOrCtrlFromPHP(t *testing.T) {
	vm := runtime.NewVM(parser.NewParser())
	php.Load(vm)
	Load(vm)

	probe := ".keys_probe.txt"
	_ = os.Remove(probe)
	t.Cleanup(func() { _ = os.Remove(probe) })

	if _, ctl := vm.LoadAndRun("test_keys_cmd.php"); ctl != nil {
		t.Fatalf("script failed: %v", ctl)
	}
	raw, err := os.ReadFile(probe)
	if err != nil {
		t.Fatalf("probe missing: %v", err)
	}
	got := string(raw)
	if got != "cmdorctrl+o" {
		t.Fatalf("Keys::cmdOrCtrl returned %q", got)
	}
}

func TestCanonicalAccelerator(t *testing.T) {
	got := canonicalAccelerator("cmdorctrl+o")
	if got != "Cmd+O" && got != "Ctrl+O" {
		t.Fatalf("unexpected canonical accelerator: %q", got)
	}
	combo := canonicalAccelerator("shift+z")
	if !strings.Contains(combo, "Shift") || !strings.Contains(strings.ToUpper(combo), "Z") {
		t.Fatalf("unexpected combo accelerator: %q", combo)
	}
}

func TestModifierPrefixFromObjectValue(t *testing.T) {
	ov := data.NewObjectValue()
	ov.SetProperty("0", data.NewStringValue("shift"))
	prefix := modifierPrefixFromValue(ov)
	if prefix != "shift+" {
		t.Fatalf("expected shift+ prefix, got %q", prefix)
	}
}
