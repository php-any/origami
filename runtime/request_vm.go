package runtime

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/perfmon"
)

// NewRequestVM creates an execution scope with request-owned mutable PHP state.
func NewRequestVM(vm data.VM) data.VM {
	var base *VM
	switch v := vm.(type) {
	case *RequestVM:
		base = v.Base
	case *VM:
		base = v
	default:
		return vm
	}
	return &RequestVM{
		Base:               base,
		addedClasses:       make(map[string]data.ClassStmt),
		addedInterfaces:    make(map[string]data.InterfaceStmt),
		addedFuncs:         make(map[string]data.FuncStmt),
		out:                newOutputState(),
		constants:          make(map[string]data.Value),
		phpFileCache:       make(map[string]struct{}),
		includeOnceResults: make(map[string]data.GetValue),
		globalVars:         make(map[string]*data.ZVal),
		initialFiles:       base.snapshotRequestFiles(),
	}
}

// RequestVM is the sole request execution VM. Base supplies registered
// language facilities and cached programs; mutable execution state stays here.
type RequestVM struct {
	Base   *VM
	parser *parser.Parser
	mu     sync.RWMutex

	addedClasses       map[string]data.ClassStmt
	addedInterfaces    map[string]data.InterfaceStmt
	addedFuncs         map[string]data.FuncStmt
	out                *outputState
	throwHandler       func(data.Control)
	lastThrow          data.Control
	request            *http.Request
	response           http.ResponseWriter
	constants          map[string]data.Value
	phpFileCache       map[string]struct{}
	initialFiles       *requestFileSnapshot
	includeOnceResults map[string]data.GetValue
	globalVars         map[string]*data.ZVal
	errorHandlers      *errorHandlerState

	globalsArray *data.ObjectValue
	sessionArray *data.ObjectValue

	call              CallState
	shutdown          shutdownQueue
	staticLocals      map[any]*data.StaticLocals
	exceptionHandlers *ExceptionHandlerState
}

func (vm *RequestVM) ScopeStaticLocals(identity any) *data.StaticLocals {
	if vm.staticLocals == nil {
		vm.staticLocals = make(map[any]*data.StaticLocals)
	}
	if store := vm.staticLocals[identity]; store != nil {
		return store
	}
	store := data.NewStaticLocals()
	vm.staticLocals[identity] = store
	return store
}

func (vm *RequestVM) AddClass(c data.ClassStmt) data.Control {
	// Declarations loaded during execution belong to this request.
	vm.addedClasses[c.GetName()] = c
	return nil
}

// AddedClasses 返回本 RequestVM 解析阶段注册的类（不合并 Base VM）。
func (vm *RequestVM) AddedClasses() []data.ClassStmt {
	out := make([]data.ClassStmt, 0, len(vm.addedClasses))
	for _, c := range vm.addedClasses {
		out = append(out, c)
	}
	return out
}

func (vm *RequestVM) registerParsedDeclarations(entry *parsedPHPFile) {
	for _, declaration := range entry.classes {
		vm.addedClasses[declaration.GetName()] = declaration
	}
	for _, declaration := range entry.interfaces {
		vm.addedInterfaces[declaration.GetName()] = declaration
	}
}

func (vm *RequestVM) AddInterface(i data.InterfaceStmt) data.Control {
	vm.addedInterfaces[i.GetName()] = i
	return nil
}

func (vm *RequestVM) AddFunc(f data.FuncStmt) data.Control {
	if vm.addedFuncs == nil {
		vm.addedFuncs = make(map[string]data.FuncStmt)
	}
	vm.addedFuncs[f.GetName()] = f
	return nil
}

func (vm *RequestVM) CreateContext(vars []data.Variable) data.Context {
	ctx := vm.Base.CreateContext(vars)
	if rctx, ok := ctx.(*Context); ok {
		rctx.SetVM(vm)
	}
	BindContextCallState(ctx, vm.requestCall())
	BindContextOutput(ctx, vm.activeOut())
	return ctx
}

// PrepareParse 克隆 Parser 并绑定到本 RequestVM，供解析阶段 GetOrLoadClass 等使用。
func (vm *RequestVM) PrepareParse(base *parser.Parser) *parser.Parser {
	p := base.Clone()
	p.SetVM(vm)
	vm.parser = p
	return p
}

