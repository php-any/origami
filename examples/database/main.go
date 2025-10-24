package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql" // 使用 mysql
	_ "github.com/mattn/go-sqlite3"    // 使用 sqlite3
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/database"
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
	database.Load(vm)
	system.Load(vm)

	// 数据库连接将在脚本中设置

	fmt.Println("Origami Database扩展示例")
	fmt.Println("按 Ctrl+C 停止程序")
	fmt.Println()

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 在goroutine中运行脚本
	done := make(chan bool, 1)
	go func() {
		_, err := vm.LoadAndRun("database.zy")
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			if !parser.InLSP {
				panic(err)
			}
		}
		done <- true
	}()

	// 等待信号或脚本完成
	select {
	case <-sigChan:
		fmt.Println("\n收到停止信号，正在关闭程序...")
		os.Exit(0)
	case <-done:
		fmt.Println("脚本执行完成")
	}
}
