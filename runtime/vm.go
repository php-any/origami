package runtime

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/utils"
)

func normalizePhpFilePath(file string) string {
	return utils.NormalizePhpFilePath(file)
}

// GetParser 返回 VM 绑定的 Parser（供嵌入层克隆解析）。
func (vm *VM) GetParser() *parser.Parser {
	return vm.parser
}

// NewVM 创建一个新的虚拟机
func NewVM(parser *parser.Parser) data.VM {
	vm := &VM{
		parser:       parser,
		loadingFiles: make(map[string]chan struct{}),
		acl: func(acl data.Control) {
			parser.ShowControl(acl)
			os.Exit(1)
		},
	}
	vm.ctx = NewContext(vm)
	parser.SetVM(vm)

	return vm
}

// VM 表示虚拟机
type VM struct {
	parser *parser.Parser
	ctx    data.Context

	// mu 保护非 map 的请求态字段；类/函数/常量等表使用 sync.Map 以支持并发读。
	mu sync.Mutex

	classMap           sync.Map // string -> data.ClassStmt
	interfaceMap       sync.Map // string -> data.InterfaceStmt
	funcMap            sync.Map // string -> data.FuncStmt
	constantMap        sync.Map // string -> data.Value
	globalVars         sync.Map // string -> *data.ZVal
	phpFileCache       sync.Map // string -> struct{}
	includeOnceResults sync.Map // string -> data.GetValue
	compiledFiles      sync.Map // string -> func() (data.GetValue, []data.Variable)
	parsedFiles        sync.Map // string -> *parsedPHPFile

	// loadingFiles 并发加载同一文件时，后续请求等待首个加载完成
	loadingFiles map[string]chan struct{}

	acl func(acl data.Control)

	// PHP 级 set_exception_handler 注册的回调
	exceptionHandler data.Value
	// 防止在异常处理回调中递归调用自身
	inExceptionHandler bool

	// PHP 级 set_error_handler 栈（restore_error_handler 弹出）
	errorHandlers []data.Value

	// PHP 级 register_shutdown_function 注册的回调列表
	shutdownCallbacks []data.Value
	shutdownRunOnce   sync.Once

	// 调用深度追踪（用于检测无限递归）
	callDepth int

	// PHP 调用栈（debug_backtrace）
	callStack     []data.CallFrame
	outputBuffers []*strings.Builder

	// 预编译文件注册表（见 compiledFiles sync.Map）
}

func (vm *VM) EnterCall() int {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.callDepth++
	return vm.callDepth
}

func (vm *VM) LeaveCall() {
	vm.mu.Lock()
	if vm.callDepth > 0 {
		vm.callDepth--
	}
	vm.mu.Unlock()
}

func (vm *VM) PushCallFrame(frame data.CallFrame) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.callStack = append(vm.callStack, frame)
}

func (vm *VM) PopCallFrame() {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if n := len(vm.callStack); n > 0 {
		vm.callStack = vm.callStack[:n-1]
	}
}

func (vm *VM) SnapshotCallStack() []data.CallFrame {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if len(vm.callStack) == 0 {
		return nil
	}
	out := make([]data.CallFrame, len(vm.callStack))
	copy(out, vm.callStack)
	return out
}

func (vm *VM) SetPhpFileCache(file string) {
	file = normalizePhpFilePath(file)
	if file == "" {
		return
	}
	vm.phpFileCache.Store(file, struct{}{})
}

func (vm *VM) GetPhpFileCache(file string) bool {
	file = normalizePhpFilePath(file)
	if file == "" {
		return false
	}
	_, ok := vm.phpFileCache.Load(file)
	return ok
}

func (vm *VM) GetIncludeOnceResult(file string) (data.GetValue, bool) {
	file = normalizePhpFilePath(file)
	if file == "" {
		return nil, false
	}
	return syncMapLoad[data.GetValue](&vm.includeOnceResults, file)
}

func (vm *VM) SetIncludeOnceResult(file string, result data.GetValue) {
	file = normalizePhpFilePath(file)
	if file == "" {
		return
	}
	syncMapStore(&vm.includeOnceResults, file, result)
}

// ClearIncludeOnceCache 清空 include_once/require_once 返回值缓存，供热重载使用。
func (vm *VM) ClearIncludeOnceCache() {
	syncMapClear(&vm.includeOnceResults)
}

