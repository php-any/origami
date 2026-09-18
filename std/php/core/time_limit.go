package core

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

var executionDeadlineNano atomic.Int64
var executionLimitSec atomic.Int64

// ExtraDeadlineExceeded / ExtraDeadlineAbort 由 php.Load 注入，用于 HTTP serve 打断卡死请求。
var ExtraDeadlineExceeded func() bool
var ExtraDeadlineAbort func()

// ExtraSetExecutionDeadline 若返回 true，表示截止时间已按请求隔离，勿再写进程级槽。
var ExtraSetExecutionDeadline func(seconds int) bool

// InHTTPRequest 为 true 时超时不能 os.Exit（会把常驻 HTTP 进程打死）。
var InHTTPRequest func() bool

// DefaultHTTPMaxExecutionTime 对齐 PHP 网页 SAPI 的 max_execution_time 默认值。
const DefaultHTTPMaxExecutionTime = 30

// SetExecutionDeadline 设置脚本最大执行截止时间（set_time_limit）。
func SetExecutionDeadline(seconds int) {
	if ExtraSetExecutionDeadline != nil && ExtraSetExecutionDeadline(seconds) {
		return
	}
	if seconds <= 0 {
		executionDeadlineNano.Store(0)
		executionLimitSec.Store(0)
		return
	}
	executionLimitSec.Store(int64(seconds))
	executionDeadlineNano.Store(time.Now().Add(time.Duration(seconds) * time.Second).UnixNano())
}

// CheckExecutionTimeLimit 在循环/语句边界检查是否超时，超时则向 stderr 输出 Fatal 并 os.Exit(1)。
func CheckExecutionTimeLimit(file string, line int) bool {
	extra := ExtraDeadlineExceeded != nil && ExtraDeadlineExceeded()
	dl := executionDeadlineNano.Load()
	phpLimit := dl != 0 && time.Now().UnixNano() >= dl
	if !extra && !phpLimit {
		return false
	}
	if file == "" {
		file = "Unknown"
	}
	if line <= 0 {
		line = 0
	}
	sec := int(executionLimitSec.Load())
	if sec <= 0 {
		sec = DefaultHTTPMaxExecutionTime
	}
	unit := ""
	if sec != 1 {
		unit = "s"
	}
	fmt.Fprintf(os.Stderr, "Fatal error: Maximum execution time of %d second%s exceeded in %s on line %d\n",
		sec, unit, file, line)
	httpReq := InHTTPRequest != nil && InHTTPRequest()
	if ExtraDeadlineAbort != nil && (extra || httpReq) {
		ExtraDeadlineAbort()
		return true
	}
	os.Exit(1)
	return true
}
