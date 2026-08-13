package compile

import (
	"fmt"
	"go/format"
	"os"
)

// writeFormattedGoFile 将 Go 源码写入文件并用 go/format 格式化。
func writeFormattedGoFile(path string, src []byte) error {
	formatted, err := format.Source(src)
	if err != nil {
		return fmt.Errorf("格式化 %s 失败: %w", path, err)
	}
	if err := os.WriteFile(path, formatted, 0644); err != nil {
		return fmt.Errorf("写入 %s 失败: %w", path, err)
	}
	return nil
}
