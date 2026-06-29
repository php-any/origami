package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"

	openai "github.com/php-any/origami-openai"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/cli"
	"github.com/php-any/origami/std/cli/annotation"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/system"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "用法: testrunner <script.php> [command]")
		fmt.Fprintln(os.Stderr, "  ./testrunner examples/run.php chat")
		fmt.Fprintln(os.Stderr, "  ./testrunner examples/run.php all")
		os.Exit(1)
	}

	p := parser.NewParser()
	vm := runtime.NewVM(p)

	std.Load(vm)
	php.Load(vm)
	system.Load(vm)
	cli.Load(vm)
	openai.Load(vm)

	// 加载入口脚本，触发 #[CliApplication] + #[Command] 注解注册
	_, ctl := vm.LoadAndRun(os.Args[1])
	if ctl != nil {
		p.ShowControl(ctl)
		os.Exit(1)
	}

	// 命令路由：os.Args[2] 是子命令名
	cmd := "all"
	if len(os.Args) > 2 {
		cmd = os.Args[2]
	}
	ctx := vm.CreateContext(nil)
	if ctl := annotation.ExecuteCommand(ctx, cmd); ctl != nil {
		p.ShowControl(ctl)
		os.Exit(1)
	}

	vm.RunShutdownCallbacks()
}
