package main

import (
	"os"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
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

	for _, name := range []string{"Illuminate\\Foundation\\Http\\Kernel", "Illuminate\\Foundation\\Console\\ServeCommand", "Illuminate\\Http\\Request", "Illuminate\\Support\\Collection", "Illuminate\\Pipeline\\Pipeline", "Symfony\\Component\\HttpFoundation\\Request", "Symfony\\Component\\HttpFoundation\\Response"} {
		class, preloaded := base.GetClass(name)
		if !preloaded {
			t.Fatalf("native component %s was not registered", name)
		}
		if _, interpreted := class.(*node.ClassStatement); interpreted {
			t.Fatalf("%s unexpectedly uses PHP before Composer", name)
		}
	}
	_, control := base.LoadAndRun("artisan")
	if exit, ok := control.(data.ExitControl); ok && exit.IsExit() && exit.GetCode() == 0 {
		control = nil
	}
	if control != nil {
		t.Fatalf("artisan package:discover: %v", control)
	}
	if thrown != nil {
		t.Fatalf("artisan package:discover throw: %v", thrown)
	}
	for _, name := range []string{"Illuminate\\Foundation\\Application", "Symfony\\Component\\Finder\\Finder"} {
		class, ctl := base.GetOrLoadClass(name)
		if ctl != nil {
			t.Fatal(ctl.AsString())
		}
		if _, official := class.(*node.ClassStatement); !official {
			t.Fatalf("%s did not use official PHP: %T", name, class)
		}
	}
	base.RunShutdownCallbacks()
}
