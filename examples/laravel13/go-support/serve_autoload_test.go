package gosupport_test

import (
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	gosupport "github.com/php-any/origami/examples/laravel13/go-support"
	"github.com/php-any/origami/examples/laravel13/go-support/requestvm"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	httplib "github.com/php-any/origami/std/net/http"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/system"
	"github.com/php-any/origami/data"
)

func TestServeAutoloadAndBootstrap(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// test cwd is laravel13 when running ./go-support/...
	if filepath.Base(root) == "go-support" {
		root = filepath.Dir(root)
	}
	autoload := filepath.Join(root, "vendor", "autoload.php")
	bootstrap := filepath.Join(root, "bootstrap", "app.php")

	p := parser.NewParser()
	base := runtime.NewVM(p).(*runtime.VM)
	base.ClearIncludeOnceCache()
	base.ClearPhpFileCache()
	base.SetThrowControl(func(data.Control) {})
	std.Load(base)
	php.Load(base)
	httplib.Load(base)
	system.Load(base)
	gosupport.Load(base)

	rec := httptest.NewRecorder()
	reqVM := requestvm.New(base, func(s string) { _, _ = io.WriteString(rec, s) })
	r := httptest.NewRequest("GET", "http://127.0.0.1/", nil)
	reqVM.BindHTTP(r, rec)

	if _, control := reqVM.LoadAndRun(autoload); control != nil {
		t.Fatalf("autoload: %v", control)
	}
	if thrown := reqVM.TakeThrow(); thrown != nil {
		t.Fatalf("autoload throw: %v", thrown)
	}

	appValue, control := reqVM.LoadAndRunFresh(bootstrap)
	if control != nil {
		t.Fatalf("bootstrap: %v", control)
	}
	if thrown := reqVM.TakeThrow(); thrown != nil {
		t.Fatalf("bootstrap throw: %v", thrown)
	}
	if _, ok := appValue.(*data.ClassValue); !ok {
		t.Fatalf("bootstrap returned %T", appValue)
	}
}
