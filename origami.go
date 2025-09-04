package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" // 使用 mysql
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/net/http"
	"github.com/php-any/origami/std/php"
	"github.com/php-any/origami/std/system"
)

func showHelp() {
	fmt.Println("折言(origami-lang) - 融合型脚本语言")
	fmt.Println()
	fmt.Println("用法: ./origami <脚本路径>")
	fmt.Println()
	fmt.Println("参数:")
	fmt.Println("  脚本路径    要执行的脚本文件路径")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  ./origami tests/run_tests.cjp")
	fmt.Println("  ./origami script.php")
	fmt.Println("  ./origami my_script.cjp")
	fmt.Println()
	fmt.Println("支持的脚本格式:")
	fmt.Println("  .cjp - 折言脚本文件")
	fmt.Println("  .php - PHP兼容脚本文件")
}

func main() {
	// 创建解析器
	p := parser.NewParser()
	// 创建全局命名空间
	p.AddScanNamespace("tests", "./tests")

	// 创建程序运行的环境
	vm := runtime.NewVM(p)
	std.Load(vm)
	php.Load(vm)
	http.Load(vm)
	system.Load(vm)

	// 检查命令行参数
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(0)
	}

	// 获取脚本路径参数
	scriptPath := os.Args[1]

	// 检查文件是否存在
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		_, _ = fmt.Fprintf(os.Stderr, "错误: 文件 '%s' 不存在\n\n", scriptPath)
		showHelp()
		os.Exit(1)
	}

	_, err := vm.LoadAndRun(scriptPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		if !parser.InLSP {
			panic(err)
		}
	}
}
