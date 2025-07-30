package http

import (
	"github.com/php-any/origami/std/tools"
	"testing"
)

func Test_newRequest(t *testing.T) {
	err := tools.Generate(&Request{}, &tools.GeneratorConfig{
		PackageName: "http",
		ClassName:   "Net\\\\Http\\\\Request",
		StructName:  "Request",
		OutputDir:   ".", // 生成到当前目录
	})
	if err != nil {
		t.Fatalf("生成代码失败: %v", err)
	}
}

func Test_newResponse(t *testing.T) {
	err := tools.Generate(&Response{}, &tools.GeneratorConfig{
		PackageName: "http",
		ClassName:   "Net\\\\Http\\\\Response",
		StructName:  "Response",
		OutputDir:   ".", // 生成到当前目录
	})
	if err != nil {
		t.Fatalf("生成代码失败: %v", err)
	}
}
