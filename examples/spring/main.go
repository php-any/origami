package main

import (
	"fmt"
	"github.com/php-any/origami/std/database/annotation"
	netAnnotation "github.com/php-any/origami/std/net/annotation"
	"os"
	"os/signal"
	"syscall"

	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/net/http"
	"github.com/php-any/origami/std/reflect"
	"github.com/php-any/origami/std/system"
)

func main() {
	p := parser.NewParser()
	p.AddScanNamespace("examples", "./")

	vm := runtime.NewVM(p)
	std.Load(vm)
	http.Load(vm)
	http.Load(vm)
	system.Load(vm)
	netAnnotation.Load(vm)
	reflect.Load(vm)
	annotation.Load(vm)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		_, err := vm.LoadAndRun("http.zy")
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			if !parser.InLSP {
				panic(err)
			}
		}
	}()

	fmt.Println("Spring 风格示例启动，按 Ctrl+C 退出")
	<-sigChan
	fmt.Println("\n收到停止信号，正在关闭服务器...")
}