// beginPhpFileLoad 开始加载文件：
// - alreadyLoaded：已在缓存中
// - wait：其他 goroutine 正在加载，等待其完成后再重试
// - finish：当前 goroutine 负责加载，完成后必须调用（无论成败）
func (vm *VM) beginPhpFileLoad(file string) (alreadyLoaded bool, wait <-chan struct{}, finish func()) {
	if _, ok := vm.phpFileCache.Load(file); ok {
		return true, nil, nil
	}
	vm.mu.Lock()
	if _, ok := vm.phpFileCache.Load(file); ok {
		vm.mu.Unlock()
		return true, nil, nil
	}
	if ch, ok := vm.loadingFiles[file]; ok {
		vm.mu.Unlock()
		return false, ch, nil
	}
	ch := make(chan struct{})
	vm.loadingFiles[file] = ch
	vm.mu.Unlock()
	return false, nil, func() {
		vm.mu.Lock()
		delete(vm.loadingFiles, file)
		close(ch)
		vm.mu.Unlock()
	}
}

// WaitPhpFileLoad 若 file 正在被其他请求加载则阻塞等待；返回 true 表示加载结束后文件已在缓存中。
func (vm *VM) WaitPhpFileLoad(file string) bool {
	file = normalizePhpFilePath(file)
	if file == "" {
		return false
	}
	for {
		if _, ok := vm.phpFileCache.Load(file); ok {
			return true
		}
		vm.mu.Lock()
		ch, loading := vm.loadingFiles[file]
		vm.mu.Unlock()
		if !loading {
			return false
		}
		<-ch
	}
}

// ClearPhpFileCache 清空已加载 PHP 文件缓存，供开发模式热重载使用。
func (vm *VM) ClearPhpFileCache() {
	syncMapClear(&vm.phpFileCache)
	vm.ClearParsedFileCache()
	vm.mu.Lock()
	vm.loadingFiles = make(map[string]chan struct{})
	vm.mu.Unlock()
}

// GlobalValue 读取 PHP 全局变量当前值。
func (vm *VM) GlobalValue(name string) data.Value {
	if zv, ok := syncMapLoad[*data.ZVal](&vm.globalVars, name); ok && zv != nil {
		return zv.Value
	}
	return nil
}

// AddNamespace 添加命名空间路径映射到类路径管理器
func (vm *VM) AddNamespace(namespace string, path string) {
	vm.parser.GetClassPathManager().AddNamespace(namespace, path)
}

func (vm *VM) SetThrowControl(fn func(acl data.Control)) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.acl = fn
}

func (vm *VM) ThrowControl(acl data.Control) {
	vm.mu.Lock()
	handler := vm.exceptionHandler
	fallback := vm.acl
	handleException := false
	if _, ok := acl.(*data.ThrowValue); ok && handler != nil && !vm.inExceptionHandler {
		vm.inExceptionHandler = true
		handleException = true
	}
	vm.mu.Unlock()

	if handleException {
		defer func() {
			vm.mu.Lock()
			vm.inExceptionHandler = false
			vm.mu.Unlock()
		}()
	}

	// 优先尝试调用用户通过 set_exception_handler 注册的 PHP 回调
	if tv, ok := acl.(*data.ThrowValue); ok && handleException {
		// 只在真正有异常对象时尝试回调
		if tv != nil && tv.Error != nil {
			// 目前仅支持 Closure/匿名函数形式的回调（*data.FuncValue）
			if fv, ok := handler.(*data.FuncValue); ok {
				vars := fv.Value.GetVariables()
				if len(vars) == 0 {
					fallback(acl)
					return
				}
				ctx := vm.CreateContext(vars)

				// PHP: handler(Throwable $e)。优先传异常实例；无 Object 时退回 ThrowValue。
				var ex data.Value = tv
				if tv.Object != nil {
					ex = tv.Object
				}

				// 可变参数 fn (...$arguments) 时，首参必须是 [ $e ]，否则 ...$arguments 无法展开。
				arg := ex
				params := fv.Value.GetParams()
				if len(params) > 0 {
					if _, variadic := params[0].(data.Parameters); variadic {
						arg = data.NewArrayValue([]data.Value{ex})
					}
				}
				_ = ctx.SetVariableValue(vars[0], arg)

				if _, hAcl := fv.Call(ctx); hAcl != nil {
					// 如果回调自身又产生未处理控制流，继续交给底层处理
					fallback(hAcl)
					return
				}
				// 回调执行完毕后直接返回，不再走默认处理
				return
			}
		}
	}

	// 默认行为：交给底层 Go 级别处理（打印并退出 / LSP 诊断等）
	fallback(acl)
}

