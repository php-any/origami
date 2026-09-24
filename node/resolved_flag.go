package node

import "sync/atomic"

// resolvedFlag 是「延迟解析已完成」的无锁快路径标记。
//
// 延迟解析节点（new / 函数调用 / 静态方法 / 静态属性）的 AST 在进程内跨请求共享，
// 而解析结果一旦确定就再也不变。此前每次执行都要抢一次 sync.Mutex，于是常驻服务里
// 所有并发请求都排在同一把互斥量上——吞吐随并发升高反而下降。
// 现在命中时只做一次 atomic.Bool 读，未命中才进锁，且每个调用点一生只进一次。
//
// 使用约定：只在解析成功后 MarkDone，失败路径不得置位，
// 以保留「解析失败下次重试」的原有语义。
type resolvedFlag struct {
	done atomic.Bool
}

// Done 报告解析结果是否已就绪；为 true 时可直接读缓存的字段（与 MarkDone 构成 acquire/release）。
func (r *resolvedFlag) Done() bool { return r.done.Load() }

// MarkDone 标记解析完成，必须在写好缓存字段之后调用。
func (r *resolvedFlag) MarkDone() { r.done.Store(true) }
