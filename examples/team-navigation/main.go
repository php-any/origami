package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	netAnnotation "github.com/php-any/origami/std/net/annotation"
	"github.com/php-any/origami/std/net/http"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/system"
	_ "modernc.org/sqlite"
)

func main() {
	// 创建解析器
	p := parser.NewParser()
	// 创建全局命名空间
	p.AddScanNamespace("App", "./src")

	// 创建程序运行的环境
	vm := runtime.NewVM(p)
	std.Load(vm)
	php.Load(vm)
	http.Load(vm)
	system.Load(vm)
	netAnnotation.Load(vm)

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("按 Ctrl+C 停止服务器")
	// 在 goroutine 中运行脚本
	go func() {
		_, err := vm.LoadAndRun("http.zy")
		if err != nil {
			p.ShowControl(err)
		}
	}()

	<-sigChan
	fmt.Println("\n收到停止信号，正在关闭服务器...")
}
