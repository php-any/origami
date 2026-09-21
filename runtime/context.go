package runtime

import (
	"context"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
)

var contextPool = sync.Pool{
	New: func() any {
		return &Context{}
	},
}

// zval 槽块按档位回收。CreateContext 每次 make([]ZVal, n) 会在 Laravel 请求里放大成百万次分配。
var zvalBlkPools [8]sync.Pool

func zvalBlkClass(n int) int {
	switch {
	case n <= 8:
		return 0
	case n <= 16:
		return 1
	case n <= 32:
		return 2
	case n <= 64:
		return 3
	case n <= 128:
		return 4
	case n <= 256:
		return 5
	case n <= 512:
		return 6
	case n <= 1024:
		return 7
	default:
		return -1
	}
}

func zvalBlkCap(class int) int {
	return 8 << class
}

func takeZValBlock(n int) []data.ZVal {
	cl := zvalBlkClass(n)
	if cl < 0 {
		return make([]data.ZVal, n)
	}
	capn := zvalBlkCap(cl)
	p := zvalBlkPools[cl].Get()
	if p == nil {
		return make([]data.ZVal, n, capn)
	}
	blk := p.([]data.ZVal)
	return blk[:n]
}

func putZValBlock(blk []data.ZVal) {
	if blk == nil {
		return
	}
	cl := zvalBlkClass(cap(blk))
	if cl < 0 || zvalBlkCap(cl) != cap(blk) {
		return
	}
	zvalBlkPools[cl].Put(blk[:cap(blk)])
}

// Context 表示运行时上下文
type Context struct {
	vm data.VM
	// 命名空间
	namespace string

	// 变量存储符号表
	variables []*data.ZVal

	// 记录本次函数/方法调用时的实参表达式列表（用于 func_get_args 等）
	callArgs []data.GetValue

	// 记录本次函数/方法调用展开后的扁平实参值列表（处理 ...$arr 展开）
	flatCallArgs []data.Value

	// 当前调用绑定的 static 局部变量存储（由 StaticVarStatement 惰性绑定）
	staticLocals *data.StaticLocals

	// retSlot 是本帧复用的 return 载体。ReturnControl 只在本帧内向上冒泡，
	// 被 Call 拆值后即失效，因此每帧一个实例就够。
	retSlot data.ReturnValue

	call    *CallState
	out     *OutputState
	pooled  bool
	escaped bool
	zvalBlk []data.ZVal
}

// BindStaticLocals 绑定函数级 static 局部变量存储
func (c *Context) BindStaticLocals(store *data.StaticLocals) {
	c.staticLocals = store
}

// StaticLocalsStore 返回绑定的 static 存储
func (c *Context) StaticLocalsStore() *data.StaticLocals {
	return c.staticLocals
}

// ReturnSlot 复用本帧的 ReturnValue，不每次 return 都分配。
func (c *Context) ReturnSlot(v data.Value) data.ReturnControl {
	c.retSlot.V = v
	return &c.retSlot
}

// NewContext 创建一个新的运行时上下文
func NewContext(vm data.VM) data.Context {
	return &Context{
		vm: vm,
	}
}

// SetNamespace 设置命名空间
func (c *Context) SetNamespace(name string) data.Context {
	c.namespace = name
	return c
}

// GetNamespace 获取命名空间
func (c *Context) GetNamespace() string {
	return c.namespace
}

// GetVariableValue 获取变量值
func (c *Context) GetVariableValue(variable data.Variable) (data.Value, data.Control) {
	if pl, ok := variable.(data.PropertyLvalue); ok {
		gv, ctl := pl.GetValue(c)
		if ctl != nil {
			return nil, ctl
		}
		if val, ok := gv.(data.Value); ok {
			return val, nil
		}
		return data.NewNullValue(), nil
	}
	return c.variables[variable.GetIndex()].Value, nil
}

func (c *Context) GetIndexValue(index int) (data.Value, bool) {
	if index < 0 || index >= len(c.variables) {
		return nil, false
	}
	return c.variables[index].Value, true
}

func (c *Context) SetIndexZVal(index int, v *data.ZVal) {
	c.variables[index] = v
}

func (c *Context) GetIndexZVal(index int) *data.ZVal {
	return c.variables[index]
}

func (c *Context) definedZValByName(name string) *data.ZVal {
	if c == nil || name == "" {
		return nil
	}
	for _, zv := range c.variables {
		if zv != nil && zv.Name == name && zv.Defined {
			return zv
		}
	}
	return nil
}

