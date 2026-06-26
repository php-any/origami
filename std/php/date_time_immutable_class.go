package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// DateTimeImmutableClass 实现 PHP DateTimeImmutable 类（简化版）
type DateTimeImmutableClass struct {
	node.Node
}

func NewDateTimeImmutableClass() *DateTimeImmutableClass {
	return &DateTimeImmutableClass{}
}

func (c *DateTimeImmutableClass) GetName() string                               { return "DateTimeImmutable" }
func (c *DateTimeImmutableClass) GetExtend() *string                            { return nil }
func (c *DateTimeImmutableClass) GetImplements() []string                       { return []string{"DateTimeInterface"} }
func (c *DateTimeImmutableClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *DateTimeImmutableClass) GetPropertyList() []data.Property              { return nil }
func (c *DateTimeImmutableClass) GetConstruct() data.Method {
	return &DateTimeImmutableConstructMethod{}
}

func (c *DateTimeImmutableClass) GetMethods() []data.Method {
	return []data.Method{
		&DateTimeImmutableConstructMethod{},
	}
}

func (c *DateTimeImmutableClass) GetMethod(name string) (data.Method, bool) {
	methods := c.GetMethods()
	for _, m := range methods {
		if m.GetName() == name {
			return m, true
		}
	}
	return nil, false
}

func (c *DateTimeImmutableClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

// DateTimeImmutableConstructMethod 构造函数
type DateTimeImmutableConstructMethod struct{}

func (m *DateTimeImmutableConstructMethod) GetName() string               { return "__construct" }
func (m *DateTimeImmutableConstructMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *DateTimeImmutableConstructMethod) GetIsStatic() bool             { return false }
func (m *DateTimeImmutableConstructMethod) GetParams() []data.GetValue    { return nil }
func (m *DateTimeImmutableConstructMethod) GetVariables() []data.Variable { return nil }
func (m *DateTimeImmutableConstructMethod) GetReturnType() data.Types     { return nil }
func (m *DateTimeImmutableConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}
