package main

import (
	"os"

	"github.com/php-any/origami/data"
	gosupport "github.com/php-any/origami/examples/laravel13/go-support"
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
	args := os.Args[1:]
	if len(args) >= 2 && args[0] == "run" {
		runScript(args[1])
		return
	}
	runArtisan(args)
}

func runScript(path string) {
	vm, p := buildVM()
	_, ctl := vm.LoadAndRun(path)
	finish(vm, p, ctl)
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
	gosupport.Load(vm)

	return vm, p
}

// runArtisan 直接执行官方 artisan（Laravel 13 原生入口），由 Origami 解释。
func runArtisan(args []string) {
	vm, p := buildVM()

	// 把 argv 传给 PHP 侧：$_SERVER['argv'] / $argv = os.Args[1:]。
	// Symfony ArgvInput 会 array_shift 掉首个脚本名，因此必须保留 "artisan"。
	os.Args = append([]string{os.Args[0], "artisan"}, args...)

	_, ctl := vm.LoadAndRun("artisan")
	finish(vm, p, ctl)
}

// finish 对齐 zy CLI：exit/die 返回 ExitControl，按退出码结束；其它 Control 才当作错误展示。
func finish(vm *runtime.VM, p *parser.Parser, ctl data.Control) {
	if ctl != nil {
		if exit, ok := ctl.(data.ExitControl); ok && exit.IsExit() {
			vm.RunShutdownCallbacks()
			if code := exit.GetCode(); code != 0 {
				os.Exit(code)
			}
			return
		}
		p.ShowControl(ctl)
		vm.RunShutdownCallbacks()
		os.Exit(1)
	}
	vm.RunShutdownCallbacks()
}