// shareZVal 让当前帧与调用者共用同一 ZVal（PHP include 语义），不按值拷贝数组。
func (c *Context) shareZVal(zv *data.ZVal) {
	if c == nil || zv == nil || zv.Name == "" {
		return
	}
	for i, cur := range c.variables {
		if cur != nil && cur.Name == zv.Name {
			c.variables[i] = zv
			return
		}
	}
	c.variables = append(c.variables, zv)
}

// SetVariableValue 设置变量值
func (c *Context) SetVariableValue(variable data.Variable, value data.Value) data.Control {
	if pl, ok := variable.(data.PropertyLvalue); ok {
		return pl.SetValue(c, value)
	}
	switch v := value.(type) {
	case *data.ReferenceValue:
		if pl, ok := v.Val.(data.PropertyLvalue); ok {
			if gzv, ok := pl.(interface {
				GetZVal(data.Context) (*data.ZVal, data.Control)
			}); ok {
				zv, ctl := gzv.GetZVal(v.Ctx)
				if ctl != nil {
					return ctl
				}
				if zv == nil {
					zv = data.NewZVal(data.NewNullValue())
				}
				c.variables[variable.GetIndex()] = zv
				return nil
			}
		}
		c.variables[variable.GetIndex()] = v.Ctx.GetIndexZVal(v.Val.GetIndex())
	case *data.ArraySlotRef:
		// &$array[] 语法：局部变量与数组元素共享 ZVal
		if v.Arr != nil && v.Idx >= 0 && v.Idx < len(v.Arr.List) {
			slot := v.Arr.List[v.Idx]
			slot.AddRefSlot()
			c.variables[variable.GetIndex()] = slot
		}
	case *data.IndexReferenceValue:
		if ie, ok := v.Expr.(*node.IndexExpression); ok {
			zv, ctl := ie.GetOrCreateZVal(v.Ctx)
			if ctl != nil {
				return ctl
			}
			if zv != nil {
				c.variables[variable.GetIndex()] = zv
			}
		}
	case *data.ArrayValue:
		zv := c.variables[variable.GetIndex()]
		zv.Value = data.CloneArrayValue(v)
		zv.Defined = true
	case *data.ObjectValue:
		// PHP 中 array 是按值赋值 + copy-on-write。
		// 在 Origami 里，关联数组可能由 ObjectValue 表示，这里也做一次结构级克隆，
		// 避免 `$b = $this->a; $b['k']=...` 反向修改到 `$this->a`（Symfony InputDefinition::$arguments 等场景）。
		zv := c.variables[variable.GetIndex()]
		zv.Value = data.CloneObjectValue(v)
		zv.Defined = true
	default:
		idx := variable.GetIndex()
		c.variables[idx].Value = value
		c.variables[idx].Defined = true
	}

	return nil
}

// CreateContext 创建函数上下文
func (c *Context) CreateContext(vars []data.Variable) data.Context {
	nc := contextPool.Get().(*Context)
	nc.vm = c.vm
	nc.namespace = ""
	nc.callArgs = nil
	nc.flatCallArgs = nil
	nc.staticLocals = nil
	nc.retSlot.V = nil
	nc.pooled = true
	nc.escaped = false
	nc.call = c.call
	nc.out = c.inheritOut()
	nc.resetVariables(vars)
	return nc
}

func (c *Context) resetVariables(vars []data.Variable) {
	c.releaseZValBlock()
	n := len(vars)
	if n == 0 {
		c.variables = c.variables[:0]
		return
	}
	// 槽位指向本帧专属 block。逃逸帧不还 block，避免闭包/引用仍握着旧 ZVal。
	if cap(c.variables) < n {
		c.variables = make([]*data.ZVal, n)
	} else {
		c.variables = c.variables[:n]
	}
	block := takeZValBlock(n)
	c.zvalBlk = block
	nullV := data.NewNullValue()
	for i := 0; i < n; i++ {
		z := &block[i]
		z.Name = vars[i].GetName()
		z.Value = nullV
		z.Defined = false
		z.RefSlotCount = 0
		z.EmptyStrKey = false
		c.variables[i] = z
	}
}

func (c *Context) releaseZValBlock() {
	if c == nil || c.zvalBlk == nil {
		return
	}
	putZValBlock(c.zvalBlk)
	c.zvalBlk = nil
}

func (c *Context) MarkEscaped() {
	if c != nil {
		c.escaped = true
	}
}

func (c *Context) IsEscaped() bool {
	return c != nil && c.escaped
}

func (c *Context) ReleasePooled() {
	if c == nil || !c.pooled || c.escaped {
		return
	}
	c.pooled = false
	c.vm = nil
	c.call = nil
	c.out = nil
	c.callArgs = nil
	c.flatCallArgs = nil
	c.staticLocals = nil
	c.retSlot.V = nil
	c.namespace = ""
	for i := range c.variables {
		c.variables[i] = nil
	}
	c.variables = c.variables[:0]
	c.releaseZValBlock()
	contextPool.Put(c)
}

