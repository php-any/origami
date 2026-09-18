package runtime

import (
	"strings"
	"sync"
	"sync/atomic"

	"github.com/php-any/origami/data"
)

const (
	outputBufMinGrow = 4096
	outputBufMaxPool = 1 << 20
)

var outputBufPool = sync.Pool{
	New: func() any {
		b := &strings.Builder{}
		b.Grow(outputBufMinGrow)
		return b
	},
}

func getOutputBuilder(grow int) *strings.Builder {
	b := outputBufPool.Get().(*strings.Builder)
	b.Reset()
	if grow > outputBufMinGrow {
		b.Grow(grow)
	}
	return b
}

func putOutputBuilder(b *strings.Builder) {
	if b == nil {
		return
	}
	if b.Cap() > outputBufMaxPool {
		return
	}
	b.Reset()
	outputBufPool.Put(b)
}

// outputLayer 是一层 PHP 输出缓冲（对应一次 ob_start）。
type outputLayer struct {
	buf       *strings.Builder
	chunkSize int
	flags     int
	name      string
	typeID    int
	handler   data.OutputBufferHandler
	started   bool
	invoking  bool
}

// OutputState 是输出缓冲栈的独立实现：自带锁，不与 VM 注册表/调用栈争用。
// 单请求内 PHP 语义仍是单线程；spawn/协程并发写同一请求栈时由 mu 串行化。
// 跨请求靠各自一份 OutputState（Context 指针或 BeginRequestOutput）隔离。
type OutputState struct {
	mu            sync.Mutex
	layers        []*outputLayer
	has           atomic.Bool
	depth         atomic.Int32
	implicitFlush atomic.Bool
	sink          func(string)
	sapiFlush     func()
	pending       atomic.Value // *ctlBox
	local         bool         // BeginRequestOutput 安装：echo 走本栈，不再查 goid
}

type outputState = OutputState

type ctlBox struct{ c data.Control }

func newOutputState() *OutputState {
	return &OutputState{}
}

// NewOutputState 创建独立输出缓冲栈，供 HTTP 请求 VM 绑定。
func NewOutputState() *OutputState {
	return newOutputState()
}

// SetSink 设置无缓冲时的最终写出目标（HTTP body / stdout）。
func (st *OutputState) SetSink(fn func(string)) {
	if st == nil {
		return
	}
	st.sink = fn
}

// SetSAPIFlush 设置 PHP flush() / ob_implicit_flush 触发的 SAPI 刷新。
func (st *OutputState) SetSAPIFlush(fn func()) {
	if st == nil {
		return
	}
	st.sapiFlush = fn
}

func (st *OutputState) syncFlag() {
	n := len(st.layers)
	st.depth.Store(int32(n))
	st.has.Store(n > 0)
}

func (st *OutputState) emitFinal(s string, fallback func(string)) {
	sink := st.sink
	if sink == nil {
		sink = fallback
	}
	if sink == nil {
		sink = data.WriteOutput
	}
	if s != "" {
		sink(s)
	}
	if st.implicitFlush.Load() {
		if f := st.sapiFlush; f != nil {
			f()
		}
	}
}

func (st *OutputState) write(s string, fallback func(string)) {
	if s == "" {
		return
	}
	if !st.has.Load() {
		st.emitFinal(s, fallback)
		return
	}
	st.mu.Lock()
	n := len(st.layers)
	if n == 0 {
		st.mu.Unlock()
		st.emitFinal(s, fallback)
		return
	}
	idx := n - 1
	for idx >= 0 && st.layers[idx].invoking {
		idx--
	}
	if idx < 0 {
		st.mu.Unlock()
		st.emitFinal(s, fallback)
		return
	}
	layer := st.layers[idx]
	layer.buf.WriteString(s)
	chunk := layer.chunkSize
	overflow := chunk > 0 && layer.buf.Len() >= chunk
	st.mu.Unlock()
	if overflow && idx == n-1 {
		st.flushCurrent(data.PHPOutputHandlerWrite, fallback)
	}
}

func (st *OutputState) start() {
	st.startSpec(data.OutputBufferStartSpec{
		Flags: data.PHPOutputHandlerStdFlags,
		Name:  "default output handler",
		Type:  data.PHPOutputHandlerInternal,
	})
}

