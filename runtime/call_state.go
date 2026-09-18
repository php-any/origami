package runtime

import (
	"context"
	"sync"
	"time"

	"github.com/php-any/origami/data"
)

// CallState 是单 goroutine / 单请求的 PHP 调用深度与 debug_backtrace 栈。
// 不内置锁：PHP-FPM 语义是请求内单线程；跨请求靠各自一份 CallState 隔离。
type CallState struct {
	Depth           int
	Stack           []data.CallFrame
	auto            bool // 按 goid 惰性创建，深度归零后可从 map 删除
	deadline        context.Context
	phpDeadlineNano int64
	phpLimitSec     int64
}

func (st *CallState) Enter() int {
	if st == nil {
		return 0
	}
	st.Depth++
	if d := st.deadline; d != nil && (st.Depth == 1 || st.Depth&0x7f == 0) {
		if d.Err() != nil {
			panic(data.ErrRequestCanceled)
		}
	}
	return st.Depth
}

func (st *CallState) Leave() {
	if st == nil || st.Depth <= 0 {
		return
	}
	st.Depth--
}

func (st *CallState) Push(frame data.CallFrame) {
	if st == nil {
		return
	}
	st.Stack = append(st.Stack, frame)
}

func (st *CallState) Pop() {
	if st == nil {
		return
	}
	if n := len(st.Stack); n > 0 {
		st.Stack = st.Stack[:n-1]
	}
}

func (st *CallState) Snapshot() []data.CallFrame {
	if st == nil || len(st.Stack) == 0 {
		return nil
	}
	out := make([]data.CallFrame, len(st.Stack))
	copy(out, st.Stack)
	return out
}

var requestCall sync.Map // uint64 goid -> *CallState

func currentRequestCallState() *CallState {
	if requestOutputActive.Load() == 0 {
		return nil
	}
	id := goid()
	if v, ok := requestCall.Load(id); ok {
		return v.(*CallState)
	}
	return nil
}

// ensureGoroutineCallState 为当前 goroutine 取得调用栈。
// 已由 BeginRequestOutput 安装则复用；否则惰性创建（auto），深度与栈都空时删除以免泄漏。
func ensureGoroutineCallState() *CallState {
	id := goid()
	if v, ok := requestCall.Load(id); ok {
		return v.(*CallState)
	}
	st := &CallState{auto: true}
	actual, _ := requestCall.LoadOrStore(id, st)
	return actual.(*CallState)
}

func releaseAutoCallState(st *CallState) {
	if st == nil || !st.auto || st.Depth > 0 || len(st.Stack) > 0 {
		return
	}
	id := goid()
	if cur, ok := requestCall.Load(id); ok && cur == st {
		requestCall.Delete(id)
	}
}

// BeginRequestDeadline 把 HTTP/验收超时绑到当前请求的 CallState。
// 必须在 BeginRequestOutput 之后调用，这样 EnterCall 才能在 PHP 热循环里打断卡死的 Handle。
func BeginRequestDeadline(ctx context.Context) (restore func()) {
	st := currentRequestCallState()
	if st == nil {
		st = ensureGoroutineCallState()
	}
	prev := st.deadline
	st.deadline = ctx
	return func() {
		st.deadline = prev
	}
}

// SetRequestPHPDeadline 把 set_time_limit / max_execution_time 写到当前 goroutine 的 CallState。
// 有请求槽时返回 true，调用方不得再改进程级截止时间（否则并发请求会互相覆盖）。
func SetRequestPHPDeadline(seconds int) bool {
	st := currentRequestCallState()
	if st == nil {
		return false
	}
	if seconds <= 0 {
		st.phpDeadlineNano = 0
		st.phpLimitSec = 0
		return true
	}
	st.phpLimitSec = int64(seconds)
	st.phpDeadlineNano = time.Now().Add(time.Duration(seconds) * time.Second).UnixNano()
	return true
}

func requestPHPDeadlineExceeded(st *CallState) bool {
	if st == nil {
		return false
	}
	n := st.phpDeadlineNano
	return n != 0 && time.Now().UnixNano() >= n
}

// BindContextCallState 把函数帧绑到请求级 CallState，使 ctx.EnterCall 与 VM.EnterCall 共用同一份栈。
func BindContextCallState(ctx data.Context, st *CallState) {
	if c, ok := ctx.(*Context); ok {
		c.call = st
	}
}

func (vm *VM) localCall() *CallState {
	if st := currentRequestCallState(); st != nil {
		return st
	}
	return ensureGoroutineCallState()
}

// RequestDeadlineExceeded 供语句/循环边界检查当前 goroutine 的 HTTP/PHP 截止时间。
func RequestDeadlineExceeded() bool {
	st := currentRequestCallState()
	if st == nil {
		return false
	}
	if st.deadline != nil && st.deadline.Err() != nil {
		return true
	}
	return requestPHPDeadlineExceeded(st)
}

// InHTTPRequest 当前 goroutine 是否处于 HTTP 请求中（有请求输出槽即可）。
func InHTTPRequest() bool {
	return currentRequestCallState() != nil
}