func (c *Context) resolveCallState() *CallState {
	if c == nil {
		return nil
	}
	if c.call != nil {
		return c.call
	}
	switch vm := c.vm.(type) {
	case *VM:
		st := vm.localCall()
		c.call = st
		return st
	case *TempVM:
		c.call = &vm.call
		return c.call
	default:
		return currentRequestCallState()
	}
}

func (c *Context) CallDepth() int {
	if st := c.resolveCallState(); st != nil {
		return st.Depth
	}
	return 0
}

func (c *Context) EnterCall() int {
	return c.resolveCallState().Enter()
}

func (c *Context) LeaveCall() {
	st := c.resolveCallState()
	st.Leave()
	releaseAutoCallState(st)
}

func (c *Context) PushCallFrame(frame data.CallFrame) {
	c.resolveCallState().Push(frame)
}

func (c *Context) PopCallFrame() {
	c.resolveCallState().Pop()
}

func (c *Context) SnapshotCallStack() []data.CallFrame {
	return c.resolveCallState().Snapshot()
}

func (c *Context) CreateBaseContext() data.Context {
	return &Context{
		vm:   c.vm,
		call: c.call,
		out:  c.inheritOut(),
	}
}

func (c *Context) inheritOut() *OutputState {
	if c != nil && c.out != nil && c.out.local {
		return c.out
	}
	if st := currentRequestOutput(); st != nil {
		return st
	}
	if c != nil {
		return c.out
	}
	return nil
}

func (c *Context) resolveOut() *OutputState {
	// 请求帧已绑 local 缓冲则不再查 goid。启动期帧钉的是 vm.out，
	// 在请求 goroutine 上 echo 时改走 TLS，但不写回共享 boot Context。
	if c != nil && c.out != nil && c.out.local {
		return c.out
	}
	if st := currentRequestOutput(); st != nil {
		return st
	}
	if c != nil && c.out != nil {
		return c.out
	}
	if c == nil || c.vm == nil {
		return nil
	}
	switch vm := c.vm.(type) {
	case *VM:
		return vm.out
	case *TempVM:
		return vm.out
	}
	return nil
}

// WriteOutput 实现 data.OutputSink：echo 走 Context 上缓存的缓冲栈，避免每次解析 goid。
func (c *Context) WriteOutput(s string) {
	st := c.resolveOut()
	if st == nil {
		if c != nil && c.vm != nil {
			if sink, ok := c.vm.(data.OutputSink); ok {
				sink.WriteOutput(s)
				return
			}
		}
		data.WriteOutput(s)
		return
	}
	st.write(s, data.WriteOutput)
}

func (c *Context) StartOutputBuffer() {
	if st := c.resolveOut(); st != nil {
		st.start()
	}
}

func (c *Context) StartOutputBufferSpec(spec data.OutputBufferStartSpec) bool {
	if st := c.resolveOut(); st != nil {
		return st.startSpec(spec)
	}
	if c != nil && c.vm != nil {
		if host, ok := c.vm.(data.OutputBufferHost); ok {
			return host.StartOutputBufferSpec(spec)
		}
	}
	return false
}

func (c *Context) CleanOutputBuffer() (string, bool) {
	if st := c.resolveOut(); st != nil {
		return st.clean()
	}
	return "", false
}

func (c *Context) FlushOutputBuffer() (string, bool) {
	st := c.resolveOut()
	if st == nil {
		return "", false
	}
	return st.flushWithFallback(data.WriteOutput)
}

func (c *Context) FlushCurrentBuffer() (string, bool) {
	st := c.resolveOut()
	if st == nil {
		return "", false
	}
	return st.flushCurrent(data.PHPOutputHandlerFlush, data.WriteOutput)
}

func (c *Context) CleanCurrentBuffer() bool {
	if st := c.resolveOut(); st != nil {
		return st.cleanCurrent()
	}
	return false
}

func (c *Context) TakeOutputControl() data.Control {
	if st := c.resolveOut(); st != nil {
		return st.takeControl()
	}
	return nil
}

func (c *Context) FlushSAPI() {
	if st := c.resolveOut(); st != nil {
		st.flushSAPI()
	}
}

func (c *Context) OutputBufferLength() (int, bool) {
	if st := c.resolveOut(); st != nil {
		return st.length()
	}
	return 0, false
}

func (c *Context) OutputBufferStatus(full bool) []data.OutputBufferStatusInfo {
	if st := c.resolveOut(); st != nil {
		return st.status(full)
	}
	return nil
}

func (c *Context) ListOutputHandlers() []string {
	if st := c.resolveOut(); st != nil {
		return st.handlers()
	}
	return nil
}