func (vm *RequestVM) LoadAndRun(file string) (data.GetValue, data.Control) {
	if vm.GetPhpFileCache(file) {
		return nil, nil
	}
	return vm.runFile(file, false)
}

func (vm *RequestVM) LoadAndRunFresh(file string) (data.GetValue, data.Control) {
	return vm.runFile(file, true)
}

func (vm *RequestVM) runFile(file string, fresh bool) (data.GetValue, data.Control) {
	_ = vm.TakeThrow()
	var program data.GetValue
	var vars []data.Variable
	var ctl data.Control
	if fresh {
		p := vm.PrepareParse(vm.Base.parser)
		program, ctl = p.ParseFile(file)
		vars = p.GetVariables()
	} else {
		program, vars, ctl = vm.Base.parseFileCachedFor(file, vm)
	}
	if ctl != nil {
		return nil, ctl
	}
	ctx := vm.CreateContext(vars)
	vm.RegisterGlobalContext(vars, ctx)
	result, ctl := program.GetValue(ctx)
	if ctl == nil {
		ctl = vm.TakeThrow()
	}
	if ctl == nil {
		vm.SetPhpFileCache(file)
	}
	return result, ctl
}

// LoadInCallerContext 在 RequestVM 上解析并执行，保证请求级类/函数注册不泄漏到 Base。
// 不设置 PhpFileCache，允许同一视图文件被多次 require。
func (vm *RequestVM) LoadInCallerContext(parent data.Context, file string) (data.GetValue, data.Control) {
	file = normalizePhpFilePath(file)

	t0 := perfmon.Now()
	program, vars, acl := vm.Base.parseFileCachedFor(file, vm)
	if acl != nil {
		return nil, acl
	}

	ctx := inheritCallerScope(parent, vm.CreateContext(vars))
	injectCallerVariables(parent, ctx, vars, includeGlobalBinder(parent, func(name string, variable data.Variable) {
		vm.bindIncludedVarToGlobal(name, variable.GetIndex(), ctx)
	}))

	result, ctrl := program.GetValue(ctx)
	perfmon.NoteInclude(file, perfmon.Since(t0))
	return result, ctrl
}

// CompileLoad 编译模式专用：仅解析并注册类/函数/接口，不执行顶层代码。
func (vm *RequestVM) CompileLoad(file string) data.Control {
	if vm.GetPhpFileCache(file) {
		return nil
	}
	_, _, ctl := vm.Base.parseFileCachedFor(file, vm)
	if ctl == nil {
		vm.SetPhpFileCache(file)
	}
	return ctl
}

func (vm *RequestVM) ParseFile(file string, object data.Value) (data.Value, data.Control) {
	program, variables, ctl := vm.Base.parseFileCachedFor(file, vm)
	if ctl != nil {
		return nil, ctl
	}
	return runTemplateFile(vm, file, program, variables, object)
}

func lookupRequestClass(added map[string]data.ClassStmt, pkg string) (data.ClassStmt, bool) {
	if c, ok := added[pkg]; ok {
		return c, true
	}
	if len(pkg) > 0 && pkg[0] == '\\' {
		if c, ok := added[pkg[1:]]; ok {
			return c, true
		}
	}
	for k, c := range added {
		if data.TypeNameEqual(k, pkg) {
			return c, true
		}
	}
	return nil, false
}

func (vm *RequestVM) GetClass(pkg string) (data.ClassStmt, bool) {
	ret, ok := vm.Base.GetClass(pkg)
	if ok {
		return ret, ok
	}
	return lookupRequestClass(vm.addedClasses, pkg)
}

func (vm *RequestVM) GetOrLoadClass(pkg string) (data.ClassStmt, data.Control) {
	if c, ok := vm.GetClass(pkg); ok {
		return c, nil
	}
	if _, ctl := vm.LoadPkg(pkg); ctl != nil {
		return nil, ctl
	}
	if c, ok := vm.GetClass(pkg); ok {
		return c, nil
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("class %s not found", pkg))
}

