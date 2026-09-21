package gosupport

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

// ProcessStart 进程进入 Origami 代码的时刻（不含 go run 编译）。
var ProcessStart = time.Now()

var firstHTTPLogged atomic.Bool

// BootLog 启动阶段耗时（stderr）。用于区分 go 编译 / artisan boot / 首请求。
func BootLog(msg string) {
	fmt.Fprintf(os.Stderr, "[origami-boot] +%s  %s\n", time.Since(ProcessStart).Truncate(time.Millisecond), msg)
}