// SetExceptionHandler 设置 PHP 级异常处理回调，返回旧的回调（如果有）
func (vm *VM) SetExceptionHandler(handler data.Value) data.Value {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	old := vm.exceptionHandler
	vm.exceptionHandler = handler
	return old
}

// GetExceptionHandler 返回当前注册的 PHP 级异常处理回调
func (vm *VM) GetExceptionHandler() data.Value {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.exceptionHandler
}

// SetErrorHandler 压入错误处理回调，返回旧的顶层回调（无则 nil）
func (vm *VM) SetErrorHandler(handler data.Value) data.Value {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	var old data.Value
	if n := len(vm.errorHandlers); n > 0 {
		old = vm.errorHandlers[n-1]
	}
	vm.errorHandlers = append(vm.errorHandlers, handler)
	return old
}

// RestoreErrorHandler 弹出当前错误处理回调；成功弹出返回 true
func (vm *VM) RestoreErrorHandler() bool {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	n := len(vm.errorHandlers)
	if n == 0 {
		return false
	}
	vm.errorHandlers = vm.errorHandlers[:n-1]
	return true
}

// GetErrorHandler 返回当前错误处理回调
func (vm *VM) GetErrorHandler() data.Value {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if n := len(vm.errorHandlers); n > 0 {
		return vm.errorHandlers[n-1]
	}
	return nil
}

func (vm *VM) AddClass(c data.ClassStmt) data.Control {
	name := c.GetName()
	if has, ok := syncMapLoad[data.ClassStmt](&vm.classMap, name); ok {
		cFrom := c.GetFrom()
		hasFrom := has.GetFrom()
		if cFrom != nil && hasFrom != nil && utils.SamePhpFile(cFrom.GetSource(), hasFrom.GetSource()) {
			return nil // 同文件重复引入，跳过
		}
		return data.NewErrorThrow(cFrom, fmt.Errorf("已存在同名的 class: %s", name))
	}
	if has, ok := syncMapLoad[data.InterfaceStmt](&vm.interfaceMap, name); ok {
		cFrom := c.GetFrom()
		hasFrom := has.GetFrom()
		if cFrom != nil && hasFrom != nil && utils.SamePhpFile(cFrom.GetSource(), hasFrom.GetSource()) {
			return nil
		}
		return data.NewErrorThrow(cFrom, fmt.Errorf("已存在同名的类或接口: %s", name))
	}
	syncMapStore(&vm.classMap, name, c)
	return nil
}

func (vm *VM) AddInterface(i data.InterfaceStmt) data.Control {
	name := i.GetName()
	if has, ok := syncMapLoad[data.ClassStmt](&vm.classMap, name); ok {
		iFrom := i.GetFrom()
		hasFrom := has.GetFrom()
		if iFrom != nil && hasFrom != nil && utils.SamePhpFile(iFrom.GetSource(), hasFrom.GetSource()) {
			return nil // 同文件不需要报错
		}
		return data.NewErrorThrow(iFrom, fmt.Errorf("已存在同名的 interface: %s", name))
	}
	if has, ok := syncMapLoad[data.InterfaceStmt](&vm.interfaceMap, name); ok {
		iFrom := i.GetFrom()
		hasFrom := has.GetFrom()
		if iFrom != nil && hasFrom != nil && utils.SamePhpFile(iFrom.GetSource(), hasFrom.GetSource()) {
			return nil // 同文件不需要报错
		}
		return data.NewErrorThrow(iFrom, fmt.Errorf("已存在同名的类或接口: %s", name))
	}
	syncMapStore(&vm.interfaceMap, name, i)
	return nil
}

func (vm *VM) findClassCaseInsensitive(name string) (data.ClassStmt, bool) {
	if v, ok := syncMapLoad[data.ClassStmt](&vm.classMap, name); ok {
		return v, true
	}
	var found data.ClassStmt
	vm.classMap.Range(func(key, value any) bool {
		if strings.EqualFold(key.(string), name) {
			found = value.(data.ClassStmt)
			return false
		}
		return true
	})
	if found != nil {
		return found, true
	}
	return nil, false
}

func (vm *VM) GetClass(pkg string) (data.ClassStmt, bool) {
	return vm.findClassCaseInsensitive(pkg)
}

