package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RecursiveIteratorInterface 实现 PHP 的 RecursiveIterator 接口
type RecursiveIteratorInterface struct {
	node.Node
}

func NewRecursiveIteratorInterface() *RecursiveIteratorInterface {
	return &RecursiveIteratorInterface{}
}

func (r *RecursiveIteratorInterface) GetName() string {
	return "RecursiveIterator"
}

func (r *RecursiveIteratorInterface) GetExtend() *string {
	// RecursiveIterator 继承自 Iterator
	iterator := "Iterator"
	return &iterator
}

func (r *RecursiveIteratorInterface) GetImplements() []string {
	return nil
}

func (r *RecursiveIteratorInterface) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (r *RecursiveIteratorInterface) GetPropertyList() []data.Property {
	return nil
}

func (r *RecursiveIteratorInterface) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(r, ctx.CreateBaseContext()), nil
}

func (r *RecursiveIteratorInterface) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "hasChildren":
		return &RecursiveIteratorHasChildren{}, true
	case "getChildren":
		return &RecursiveIteratorGetChildren{}, true
	}
	return nil, false
}

func (r *RecursiveIteratorInterface) GetMethods() []data.Method {
	return []data.Method{
		&RecursiveIteratorHasChildren{},
		&RecursiveIteratorGetChildren{},
	}
}

func (r *RecursiveIteratorInterface) GetConstruct() data.Method {
	return nil
}

// RecursiveIteratorHasChildren hasChildren 方法
type RecursiveIteratorHasChildren struct{}

func (m *RecursiveIteratorHasChildren) GetName() string {
	return "hasChildren"
}

func (m *RecursiveIteratorHasChildren) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *RecursiveIteratorHasChildren) GetIsStatic() bool {
	return false
}

func (m *RecursiveIteratorHasChildren) GetVariables() []data.Variable {
	return nil
}

func (m *RecursiveIteratorHasChildren) GetReturnType() data.Types {
	return nil
}

func (m *RecursiveIteratorHasChildren) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *RecursiveIteratorHasChildren) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}

// RecursiveIteratorGetChildren getChildren 方法
type RecursiveIteratorGetChildren struct{}

func (m *RecursiveIteratorGetChildren) GetName() string {
	return "getChildren"
}

func (m *RecursiveIteratorGetChildren) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *RecursiveIteratorGetChildren) GetIsStatic() bool {
	return false
}

func (m *RecursiveIteratorGetChildren) GetVariables() []data.Variable {
	return nil
}

func (m *RecursiveIteratorGetChildren) GetReturnType() data.Types {
	return nil
}

func (m *RecursiveIteratorGetChildren) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *RecursiveIteratorGetChildren) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}
