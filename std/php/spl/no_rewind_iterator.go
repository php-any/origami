package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NoRewindIteratorClass 实现 PHP NoRewindIterator（rewind 不调用内部迭代器�?
type NoRewindIteratorClass struct {
	node.Node
}

func NewNoRewindIteratorClass() *NoRewindIteratorClass {
	return &NoRewindIteratorClass{}
}

func (c *NoRewindIteratorClass) GetName() string { return "NoRewindIterator" }
func (c *NoRewindIteratorClass) GetExtend() *string {
	parent := "IteratorIterator"
	return &parent
}
func (c *NoRewindIteratorClass) GetImplements() []string { return nil }
func (c *NoRewindIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}
func (c *NoRewindIteratorClass) GetPropertyList() []data.Property { return nil }
func (c *NoRewindIteratorClass) GetConstruct() data.Method        { return &NRIConstructMethod{} }

func (c *NoRewindIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.SetProperty(iiValidKey, data.NewBoolValue(false))
	cv.SetProperty(iiCurValKey, data.NewNullValue())
	cv.SetProperty(iiCurKeyKey, data.NewNullValue())
	return cv, nil
}

func (c *NoRewindIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &NRIConstructMethod{}, true
	case "rewind":
		return &NRIRewindMethod{}, true
	}
	return nil, false
}

func (c *NoRewindIteratorClass) GetMethods() []data.Method {
	return []data.Method{&NRIConstructMethod{}, &NRIRewindMethod{}}
}

type NRIConstructMethod struct{}

func (m *NRIConstructMethod) GetName() string            { return "__construct" }
func (m *NRIConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *NRIConstructMethod) GetIsStatic() bool          { return false }
func (m *NRIConstructMethod) GetReturnType() data.Types  { return nil }
func (m *NRIConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "iterator", 0, nil, data.NewBaseType("Iterator")),
	}
}
func (m *NRIConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "iterator", 0, data.NewBaseType("Iterator")),
	}
}
func (m *NRIConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	it, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, nil
	}
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	iiSetInner(cv, it)
	iiSyncFromInner(cv)
	return nil, nil
}

type NRIRewindMethod struct{}

func (m *NRIRewindMethod) GetName() string               { return "rewind" }
func (m *NRIRewindMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *NRIRewindMethod) GetIsStatic() bool             { return false }
func (m *NRIRewindMethod) GetReturnType() data.Types     { return nil }
func (m *NRIRewindMethod) GetParams() []data.GetValue    { return nil }
func (m *NRIRewindMethod) GetVariables() []data.Variable { return nil }
func (m *NRIRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// NoRewindIterator::rewind 为空操作
	return nil, nil
}
