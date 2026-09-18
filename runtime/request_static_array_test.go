package runtime

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
)

func TestStaticArrayAppendDoesNotLeakAcrossHTTPRequests(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "StaticArrCow_Box.php")
	src := "<?php\nclass StaticArrCow_Box { public static array $items = []; }\n"
	if err := os.WriteFile(file, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}

	vm := NewVM(parser.NewParser()).(*VM)
	if _, ctl := vm.LoadAndRun(file); ctl != nil {
		t.Fatalf("LoadAndRun: %v", ctl)
	}
	cls, ok := vm.GetClass("StaticArrCow_Box")
	if !ok {
		t.Fatal("class not found")
	}
	getter, ok := cls.(data.GetStaticProperty)
	if !ok {
		t.Fatal("no GetStaticProperty")
	}

	push := func() {
		v, ok := getter.GetStaticProperty("items")
		if !ok {
			t.Fatal("items missing")
		}
		arr, ok := v.(*data.ArrayValue)
		if !ok {
			t.Fatalf("items type %T", v)
		}
		arr.AppendValue(data.NewStringValue("Login"))
	}
	count := func() int {
		v, ok := getter.GetStaticProperty("items")
		if !ok {
			return -1
		}
		arr, ok := v.(*data.ArrayValue)
		if !ok {
			return -1
		}
		return len(arr.List)
	}

	restore := BeginRequestOutput()
	push()
	if n := count(); n != 1 {
		restore()
		t.Fatalf("request1 count=%d", n)
	}
	restore()

	restore = BeginRequestOutput()
	defer restore()
	if n := count(); n != 0 {
		t.Fatalf("request2 leaked static array count=%d", n)
	}
}
