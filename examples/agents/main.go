package main

import (
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"

	openai "github.com/php-any/origami-openai"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	cliannotation "github.com/php-any/origami/std/cli/annotation"
	netannotation "github.com/php-any/origami/std/net/annotation"
	httplib "github.com/php-any/origami/std/net/http"
	"github.com/php-any/origami/std/net/websocket"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/system"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	if os.Args[1] == "dev" {
		port := 8080
		if len(os.Args) > 2 {
			if p, err := strconv.Atoi(os.Args[2]); err == nil {
				port = p
			}
		}
		if err := runDevServer(port); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	runScript(os.Args[1], os.Args[2:])
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "用法:")
	fmt.Fprintln(os.Stderr, "  agents dev [port]        # Web 开发服务（热更新）")
	fmt.Fprintln(os.Stderr, "  agents index.php         # Web 生产服务")
	fmt.Fprintln(os.Stderr, "  agents run.php <command> # CLI 命令（pipeline/debate/handoff）")
}

// buildVM 组装运行时：与标准 zy 二进制一致，额外注册 OpenAI 扩展。
// 命名空间解析交给 PHP 层的 autoload.php（spl_autoload_register）。
func buildVM() (*runtime.VM, *parser.Parser) {
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)

	std.Load(vm)
	php.Load(vm)
	httplib.Load(vm)
	websocket.Load(vm)
	netannotation.Load(vm)
	system.Load(vm)
	openai.Load(vm)

	return vm, p
}

// runScript 执行脚本。若脚本后带有参数，则将首个参数作为 CLI 命令名分发。
func runScript(script string, rest []string) {
	vm, p := buildVM()

	if _, ctl := vm.LoadAndRun(script); ctl != nil {
		p.ShowControl(ctl)
		os.Exit(1)
	}

	if len(rest) > 0 {
		ctx := vm.CreateContext(nil)
		if ctl := cliannotation.ExecuteCommand(ctx, rest[0]); ctl != nil {
			p.ShowControl(ctl)
			os.Exit(1)
		}
	}

	vm.RunShutdownCallbacks()
}
