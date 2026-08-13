package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"

	"github.com/php-any/origami/cmd"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std"
	netannotation "github.com/php-any/origami/std/net/annotation"
	"github.com/php-any/origami/std/net/http"
	"github.com/php-any/origami/std/net/websocket"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/system"
)

func init() {
	cmd.SetRuntimeLoader(func(vm data.VM) {
		std.Load(vm)
		php.Load(vm)
		http.Load(vm)
		websocket.Load(vm)
		netannotation.Load(vm)
		system.Load(vm)
	})
}

func main() {
	if len(os.Args) > 1 && cmd.IsDirectScriptArg(os.Args[1]) {
		if err := cmd.RunScriptFile(os.Args[1]); err != nil {
			os.Exit(1)
		}
		return
	}

	if len(os.Args) == 1 {
		if err := cmd.RootHelp(); err != nil {
			os.Exit(1)
		}
		return
	}

	cmd.Execute()
}
