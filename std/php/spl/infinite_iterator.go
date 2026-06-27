package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// InfiniteIteratorClass 实现 PHP InfiniteIterator（无限重复内部迭代器�?
type InfiniteIteratorClass struct {
	node.Node
}

func NewInfiniteIteratorClass() *InfiniteIteratorClass {
	return &InfiniteIteratorClass{}
}

func (c *InfiniteIteratorClass) GetName() string { return "InfiniteIterator" }
func (c *InfiniteIteratorClass) GetExtend() *string {
	parent := "IteratorIterator"
	return &parent
}
func (c *InfiniteIteratorClass) GetImplements() []string { return nil }
func (c *InfiniteIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}
func (c *InfiniteIteratorClass) GetPropertyList() []data.Property { return nil }
func (c *InfiniteIteratorClass) GetConstruct() data.Method        { return &InfIConstructMethod{} }

func (c *InfiniteIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.SetProperty(iiValidKey, data.NewBoolValue(false))
	cv.SetProperty(iiCurValKey, data.NewNullValue())
	cv.SetProperty(iiCurKeyKey, data.NewNullValue())
	return cv, nil
}

func (c *InfiniteIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &InfIConstructMethod{}, true
	case "next":
		return &InfINextMethod{}, true
	}
	return nil, false
}

func (c *InfiniteIteratorClass) GetMethods() []data.Method {
	return []data.Method{&InfIConstructMethod{}, &InfINextMethod{}}
}

type InfIConstructMethod struct{}

func (m *InfIConstructMethod) GetName() string            { return "__construct" }
func (m *InfIConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *InfIConstructMethod) GetIsStatic() bool          { return false }
func (m *InfIConstructMethod) GetReturnType() data.Types  { return nil }
func (m *InfIConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "iterator", 0, nil, data.NewBaseType("Iterator")),
	}
}
func (m *InfIConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "iterator", 0, data.NewBaseType("Iterator")),
	}
}
func (m *InfIConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	it, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, nil
	}
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	iiSetInner(cv, it)
	iiSetValid(cv, false)
	return nil, nil
}

type InfINextMethod struct{}

func (m *InfINextMethod) GetName() string               { return "next" }
func (m *InfINextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *InfINextMethod) GetIsStatic() bool             { return false }
func (m *InfINextMethod) GetReturnType() data.Types     { return nil }
func (m *InfINextMethod) GetParams() []data.GetValue    { return nil }
func (m *InfINextMethod) GetVariables() []data.Variable { return nil }
func (m *InfINextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	inner := iiGetInner(cv)
	iiCallInnerMethod(inner, "next")
	if !iiInnerValid(inner) {
		iiCallInnerMethod(inner, "rewind")
	}
	iiSyncFromInner(cv)
	return nil, nil
}
