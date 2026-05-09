package data

import "fmt"

// OutputWriter 是输出写入函数类型
type OutputWriter func(string)

// DefaultOutputWriter 默认输出到 stdout（供 ob_* 函数恢复时使用）
func DefaultOutputWriter(s string) {
	fmt.Print(s)
}

// WriteOutput 是当前的输出写入函数
// 由 ob_start/ob_get_clean 切换输出目标
var WriteOutput OutputWriter = DefaultOutputWriter

// ResetOutputWriter 恢复默认输出（直接 stdout）
func ResetOutputWriter() {
	WriteOutput = DefaultOutputWriter
}
