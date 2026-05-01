package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

type ReflectionPropertyClass struct {
	node.Node
}

func (c *ReflectionPropertyClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *ReflectionPropertyClass) GetName() string                               { return "ReflectionProperty" }
func (c *ReflectionPropertyClass) GetExtend() *string                            { return nil }
func (c *ReflectionPropertyClass) GetImplements() []string                       { return nil }
func (c *ReflectionPropertyClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *ReflectionPropertyClass) GetPropertyList() []data.Property              { return nil }
func (c *ReflectionPropertyClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case token.ConstructName:
		return &ReflectionPropertyConstructMethod{}, true
	}
	return nil, false
}
func (c *ReflectionPropertyClass) GetMethods() []data.Method {
	return []data.Method{&ReflectionPropertyConstructMethod{}}
}
func (c *ReflectionPropertyClass) GetConstruct() data.Method {
	return &ReflectionPropertyConstructMethod{}
}

type ReflectionPropertyConstructMethod struct{}

func (m *ReflectionPropertyConstructMethod) GetName() string            { return "__construct" }
func (m *ReflectionPropertyConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ReflectionPropertyConstructMethod) GetIsStatic() bool          { return false }
func (m *ReflectionPropertyConstructMethod) GetReturnType() data.Types  { return nil }
func (m *ReflectionPropertyConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "class", 0, nil, nil),
		node.NewParameter(nil, "property", 1, nil, nil),
	}
}
func (m *ReflectionPropertyConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "class", 0, data.Mixed{}),
		node.NewVariable(nil, "property", 1, data.Mixed{}),
	}
}
func (m *ReflectionPropertyConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}
