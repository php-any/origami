package requestvm

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/runtime"
)

// LaravelRequestVM 包装共享 runtime.VM：类/函数全部委托 base，仅隔离本请求的
// output sink 与 ob_* 缓冲栈。不使用 runtime.TempVM。
type LaravelRequestVM struct {
	base   *runtime.VM
	output data.OutputWriter
	ob     obStack

	mu        sync.RWMutex
	request   *http.Request
	response  http.ResponseWriter
	lastThrow data.Control

	callDepth int
	callStack []data.CallFrame
}

// New 创建请求级 VM 包装。output 为本请求 HTTP body 写入目标。
func New(base *runtime.VM, output data.OutputWriter) *LaravelRequestVM {
	if output == nil {
		output = data.DefaultOutputWriter
	}
	return &LaravelRequestVM{
		base:   base,
		output: output,
		ob:     obStack{},
	}
}

// Base 返回共享底层 VM。
func (v *LaravelRequestVM) Base() *runtime.VM { return v.base }

// BindHTTP 绑定本请求的 Go Request/ResponseWriter。
func (v *LaravelRequestVM) BindHTTP(req *http.Request, resp http.ResponseWriter) {
	v.mu.Lock()
	v.request = req
	v.response = resp
	v.mu.Unlock()
}

// HTTPRequest 返回本请求的 *http.Request。
func (v *LaravelRequestVM) HTTPRequest() *http.Request {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.request
}

// HTTPResponseWriter 返回本请求的 ResponseWriter。
func (v *LaravelRequestVM) HTTPResponseWriter() http.ResponseWriter {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.response
}

// WriteOutput 实现 data.OutputSink：有 ob 层时写入缓冲，否则写入本请求 output。
func (v *LaravelRequestVM) WriteOutput(s string) {
	v.mu.Lock()
	if v.ob.write(s) {
		v.mu.Unlock()
		return
	}
	v.mu.Unlock()
	v.output(s)
}

func (v *LaravelRequestVM) StartOutputBuffer() {
	v.mu.Lock()
	v.ob.push()
	v.mu.Unlock()
}

func (v *LaravelRequestVM) CleanOutputBuffer() (string, bool) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.ob.level() == 0 {
		return "", false
	}
	return v.ob.pop(), true
}

func (v *LaravelRequestVM) OutputBufferContents() (string, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	if v.ob.level() == 0 {
		return "", false
	}
	return v.ob.contents(), true
}

func (v *LaravelRequestVM) OutputBufferLevel() int {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.ob.level()
}

func (v *LaravelRequestVM) AddClass(c data.ClassStmt) data.Control {
	return v.base.AddClass(c)
}
func (v *LaravelRequestVM) GetClass(pkg string) (data.ClassStmt, bool) {
	return v.base.GetClass(pkg)
}
func (v *LaravelRequestVM) GetOrLoadClass(pkg string) (data.ClassStmt, data.Control) {
	return v.base.GetOrLoadClass(pkg)
}
func (v *LaravelRequestVM) LoadPkg(pkg string) (data.GetValue, data.Control) {
	return v.base.LoadPkg(pkg)
}
func (v *LaravelRequestVM) AddInterface(i data.InterfaceStmt) data.Control {
	return v.base.AddInterface(i)
}
func (v *LaravelRequestVM) GetInterface(pkg string) (data.InterfaceStmt, bool) {
	return v.base.GetInterface(pkg)
}
func (v *LaravelRequestVM) GetOrLoadInterface(pkg string) (data.InterfaceStmt, data.Control) {
	return v.base.GetOrLoadInterface(pkg)
}
func (v *LaravelRequestVM) AddFunc(f data.FuncStmt) data.Control {
	return v.base.AddFunc(f)
}
func (v *LaravelRequestVM) GetFunc(pkg string) (data.FuncStmt, bool) {
	return v.base.GetFunc(pkg)
}
func (v *LaravelRequestVM) RegisterFunction(name string, fn interface{}) data.Control {
	return v.base.RegisterFunction(name, fn)
}
func (v *LaravelRequestVM) RegisterReflectClass(name string, instance interface{}) data.Control {
	return v.base.RegisterReflectClass(name, instance)
}

func (v *LaravelRequestVM) CreateContext(vars []data.Variable) data.Context {
	ctx := v.base.CreateContext(vars)
	if rctx, ok := ctx.(interface{ SetVM(data.VM) }); ok {
		rctx.SetVM(v)
	}
	return ctx
}

func (v *LaravelRequestVM) SetThrowControl(fn func(acl data.Control)) {
	// 请求级 VM 自行捕获 throw，避免委托 base 的 os.Exit。
}
func (v *LaravelRequestVM) ThrowControl(acl data.Control) {
	v.mu.Lock()
	v.lastThrow = acl
	v.mu.Unlock()
}

func (v *LaravelRequestVM) takeThrow() data.Control {
	v.mu.Lock()
	defer v.mu.Unlock()
	acl := v.lastThrow
	v.lastThrow = nil
	return acl
}

// TakeThrow 取出并清空最近一次被吞掉的 throw（供 HTTP 层展示）。
func (v *LaravelRequestVM) TakeThrow() data.Control {
	return v.takeThrow()
}

func (v *LaravelRequestVM) LoadAndRun(file string) (data.GetValue, data.Control) {
	return v.runFile(file, false)
}

// LoadAndRunFresh 强制重新读盘并解析（热重载/调试）。
func (v *LaravelRequestVM) LoadAndRunFresh(file string) (data.GetValue, data.Control) {
	return v.runFile(file, true)
}