func (vm *VM) GetOrLoadClass(pkg string) (data.ClassStmt, data.Control) {
	if len(pkg) == 0 {
		return nil, nil
	}
	if pkg[0:1] == "\\" {
		pkg = pkg[1:]
	}

	if v, ok := vm.findClassCaseInsensitive(pkg); ok {
		return v, nil
	}

	acl := vm.parser.GetClassPathManager().LoadClass(pkg, vm.parser)
	if acl != nil {
		return nil, acl
	}

	if v, ok := vm.findClassCaseInsensitive(pkg); ok {
		return v, nil
	}
	// 兼容调用方把 interface 误当 class 查询的场景（例如第三方框架类型判断链），
	// 此处不抛错，让上层按“非 class”分支继续处理。
	if _, ok := vm.lookupInterface(pkg); ok {
		return nil, nil
	}

	return nil, utils.NewThrowf("找不到 %s; class 定义需要和文件名称一致才能自动加载", pkg)
}

func (vm *VM) lookupInterface(pkg string) (data.InterfaceStmt, bool) {
	return syncMapLoad[data.InterfaceStmt](&vm.interfaceMap, pkg)
}

func (vm *VM) LoadPkg(pkg string) (data.GetValue, data.Control) {
	if len(pkg) == 0 {
		return nil, nil
	}
	if pkg[0:1] == "\\" {
		temp := pkg[1:]
		if c, ok := vm.findClassCaseInsensitive(temp); ok {
			return c, nil
		}
		if c, ok := vm.lookupInterface(temp); ok {
			return c, nil
		}
	}

	if c, ok := vm.findClassCaseInsensitive(pkg); ok {
		return c, nil
	}
	if c, ok := vm.lookupInterface(pkg); ok {
		return c, nil
	}

	acl := vm.parser.GetClassPathManager().LoadClass(pkg, vm.parser)
	if acl != nil {
		return nil, acl
	}
	if c, ok := vm.findClassCaseInsensitive(pkg); ok {
		return c, nil
	}
	if c, ok := vm.lookupInterface(pkg); ok {
		return c, nil
	}

	return nil, nil
}

func (vm *VM) GetInterface(pkg string) (data.InterfaceStmt, bool) {
	return vm.lookupInterface(pkg)
}

func (vm *VM) GetOrLoadInterface(pkg string) (data.InterfaceStmt, data.Control) {
	if len(pkg) == 0 {
		return nil, nil
	}
	if pkg[0:1] == "\\" {
		pkg = pkg[1:]
	}

	if inf, ok := vm.lookupInterface(pkg); ok {
		return inf, nil
	}

	// 使用 LoadClass 来加载接口（接口和类使用相同的加载机制）
	acl := vm.parser.GetClassPathManager().LoadClass(pkg, vm.parser)
	if acl != nil {
		return nil, acl
	}

	if inf, ok := vm.lookupInterface(pkg); ok {
		return inf, nil
	}

	return nil, utils.NewThrowf("找不到 %s; interface 定义需要和文件名称一致才能自动加载", pkg)
}

func (vm *VM) AddFunc(f data.FuncStmt) data.Control {
	name := f.GetName()
	if _, ok := syncMapLoad[data.FuncStmt](&vm.funcMap, name); ok {
		switch ff := f.(type) {
		case node.GetFrom:
			return data.NewErrorThrow(ff.GetFrom(), fmt.Errorf("已存在同名的 function: %s", name))
		default:
			return utils.NewThrowf("已存在同名的 function: %s", name)
		}
	}
	syncMapStore(&vm.funcMap, name, f)
	return nil
}

func (vm *VM) GetFunc(pkg string) (data.FuncStmt, bool) {
	if v, ok := syncMapLoad[data.FuncStmt](&vm.funcMap, pkg); ok {
		return v, true
	}
	if len(pkg) > 0 && pkg[0:1] == "\\" {
		return syncMapLoad[data.FuncStmt](&vm.funcMap, pkg[1:])
	}
	return nil, false
}

// AllFuncs 返回 VM 中已注册的全部函数（按名称排序）。
func (vm *VM) AllFuncs() []data.FuncStmt {
	funcs := make([]data.FuncStmt, 0)
	vm.funcMap.Range(func(_, value any) bool {
		funcs = append(funcs, value.(data.FuncStmt))
		return true
	})
	sort.Slice(funcs, func(i, j int) bool {
		return funcs[i].GetName() < funcs[j].GetName()
	})
	return funcs
}

