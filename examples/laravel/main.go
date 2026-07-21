package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"

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
	runArtisan(os.Args[1:])
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
