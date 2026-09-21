package main

import (
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"

	"github.com/php-any/origami/data"
	gosupport "github.com/php-any/origami/examples/laravel13/go-support"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/perfmon"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	netannotation "github.com/php-any/origami/std/net/annotation"
	httplib "github.com/php-any/origami/std/net/http"
	"github.com/php-any/origami/std/net/websocket"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/system"
	"github.com/php-any/origami/std/vendoraccel"
)

var stopPerf = func() {}

func main() {
	stopPerf = perfmon.Start()
	defer stopPerf()

	args := os.Args[1:]
	if len(args) >= 2 && args[0] == "run" {
		runScript(args[1])
		return
	}
	if len(args) >= 1 && args[0] == "serve" {
		runServe()
		return
	}
	runArtisan(args)
}

func runServe() {
	node.ResetSuperglobals()
	t0 := time.Now()
	gosupport.BootLog("process entered main (HTTP serve)")
	vm, p := buildVM()
	gosupport.BootLog("buildVM+std+vendoraccel done (" + time.Since(t0).Truncate(time.Millisecond).String() + ")")
	if err := gosupport.RunHTTPServe(vm); err != nil {
		fmt.Fprintf(os.Stderr, "serve: %v\n", err)
		p.ShowControl(data.NewErrorThrow(nil, err))
		vm.RunShutdownCallbacks()
		stopPerf()
		os.Exit(1)
	}
	vm.RunShutdownCallbacks()
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
	vendoraccel.Load(vm)
	gosupport.Load(vm)

	return vm, p
}

// runArtisan 直接执行官方 artisan（Laravel 13 原生入口），由 Origami 解释。
func runArtisan(args []string) {
	// Symfony ArgvInput 会 array_shift 掉首个脚本名，因此必须保留 "artisan"。
	// 必须在 buildVM / putenv / $_SERVER 初始化之前改写，否则 $_SERVER['argv'] 会少命令名，
	// `serve` 会被当成默认的 list。
	os.Args = append([]string{os.Args[0], "artisan"}, args...)
	node.ResetSuperglobals()

	t0 := time.Now()
	gosupport.BootLog("process entered main")
	vm, p := buildVM()
	gosupport.BootLog("buildVM+std+vendoraccel done (" + time.Since(t0).Truncate(time.Millisecond).String() + ")")

	if vendoraccel.ShouldWarmup() {
		if root, err := os.Getwd(); err == nil {
			vendoraccel.WarmupVendorClassmap(vm, root)
		}
	}

	_, ctl := vm.LoadAndRun("artisan")
	finish(vm, p, ctl)
}

// finish 对齐 zy CLI：exit/die 返回 ExitControl，按退出码结束；其它 Control 才当作错误展示。
func finish(vm *runtime.VM, p *parser.Parser, ctl data.Control) {
	if ctl != nil {
		if exit, ok := ctl.(data.ExitControl); ok && exit.IsExit() {
			vm.RunShutdownCallbacks()
			if code := exit.GetCode(); code != 0 {
				stopPerf()
				os.Exit(code)
			}
			return
		}
		p.ShowControl(ctl)
		vm.RunShutdownCallbacks()
		stopPerf()
		os.Exit(1)
	}
	vm.RunShutdownCallbacks()
}
