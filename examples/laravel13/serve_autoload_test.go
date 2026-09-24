package main

import (
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/laravel/httpkernel"
	"github.com/php-any/origami/std/laravel/serve"
	httplib "github.com/php-any/origami/std/net/http"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/php/fpm"
	"github.com/php-any/origami/std/system"
)

// TestServeAutoloadAndBootstrap 需要在 examples/laravel13 目录下运行（依赖 vendor/ 与 bootstrap/app.php）。
func TestServeAutoloadAndBootstrap(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	autoload := filepath.Join(root, "vendor", "autoload.php")
	bootstrap := filepath.Join(root, "bootstrap", "app.php")
	if _, err := os.Stat(autoload); err != nil {
		t.Skipf("需要 examples/laravel13 作为工作目录（composer install 后）: %v", err)
	}

	p := parser.NewParser()
	base := runtime.NewVM(p).(*runtime.VM)
	base.ClearIncludeOnceCache()
	base.ClearPhpFileCache()
	base.SetThrowControl(func(data.Control) {})
	std.Load(base)
	php.Load(base)
	httplib.Load(base)
	system.Load(base)
	// 与 serve 路径一致：应用侧 HTTP 桥 + Go 版 ServeCommand（其余 Illuminate 类走 vendor autoload）。
	httpkernel.Load(base)
	serve.Load(base)

	rec := httptest.NewRecorder()
	reqVM := fpm.New(base, func(s string) { _, _ = io.WriteString(rec, s) })
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
