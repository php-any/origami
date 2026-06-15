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

// NewVM 创建一个新的虚拟机
func NewVM(parser *parser.Parser) data.VM {
	vm := &VM{
		parser:        parser,
		classMap:      make(map[string]data.ClassStmt),
		interfaceMap:  make(map[string]data.InterfaceStmt),
		funcMap:       make(map[string]data.FuncStmt),
		constantMap:   make(map[string]data.Value),
		globalVars:    make(map[string]*data.ZVal),
		phpFileCache:  make(map[string]struct{}),
		compiledFiles: make(map[string]func() (data.GetValue, []data.Variable)),
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

	mu           sync.RWMutex
	classMap     map[string]data.ClassStmt
	interfaceMap map[string]data.InterfaceStmt
	funcMap      map[string]data.FuncStmt
	constantMap  map[string]data.Value // 全局常量映射
	globalVars   map[string]*data.ZVal // 全局变量 ZVal 映射

	// 已引入/加载过的 PHP 文件缓存
	phpFileCache map[string]struct{}

	acl func(acl data.Control)

	// PHP 级 set_exception_handler 注册的回调
	exceptionHandler data.Value
	// 防止在异常处理回调中递归调用自身
	inExceptionHandler bool

	// PHP 级 register_shutdown_function 注册的回调列表
	shutdownCallbacks []data.Value
	shutdownRunOnce   sync.Once

	// 调用深度追踪（用于检测无限递归）
	callDepth int

	// 注解 @Controller 注册的 HTTP 路由（flash 引导后写入此处）
	httpRoutes []Route

	// 预编译文件注册表
	compiledFiles map[string]func() (data.GetValue, []data.Variable)
}

func (vm *VM) EnterCall() int {
	vm.callDepth++
	return vm.callDepth
}

func (vm *VM) LeaveCall() {
	vm.callDepth--
}

func (vm *VM) SetPhpFileCache(file string) {
	file = normalizePhpFilePath(file)
	if file == "" {
		return
	}
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.phpFileCache[file] = struct{}{}
}

func (vm *VM) GetPhpFileCache(file string) bool {
	file = normalizePhpFilePath(file)
	if file == "" {
		return false
	}
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	_, ok := vm.phpFileCache[file]
	return ok
}

// AddNamespace 添加命名空间路径映射到类路径管理器
func (vm *VM) AddNamespace(namespace string, path string) {
	vm.parser.GetClassPathManager().AddNamespace(namespace, path)
}

func (vm *VM) SetThrowControl(fn func(acl data.Control)) {
	vm.acl = fn
}

func (vm *VM) ThrowControl(acl data.Control) {
	vm.acl(acl) // TODO 临时调试
	// 优先尝试调用用户通过 set_exception_handler 注册的 PHP 回调
	if tv, ok := acl.(*data.ThrowValue); ok && vm.exceptionHandler != nil && !vm.inExceptionHandler {
		// 只在真正有异常对象时尝试回调
		if tv != nil && tv.Error != nil {
			// 仅处理一次，避免回调内部再次抛出未捕获异常导致无限递归
			vm.inExceptionHandler = true
			defer func() { vm.inExceptionHandler = false }()

			// 目前仅支持 Closure/匿名函数形式的回调（*data.FuncValue）
			if fv, ok := vm.exceptionHandler.(*data.FuncValue); ok {
				// 使用回调自身的变量列表创建上下文
				vars := fv.Value.GetVariables()
				ctx := vm.CreateContext(vars)

				_ = ctx.SetVariableValue(vars[0], acl)

				if _, hAcl := fv.Call(ctx); hAcl != nil {
					// 如果回调自身又产生未处理控制流，继续交给底层处理
					vm.acl(hAcl)
					return
				}
				// 回调执行完毕后直接返回，不再走默认处理
				return
			}
		}
	}

	// 默认行为：交给底层 Go 级别处理（打印并退出 / LSP 诊断等）
	vm.acl(acl)
}

// SetExceptionHandler 设置 PHP 级异常处理回调，返回旧的回调（如果有）
func (vm *VM) SetExceptionHandler(handler data.Value) data.Value {
	old := vm.exceptionHandler
	vm.exceptionHandler = handler
	return old
}

// GetExceptionHandler 返回当前注册的 PHP 级异常处理回调
func (vm *VM) GetExceptionHandler() data.Value {
	return vm.exceptionHandler
}

func (vm *VM) AddClass(c data.ClassStmt) data.Control {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	// 检查 interfaceMap、classMap 中是否已存在
	if has, ok := vm.classMap[c.GetName()]; ok {
		cFrom := c.GetFrom()
		hasFrom := has.GetFrom()
		if cFrom != nil && hasFrom != nil && utils.SamePhpFile(cFrom.GetSource(), hasFrom.GetSource()) {
			return nil // 同文件重复引入，跳过
		}
		return data.NewErrorThrow(cFrom, fmt.Errorf("已存在同名的 class: %s", c.GetName()))
	}
	if has, ok := vm.interfaceMap[c.GetName()]; ok {
		cFrom := c.GetFrom()
		hasFrom := has.GetFrom()
		if cFrom != nil && hasFrom != nil && utils.SamePhpFile(cFrom.GetSource(), hasFrom.GetSource()) {
			return nil
		}
		return data.NewErrorThrow(cFrom, fmt.Errorf("已存在同名的类或接口: %s", c.GetName()))
	}
	vm.classMap[c.GetName()] = c
	return nil
}

func (vm *VM) AddInterface(i data.InterfaceStmt) data.Control {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	// 检查 interfaceMap、classMap 中是否已存在
	if has, ok := vm.classMap[i.GetName()]; ok {
		iFrom := i.GetFrom()
		hasFrom := has.GetFrom()
		if iFrom != nil && hasFrom != nil && utils.SamePhpFile(iFrom.GetSource(), hasFrom.GetSource()) {
			return nil // 同文件不需要报错
		}
		return data.NewErrorThrow(iFrom, fmt.Errorf("已存在同名的 interface: %s", i.GetName()))
	}
	if has, ok := vm.interfaceMap[i.GetName()]; ok {
		iFrom := i.GetFrom()
		hasFrom := has.GetFrom()
		if iFrom != nil && hasFrom != nil && utils.SamePhpFile(iFrom.GetSource(), hasFrom.GetSource()) {
			return nil // 同文件不需要报错
		}
		return data.NewErrorThrow(iFrom, fmt.Errorf("已存在同名的类或接口: %s", i.GetName()))
	}

	vm.interfaceMap[i.GetName()] = i
	return nil
}

func (vm *VM) findClassCaseInsensitive(name string) (data.ClassStmt, bool) {
	if v, ok := vm.classMap[name]; ok {
		return v, true
	}
	for k, v := range vm.classMap {
		if strings.EqualFold(k, name) {
			return v, true
		}
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

	return nil, utils.NewThrowf("找不到 %s; class 定义需要和文件名称一致才能自动加载", pkg)
}

func (vm *VM) LoadPkg(pkg string) (data.GetValue, data.Control) {
	if len(pkg) == 0 {
		return nil, nil
	}
	if pkg[0:1] == "\\" {
		temp := pkg[1:]
		if c, ok := vm.classMap[temp]; ok {
			return c, nil
		}
		if c, ok := vm.interfaceMap[temp]; ok {
			return c, nil
		}
	}

	if c, ok := vm.classMap[pkg]; ok {
		return c, nil
	}
	if c, ok := vm.interfaceMap[pkg]; ok {
		return c, nil
	}

	acl := vm.parser.GetClassPathManager().LoadClass(pkg, vm.parser)
	if acl != nil {
		return nil, acl
	}
	if c, ok := vm.classMap[pkg]; ok {
		return c, nil
	}
	if c, ok := vm.interfaceMap[pkg]; ok {
		return c, nil
	}

	return nil, nil
}

func (vm *VM) GetInterface(pkg string) (data.InterfaceStmt, bool) {
	if inf, ok := vm.interfaceMap[pkg]; ok {
		return inf, true
	}

	return nil, false
}

func (vm *VM) GetOrLoadInterface(pkg string) (data.InterfaceStmt, data.Control) {
	if len(pkg) == 0 {
		return nil, nil
	}
	if pkg[0:1] == "\\" {
		pkg = pkg[1:]
	}

	if inf, ok := vm.interfaceMap[pkg]; ok {
		return inf, nil
	}

	// 使用 LoadClass 来加载接口（接口和类使用相同的加载机制）
	acl := vm.parser.GetClassPathManager().LoadClass(pkg, vm.parser)
	if acl != nil {
		return nil, acl
	}

	if inf, ok := vm.interfaceMap[pkg]; ok {
		return inf, nil
	}

	return nil, utils.NewThrowf("找不到 %s; interface 定义需要和文件名称一致才能自动加载", pkg)
}

func (vm *VM) AddFunc(f data.FuncStmt) data.Control {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	if _, ok := vm.funcMap[f.GetName()]; ok {
		switch ff := f.(type) {
		case node.GetFrom:
			return data.NewErrorThrow(ff.GetFrom(), fmt.Errorf("已存在同名的 function: %s", f.GetName()))
		default:
			return utils.NewThrowf("已存在同名的 function: %s", f.GetName())
		}
	}

	vm.funcMap[f.GetName()] = f
	return nil
}
func (vm *VM) GetFunc(pkg string) (data.FuncStmt, bool) {
	if v, ok := vm.funcMap[pkg]; ok {
		return v, true
	} else if len(pkg) > 0 && pkg[0:1] == "\\" {
		if v, ok := vm.funcMap[pkg[1:]]; ok {
			return v, true
		}
	}
	return nil, false
}

// AllFuncs 返回 VM 中已注册的全部函数（按名称排序）。
func (vm *VM) AllFuncs() []data.FuncStmt {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	funcs := make([]data.FuncStmt, 0, len(vm.funcMap))
	for _, f := range vm.funcMap {
		funcs = append(funcs, f)
	}
	sort.Slice(funcs, func(i, j int) bool {
		return funcs[i].GetName() < funcs[j].GetName()
	})
	return funcs
}

// AllClasses 返回 VM 中已注册的全部类（按名称排序）。
func (vm *VM) AllClasses() []data.ClassStmt {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	classes := make([]data.ClassStmt, 0, len(vm.classMap))
	for _, c := range vm.classMap {
		classes = append(classes, c)
	}
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
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.compiledFiles[file] = fn
}

func (vm *VM) RunCompiledFile(file string) (data.GetValue, data.Control) {
	file = normalizePhpFilePath(file)
	if vm.GetPhpFileCache(file) {
		return nil, nil
	}
	vm.SetPhpFileCache(file)

	vm.mu.RLock()
	fn, ok := vm.compiledFiles[file]
	vm.mu.RUnlock()
	if !ok {
		return nil, utils.NewThrowf("run_php_file: 未找到预编译文件 %s", file)
	}
	program, vars := fn()
	ctx := vm.CreateContext(vars)
	vm.RegisterGlobalContext(vars, ctx)
	result, ctrl := program.GetValue(ctx)
	if data.FlushAllBuffersFn != nil {
		data.FlushAllBuffersFn()
	}
	return result, ctrl
}

func (vm *VM) LoadAndRun(file string) (data.GetValue, data.Control) {
	file = normalizePhpFilePath(file)
	if vm.GetPhpFileCache(file) {
		return nil, nil
	}
	vm.SetPhpFileCache(file)

	data.ResetUserOutput()
	// 解析文件
	p := vm.parser.Clone()

	program, acl := p.ParseFile(file)
	if acl != nil {
		return nil, acl
	}

	vars := p.GetVariables()
	ctx := vm.CreateContext(vars)
	// 将顶层变量注册到全局变量表，供 global 语句使用
	vm.RegisterGlobalContext(vars, ctx)
	result, ctrl := program.GetValue(ctx)

	if data.FlushAllBuffersFn != nil {
		data.FlushAllBuffersFn()
	}

	return result, ctrl
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
	vm.mu.Lock()
	defer vm.mu.Unlock()

	// 如果常量已存在，不允许重新定义
	if _, ok := vm.constantMap[name]; ok {
		return utils.NewThrowf("常量 %s 已经定义，不能重新定义", name)
	}

	vm.constantMap[name] = value
	return nil
}

// GetConstant 获取全局常量
func (vm *VM) GetConstant(name string) (data.Value, bool) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	if len(name) > 0 && name[0:1] == "\\" {
		name = name[1:]
	}

	value, ok := vm.constantMap[name]
	return value, ok
}

// EnsureGlobalZVal 获取或创建全局变量的 ZVal
// 如果该全局变量不存在，则创建一个初始值为 null 的 ZVal
func (vm *VM) EnsureGlobalZVal(name string) *data.ZVal {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if zv, ok := vm.globalVars[name]; ok {
		return zv
	}
	zv := data.NewZVal(data.NewNullValue())
	vm.globalVars[name] = zv
	return zv
}

// RegisterGlobalContext 将顶层 ctx 中的变量注册到全局变量表
func (vm *VM) RegisterGlobalContext(vars []data.Variable, ctx data.Context) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	for _, v := range vars {
		if v == nil {
			continue
		}
		name := v.GetName()
		if _, exists := vm.globalVars[name]; !exists {
			zv := ctx.GetIndexZVal(v.GetIndex())
			if zv != nil {
				vm.globalVars[name] = zv
			}
		}
	}
}