func (vm *RequestVM) LoadPkg(pkg string) (data.GetValue, data.Control) {
	if c, ok := vm.GetClass(pkg); ok {
		return c, nil
	}
	if c, ok := vm.GetInterface(pkg); ok {
		return c, nil
	}
	if vm.parser == nil {
		vm.PrepareParse(vm.Base.parser)
	}
	if ctl := vm.parser.GetClassPathManager().LoadClass(pkg, vm.parser); ctl != nil {
		return nil, ctl
	}
	if c, ok := vm.GetClass(pkg); ok {
		return c, nil
	}
	if c, ok := vm.GetInterface(pkg); ok {
		return c, nil
	}
	return nil, nil
}

func (vm *RequestVM) GetInterface(pkg string) (data.InterfaceStmt, bool) {
	ret, ok := vm.Base.GetInterface(pkg)
	if ok {
		return ret, ok
	}
	if c, ok := vm.addedInterfaces[pkg]; ok {
		return c, true
	}
	for name, c := range vm.addedInterfaces {
		if data.TypeNameEqual(name, pkg) {
			return c, true
		}
	}
	return nil, false
}

func (vm *RequestVM) GetOrLoadInterface(pkg string) (data.InterfaceStmt, data.Control) {
	if c, ok := vm.GetInterface(pkg); ok {
		return c, nil
	}
	if _, ctl := vm.LoadPkg(pkg); ctl != nil {
		return nil, ctl
	}
	if c, ok := vm.GetInterface(pkg); ok {
		return c, nil
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("interface %s not found", pkg))
}

func (vm *RequestVM) GetFunc(pkg string) (data.FuncStmt, bool) {
	if f, ok := vm.addedFuncs[pkg]; ok {
		return f, true
	}
	return vm.Base.GetFunc(pkg)
}
func (vm *RequestVM) RegisterFunction(name string, fn interface{}) data.Control {
	return vm.Base.RegisterFunction(name, fn)
}
func (vm *RequestVM) RegisterReflectClass(name string, instance interface{}) data.Control {
	return vm.Base.RegisterReflectClass(name, instance)
}
func (vm *RequestVM) SetThrowControl(fn func(acl data.Control)) {
	vm.mu.Lock()
	vm.throwHandler = fn
	vm.mu.Unlock()
}
func (vm *RequestVM) ThrowControl(acl data.Control) {
	vm.mu.RLock()
	handler := vm.throwHandler
	vm.mu.RUnlock()
	if handler != nil {
		handler(acl)
		return
	}
	vm.lastThrow = acl
}
func (vm *RequestVM) SetPhpFileCache(file string) {
	vm.phpFileCache[normalizePhpFilePath(file)] = struct{}{}
}
func (vm *RequestVM) GetPhpFileCache(file string) bool {
	file = normalizePhpFilePath(file)
	_, found := vm.phpFileCache[file]
	if !found {
		_, found = vm.initialFiles.loaded[file]
	}
	return found
}
func (vm *RequestVM) GetIncludeOnceResult(file string) (data.GetValue, bool) {
	file = normalizePhpFilePath(file)
	result, found := vm.includeOnceResults[file]
	if !found {
		result, found = vm.initialFiles.results[file]
	}
	return result, found
}
func (vm *RequestVM) SetIncludeOnceResult(file string, result data.GetValue) {
	vm.includeOnceResults[normalizePhpFilePath(file)] = result
}

// AddNamespace 添加命名空间路径映射（委托给 Base VM，因为命名空间映射是全局共享的）
func (vm *RequestVM) AddNamespace(namespace string, path string) {
	vm.Base.AddNamespace(namespace, path)
}

func (vm *RequestVM) SetConstant(name string, value data.Value) data.Control {
	if _, found := vm.constants[name]; found {
		return data.NewErrorThrow(nil, fmt.Errorf("constant %s already defined", name))
	}
	vm.constants[name] = value
	return nil
}
func (vm *RequestVM) GetConstant(name string) (data.Value, bool) {
	if value, found := vm.constants[name]; found {
		return value, true
	}
	return vm.Base.GetConstant(name)
}
func (vm *RequestVM) EnsureGlobalZVal(name string) *data.ZVal {
	if slot, found := vm.globalVars[name]; found {
		return slot
	}
	slot := data.NewZVal(data.NewNullValue())
	vm.globalVars[name] = slot
	return slot
}
func (vm *RequestVM) RegisterGlobalContext(vars []data.Variable, ctx data.Context) {
	for _, variable := range vars {
		if variable == nil || variable.GetName() == "" {
			continue
		}
		vm.bindIncludedVarToGlobal(variable.GetName(), variable.GetIndex(), ctx)
	}
}
func (vm *RequestVM) bindIncludedVarToGlobal(name string, index int, ctx data.Context) {
	if slot, found := vm.globalVars[name]; found {
		ctx.SetIndexZVal(index, slot)
		return
	}
	vm.globalVars[name] = ctx.GetIndexZVal(index)
}

