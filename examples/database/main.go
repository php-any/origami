package main

import (
	"fmt"

	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/database"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/system"
	_ "modernc.org/sqlite"
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

	_, err := vm.LoadAndRun("database.zy")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		if !parser.InLSP {
			panic(err)
		}
	}
}
