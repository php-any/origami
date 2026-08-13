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

// OutputBufferStatusInfo 描述一个输出缓冲层的基本状态（对应 PHP ob_get_status 的元素）。
type OutputBufferStatusInfo struct {
	Level      int    // 层级（1 起）
	Type       int    // 缓冲类型，PHP 中 PHPTAL/内部缓冲类型，这里统一为 1（PHP_OUTPUT_HANDLER_INTERNAL）
	Flags      int    // 处理标志
	ChunkSize  int    // 块大小（未指定为 0）
	BufferSize int    // 当前缓冲字节数
	Name       string // 处理器名称
}

// OutputBufferHost 为 VM 提供请求级 PHP 输出缓冲能力。
type OutputBufferHost interface {
	StartOutputBuffer()
	CleanOutputBuffer() (string, bool)
	OutputBufferContents() (string, bool)
	OutputBufferLevel() int

	// 以下为完整 ob_* 系列支持所需的能力。
	// FlushOutputBuffer 弹出并返回栈顶缓冲内容，并把内容写出到上一层（或最终输出），对应 ob_end_flush / ob_get_flush。
	FlushOutputBuffer() (string, bool)
	// CleanCurrentBuffer 清空栈顶缓冲内容但不结束缓冲，对应 ob_clean。
	CleanCurrentBuffer() bool
	// OutputBufferLength 返回栈顶缓冲的字节长度（无缓冲返回 false），对应 ob_get_length。
	OutputBufferLength() (int, bool)
	// OutputBufferStatus 返回缓冲层状态列表；full=true 返回全部层，否则仅最顶层，对应 ob_get_status。
	OutputBufferStatus(full bool) []OutputBufferStatusInfo
	// ListOutputHandlers 返回所有激活缓冲的处理器名，对应 ob_list_handlers。
	ListOutputHandlers() []string
	// SetImplicitFlush 设置/清除隐式刷新标志，对应 ob_implicit_flush。
	SetImplicitFlush(on bool)
	// IsImplicitFlush 返回当前隐式刷新标志。
	IsImplicitFlush() bool
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