// EnsureGlobalsArray 本 RequestVM 独立的 $GLOBALS（热重载/请求包装不串态）。
func (vm *RequestVM) EnsureGlobalsArray() *data.ObjectValue {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if vm.globalsArray == nil {
		vm.globalsArray = data.NewObjectValue()
	}
	return vm.globalsArray
}

// EnsureSessionArray 本 RequestVM 独立的 $_SESSION。
func (vm *RequestVM) EnsureSessionArray() *data.ObjectValue {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if vm.sessionArray == nil {
		vm.sessionArray = data.NewObjectValue()
	}
	return vm.sessionArray
}

func (vm *RequestVM) AddShutdownCallback(cb data.Value) {
	if st := vm.requestCall(); st != &vm.call {
		st.shutdown.add(cb)
		return
	}
	vm.shutdown.add(cb)
}

func (vm *RequestVM) RunShutdownCallbacks() {
	if st := vm.requestCall(); st != &vm.call {
		st.shutdown.run(vm)
		return
	}
	vm.shutdown.run(vm)
}

func (vm *RequestVM) requestCall() *CallState {
	if vm.out.local {
		return &vm.call
	}
	if st := currentRequestCallState(); st != nil {
		return st
	}
	return &vm.call
}

func (vm *RequestVM) EnterCall() int {
	return vm.requestCall().Enter()
}

func (vm *RequestVM) LeaveCall() {
	st := vm.requestCall()
	st.Leave()
	releaseAutoCallState(st)
}

func (vm *RequestVM) PushCallFrame(frame data.CallFrame) {
	vm.requestCall().Push(frame)
}

func (vm *RequestVM) PopCallFrame() {
	vm.requestCall().Pop()
}

func (vm *RequestVM) SnapshotCallStack() []data.CallFrame {
	return vm.requestCall().Snapshot()
}

// RegisterCompiledFile 注册预编译的文件 AST（委托给 Base VM）
func (vm *RequestVM) RegisterCompiledFile(file string, fn func() (data.GetValue, []data.Variable)) {
	vm.Base.RegisterCompiledFile(file, fn)
}

// RunCompiledFile 从共享缓存取得程序，在当前请求中执行。
func (vm *RequestVM) RunCompiledFile(file string) (data.GetValue, data.Control) {
	file = normalizePhpFilePath(file)
	if vm.GetPhpFileCache(file) {
		return nil, nil
	}
	build, found := syncMapLoad[func() (data.GetValue, []data.Variable)](&vm.Base.compiledFiles, file)
	if !found {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("compiled file %s not found", file))
	}
	program, vars := build()
	ctx := vm.CreateContext(vars)
	vm.RegisterGlobalContext(vars, ctx)
	result, ctl := program.GetValue(ctx)
	if ctl == nil {
		vm.SetPhpFileCache(file)
	}
	return result, ctl
}

func (vm *RequestVM) TakeThrow() data.Control {
	ctl := vm.lastThrow
	vm.lastThrow = nil
	return ctl
}
func (vm *RequestVM) SetOutputWriter(write data.OutputWriter) {
	if write == nil {
		write = data.DefaultOutputWriter
	}
	vm.out.local = true
	vm.out.call = &vm.call
	vm.out.SetSink(write)
}
func (vm *RequestVM) BindHTTP(req *http.Request, response http.ResponseWriter) {
	vm.request, vm.response = req, response
	if flush, ok := response.(http.Flusher); ok {
		vm.activeOut().SetSAPIFlush(flush.Flush)
	}
}
func (vm *RequestVM) HTTPRequest() *http.Request              { return vm.request }
func (vm *RequestVM) HTTPResponseWriter() http.ResponseWriter { return vm.response }

var _ data.VM = (*RequestVM)(nil)
