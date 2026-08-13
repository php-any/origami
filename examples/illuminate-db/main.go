package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	netannotation "github.com/php-any/origami/std/net/annotation"
	httplib "github.com/php-any/origami/std/net/http"
	"github.com/php-any/origami/std/net/websocket"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/system"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: ./illuminate-db <script.php>")
		fmt.Println("示例: ./illuminate-db examples/01_basic_connection.php")
		os.Exit(1)
	}

	vm, p := buildVM()
	path := os.Args[1]

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
	system.Load(vm)

	return vm, p
}

var _ = data.VM(nil)
