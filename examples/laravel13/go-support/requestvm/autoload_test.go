package requestvm

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	rt "github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/system"
)

func laravelRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller failed")
	}
	// .../go-support/requestvm -> laravel13
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func TestLoadComposerAutoloadViaRequestVM(t *testing.T) {
	root := laravelRoot(t)
	autoload := filepath.Join(root, "vendor", "autoload.php")

	p := parser.NewParser()
	base := rt.NewVM(p).(*rt.VM)
	base.ClearIncludeOnceCache()
	base.ClearPhpFileCache()
	base.SetThrowControl(func(data.Control) {})
	std.Load(base)
	php.Load(base)
	system.Load(base)

	req := New(base, data.DefaultOutputWriter)
	result, control := req.LoadAndRun(autoload)
	if thrown := req.TakeThrow(); control == nil && thrown != nil {
		control = thrown
	}
	if control != nil {
		t.Fatalf("autoload control: %v", control)
	}
	if result == nil {
		t.Fatal("autoload returned nil")
	}
	if _, ok := base.GetFunc("Illuminate\\Filesystem\\join_paths"); !ok {
		t.Fatal("composer autoload did not register Illuminate\\Filesystem\\join_paths")
	}
	found := false
	for _, c := range base.AllClasses() {
		if len(c.GetName()) >= 22 && c.GetName()[:22] == "ComposerAutoloaderInit" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("ComposerAutoloaderInit missing after RequestVM LoadAndRun")
	}
}

func TestLoadComposerAutoloadViaRequestVMWithLaravelLoad(t *testing.T) {
	root := laravelRoot(t)
	autoload := filepath.Join(root, "vendor", "autoload.php")

	p := parser.NewParser()
	base := rt.NewVM(p).(*rt.VM)
	base.ClearIncludeOnceCache()
	base.ClearPhpFileCache()
	base.SetThrowControl(func(data.Control) {})
	std.Load(base)
	php.Load(base)
	system.Load(base)

	req := New(base, data.DefaultOutputWriter)
	result, control := req.LoadAndRun(autoload)
	if thrown := req.TakeThrow(); control == nil && thrown != nil {
		control = thrown
	}
	if control != nil {
		t.Fatalf("autoload control: %v", control)
	}
	if result == nil {
		t.Fatal("autoload returned nil")
	}
}