func (st *OutputState) startSpec(spec data.OutputBufferStartSpec) bool {
	flags := spec.Flags
	name := spec.Name
	if name == "" {
		if spec.Handler != nil {
			name = "user output handler"
		} else {
			name = "default output handler"
		}
	}
	chunk := spec.ChunkSize
	if chunk < 0 {
		chunk = 0
	}
	grow := outputBufMinGrow
	if chunk > grow {
		grow = chunk
	}
	layer := &outputLayer{
		buf:       getOutputBuilder(grow),
		chunkSize: chunk,
		flags:     flags,
		name:      name,
		typeID:    spec.Type,
		handler:   spec.Handler,
	}
	st.mu.Lock()
	st.layers = append(st.layers, layer)
	st.syncFlag()
	st.mu.Unlock()
	return true
}

func (st *OutputState) runHandler(layer *outputLayer, handler data.OutputBufferHandler, content string, phase int) (string, bool) {
	if handler == nil {
		if layer != nil {
			st.mu.Lock()
			layer.invoking = false
			st.mu.Unlock()
		}
		return content, true
	}
	rewritten, ok, ctl := handler(content, phase)
	if layer != nil {
		st.mu.Lock()
		layer.invoking = false
		st.mu.Unlock()
	}
	if ctl != nil {
		st.setControl(ctl)
		return content, false
	}
	if !ok {
		return content, false
	}
	return rewritten, true
}

func (st *OutputState) setControl(c data.Control) {
	if c == nil {
		return
	}
	st.pending.Store(&ctlBox{c: c})
}

func (st *OutputState) takeControl() data.Control {
	v := st.pending.Swap((*ctlBox)(nil))
	if v == nil {
		return nil
	}
	box, ok := v.(*ctlBox)
	if !ok || box == nil {
		return nil
	}
	return box.c
}

func (st *OutputState) popLayer() *outputLayer {
	n := len(st.layers)
	if n == 0 {
		return nil
	}
	layer := st.layers[n-1]
	st.layers[n-1] = nil
	st.layers = st.layers[:n-1]
	st.syncFlag()
	return layer
}

func (st *OutputState) releaseLayer(layer *outputLayer) {
	if layer == nil {
		return
	}
	putOutputBuilder(layer.buf)
	layer.buf = nil
	layer.handler = nil
}

// cleanPop 对应 ob_get_clean / ob_end_clean：弹出栈顶，调用 CLEAN|FINAL handler，不写出。
func (st *OutputState) cleanPop() (string, bool) {
	st.mu.Lock()
	layer := st.popLayer()
	if layer == nil {
		st.mu.Unlock()
		return "", false
	}
	if layer.flags&data.PHPOutputHandlerRemovable == 0 {
		st.layers = append(st.layers, layer)
		st.syncFlag()
		st.mu.Unlock()
		return "", false
	}
	content := layer.buf.String()
	phase := data.PHPOutputHandlerClean | data.PHPOutputHandlerFinal
	if !layer.started {
		phase |= data.PHPOutputHandlerStart
	}
	handler := layer.handler
	st.mu.Unlock()
	out, ok := st.runHandler(nil, handler, content, phase)
	st.releaseLayer(layer)
	if !ok {
		return "", false
	}
	return out, true
}

func (st *OutputState) clean() (string, bool) {
	return st.cleanPop()
}

// flushEnd 对应 ob_end_flush / ob_get_flush：弹出、FINAL handler、写出到上一层。
func (st *OutputState) flushEnd(fallback func(string)) (string, bool) {
	st.mu.Lock()
	layer := st.popLayer()
	if layer == nil {
		st.mu.Unlock()
		return "", false
	}
	if layer.flags&data.PHPOutputHandlerRemovable == 0 {
		st.layers = append(st.layers, layer)
		st.syncFlag()
		st.mu.Unlock()
		return "", false
	}
	content := layer.buf.String()
	phase := data.PHPOutputHandlerFinal
	if !layer.started {
		phase |= data.PHPOutputHandlerStart
	}
	handler := layer.handler
	st.mu.Unlock()
	out, ok := st.runHandler(nil, handler, content, phase)
	st.releaseLayer(layer)
	if !ok {
		return "", false
	}
	st.write(out, fallback)
	return out, true
}

func (st *OutputState) flush() (string, bool) {
	return st.flushEnd(nil)
}

func (st *OutputState) flushWithFallback(fallback func(string)) (string, bool) {
	return st.flushEnd(fallback)
}

