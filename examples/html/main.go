package main

import (
	"fmt"
	"os"

	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
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
	system.Load(vm)

	fmt.Println("Origami HTML解析器示例")
	fmt.Println("按 Ctrl+C 停止程序")
	fmt.Println()

	// 获取要执行的文件名
	filename := "html_features.zy"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	_, err := vm.LoadAndRun(filename)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		if !parser.InLSP {
			panic(err)
		}
	}
}
