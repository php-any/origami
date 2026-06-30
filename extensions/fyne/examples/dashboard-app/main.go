package main

import (
	"fmt"
	"os"

	fyne "github.com/php-any/origami-fyne"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/system"

	dash "github.com/php-any/origami-fyne/examples/dashboard-app/build/gen"
)

func main() {
	p := parser.NewParser()
	vm := runtime.NewVM(p)

	std.Load(vm)
	php.Load(vm)
	system.Load(vm)
	fyne.Load(vm)

	dash.Register(vm)

	_, ctrl := vm.RunCompiledFile(dash.EntryPath)
	if data.FlushAllBuffersFn != nil {
		data.FlushAllBuffersFn()
	}
	if ctrl != nil {
		fmt.Fprintf(os.Stderr, "run failed\n")
		p.ShowControl(ctrl)
		os.Exit(1)
	}
	vm.RunShutdownCallbacks()
}
