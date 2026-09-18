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

// PHP output handler 操作标志（传给回调第二参数 $phase），对齐 php-src main/php_output.h。
const (
	PHPOutputHandlerWrite = 0x00
	PHPOutputHandlerStart = 0x01
	PHPOutputHandlerClean = 0x02
	PHPOutputHandlerFlush = 0x04
	PHPOutputHandlerFinal = 0x08
	PHPOutputHandlerCont  = PHPOutputHandlerWrite
	PHPOutputHandlerEnd   = PHPOutputHandlerFinal
)

// PHP output handler 能力标志（ob_start 第三参数 $flags）。
const (
	PHPOutputHandlerCleanable = 0x0010
	PHPOutputHandlerFlushable = 0x0020
	PHPOutputHandlerRemovable = 0x0040
	PHPOutputHandlerStdFlags  = 0x0070
)

// PHP output handler 类型（ob_get_status()['type']）。
const (
	PHPOutputHandlerInternal = 0
	PHPOutputHandlerUser     = 1
)

// OutputBufferHandler 是一层输出缓冲的用户回调。
// phase 为 PHP_OUTPUT_HANDLER_* 位掩码；ok=false 表示该次 ob_* 失败。
// 回调返回 PHP false 时应由包装器改写为 rewritten=原始缓冲、ok=true，以便原样放行。
type OutputBufferHandler func(buffer string, phase int) (rewritten string, ok bool, ctl Control)

// OutputBufferStartSpec 描述 ob_start($callback, $chunk_size, $flags)。
type OutputBufferStartSpec struct {
	Handler   OutputBufferHandler
	ChunkSize int
	Flags     int
	Name      string
	Type      int
}

// OutputBufferStatusInfo 描述一个输出缓冲层的基本状态（对应 PHP ob_get_status 的元素）。
type OutputBufferStatusInfo struct {
	Level      int    // 层级（1 起，对齐 ob_get_level）
	Type       int    // PHP_OUTPUT_HANDLER_INTERNAL / USER
	Flags      int    // 处理标志
	ChunkSize  int    // 块大小（未指定为 0）
	BufferSize int    // 已分配缓冲容量（字节）
	BufferUsed int    // 当前已用字节
	Name       string // 处理器名称
}

// OutputBufferHost 为 VM / Context 提供请求级 PHP 输出缓冲能力。
type OutputBufferHost interface {
	StartOutputBuffer()
	StartOutputBufferSpec(spec OutputBufferStartSpec) bool
	CleanOutputBuffer() (string, bool)
	OutputBufferContents() (string, bool)
	OutputBufferLevel() int

	// FlushOutputBuffer 弹出栈顶并把（经 handler 改写后的）内容写出到上一层或最终输出。
	FlushOutputBuffer() (string, bool)
	// FlushCurrentBuffer 把栈顶内容写出到上一层但保留该层（ob_flush）。
	FlushCurrentBuffer() (string, bool)
	// CleanCurrentBuffer 清空栈顶内容但不结束缓冲（ob_clean）。
	CleanCurrentBuffer() bool
	// OutputBufferLength 返回栈顶缓冲的已用字节长度（无缓冲返回 false）。
	OutputBufferLength() (int, bool)
	// OutputBufferStatus 返回缓冲层状态列表；full=true 返回全部层，否则仅最顶层。
	OutputBufferStatus(full bool) []OutputBufferStatusInfo
	// ListOutputHandlers 返回所有激活缓冲的处理器名。
	ListOutputHandlers() []string
	SetImplicitFlush(on bool)
	IsImplicitFlush() bool
	// TakeOutputControl 取出 handler 回调产生的 throw/return 控制流（无则 nil）。
	TakeOutputControl() Control
	// FlushSAPI 对齐 PHP flush()：刷新底层 SAPI/HTTP，不弹出 ob 栈。
	FlushSAPI()
}

// EmitOutput 优先写到当前 Context（请求级缓冲指针），再回退 VM OutputSink。
func EmitOutput(ctx Context, s string) Control {
	if ctx != nil {
		if sink, ok := ctx.(OutputSink); ok {
			sink.WriteOutput(s)
			if host, ok := ctx.(OutputBufferHost); ok {
				return host.TakeOutputControl()
			}
			return nil
		}
		if vm := ctx.GetVM(); vm != nil {
			if sink, ok := vm.(OutputSink); ok {
				sink.WriteOutput(s)
				if host, ok := vm.(OutputBufferHost); ok {
					return host.TakeOutputControl()
				}
				return nil
			}
		}
	}
	WriteOutput(s)
	return nil
}

// CompileMode 编译模式标记。
// 设为 true 时，注解构造函数应跳过有副作用的操作（扫描目录、初始化数据库、调用 boot 等），
// 供 compile 子命令在纯解析阶段使用。
var CompileMode bool
