package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"

	gosupport "github.com/php-any/origami/examples/laravel/go-support"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	cliannotation "github.com/php-any/origami/std/cli/annotation"
	netannotation "github.com/php-any/origami/std/net/annotation"
	httplib "github.com/php-any/origami/std/net/http"
	"github.com/php-any/origami/std/net/websocket"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/system"
)

func main() {
	args := os.Args[1:]
	// ./laravel run path/to/script.php — 加载 go-support 后直接跑脚本（用于 illuminate 冒烟）
	if len(args) >= 2 && args[0] == "run" {
		runScript(args[1])
		return
	}
	runArtisan(args)
}

func runScript(path string) {
	vm, p := buildVM()
	if _, ctl := vm.LoadAndRun(path); ctl != nil {
		p.ShowControl(ctl)
		os.Exit(1)
	}
	vm.RunShutdownCallbacks()
}

func buildVM() (*runtime.VM, *parser.Parser) {
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)

	std.Load(vm)
	php.Load(vm)
	httplib.Load(vm)
	websocket.Load(vm)
	netannotation.Load(vm)
	cliannotation.Load(vm)
	system.Load(vm)
	gosupport.Load(vm)

	return vm, p
}

func runArtisan(args []string) {
	vm, p := buildVM()

	if _, ctl := vm.LoadAndRun("artisan.php"); ctl != nil {
		p.ShowControl(ctl)
		os.Exit(1)
	}

	cmd := "list"
	if len(args) > 0 {
		switch args[0] {
		case "help", "-h", "--help":
			cmd = "list"
		default:
			cmd = args[0]
		}
	}

	ctx := vm.CreateContext(nil)
	if ctl := cliannotation.ExecuteCommand(ctx, cmd); ctl != nil {
		p.ShowControl(ctl)
		os.Exit(1)
	}

	vm.RunShutdownCallbacks()
}
