package data

import "fmt"

// userOutputEmitted 表示本次请求是否已向 stdout 写出用户可见内容（echo/var_dump 等）
var userOutputEmitted bool

// MarkUserOutput 标记已有用户输出（用于 Fatal 前空行等格式）
func MarkUserOutput() {
	userOutputEmitted = true
}

// HasUserOutput 是否已有用户输出
func HasUserOutput() bool {
	return userOutputEmitted
}

// ResetUserOutput 重置用户输出标记（每个脚本执行前调用）
func ResetUserOutput() {
	userOutputEmitted = false
}

// OutputWriter 是输出写入函数类型
type OutputWriter func(string)

// DefaultOutputWriter 默认输出到 stdout（供 ob_* 函数恢复时使用）
func DefaultOutputWriter(s string) {
	MarkUserOutput()
	fmt.Print(s)
}

// WriteOutput 是当前的输出写入函数
// 由 ob_start/ob_get_clean 切换输出目标
var WriteOutput OutputWriter = DefaultOutputWriter

// ResetOutputWriter 恢复默认输出（直接 stdout）
func ResetOutputWriter() {
	WriteOutput = DefaultOutputWriter
}
