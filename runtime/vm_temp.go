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
		return &TempVM{Base: v, Cache: make([]Route, 0)}
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

func (t *TempVM) AddClass(c data.ClassStmt) data.Control {
	// 仅注册到临时 VM 的映射中（请求级生效）
	if t.addedClasses == nil {
		t.addedClasses = make(map[string]data.ClassStmt)
	}
	t.addedClasses[c.GetName()] = c
	return nil
}

func (t *TempVM) AddInterface(i data.InterfaceStmt) data.Control {
	if t.addedInterfaces == nil {
		t.addedInterfaces = make(map[string]data.InterfaceStmt)
	}
	t.addedInterfaces[i.GetName()] = i
	return nil
}

func (t *TempVM) AddFunc(f data.FuncStmt) data.Control {
	if t.addedFuncs == nil {
		t.addedFuncs = make(map[string]data.FuncStmt)
	}
	t.addedFuncs[f.GetName()] = f
	return nil
}

func (t *TempVM) CreateContext(vars []data.Variable) data.Context {
	ctx := t.Base.CreateContext(vars)
	if rctx, ok := ctx.(*Context); ok {
		rctx.SetVM(t)
	}
	return ctx
}

func (t *TempVM) LoadAndRun(file string) (data.GetValue, data.Control) {
	p := t.Base.parser.Clone()
	p.SetVM(t)
	t.parser = p

	program, acl := p.ParseFile(file)
	if acl != nil {
		return nil, acl
	}
	return program.GetValue(t.CreateContext(p.GetVariables()))
}

func (t *TempVM) ParseFile(file string, data data.Value) (data.Value, data.Control) {
	return t.Base.ParseFile(file, data)
}

func (t *TempVM) GetClass(pkg string) (data.ClassStmt, bool) {
	ret, ok := t.Base.GetClass(pkg)
	if ok {
		return ret, ok
	}
	if c, ok := t.addedClasses[pkg]; ok {
		return c, true
	}
	return nil, false
}

func (t *TempVM) GetOrLoadClass(pkg string) (data.ClassStmt, data.Control) {
	c, ok := t.Base.GetClass(pkg)
	if ok {
		return c, nil
	}
	// 优先从本请求新增的类中查找
	if t.addedClasses != nil {
		if c, ok := t.addedClasses[pkg]; ok {
			return c, nil
		} else {
			acl := t.parser.GetClassPathManager().LoadClass(pkg, t.parser)
			if acl != nil {
				return nil, acl
			}
		}
		if c, ok := t.addedClasses[pkg]; ok {
			return c, nil
		}
	}

	return nil, data.NewErrorThrow(nil, fmt.Errorf("class %s not found", pkg))
}

func (t *TempVM) GetInterface(pkg string) (data.InterfaceStmt, bool) {
	// 优先从本请求新增的接口中查找
	if t.addedInterfaces != nil {
		if i, ok := t.addedInterfaces[pkg]; ok {
			return i, true
		}
	}
	// 再从 Base VM 查找
	return t.Base.GetInterface(pkg)
}

func (t *TempVM) GetFunc(pkg string) (data.FuncStmt, bool) {
	// 优先从本请求新增的函数中查找
	if t.addedFuncs != nil {
		if f, ok := t.addedFuncs[pkg]; ok {
			return f, true
		}
	}
	// 再从 Base VM 查找
	return t.Base.GetFunc(pkg)
}
func (t *TempVM) RegisterFunction(name string, fn interface{}) data.Control {
	return t.Base.RegisterFunction(name, fn)
}
func (t *TempVM) RegisterReflectClass(name string, instance interface{}) data.Control {
	return t.Base.RegisterReflectClass(name, instance)
}
func (t *TempVM) SetThrowControl(fn func(acl data.Control))    { t.Base.SetThrowControl(fn) }
func (t *TempVM) ThrowControl(acl data.Control)                { t.Base.ThrowControl(acl) }
func (t *TempVM) SetClassPathCache(name, path string)          { t.Base.SetClassPathCache(name, path) }
func (t *TempVM) GetClassPathCache(name string) (string, bool) { return t.Base.GetClassPathCache(name) }
