package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql" // 使用 mysql
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/net/http"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/system"
)

func main() {
	// 创建解析器
	p := parser.NewParser()
	// 创建全局命名空间
	p.AddScanNamespace("examples", "./")

	// 创建程序运行的环境
	vm := runtime.NewVM(p)
	std.Load(vm)
	php.Load(vm)
	http.Load(vm)
	system.Load(vm)

	fmt.Println()

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 在goroutine中运行脚本
	go func() {
		_, err := vm.LoadAndRun("http.zy")
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			if !parser.InLSP {
				panic(err)
			}
		}
	}()

	fmt.Println("Origami HTTP扩展示例")
	fmt.Println("按 Ctrl+C 停止服务器")
	// 等待信号
	<-sigChan
	fmt.Println("\n收到停止信号，正在关闭服务器...")
	os.Exit(0)
}
