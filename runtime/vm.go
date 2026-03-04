package runtime

import (
	"fmt"
	"os"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/utils"
)

// NewVM 创建一个新的虚拟机
func NewVM(parser *parser.Parser) data.VM {
	vm := &VM{
		parser:       parser,
		classMap:     make(map[string]data.ClassStmt),
		interfaceMap: make(map[string]data.InterfaceStmt),
		funcMap:      make(map[string]data.FuncStmt),
		constantMap:  make(map[string]data.Value),
		classPathMap: make(map[string]string),
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

	// 类解释过程中的缓存, 用于支持循环依赖
	classPathMap map[string]string

	acl func(acl data.Control)

	// PHP 级 set_exception_handler 注册的回调
	exceptionHandler data.Value
	// 防止在异常处理回调中递归调用自身
	inExceptionHandler bool
}

func (vm *VM) SetClassPathCache(name string, path string) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	vm.classPathMap[name] = path
}

func (vm *VM) GetClassPathCache(name string) (string, bool) {
	path, ok := vm.classPathMap[name]
	return path, ok
}

func (vm *VM) SetThrowControl(fn func(acl data.Control)) {
	vm.acl = fn
}

func (vm *VM) ThrowControl(acl data.Control) {
	// vm.acl(acl)
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
	if _, ok := vm.classMap[c.GetName()]; ok {
		return data.NewErrorThrow(c.GetFrom(), fmt.Errorf("已存在同名的 class: %s", c.GetName()))
	}
	if _, ok := vm.interfaceMap[c.GetName()]; ok {
		return data.NewErrorThrow(c.GetFrom(), fmt.Errorf("已存在同名的类或接口: %s", c.GetName()))
	}
	vm.classMap[c.GetName()] = c
	return nil
}

func (vm *VM) AddInterface(i data.InterfaceStmt) data.Control {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	// 检查 interfaceMap、classMap 中是否已存在
	if has, ok := vm.classMap[i.GetName()]; ok {
		if i.GetFrom().GetSource() == has.GetFrom().GetSource() {
			return nil // 同文件不需要报错
		}
		return data.NewErrorThrow(i.GetFrom(), fmt.Errorf("已存在同名的 interface: %s", i.GetName()))
	}
	if has, ok := vm.interfaceMap[i.GetName()]; ok {
		if i.GetFrom().GetSource() == has.GetFrom().GetSource() {
			return nil // 同文件不需要报错
		}
		return data.NewErrorThrow(i.GetFrom(), fmt.Errorf("已存在同名的类或接口: %s", i.GetName()))
	}

	vm.interfaceMap[i.GetName()] = i
	return nil
}

func (vm *VM) GetClass(pkg string) (data.ClassStmt, bool) {
	if v, ok := vm.classMap[pkg]; ok {
		return v, true
	}

	return nil, false
}

func (vm *VM) GetOrLoadClass(pkg string) (data.ClassStmt, data.Control) {
	if pkg[0:1] == "\\" {
		pkg = pkg[1:]
	}

	if v, ok := vm.classMap[pkg]; ok {
		return v, nil
	}

	acl := vm.parser.GetClassPathManager().LoadClass(pkg, vm.parser)
	if acl != nil {
		return nil, acl
	}

	if v, ok := vm.classMap[pkg]; ok {
		return v, nil
	}

	return nil, utils.NewThrowf("找不到 %s; class 定义需要和文件名称一致才能自动加载", pkg)
}

func (vm *VM) LoadPkg(pkg string) (data.GetValue, data.Control) {
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

	// 先检查是否存在对应的物理文件；存在则按正常流程加载
	if _, ok := vm.parser.GetClassPathManager().FindClassFile(pkg); ok {
		if acl := vm.parser.GetClassPathManager().LoadClass(pkg, vm.parser); acl != nil {
			return nil, acl
		}
		if c, ok := vm.classMap[pkg]; ok {
			return c, nil
		}
		if c, ok := vm.interfaceMap[pkg]; ok {
			return c, nil
		}
	}

	// 若未找到物理文件，则仍尝试通过 autoload 机制加载（例如 Composer 的 PSR-4），
	// LoadClass 在未找到类时会返回一个 “类不存在或无法加载” 的错误，这里按约定吞掉，
	// 将“未找到类/接口”视为正常情况并返回 nil。
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

func (vm *VM) CreateContext(vars []data.Variable) data.Context {
	return vm.ctx.CreateContext(vars)
}

func (vm *VM) LoadAndRun(file string) (data.GetValue, data.Control) {
	// 解析文件
	p := vm.parser.Clone()

	program, acl := p.ParseFile(file)
	if acl != nil {
		return nil, acl
	}

	return program.GetValue(vm.CreateContext(p.GetVariables()))
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
		for name, value := range v.GetProperties() {
			for _, variable := range varList {
				if variable.GetName() == name {
					variable.SetValue(ctx, value)
				}
			}
		}
	case *data.ClassValue:
		for name, value := range v.GetProperties() {
			for _, variable := range varList {
				if variable.GetName() == name {
					variable.SetValue(ctx, value)
				}
			}
		}
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
