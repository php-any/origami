package main

import (
	"fmt"
	"github.com/php-any/generator"
	"net/http"
	"os"
)

func main() {
	array := []any{
		http.NewServeMux,
	}

	outRoot := "generator_build"
	for _, elem := range array {
		if err := generator.GenerateFromConstructor(elem, generator.GenOptions{OutputRoot: outRoot}); err != nil {
			fmt.Fprintln(os.Stderr, "生成失败:", err)
			continue
		}
		fmt.Println("生成完成 ->", outRoot)
	}
}