// flushCurrent 对应 ob_flush：保留层级，内容经 handler 后冒泡到上一层。
func (st *OutputState) flushCurrent(phase int, fallback func(string)) (string, bool) {
	st.mu.Lock()
	n := len(st.layers)
	if n == 0 {
		st.mu.Unlock()
		return "", false
	}
	layer := st.layers[n-1]
	if phase&data.PHPOutputHandlerFlush != 0 && layer.flags&data.PHPOutputHandlerFlushable == 0 {
		st.mu.Unlock()
		return "", false
	}
	content := layer.buf.String()
	layer.buf.Reset()
	if !layer.started {
		phase |= data.PHPOutputHandlerStart
		layer.started = true
	}
	handler := layer.handler
	// 保持 invoking，使 write() 跳过本层、把内容冒泡到上一层/SAPI。
	layer.invoking = true
	st.mu.Unlock()
	out, ok := st.runHandler(nil, handler, content, phase)
	if ok {
		st.write(out, fallback)
	}
	st.mu.Lock()
	layer.invoking = false
	st.mu.Unlock()
	return out, ok
}

func (st *OutputState) cleanCurrent() bool {
	st.mu.Lock()
	n := len(st.layers)
	if n == 0 {
		st.mu.Unlock()
		return false
	}
	layer := st.layers[n-1]
	if layer.flags&data.PHPOutputHandlerCleanable == 0 {
		st.mu.Unlock()
		return false
	}
	content := layer.buf.String()
	layer.buf.Reset()
	phase := data.PHPOutputHandlerClean
	if !layer.started {
		phase |= data.PHPOutputHandlerStart
		layer.started = true
	}
	handler := layer.handler
	layer.invoking = handler != nil
	st.mu.Unlock()
	_, ok := st.runHandler(layer, handler, content, phase)
	return ok
}

func (st *OutputState) length() (int, bool) {
	st.mu.Lock()
	defer st.mu.Unlock()
	n := len(st.layers)
	if n == 0 {
		return 0, false
	}
	return st.layers[n-1].buf.Len(), true
}

func (st *OutputState) contents() (string, bool) {
	st.mu.Lock()
	defer st.mu.Unlock()
	n := len(st.layers)
	if n == 0 {
		return "", false
	}
	return st.layers[n-1].buf.String(), true
}

func (st *OutputState) level() int {
	return int(st.depth.Load())
}

func (st *OutputState) status(full bool) []data.OutputBufferStatusInfo {
	st.mu.Lock()
	defer st.mu.Unlock()
	n := len(st.layers)
	if n == 0 {
		return nil
	}
	start := 0
	if !full {
		start = n - 1
	}
	result := make([]data.OutputBufferStatusInfo, 0, n-start)
	for i := start; i < n; i++ {
		layer := st.layers[i]
		result = append(result, data.OutputBufferStatusInfo{
			Level:      i + 1,
			Type:       layer.typeID,
			Flags:      layer.flags,
			ChunkSize:  layer.chunkSize,
			BufferSize: layer.buf.Cap(),
			BufferUsed: layer.buf.Len(),
			Name:       layer.name,
		})
	}
	return result
}

func (st *OutputState) handlers() []string {
	st.mu.Lock()
	defer st.mu.Unlock()
	n := len(st.layers)
	if n == 0 {
		return nil
	}
	result := make([]string, n)
	for i, layer := range st.layers {
		result[i] = layer.name
	}
	return result
}

func (st *OutputState) setImplicitFlush(on bool) {
	st.implicitFlush.Store(on)
}

func (st *OutputState) isImplicitFlush() bool {
	return st.implicitFlush.Load()
}

func (st *OutputState) flushSAPI() {
	if st == nil {
		return
	}
	if f := st.sapiFlush; f != nil {
		f()
	}
}

// 下列导出方法供 std/php/fpm.RequestVM 等跨包宿主委托，语义与内部方法一致。

