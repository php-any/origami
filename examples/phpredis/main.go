// Run from the repository root: go run -mod=mod ./examples/phpredis
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/php/redis"
)

func main() {
	vm := runtime.NewVM(parser.NewParser()).(*runtime.VM)
	std.Load(vm)
	php.Load(vm)
	redis.Load(vm)
	script := "examples/phpredis/examples.php"
	if len(os.Args) > 1 {
		script = os.Args[1]
	}
	path, err := filepath.Abs(script)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if _, ctl := vm.LoadAndRun(path); ctl != nil {
		fmt.Fprintf(os.Stderr, "PHP execution failed: %v\n", ctl)
		os.Exit(1)
	}
}
