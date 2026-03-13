package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// OuterIteratorInterface 实现 PHP 的 OuterIterator 接口
type OuterIteratorInterface struct {
	node.Node
}

func NewOuterIteratorInterface() *OuterIteratorInterface {
	return &OuterIteratorInterface{}
}

func (o *OuterIteratorInterface) GetName() string {
	return "OuterIterator"
}

func (o *OuterIteratorInterface) GetExtend() *string {
	// OuterIterator 继承自 Iterator
	iterator := "Iterator"
	return &iterator
}

func (o *OuterIteratorInterface) GetImplements() []string {
	return nil
}

func (o *OuterIteratorInterface) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (o *OuterIteratorInterface) GetPropertyList() []data.Property {
	return nil
}

func (o *OuterIteratorInterface) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(o, ctx.CreateBaseContext()), nil
}

func (o *OuterIteratorInterface) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "getInnerIterator":
		return &OuterIteratorGetInnerIterator{}, true
	}
	return nil, false
}

func (o *OuterIteratorInterface) GetMethods() []data.Method {
	return []data.Method{
		&OuterIteratorGetInnerIterator{},
	}
}

func (o *OuterIteratorInterface) GetConstruct() data.Method {
	return nil
}

// OuterIteratorGetInnerIterator getInnerIterator 方法
type OuterIteratorGetInnerIterator struct{}

func (m *OuterIteratorGetInnerIterator) GetName() string {
	return "getInnerIterator"
}

func (m *OuterIteratorGetInnerIterator) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *OuterIteratorGetInnerIterator) GetIsStatic() bool {
	return false
}

func (m *OuterIteratorGetInnerIterator) GetVariables() []data.Variable {
	return nil
}

func (m *OuterIteratorGetInnerIterator) GetReturnType() data.Types {
	return nil
}

func (m *OuterIteratorGetInnerIterator) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *OuterIteratorGetInnerIterator) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 这个方法应该由实现类自己提供具体实现
	// RecursiveIteratorIterator 已经实现了这个方法
	return nil, nil
}