// AllClasses 返回 VM 中已注册的全部类（按名称排序）。
func (vm *VM) AllClasses() []data.ClassStmt {
	classes := make([]data.ClassStmt, 0)
	vm.classMap.Range(func(_, value any) bool {
		classes = append(classes, value.(data.ClassStmt))
		return true
	})
	sort.Slice(classes, func(i, j int) bool {
		return classes[i].GetName() < classes[j].GetName()
	})
	return classes
}

func (vm *VM) CreateContext(vars []data.Variable) data.Context {
	return vm.ctx.CreateContext(vars)
}

// EvalCode 执行 eval() 传入的 PHP 代码（在当前上下文中）
func (vm *VM) EvalCode(code string, ctx data.Context, evalFrom data.From) (data.GetValue, data.Control) {
	p := vm.parser.Clone()
	parentFile := ""
	parentLine := 0
	if evalFrom != nil {
		parentFile = evalFrom.GetSource()
		parentLine, _ = evalFrom.GetStartPosition()
		parentLine++
	}
	evalPath := fmt.Sprintf("%s(%d) : eval()'d code", parentFile, parentLine)
	src := strings.TrimSpace(code)
	if !strings.HasPrefix(src, "<?") {
		src = "<?php\n" + src
	}
	program, acl := p.ParseString(src, evalPath)
	if acl != nil {
		return nil, acl
	}
	return program.GetValue(ctx)
}

func (vm *VM) RegisterCompiledFile(file string, fn func() (data.GetValue, []data.Variable)) {
	file = normalizePhpFilePath(file)
	syncMapStore(&vm.compiledFiles, file, fn)
}

func (vm *VM) RunCompiledFile(file string) (data.GetValue, data.Control) {
	file = normalizePhpFilePath(file)
	for {
		loaded, wait, finish := vm.beginPhpFileLoad(file)
		if loaded {
			return nil, nil
		}
		if wait != nil {
			<-wait
			continue
		}

		fn, ok := syncMapLoad[func() (data.GetValue, []data.Variable)](&vm.compiledFiles, file)
		if !ok {
			finish()
			return nil, utils.NewThrowf("run_php_file: 未找到预编译文件 %s", file)
		}
		program, vars := fn()
		ctx := vm.CreateContext(vars)
		vm.RegisterGlobalContext(vars, ctx)
		result, ctrl := program.GetValue(ctx)
		if ctrl == nil {
			vm.SetPhpFileCache(file)
		}
		finish()
		return result, ctrl
	}
}

func (vm *VM) LoadAndRun(file string) (data.GetValue, data.Control) {
	file = normalizePhpFilePath(file)
	for {
		loaded, wait, finish := vm.beginPhpFileLoad(file)
		if loaded {
			return nil, nil
		}
		if wait != nil {
			<-wait
			continue
		}

		data.ResetUserOutput()
		p := vm.parser.Clone()

		program, acl := p.ParseFile(file)
		if acl != nil {
			finish()
			return nil, acl
		}

		vars := p.GetVariables()
		ctx := vm.CreateContext(vars)
		vm.RegisterGlobalContext(vars, ctx)
		result, ctrl := program.GetValue(ctx)

		if ctrl == nil {
			vm.SetPhpFileCache(file)
		}
		finish()
		return result, ctrl
	}
}

// LoadInCallerContext 在独立文件作用域中执行，但按名注入调用者已有变量（含 extract 动态槽）。
// 不设置 PhpFileCache，允许同一视图文件被多次 require（PhpEngine）。
//
// 对调用者中不存在的变量，与 $GLOBALS 对齐：顶层 include（如 tests/run_tests.php）
// 里的赋值可被函数内 global 关键字看到。
func (vm *VM) LoadInCallerContext(parent data.Context, file string) (data.GetValue, data.Control) {
	file = normalizePhpFilePath(file)

	program, vars, acl := vm.ParseFileCached(file)
	if acl != nil {
		return nil, acl
	}

	ctx := vm.CreateContext(vars)

	for _, variable := range vars {
		name := variable.GetName()
		if name == "" {
			continue
		}
		if val, ok := parent.GetVariableByName(name); ok && val != nil {
			// 调用者已有（extract / 闭包 use）：优先注入，不覆盖为全局槽
			if ctl := variable.SetValue(ctx, val); ctl != nil {
				return nil, ctl
			}
			continue
		}
		// 顶层 include：共享/注册全局 ZVal，供 global 关键字使用
		vm.bindIncludedVarToGlobal(name, variable.GetIndex(), ctx)
	}

	result, ctrl := program.GetValue(ctx)
	return result, ctrl
}

