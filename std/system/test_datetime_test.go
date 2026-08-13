package system

import (
	"testing"

	"github.com/php-any/origami/std/tools"
)

func TestDateTimeWrapper(t *testing.T) {
	err := tools.Generate(newDateTime(), &tools.GeneratorConfig{
		PackageName: "system",
		ClassName:   "System\\\\DateTime",
		StructName:  "DateTime",
		OutputDir:   ".", // 生成到当前目录
	})
	if err != nil {
		t.Fatalf("生成代码失败: %v", err)
	}
}
