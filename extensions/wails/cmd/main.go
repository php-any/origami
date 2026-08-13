package main

import (
	"fmt"
	"os"

	wails "github.com/php-any/origami-wails"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/system"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "用法: wailsrunner <script.php>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "构建（须在 extensions/wails 目录下）:")
		fmt.Fprintln(os.Stderr, "  go build -o wailsrunner ./cmd/")
		fmt.Fprintln(os.Stderr, "  或: make build")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "示例:")
		fmt.Fprintln(os.Stderr, "  ./wailsrunner examples/hello.php")
		fmt.Fprintln(os.Stderr, "  ./wailsrunner examples/chat_demo.php")
		os.Exit(1)
	}

	p := parser.NewParser()
	vm := runtime.NewVM(p)

	std.Load(vm)
	php.Load(vm)
	system.Load(vm)
	wails.Load(vm)

	// 加载并运行入口脚本（Wails 示例为直接运行脚本，无 CLI 命令路由）
	_, ctl := vm.LoadAndRun(os.Args[1])
	if ctl != nil {
		p.ShowControl(ctl)
		os.Exit(1)
	}

	vm.RunShutdownCallbacks()
}