// bindIncludedVarToGlobal 将被引入文件的变量槽与 $GLOBALS 对齐。
// 若全局已有同名 ZVal 则复用；否则把当前槽注册进全局表。
func (vm *VM) bindIncludedVarToGlobal(name string, index int, ctx data.Context) {
	if existing, ok := syncMapLoad[*data.ZVal](&vm.globalVars, name); ok && existing != nil {
		ctx.SetIndexZVal(index, existing)
		return
	}
	zv := ctx.GetIndexZVal(index)
	if zv != nil {
		vm.globalVars.Store(name, zv)
	}
}

// CompileLoad 编译模式专用：仅解析文件并注册类/函数/接口，不执行顶层代码。
func (vm *VM) CompileLoad(file string) data.Control {
	file = normalizePhpFilePath(file)
	for {
		loaded, wait, finish := vm.beginPhpFileLoad(file)
		if loaded {
			return nil
		}
		if wait != nil {
			<-wait
			continue
		}

		p := vm.parser.Clone()
		_, acl := p.ParseFile(file)
		if acl == nil {
			vm.SetPhpFileCache(file)
		}
		finish()
		return acl
	}
}

func bindTemplateVariables(ctx data.Context, varList []data.Variable, props map[string]data.Value) {
	for name, value := range props {
		for _, variable := range varList {
			if variable.GetName() == name {
				variable.SetValue(ctx, value)
			}
		}
	}
}

func templatePropsFromArray(arr *data.ArrayValue) map[string]data.Value {
	props := make(map[string]data.Value)
	for _, z := range arr.List {
		if z == nil || z.Name == "" {
			continue
		}
		props[z.Name] = z.Value
	}
	return props
}

func (vm *VM) ParseFile(file string, object data.Value) (data.Value, data.Control) {
	// 解析文件
	p := vm.parser.Clone()

	program, acl := p.ParseFile(file)
	if acl != nil {
		return nil, acl
	}

	varList := p.GetVariables()
	ctx := vm.CreateContext(varList)
	switch v := object.(type) {
	case *data.ObjectValue:
		bindTemplateVariables(ctx, varList, v.GetProperties())
	case *data.ClassValue:
		bindTemplateVariables(ctx, varList, v.GetProperties())
	case *data.ArrayValue:
		// 关联数组（字符串键）与 object 一样按名注入；纯数字下标不映射到模板变量
		bindTemplateVariables(ctx, varList, templatePropsFromArray(v))
	default:
		return nil, utils.NewThrowf("DIY解析文件无法设置指定值到文件域, file(%s)", file)
	}

	v, acl := program.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	if vv, ok := v.(data.Value); ok {
		return vv, nil
	}

	return data.NewNullValue(), nil
}

// SetConstant 设置全局常量
func (vm *VM) SetConstant(name string, value data.Value) data.Control {
	if _, ok := syncMapLoad[data.Value](&vm.constantMap, name); ok {
		return utils.NewThrowf("常量 %s 已经定义，不能重新定义", name)
	}
	syncMapStore(&vm.constantMap, name, value)
	return nil
}

// GetConstant 获取全局常量
func (vm *VM) GetConstant(name string) (data.Value, bool) {
	if len(name) > 0 && name[0:1] == "\\" {
		name = name[1:]
	}
	return syncMapLoad[data.Value](&vm.constantMap, name)
}

// EnsureGlobalZVal 获取或创建全局变量的 ZVal
// 如果该全局变量不存在，则创建一个初始值为 null 的 ZVal
func (vm *VM) EnsureGlobalZVal(name string) *data.ZVal {
	if zv, ok := syncMapLoad[*data.ZVal](&vm.globalVars, name); ok {
		return zv
	}
	zv := data.NewZVal(data.NewNullValue())
	actual, _ := vm.globalVars.LoadOrStore(name, zv)
	return actual.(*data.ZVal)
}

// RegisterGlobalContext 将顶层 ctx 中的变量注册到全局变量表
func (vm *VM) RegisterGlobalContext(vars []data.Variable, ctx data.Context) {
	for _, v := range vars {
		if v == nil {
			continue
		}
		name := v.GetName()
		if name == "" {
			continue
		}
		if _, exists := vm.globalVars.Load(name); exists {
			continue
		}
		if zv := ctx.GetIndexZVal(v.GetIndex()); zv != nil {
			vm.globalVars.Store(name, zv)
		}
	}
}
