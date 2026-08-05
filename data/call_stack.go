package data

// CallFrame 表示 PHP 调用栈中的一帧（供 debug_backtrace）
type CallFrame struct {
	File     string
	Line     int
	Function string
	Class    string
	Type     string // "->" 或 "::"
}

// CallStackTracker 可选接口：VM 实现后 debug_backtrace 可返回真实 PHP 调用栈
type CallStackTracker interface {
	PushCallFrame(frame CallFrame)
	PopCallFrame()
	SnapshotCallStack() []CallFrame
}
