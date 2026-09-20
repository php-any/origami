package runtime

import (
	"fmt"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/perfmon"
)

// NewTempVM 根据给定 VM 创建/返回一个临时 VM 实例
func NewTempVM(vm data.VM) data.VM {
	var base *VM
	switch v := vm.(type) {
	case *TempVM:
		base = v.Base
	case *VM:
		base = v
	default:
		return vm
	}
	return &TempVM{
		Base:            base,
		addedClasses:    make(map[string]data.ClassStmt),
		addedInterfaces: make(map[string]data.InterfaceStmt),
		addedFuncs:      make(map[string]data.FuncStmt),
		out:             newOutputState(),
	}
}

// TempVM 用于模拟 php-fpm 请求级生效的 VM（热重载）
// 确保解析阶段（Parser）也绑定到 TempVM
type TempVM struct {
	Base   *VM
	parser *parser.Parser
	mu     sync.RWMutex

	addedClasses    map[string]data.ClassStmt
	addedInterfaces map[string]data.InterfaceStmt
	addedFuncs      map[string]data.FuncStmt
	out             *outputState
	throwHandler    func(data.Control)

	globalsArray *data.ObjectValue
	sessionArray *data.ObjectValue

	call CallState
}

func (vm *TempVM) AddClass(c data.ClassStmt) data.Control {
	// 仅注册到临时 VM 的映射中（请求级生效）
	vm.addedClasses[c.GetName()] = c
	return nil
}

// AddedClasses 返回本 TempVM 解析阶段注册的类（不合并 Base VM）。
func (vm *TempVM) AddedClasses() []data.ClassStmt {
	out := make([]data.ClassStmt, 0, len(vm.addedClasses))
	for _, c := range vm.addedClasses {
		out = append(out, c)
	}
	return out
}

func (vm *TempVM) AddInterface(i data.InterfaceStmt) data.Control {
	vm.addedInterfaces[i.GetName()] = i
	return nil
}

func (vm *TempVM) AddFunc(f data.FuncStmt) data.Control {
	if vm.addedFuncs == nil {
		vm.addedFuncs = make(map[string]data.FuncStmt)
	}
	vm.addedFuncs[f.GetName()] = f
	return nil
}

func (vm *TempVM) CreateContext(vars []data.Variable) data.Context {
	ctx := vm.Base.CreateContext(vars)
	if rctx, ok := ctx.(*Context); ok {
		rctx.SetVM(vm)
	}
	if st := currentRequestCallState(); st != nil {
		BindContextCallState(ctx, st)
	} else {
		BindContextCallState(ctx, &vm.call)
	}
	BindContextOutput(ctx, vm.activeOut())
	return ctx
}

// PrepareParse 克隆 Parser 并绑定到本 TempVM，供解析阶段 GetOrLoadClass 等使用。
func (vm *TempVM) PrepareParse(base *parser.Parser) *parser.Parser {
	p := base.Clone()
	p.SetVM(vm)
	vm.parser = p
	return p
}

func (vm *TempVM) LoadAndRun(file string) (data.GetValue, data.Control) {
	file = normalizePhpFilePath(file)
	if vm.Base.GetPhpFileCache(file) {
		return nil, nil
	}

	p := vm.PrepareParse(vm.Base.parser)

	program, acl := p.ParseFile(file)
	if acl != nil {
		return nil, acl
	}
	result, ctrl := program.GetValue(vm.CreateContext(p.GetVariables()))
	if ctrl == nil {
		vm.Base.SetPhpFileCache(file)
	}
	return result, ctrl
}

// LoadInCallerContext 在 TempVM 上解析并执行，保证请求级类/函数注册不泄漏到 Base。
// 不设置 PhpFileCache，允许同一视图文件被多次 require。
func (vm *TempVM) LoadInCallerContext(parent data.Context, file string) (data.GetValue, data.Control) {
	file = normalizePhpFilePath(file)

	t0 := perfmon.Now()
	program, vars, acl := vm.Base.ParseFileCached(file)
	if acl != nil {
		return nil, acl
	}

	ctx := inheritCallerScope(parent, vm.CreateContext(vars))
	injectCallerVariables(parent, ctx, vars, includeGlobalBinder(parent, func(name string, variable data.Variable) {
		vm.Base.bindIncludedVarToGlobal(name, variable.GetIndex(), ctx)
	}))

	result, ctrl := program.GetValue(ctx)
	perfmon.NoteInclude(file, perfmon.Since(t0))
	return result, ctrl
}

// CompileLoad 编译模式专用：仅解析并注册类/函数/接口，不执行顶层代码。
func (vm *TempVM) CompileLoad(file string) data.Control {
	return vm.Base.CompileLoad(file)
}

