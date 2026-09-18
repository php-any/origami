//go:build !origamidebug

package perfmon

import "time"

// Enabled 报告当前二进制是否编入了性能监控。
func Enabled() bool { return false }

// Start 在生产构建中为零操作。
func Start() func() { return func() {} }

// Now 在生产构建中不读时钟。
func Now() time.Time { return time.Time{} }

// Since 在生产构建中恒为 0。
func Since(time.Time) time.Duration { return 0 }

// NoteParse 在生产构建中为零操作。
func NoteParse(string, bool, time.Duration) {}

// NoteClassLoad 在生产构建中为零操作（仅 autoload miss 路径会调用）。
func NoteClassLoad(string, time.Duration) {}

// NoteInclude 在生产构建中为零操作。
func NoteInclude(string, time.Duration) {}

// NoteFileRun 在生产构建中为零操作。
func NoteFileRun(string, time.Duration) {}

// Span 是一次请求/脚本的计时段；生产构建里 BeginRequest 返回 nil。
type Span struct{}

// BeginRequest 在生产构建中返回 nil。
func BeginRequest(string, string) *Span { return nil }

// End 允许 nil 接收者，生产路径为零操作。
func (*Span) End(string) {}
