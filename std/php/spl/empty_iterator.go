package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// EmptyIteratorClass 实现 PHP SPL �?EmptyIterator
type EmptyIteratorClass struct {
	node.Node
}

func NewEmptyIteratorClass() *EmptyIteratorClass {
	return &EmptyIteratorClass{}
}

func (c *EmptyIteratorClass) GetName() string    { return "EmptyIterator" }
func (c *EmptyIteratorClass) GetExtend() *string { return nil }
func (c *EmptyIteratorClass) GetImplements() []string {
	return []string{"Iterator"}
}
func (c *EmptyIteratorClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *EmptyIteratorClass) GetPropertyList() []data.Property              { return nil }
func (c *EmptyIteratorClass) GetConstruct() data.Method                     { return nil }
func (c *EmptyIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *EmptyIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "rewind":
		return &EmptyIteratorRewindMethod{}, true
	case "valid":
		return &EmptyIteratorValidMethod{}, true
	case "current":
		return &EmptyIteratorCurrentMethod{}, true
	case "key":
		return &EmptyIteratorKeyMethod{}, true
	case "next":
		return &EmptyIteratorNextMethod{}, true
	}
	return nil, false
}

func (c *EmptyIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&EmptyIteratorRewindMethod{},
		&EmptyIteratorValidMethod{},
		&EmptyIteratorCurrentMethod{},
		&EmptyIteratorKeyMethod{},
		&EmptyIteratorNextMethod{},
	}
}

type EmptyIteratorRewindMethod struct{}

func (m *EmptyIteratorRewindMethod) GetName() string            { return "rewind" }
func (m *EmptyIteratorRewindMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *EmptyIteratorRewindMethod) GetIsStatic() bool          { return false }
func (m *EmptyIteratorRewindMethod) GetReturnType() data.Types  { return nil }
func (m *EmptyIteratorRewindMethod) GetParams() []data.GetValue { return nil }
func (m *EmptyIteratorRewindMethod) GetVariables() []data.Variable {
	return nil
}
func (m *EmptyIteratorRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}

type EmptyIteratorValidMethod struct{}

func (m *EmptyIteratorValidMethod) GetName() string            { return "valid" }
func (m *EmptyIteratorValidMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *EmptyIteratorValidMethod) GetIsStatic() bool          { return false }
func (m *EmptyIteratorValidMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *EmptyIteratorValidMethod) GetParams() []data.GetValue { return nil }
func (m *EmptyIteratorValidMethod) GetVariables() []data.Variable {
	return nil
}
func (m *EmptyIteratorValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(false), nil
}

type EmptyIteratorCurrentMethod struct{}

func (m *EmptyIteratorCurrentMethod) GetName() string            { return "current" }
func (m *EmptyIteratorCurrentMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *EmptyIteratorCurrentMethod) GetIsStatic() bool          { return false }
func (m *EmptyIteratorCurrentMethod) GetReturnType() data.Types  { return data.Mixed{} }
func (m *EmptyIteratorCurrentMethod) GetParams() []data.GetValue { return nil }
func (m *EmptyIteratorCurrentMethod) GetVariables() []data.Variable {
	return nil
}
func (m *EmptyIteratorCurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}

type EmptyIteratorKeyMethod struct{}

func (m *EmptyIteratorKeyMethod) GetName() string            { return "key" }
func (m *EmptyIteratorKeyMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *EmptyIteratorKeyMethod) GetIsStatic() bool          { return false }
func (m *EmptyIteratorKeyMethod) GetReturnType() data.Types  { return data.Mixed{} }
func (m *EmptyIteratorKeyMethod) GetParams() []data.GetValue { return nil }
func (m *EmptyIteratorKeyMethod) GetVariables() []data.Variable {
	return nil
}
func (m *EmptyIteratorKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}

type EmptyIteratorNextMethod struct{}

func (m *EmptyIteratorNextMethod) GetName() string            { return "next" }
func (m *EmptyIteratorNextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *EmptyIteratorNextMethod) GetIsStatic() bool          { return false }
func (m *EmptyIteratorNextMethod) GetReturnType() data.Types  { return nil }
func (m *EmptyIteratorNextMethod) GetParams() []data.GetValue { return nil }
func (m *EmptyIteratorNextMethod) GetVariables() []data.Variable {
	return nil
}
func (m *EmptyIteratorNextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}