func (st *OutputState) WriteTo(s string, fallback func(string)) { st.write(s, fallback) }
func (st *OutputState) StartDefault()                           { st.start() }
func (st *OutputState) StartSpec(spec data.OutputBufferStartSpec) bool {
	return st.startSpec(spec)
}
func (st *OutputState) Clean() (string, bool)    { return st.clean() }
func (st *OutputState) Contents() (string, bool) { return st.contents() }
func (st *OutputState) Level() int               { return st.level() }
func (st *OutputState) FlushEnd(fallback func(string)) (string, bool) {
	return st.flushEnd(fallback)
}
func (st *OutputState) FlushCurrent(fallback func(string)) (string, bool) {
	return st.flushCurrent(data.PHPOutputHandlerFlush, fallback)
}
func (st *OutputState) CleanCurrent() bool { return st.cleanCurrent() }
func (st *OutputState) Length() (int, bool) {
	return st.length()
}
func (st *OutputState) Status(full bool) []data.OutputBufferStatusInfo {
	return st.status(full)
}
func (st *OutputState) Handlers() []string        { return st.handlers() }
func (st *OutputState) SetImplicitFlush(on bool)  { st.setImplicitFlush(on) }
func (st *OutputState) IsImplicitFlush() bool     { return st.isImplicitFlush() }
func (st *OutputState) TakeControl() data.Control { return st.takeControl() }
func (st *OutputState) FlushSAPI()                { st.flushSAPI() }

// --- 请求级输出缓冲 ---
// 按 goroutine 隔离（map 以廉价 goid 为键）。热路径优先用 Context 上继承的指针；
// 启动期帧在请求 goroutine 上 echo 时再查一次当前 goroutine 的槽，禁止进程级单槽。

func init() {
	data.RequestGoid = goid
}

var (
	requestOutput       sync.Map // uint64 goid -> *OutputState
	requestOutputActive atomic.Int32
)

// BeginRequestOutput 为当前 goroutine 安装独立输出缓冲栈，返回 restore（务必 defer）。
// Laravel serve 等多请求共享 runtime.VM 时必须调用，否则 ob_/echo 会串请求。
func BeginRequestOutput() (restore func()) {
	id := goid()
	st := newOutputState()
	st.local = true
	prevOut, _ := requestOutput.Swap(id, st)
	callSt := &CallState{}
	prevCall, _ := requestCall.Swap(id, callSt)
	requestOutputActive.Add(1)
	data.BeginRequestStaticOverlay()
	return func() {
		data.EndRequestStaticOverlay()
		if prevOut != nil {
			requestOutput.Store(id, prevOut)
		} else {
			requestOutput.Delete(id)
		}
		if prevCall != nil {
			requestCall.Store(id, prevCall)
		} else {
			requestCall.Delete(id)
		}
		requestOutputActive.Add(-1)
	}
}

// MuteRequestStdout 让当前请求无 ob 的 echo 不再落到进程 stdout。
// Laravel Kernel 把页面放在 Response 里由 SendResponseTo 写给浏览器；
// 若 emitFinal 仍 fallback 到 data.WriteOutput，同一份 HTML 会同时出现在浏览器和控制台。
func MuteRequestStdout() {
	if st := currentRequestOutput(); st != nil {
		st.SetSink(func(string) {})
	}
}

// StartRequestOutputBuffer 对齐 Octane Worker：请求入口 ob_start()，
// Blade/echo 漏出 View 内层缓冲时仍留在本请求栈上，而不是 stdout。
func StartRequestOutputBuffer() {
	if st := currentRequestOutput(); st != nil {
		st.start()
	}
}

// TakeRequestOutput 对齐 Octane ob_get_contents + ob_end_clean。
// 从外层到内层拼残留输出，供 SwooleClient 那样先写 leftover 再写 Response。
func TakeRequestOutput() string {
	st := currentRequestOutput()
	if st == nil {
		return ""
	}
	var parts []string
	for st.level() > 0 {
		s, ok := st.clean()
		if !ok {
			break
		}
		parts = append(parts, s)
	}
	if len(parts) == 0 {
		return ""
	}
	var b strings.Builder
	for i := len(parts) - 1; i >= 0; i-- {
		b.WriteString(parts[i])
	}
	return b.String()
}

func currentRequestOutput() *OutputState {
	if requestOutputActive.Load() == 0 {
		return nil
	}
	id := goid()
	if v, ok := requestOutput.Load(id); ok {
		return v.(*OutputState)
	}
	return nil
}

// BindContextOutput 把请求级缓冲栈绑到 Context，后续 echo 不再查 goid。
func BindContextOutput(ctx data.Context, st *OutputState) {
	if c, ok := ctx.(*Context); ok && st != nil {
		c.out = st
	}
}