func (v *LaravelRequestVM) runFile(file string, fresh bool) (data.GetValue, data.Control) {
	_ = v.takeThrow()

	var program data.GetValue
	var vars []data.Variable
	var acl data.Control

	if fresh {
		p := v.base.GetParser().Clone()
		p.SetVM(v)
		program, acl = p.ParseFile(file)
		if acl != nil {
			return nil, acl
		}
		vars = p.GetVariables()
	} else {
		program, vars, acl = v.base.ParseFileCached(file)
		if acl != nil {
			return nil, acl
		}
	}

	ctx := v.CreateContext(vars)
	v.base.RegisterGlobalContext(vars, ctx)
	result, ctrl := program.GetValue(ctx)
	if ctrl == nil {
		if thrown := v.takeThrow(); thrown != nil {
			return result, thrown
		}
	}
	return result, ctrl
}

func (v *LaravelRequestVM) LoadInCallerContext(parent data.Context, file string) (data.GetValue, data.Control) {
	program, vars, acl := v.base.ParseFileCached(file)
	if acl != nil {
		return nil, acl
	}
	ctx := v.CreateContext(vars)
	for _, variable := range vars {
		name := variable.GetName()
		if name == "" {
			continue
		}
		if val, ok := parent.GetVariableByName(name); ok && val != nil {
			if ctl := variable.SetValue(ctx, val); ctl != nil {
				return nil, ctl
			}
			continue
		}
		ctx.SetIndexZVal(variable.GetIndex(), v.EnsureGlobalZVal(name))
	}
	return program.GetValue(ctx)
}

func (v *LaravelRequestVM) CompileLoad(file string) data.Control {
	return v.base.CompileLoad(file)
}
func (v *LaravelRequestVM) RegisterCompiledFile(file string, fn func() (data.GetValue, []data.Variable)) {
	v.base.RegisterCompiledFile(file, fn)
}
func (v *LaravelRequestVM) RunCompiledFile(file string) (data.GetValue, data.Control) {
	return v.base.RunCompiledFile(file)
}
func (v *LaravelRequestVM) ParseFile(file string, value data.Value) (data.Value, data.Control) {
	return v.base.ParseFile(file, value)
}
func (v *LaravelRequestVM) SetPhpFileCache(file string)      { v.base.SetPhpFileCache(file) }
func (v *LaravelRequestVM) GetPhpFileCache(file string) bool { return v.base.GetPhpFileCache(file) }
func (v *LaravelRequestVM) GetIncludeOnceResult(file string) (data.GetValue, bool) {
	return v.base.GetIncludeOnceResult(file)
}
func (v *LaravelRequestVM) SetIncludeOnceResult(file string, result data.GetValue) {
	v.base.SetIncludeOnceResult(file, result)
}
func (v *LaravelRequestVM) AddNamespace(namespace string, path string) {
	v.base.AddNamespace(namespace, path)
}
func (v *LaravelRequestVM) EnterCall() int {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.callDepth++
	return v.callDepth
}

func (v *LaravelRequestVM) LeaveCall() {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.callDepth > 0 {
		v.callDepth--
	}
}
func (v *LaravelRequestVM) SetConstant(name string, value data.Value) data.Control {
	return v.base.SetConstant(name, value)
}
func (v *LaravelRequestVM) GetConstant(name string) (data.Value, bool) {
	return v.base.GetConstant(name)
}
func (v *LaravelRequestVM) EnsureGlobalZVal(name string) *data.ZVal {
	return v.base.EnsureGlobalZVal(name)
}
func (v *LaravelRequestVM) SetExceptionHandler(handler data.Value) data.Value {
	return v.base.SetExceptionHandler(handler)
}
func (v *LaravelRequestVM) GetExceptionHandler() data.Value {
	return v.base.GetExceptionHandler()
}
func (v *LaravelRequestVM) SetErrorHandler(handler data.Value) data.Value {
	return v.base.SetErrorHandler(handler)
}
func (v *LaravelRequestVM) RestoreErrorHandler() bool {
	return v.base.RestoreErrorHandler()
}
func (v *LaravelRequestVM) GetErrorHandler() data.Value {
	return v.base.GetErrorHandler()
}
func (v *LaravelRequestVM) AddShutdownCallback(cb data.Value) {
	v.base.AddShutdownCallback(cb)
}
func (v *LaravelRequestVM) RunShutdownCallbacks() {
	v.base.RunShutdownCallbacks()
}

func (v *LaravelRequestVM) PushCallFrame(frame data.CallFrame) {
	v.mu.Lock()
	v.callStack = append(v.callStack, frame)
	v.mu.Unlock()
}

func (v *LaravelRequestVM) PopCallFrame() {
	v.mu.Lock()
	if n := len(v.callStack); n > 0 {
		v.callStack = v.callStack[:n-1]
	}
	v.mu.Unlock()
}

func (v *LaravelRequestVM) SnapshotCallStack() []data.CallFrame {
	v.mu.RLock()
	defer v.mu.RUnlock()
	if len(v.callStack) == 0 {
		return nil
	}
	out := make([]data.CallFrame, len(v.callStack))
	copy(out, v.callStack)
	return out
}

// FromContext 从执行上下文取出 LaravelRequestVM。
func FromContext(ctx data.Context) (*LaravelRequestVM, error) {
	if ctx == nil {
		return nil, fmt.Errorf("requestvm: nil context")
	}
	vm, ok := ctx.GetVM().(*LaravelRequestVM)
	if !ok || vm == nil {
		return nil, fmt.Errorf("requestvm: current VM is not LaravelRequestVM")
	}
	return vm, nil
}

var (
	_ data.VM               = (*LaravelRequestVM)(nil)
	_ data.OutputSink       = (*LaravelRequestVM)(nil)
	_ data.OutputBufferHost = (*LaravelRequestVM)(nil)
	_ data.CallStackTracker = (*LaravelRequestVM)(nil)
)
