package runtime

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
)

// NewTempVM 根据给定 VM 创建/返回一个临时 VM 实例
func NewTempVM(vm data.VM) data.VM {
	switch v := vm.(type) {
	case *TempVM:
		return v
	case *VM:
		return &TempVM{
			Base:            v,
			addedClasses:    make(map[string]data.ClassStmt),
			addedInterfaces: make(map[string]data.InterfaceStmt),
			addedFuncs:      make(map[string]data.FuncStmt),
			Cache:           make([]Route, 0),
		}
	default:
		return vm
	}
}

type Route struct {
	Method string
	Path   string
	Target data.Method
}

// TempVM 用于模拟 php-fpm 请求级生效的 VM（热重载）
// 确保解析阶段（Parser）也绑定到 TempVM
type TempVM struct {
	Base   *VM
	parser *parser.Parser

	addedClasses    map[string]data.ClassStmt
	addedInterfaces map[string]data.InterfaceStmt
	addedFuncs      map[string]data.FuncStmt
	Cache           []Route
}

func (vm *TempVM) AddClass(c data.ClassStmt) data.Control {
	// 仅注册到临时 VM 的映射中（请求级生效）
	vm.addedClasses[c.GetName()] = c
	return nil
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
	return ctx
}

func (vm *TempVM) LoadAndRun(file string) (data.GetValue, data.Control) {
	p := vm.Base.parser.Clone()
	p.SetVM(vm)
	vm.parser = p

	program, acl := p.ParseFile(file)
	if acl != nil {
		return nil, acl
	}
	return program.GetValue(vm.CreateContext(p.GetVariables()))
}

func (vm *TempVM) ParseFile(file string, data data.Value) (data.Value, data.Control) {
	return vm.Base.ParseFile(file, data)
}

func (vm *TempVM) GetClass(pkg string) (data.ClassStmt, bool) {
	ret, ok := vm.Base.GetClass(pkg)
	if ok {
		return ret, ok
	}
	if c, ok := vm.addedClasses[pkg]; ok {
		return c, true
	}
	return nil, false
}

func (vm *TempVM) GetOrLoadClass(pkg string) (data.ClassStmt, data.Control) {
	c, ok := vm.Base.GetClass(pkg)
	if ok {
		return c, nil
	}
	// 优先从本请求新增的类中查找
	if c, ok := vm.addedClasses[pkg]; ok {
		return c, nil
	} else {
		acl := vm.parser.GetClassPathManager().LoadClass(pkg, vm.parser)
		if acl != nil {
			return nil, acl
		}
	}
	if c, ok := vm.addedClasses[pkg]; ok {
		return c, nil
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
func (vm *TempVM) SetThrowControl(fn func(acl data.Control)) { vm.Base.SetThrowControl(fn) }
func (vm *TempVM) ThrowControl(acl data.Control)             { vm.Base.ThrowControl(acl) }
func (vm *TempVM) SetClassPathCache(name, path string)       { vm.Base.SetClassPathCache(name, path) }
func (vm *TempVM) GetClassPathCache(name string) (string, bool) {
	return vm.Base.GetClassPathCache(name)
}

// SetConstant 设置全局常量（委托给 Base VM，因为常量是全局的）
func (vm *TempVM) SetConstant(name string, value data.Value) data.Control {
	return vm.Base.SetConstant(name, value)
}

// GetConstant 获取全局常量（从 Base VM 获取，因为常量是全局的）
func (vm *TempVM) GetConstant(name string) (data.Value, bool) {
	return vm.Base.GetConstant(name)
}

// EnsureGlobalZVal 委托给底层 VM，全局变量是全局共享的
func (vm *TempVM) EnsureGlobalZVal(name string) *data.ZVal {
	return vm.Base.EnsureGlobalZVal(name)
}

// SetExceptionHandler 委托到底层 VM，确保异常处理在整个进程内全局生效
func (vm *TempVM) SetExceptionHandler(handler data.Value) data.Value {
	return vm.Base.SetExceptionHandler(handler)
}

// GetExceptionHandler 从底层 VM 获取当前异常处理回调
func (vm *TempVM) GetExceptionHandler() data.Value {
	return vm.Base.GetExceptionHandler()
}

func (vm *TempVM) EnterCall() int {
	return vm.Base.EnterCall()
}

func (vm *TempVM) LeaveCall() {
	vm.Base.LeaveCall()
}