func (vm *TempVM) ParseFile(file string, data data.Value) (data.Value, data.Control) {
	return vm.Base.ParseFile(file, data)
}

func lookupTempClass(added map[string]data.ClassStmt, pkg string) (data.ClassStmt, bool) {
	if c, ok := added[pkg]; ok {
		return c, true
	}
	if len(pkg) > 0 && pkg[0] == '\\' {
		if c, ok := added[pkg[1:]]; ok {
			return c, true
		}
	}
	want := phpIdentKey(pkg)
	for k, c := range added {
		if phpIdentKey(k) == want {
			return c, true
		}
	}
	return nil, false
}

func (vm *TempVM) GetClass(pkg string) (data.ClassStmt, bool) {
	ret, ok := vm.Base.GetClass(pkg)
	if ok {
		return ret, ok
	}
	return lookupTempClass(vm.addedClasses, pkg)
}

func (vm *TempVM) GetOrLoadClass(pkg string) (data.ClassStmt, data.Control) {
	c, ok := vm.Base.GetClass(pkg)
	if ok {
		return c, nil
	}
	// 优先从本请求新增的类中查找
	if c, ok := lookupTempClass(vm.addedClasses, pkg); ok {
		return c, nil
	}
	if vm.parser == nil {
		return vm.Base.GetOrLoadClass(pkg)
	}
	acl := vm.parser.GetClassPathManager().LoadClass(pkg, vm.parser)
	if acl != nil {
		return nil, acl
	}
	if c, ok := lookupTempClass(vm.addedClasses, pkg); ok {
		return c, nil
	}
	// 编译模式下类可能注册到 baseVM（CompileLoad 不经过 TempVM）
	if data.CompileMode {
		if c, ok := vm.Base.GetClass(pkg); ok {
			return c, nil
		}
	}

	return nil, data.NewErrorThrow(nil, fmt.Errorf("class %s not found", pkg))
}

func (vm *TempVM) LoadPkg(pkg string) (data.GetValue, data.Control) {
	if c, ok := vm.addedClasses[pkg]; ok {
		return c, nil
	}
	if c, ok := vm.addedInterfaces[pkg]; ok {
		return c, nil
	}
	c, acl := vm.Base.LoadPkg(pkg)
	if acl != nil {
		return nil, acl
	}
	if c != nil {
		return c, nil
	}
	if vm.parser == nil {
		return nil, nil
	}
	_, ok := vm.parser.GetClassPathManager().FindClassFile(pkg)
	if !ok {
		return nil, nil
	}
	if acl = vm.parser.GetClassPathManager().LoadClass(pkg, vm.parser); acl != nil {
		return nil, acl
	}
	if c, ok := vm.addedClasses[pkg]; ok {
		return c, nil
	}
	if c, ok := vm.addedInterfaces[pkg]; ok {
		return c, nil
	}
	// 编译模式下类可能注册到 baseVM（CompileLoad 不经过 TempVM）
	if data.CompileMode {
		if c, ok := vm.Base.GetClass(pkg); ok {
			return c, nil
		}
		if c, ok := vm.Base.GetInterface(pkg); ok {
			return c, nil
		}
	}
	return nil, nil
}

func (vm *TempVM) GetInterface(pkg string) (data.InterfaceStmt, bool) {
	ret, ok := vm.Base.GetInterface(pkg)
	if ok {
		return ret, ok
	}
	if c, ok := vm.addedInterfaces[pkg]; ok {
		return c, true
	}
	return nil, false
}

func (vm *TempVM) GetOrLoadInterface(pkg string) (data.InterfaceStmt, data.Control) {
	// 优先从本请求新增的接口中查找
	if c, ok := vm.addedInterfaces[pkg]; ok {
		return c, nil
	}

	// 从 Base VM 查找或加载
	ret, acl := vm.Base.GetOrLoadInterface(pkg)
	if acl == nil && ret != nil {
		return ret, nil
	}
	if acl != nil {
		return nil, acl
	}

	// 如果 Base 加载后，再次检查新增的接口
	if c, ok := vm.addedInterfaces[pkg]; ok {
		return c, nil
	}

	return nil, data.NewErrorThrow(nil, fmt.Errorf("interface %s not found", pkg))
}

