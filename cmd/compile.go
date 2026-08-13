package cmd

import (
	"github.com/php-any/origami/cmd/compile"
	"github.com/php-any/origami/data"
)

var compileCmd = compile.NewCommand(func(vm data.VM) {
	if runtimeLoader != nil {
		runtimeLoader(vm)
	}
})
