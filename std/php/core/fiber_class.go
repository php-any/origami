package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// FiberClass 表示 PHP 8.1+ 的 Fiber 类的简化桩
type FiberClass struct {
	node.Node
}

func NewFiberClass() *FiberClass {
	return &FiberClass{}
}

func (f *FiberClass) GetName() string {
	return "Fiber"
}

func (f *FiberClass) GetExtend() *string {
	return nil
}

func (f *FiberClass) GetImplements() []string {
	return nil
}

func (f *FiberClass) GetMethods() []data.Method {
	return []data.Method{
		&FiberGetCurrentMethod{},
	}
}

func (f *FiberClass) GetMethod(name string) (data.Method, bool) {
	methods := f.GetMethods()
	for _, method := range methods {
		if method.GetName() == name {
			return method, true
		}
	}
	return nil, false
}

func (f *FiberClass) GetConstruct() data.Method {
	return nil
}

func (f *FiberClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (f *FiberClass) GetPropertyList() []data.Property {
	return nil
}

func (f *FiberClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(f, ctx.CreateBaseContext()), nil
}

// FiberGetCurrentMethod 静态方法 getCurrent
type FiberGetCurrentMethod struct{}

func (m *FiberGetCurrentMethod) GetName() string               { return "getCurrent" }
func (m *FiberGetCurrentMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *FiberGetCurrentMethod) GetIsStatic() bool             { return true }
func (m *FiberGetCurrentMethod) GetParams() []data.GetValue    { return nil }
func (m *FiberGetCurrentMethod) GetVariables() []data.Variable { return nil }
func (m *FiberGetCurrentMethod) GetReturnType() data.Types     { return nil }
func (m *FiberGetCurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}