func (vm *TempVM) GetFunc(pkg string) (data.FuncStmt, bool) {
	if f, ok := vm.addedFuncs[pkg]; ok {
		return f, true
	}
	return vm.Base.GetFunc(pkg)
}
func (vm *TempVM) RegisterFunction(name string, fn interface{}) data.Control {
	return vm.Base.RegisterFunction(name, fn)
}
func (vm *TempVM) RegisterReflectClass(name string, instance interface{}) data.Control {
	return vm.Base.RegisterReflectClass(name, instance)
}
func (vm *TempVM) SetThrowControl(fn func(acl data.Control)) {
	vm.mu.Lock()
	vm.throwHandler = fn
	vm.mu.Unlock()
}
func (vm *TempVM) ThrowControl(acl data.Control) {
	vm.mu.RLock()
	handler := vm.throwHandler
	vm.mu.RUnlock()
	if handler != nil {
		handler(acl)
		return
	}
	panic(acl)
}
func (vm *TempVM) SetPhpFileCache(file string) { vm.Base.SetPhpFileCache(file) }
func (vm *TempVM) GetPhpFileCache(file string) bool {
	return vm.Base.GetPhpFileCache(file)
}

func (vm *TempVM) GetIncludeOnceResult(file string) (data.GetValue, bool) {
	return vm.Base.GetIncludeOnceResult(file)
}

func (vm *TempVM) SetIncludeOnceResult(file string, result data.GetValue) {
	vm.Base.SetIncludeOnceResult(file, result)
}

// AddNamespace 添加命名空间路径映射（委托给 Base VM，因为命名空间映射是全局共享的）
func (vm *TempVM) AddNamespace(namespace string, path string) {
	vm.Base.AddNamespace(namespace, path)
}

// SetConstant 设置全局常量（委托给 Base VM，因为常量是全局的）
func (vm *TempVM) SetConstant(name string, value data.Value) data.Control {
	return vm.Base.SetConstant(name, value)
}

// GetConstant 获取全局常量（从 Base VM 获取，因为常量是全局的）
func (vm *TempVM) GetConstant(name string) (data.Value, bool) {
	return vm.Base.GetConstant(name)
}

// EnsureGlobalZVal 委托给 Base；PHP global 是共享语言状态，不属于请求输出通道。
func (vm *TempVM) EnsureGlobalZVal(name string) *data.ZVal {
	return vm.Base.EnsureGlobalZVal(name)
}

// EnsureGlobalsArray 本 TempVM 独立的 $GLOBALS（热重载/请求包装不串态）。
func (vm *TempVM) EnsureGlobalsArray() *data.ObjectValue {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if vm.globalsArray == nil {
		vm.globalsArray = data.NewObjectValue()
	}
	return vm.globalsArray
}

// EnsureSessionArray 本 TempVM 独立的 $_SESSION。
func (vm *TempVM) EnsureSessionArray() *data.ObjectValue {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if vm.sessionArray == nil {
		vm.sessionArray = data.NewObjectValue()
	}
	return vm.sessionArray
}

func (vm *TempVM) SetExceptionHandler(handler data.Value) data.Value {
	return vm.Base.SetExceptionHandler(handler)
}

func (vm *TempVM) GetExceptionHandler() data.Value {
	return vm.Base.GetExceptionHandler()
}

func (vm *TempVM) SetErrorHandler(handler data.Value) data.Value {
	return vm.Base.SetErrorHandler(handler)
}

func (vm *TempVM) RestoreErrorHandler() bool {
	return vm.Base.RestoreErrorHandler()
}

func (vm *TempVM) GetErrorHandler() data.Value {
	return vm.Base.GetErrorHandler()
}

func (vm *TempVM) AddShutdownCallback(cb data.Value) {
	vm.Base.AddShutdownCallback(cb)
}

func (vm *TempVM) RunShutdownCallbacks() {
	vm.Base.RunShutdownCallbacks()
}

func (vm *TempVM) requestCall() *CallState {
	if st := currentRequestCallState(); st != nil {
		return st
	}
	return &vm.call
}

func (vm *TempVM) EnterCall() int {
	return vm.requestCall().Enter()
}

func (vm *TempVM) LeaveCall() {
	st := vm.requestCall()
	st.Leave()
	releaseAutoCallState(st)
}

func (vm *TempVM) PushCallFrame(frame data.CallFrame) {
	vm.requestCall().Push(frame)
}

func (vm *TempVM) PopCallFrame() {
	vm.requestCall().Pop()
}

func (vm *TempVM) SnapshotCallStack() []data.CallFrame {
	return vm.requestCall().Snapshot()
}

// RegisterCompiledFile 注册预编译的文件 AST（委托给 Base VM）
func (vm *TempVM) RegisterCompiledFile(file string, fn func() (data.GetValue, []data.Variable)) {
	vm.Base.RegisterCompiledFile(file, fn)
}

// RunCompiledFile 执行预编译文件（委托给 Base VM）
func (vm *TempVM) RunCompiledFile(file string) (data.GetValue, data.Control) {
	return vm.Base.RunCompiledFile(file)
}
