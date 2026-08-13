package attribute

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// OverrideClass PHP 8.3+ 原生 Override 属性类
type OverrideClass struct {
	node.Node
	construct data.Method
}

func NewOverrideClass() *OverrideClass {
	return &OverrideClass{
		construct: &OverrideConstructMethod{},
	}
}

func (c *OverrideClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *OverrideClass) GetName() string         { return "Override" }
func (c *OverrideClass) GetExtend() *string      { return nil }
func (c *OverrideClass) GetImplements() []string { return nil }
func (c *OverrideClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}
func (c *OverrideClass) GetPropertyList() []data.Property { return []data.Property{} }

func (c *OverrideClass) GetMethod(name string) (data.Method, bool) {
	if name == token.ConstructName {
		return c.construct, true
	}
	return nil, false
}

func (c *OverrideClass) GetMethods() []data.Method {
	return []data.Method{c.construct}
}

func (c *OverrideClass) GetConstruct() data.Method { return c.construct }

type OverrideConstructMethod struct{}

func (m *OverrideConstructMethod) GetName() string            { return token.ConstructName }
func (m *OverrideConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *OverrideConstructMethod) GetIsStatic() bool          { return false }
func (m *OverrideConstructMethod) GetParams() []data.GetValue { return []data.GetValue{} }
func (m *OverrideConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (m *OverrideConstructMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
func (m *OverrideConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}
