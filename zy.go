package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"

	"github.com/php-any/origami/cmd"
)

func init() {
	cmd.SetRuntimeLoader(loadDefaultRuntime)
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
