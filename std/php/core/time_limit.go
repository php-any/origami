package core

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var timeLimitMu sync.Mutex
var executionDeadline time.Time
var executionLimitSec int

// SetExecutionDeadline 设置脚本最大执行截止时间（set_time_limit）。
func SetExecutionDeadline(seconds int) {
	timeLimitMu.Lock()
	defer timeLimitMu.Unlock()
	if seconds <= 0 {
		executionDeadline = time.Time{}
		executionLimitSec = 0
		return
	}
	executionLimitSec = seconds
	executionDeadline = time.Now().Add(time.Duration(seconds) * time.Second)
}

// CheckExecutionTimeLimit 在循环/语句边界检查是否超时，超时则向 stderr 输出 Fatal 并 os.Exit(1)。
func CheckExecutionTimeLimit(file string, line int) bool {
	timeLimitMu.Lock()
	deadline := executionDeadline
	timeLimitMu.Unlock()
	if deadline.IsZero() || time.Now().Before(deadline) {
		return false
	}
	if file == "" {
		file = "Unknown"
	}
	if line <= 0 {
		line = 0
	}
	sec := executionLimitSec
	if sec <= 0 {
		sec = 1
	}
	unit := ""
	if sec != 1 {
		unit = "s"
	}
	fmt.Fprintf(os.Stderr, "Fatal error: Maximum execution time of %d second%s exceeded in %s on line %d\n",
		sec, unit, file, line)
	os.Exit(1)
	return true
}
