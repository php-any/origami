package main

import (
	"os"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/httpkernel"
)

// TestServeAutoloadAndBootstrap 需要在 examples/laravel13 目录下运行（依赖 vendor/ 与 bootstrap/app.php）。
func TestServeAutoloadAndBootstrap(t *testing.T) {
	if _, err := os.Stat("vendor/autoload.php"); err != nil {
		t.Skipf("需要 examples/laravel13 作为工作目录（composer install 后）: %v", err)
	}

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{oldArgs[0], "artisan", "package:discover", "--ansi"}
	node.ResetSuperglobals()

	base, _ := buildVM()
	var thrown data.Control
	base.SetThrowControl(func(control data.Control) { thrown = control })

	class, control := base.GetOrLoadClass("Illuminate\\Foundation\\Http\\Kernel")
	if control != nil {
		t.Fatalf("resolve preloaded HTTP kernel class: %v", control)
	}
	if _, ok := class.(*httpkernel.KernelClass); !ok {
		t.Fatalf("HTTP kernel class = %T, want preloaded Go class", class)
	}

	_, control = base.LoadAndRun("artisan")
	if exit, ok := control.(data.ExitControl); ok && exit.IsExit() && exit.GetCode() == 0 {
		control = nil
	}
	if control != nil {
		t.Fatalf("artisan package:discover: %v", control)
	}
	if thrown != nil {
		t.Fatalf("artisan package:discover throw: %v", thrown)
	}
	base.RunShutdownCallbacks()
}
