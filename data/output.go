package data

import (
	"fmt"
	"sync/atomic"
)

// userOutputEmitted 表示当前进程是否已向 stdout 写出用户可见内容。
// 它只影响 CLI Fatal 输出前的换行格式，不承载请求状态。
var userOutputEmitted atomic.Bool

// MarkUserOutput 标记已有用户输出（用于 Fatal 前空行等格式）
func MarkUserOutput() {
	userOutputEmitted.Store(true)
}

// HasUserOutput 是否已有用户输出
func HasUserOutput() bool {
	return userOutputEmitted.Load()
}

// ResetUserOutput 重置 CLI 脚本输出标记。
func ResetUserOutput() {
	userOutputEmitted.Store(false)
}

// OutputWriter 是输出写入函数类型
type OutputWriter func(string)

// DefaultOutputWriter 默认输出到 stdout
func DefaultOutputWriter(s string) {
	MarkUserOutput()
	fmt.Print(s)
}

// WriteOutput 是不可替换的语言默认输出快路径。
func WriteOutput(s string) { DefaultOutputWriter(s) }

// OutputSink 可选接口：VM 实现后 echo/HTML/Response 走请求级输出，默认路径仍用 WriteOutput。
type OutputSink interface {
	WriteOutput(s string)
}

// OutputBufferHost 为 VM 提供请求级 PHP 输出缓冲能力。
type OutputBufferHost interface {
	StartOutputBuffer()
	CleanOutputBuffer() (string, bool)
	OutputBufferContents() (string, bool)
	OutputBufferLevel() int
}

// EmitOutput 优先写到当前 VM 的 OutputSink，否则走语言默认 WriteOutput。
func EmitOutput(ctx Context, s string) {
	if ctx != nil {
		if sink, ok := ctx.GetVM().(OutputSink); ok {
			sink.WriteOutput(s)
			return
		}
	}
	WriteOutput(s)
}

// CompileMode 编译模式标记。
// 设为 true 时，注解构造函数应跳过有副作用的操作（扫描目录、初始化数据库、调用 boot 等），
// 供 compile 子命令在纯解析阶段使用。
var CompileMode bool