func (c *Context) SetImplicitFlush(on bool) {
	if st := c.resolveOut(); st != nil {
		st.setImplicitFlush(on)
	}
}

func (c *Context) IsImplicitFlush() bool {
	if st := c.resolveOut(); st != nil {
		return st.isImplicitFlush()
	}
	return false
}

func (c *Context) OutputBufferContents() (string, bool) {
	if st := c.resolveOut(); st != nil {
		return st.contents()
	}
	return "", false
}

func (c *Context) OutputBufferLevel() int {
	if st := c.resolveOut(); st != nil {
		return st.level()
	}
	return 0
}

func (c *Context) GetVM() data.VM {
	return c.vm
}

func (c *Context) GoContext() context.Context {
	return context.Background()
}

// SetVM 替换当前 Context 所绑定的 VM
func (c *Context) SetVM(vm data.VM) {
	if vm == nil {
		return
	}
	c.vm = vm
}

// SetCallArgs 记录本次调用时传入的参数表达式列表
func (c *Context) SetCallArgs(args []data.GetValue) {
	c.callArgs = args
}

// GetCallArgs 获取本次调用时传入的参数表达式列表
func (c *Context) GetCallArgs() []data.GetValue {
	return c.callArgs
}

// SetFlatCallArgs 记录本次调用展开后的扁平实参值列表（处理 ...$arr 展开）
func (c *Context) SetFlatCallArgs(values []data.Value) {
	c.flatCallArgs = values
}

// GetFlatCallArgs 获取本次调用展开后的扁平实参值列表
func (c *Context) GetFlatCallArgs() []data.Value {
	return c.flatCallArgs
}

func makeSliceVariable(i int) []*data.ZVal {
	if i <= 0 {
		return nil
	}
	l := make([]*data.ZVal, i)
	block := make([]data.ZVal, i)
	nullV := data.NewNullValue()
	for j := range l {
		block[j].Value = nullV
		l[j] = &block[j]
	}
	return l
}

// makeSliceVariableWithNames 创建带变量名的 ZVal 切片
func makeSliceVariableWithNames(vars []data.Variable) []*data.ZVal {
	n := len(vars)
	if n == 0 {
		return nil
	}
	l := make([]*data.ZVal, n)
	block := make([]data.ZVal, n)
	nullV := data.NewNullValue()
	for i := range l {
		if vars[i] != nil {
			block[i].Name = vars[i].GetName()
		}
		block[i].Value = nullV
		l[i] = &block[i]
	}
	return l
}

// SetVariableByName 通过变量名设置变量值，用于 extract 等动态赋值场景。
// 若该名称尚无槽位，则追加新槽（Laravel getRequire: extract 注入的键在闭包体内未必被静态引用）。
func (c *Context) SetVariableByName(name string, value data.Value) {
	for _, zv := range c.variables {
		if zv != nil && zv.Name == name {
			switch v := value.(type) {
			case *data.ArrayValue:
				zv.Value = data.CloneArrayValue(v)
			case *data.ObjectValue:
				zv.Value = data.CloneObjectValue(v)
			default:
				zv.Value = value
			}
			zv.Defined = true
			return
		}
	}
	var stored data.Value = value
	switch v := value.(type) {
	case *data.ArrayValue:
		stored = data.CloneArrayValue(v)
	case *data.ObjectValue:
		stored = data.CloneObjectValue(v)
	}
	c.variables = append(c.variables, data.NewNamedZVal(name, stored))
}

// GetVariableByName 通过变量名读取变量值
func (c *Context) GetVariableByName(name string) (data.Value, bool) {
	for _, zv := range c.variables {
		if zv != nil && zv.Name == name {
			return zv.Value, true
		}
	}
	return nil, false
}

// HasVariableByName 检查调用者上下文中是否已存在指定名称的变量。
// 仅「已赋值」的槽算存在；符号表预留但未赋值的槽不算（对齐 PHP EXTR_SKIP）。
func (c *Context) HasVariableByName(name string) bool {
	for _, zv := range c.variables {
		if zv != nil && zv.Name == name && zv.Defined {
			return true
		}
	}
	return false
}

// GetDefinedVariables 返回当前作用域已赋值变量的快照（对齐 get_defined_vars）。
func (c *Context) GetDefinedVariables() map[string]data.Value {
	result := make(map[string]data.Value)
	for _, zv := range c.variables {
		if zv != nil && zv.Name != "" && zv.Defined {
			result[zv.Name] = zv.Value
		}
	}
	return result
}

// NewContextToDo 不实现具体功能的上下文
func NewContextToDo() data.Context {
	vm := NewVM(&parser.Parser{})
	return vm.CreateContext([]data.Variable{})
}
