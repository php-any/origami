package runtime

import (
	"errors"
	"fmt"
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
		classPathMap: make(map[string]string),
		acl: func(acl data.Control) {
			parser.ShowControl(acl)
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

	// 类解释过程中的缓存, 用于支持循环依赖
	classPathMap map[string]string

	acl func(acl data.Control)
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
	vm.acl(acl)
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

	return nil, data.NewErrorThrow(nil, errors.New("找不到 class; class 定义需要和文件名称一致才能自动加载"))
}

func (vm *VM) LoadPkg(pkg string) (data.GetValue, data.Control) {
	if c, ok := vm.classMap[pkg]; ok {
		return c, nil
	}
	if c, ok := vm.interfaceMap[pkg]; ok {
		return c, nil
	}

	_, ok := vm.parser.GetClassPathManager().FindClassFile(pkg)
	if !ok {
		return nil, nil
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
