package main

import (
	"fmt"
	"os"

	"github.com/php-any/origami/parser"
)

func main() {
	// 1) 已知样例：验证二进制是否加载了最新的 preprocessor 代码
	sample := "<?php if(a): echo 1; endif; ?><?php foreach($x as $y): endforeach; ?>"
	convertedSample := parser.ConvertAltPHPSyntaxForTest("sample.php", sample)
	fmt.Printf("[SAMPLE INPUT ] %q\n", sample)
	fmt.Printf("[SAMPLE OUTPUT] %q\n", convertedSample)

	// 2) 真实文件
	file := os.Args[1]
	content, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	code := string(content)
	fmt.Printf("[REAL] len=%d hasControlColon=%v\n", len(code), parser.HasControlColonForTest(code))
	converted := parser.ConvertAltPHPSyntaxForTest(file, code)
	if converted == code {
		fmt.Println("[REAL] NO CHANGE (conversion is no-op)")
	} else {
		fmt.Println("[REAL] CONVERTED (differs from input)")
	}
	os.Stdout.WriteString(converted)
	fmt.Println("\n=== END CONVERTED ===")
}
