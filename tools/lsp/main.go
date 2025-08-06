package main

import (
	"context"
	"log"
	"os"

	"github.com/php-any/origami/tools/lsp/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()

	// 创建 LSP 服务器
	srv := server.NewServer()

	// 启动服务器
	if err := srv.Start(ctx, os.Stdin, os.Stdout); err != nil {
		log.Fatalf("LSP 服务器启动失败: %v", err)
	}
}
