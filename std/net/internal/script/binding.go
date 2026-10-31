// Package script contains native binding helpers shared by network modules.
package script

import (
	"context"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/utils"
)

type Method struct {
	Name string
	Args []string
	Run  func(data.Context) (data.GetValue, error)
}

func (m *Method) GetName() string            { return m.Name }
func (m *Method) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *Method) GetIsStatic() bool          { return false }
func (m *Method) GetReturnType() data.Types  { return data.NewBaseType("mixed") }
func (m *Method) GetParams() []data.GetValue {
	a := []data.GetValue{}
	for i, n := range m.Args {
		a = append(a, node.NewParameter(nil, n, i, nil, nil))
	}
	return a
}
func (m *Method) GetVariables() []data.Variable {
	a := []data.Variable{}
	for i, n := range m.Args {
		a = append(a, node.NewVariable(nil, n, i, nil))
	}
	return a
}
func (m *Method) Call(c data.Context) (data.GetValue, data.Control) {
	v, e := m.Run(c)
	if e != nil {
		return nil, utils.NewThrow(e)
	}
	if v == nil {
		v = data.NewNullValue()
	}
	return v, nil
}

type Class struct {
	node.Node
	Name    string
	Factory func() data.ClassStmt
	Methods []data.Method
	Source  any
}

func (c *Class) GetName() string         { return c.Name }
func (c *Class) GetExtend() *string      { return nil }
func (c *Class) GetImplements() []string { return nil }
func (c *Class) AsString() string        { return c.Name + "{}" }
func (c *Class) GetSource() any          { return c.Source }
func (c *Class) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if c.Factory == nil {
		return nil, utils.NewThrowf("%s cannot be constructed directly", c.Name)
	}
	return data.NewProxyValue(c.Factory(), ctx.CreateBaseContext()), nil
}
func (c *Class) GetMethod(n string) (data.Method, bool) {
	for _, m := range c.Methods {
		if m.GetName() == n {
			return m, true
		}
	}
	return nil, false
}
func (c *Class) GetMethods() []data.Method                { return c.Methods }
func (c *Class) GetConstruct() data.Method                { m, _ := c.GetMethod("__construct"); return m }
func (c *Class) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *Class) GetPropertyList() []data.Property         { return nil }
func Options(ctx data.Context, index int, defaults map[string]int) (map[string]int, error) {
	v, ok := ctx.GetIndexValue(index)
	if !ok {
		return nil, fmt.Errorf("options required")
	}
	a, ok := v.(*data.ArrayValue)
	if !ok {
		return nil, fmt.Errorf("options must be an array")
	}
	out := map[string]int{}
	for k, v := range defaults {
		out[k] = v
	}
	for _, z := range a.Range() {
		if _, ok := out[z.Name]; !ok {
			return nil, fmt.Errorf("unknown option %s", z.Name)
		}
		v, ok := z.ReadValue().(*data.IntValue)
		if !ok {
			return nil, fmt.Errorf("%s must be an integer", z.Name)
		}
		n, e := utils.Convert[int](v)
		if e != nil || n < 1 || n > 1<<30 {
			return nil, fmt.Errorf("%s must be between 1 and 1073741824", z.Name)
		}
		out[z.Name] = n
	}
	return out, nil
}

// Invoke creates the same isolated language scope used by HTTP, with cancellation.
func Invoke(base data.VM, fn data.FuncStmt, goctx context.Context, args func(data.Context) []data.Value) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("callback panic: %v", p)
		}
	}()
	vm, ok := runtime.NewRequestVM(base).(*runtime.RequestVM)
	if !ok {
		return fmt.Errorf("request VM required")
	}
	restore := vm.Activate()
	defer restore()
	defer vm.RunShutdownCallbacks()
	vm.BindContext(goctx)
	c := vm.CreateContext(fn.GetVariables())
	if ctl := data.BindDeclaredArgs(c, fn, args(c)); ctl != nil {
		return fmt.Errorf("callback arguments: %s", ctl.AsString())
	}
	_, ctl := fn.Call(c)
	if ctl != nil {
		return fmt.Errorf("callback failed: %s", ctl.AsString())
	}
	return nil
}

// boundCallback preserves Closure::bindTo's object and visibility scope.
type boundCallback struct {
	data.FuncStmt
	bound *data.BoundFuncValue
}

func (f *boundCallback) Call(ctx data.Context) (data.GetValue, data.Control) {
	return f.bound.Call(ctx)
}

func Callback(ctx data.Context, index int) (data.FuncStmt, error) {
	value, ok := ctx.GetIndexValue(index)
	if ok {
		switch f := value.(type) {
		case *data.FuncValue:
			if f != nil && f.Value != nil {
				return f.Value, nil
			}
		case *data.BoundFuncValue:
			if f != nil && f.Value != nil {
				return &boundCallback{FuncStmt: f.Value, bound: f}, nil
			}
		}
	}
	return nil, fmt.Errorf("callback must be a closure")
}
