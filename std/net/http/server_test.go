package http

import (
	"github.com/php-any/origami/std/tools"
	"testing"
)

func Test_newServer(t *testing.T) {
	err := tools.Generate(&Server{}, &tools.GeneratorConfig{
		PackageName: "http",
		ClassName:   "Net\\\\Http\\\\Server",
		StructName:  "Server",
		OutputDir:   ".", // 生成到当前目录
	})
	if err != nil {
		t.Fatalf("生成代码失败: %v", err)
	}
}
