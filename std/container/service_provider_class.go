package container

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewServiceProviderClass() data.ClassStmt {
	return &ServiceProviderClass{
		registerMethod: &ServiceProviderRegisterMethod{},
		bootMethod:     &ServiceProviderBootMethod{},
	}
}

type ServiceProviderClass struct {
	node.Node
	registerMethod data.Method
	bootMethod     data.Method
}

func (s *ServiceProviderClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(s, ctx), nil
}

func (s *ServiceProviderClass) GetName() string { return "Container\\ServiceProvider" }

func (s *ServiceProviderClass) GetExtend() *string { return nil }

func (s *ServiceProviderClass) GetImplements() []string { return nil }

func (s *ServiceProviderClass) GetProperty(name string) (data.Property, bool) {
	if name == "container" {
		return node.NewProperty(nil, "container", "protected", false, nil), true
	}
	return nil, false
}

func (s *ServiceProviderClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "container", "protected", false, nil),
	}
}

func (s *ServiceProviderClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "register":
		return s.registerMethod, true
	case "boot":
		return s.bootMethod, true
	}
	return nil, false
}

func (s *ServiceProviderClass) GetMethods() []data.Method {
	return []data.Method{s.registerMethod, s.bootMethod}
}

func (s *ServiceProviderClass) GetConstruct() data.Method { return nil }

type ServiceProviderRegisterMethod struct{}

func (m *ServiceProviderRegisterMethod) GetName() string               { return "register" }
func (m *ServiceProviderRegisterMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ServiceProviderRegisterMethod) GetIsStatic() bool             { return false }
func (m *ServiceProviderRegisterMethod) GetParams() []data.GetValue    { return nil }
func (m *ServiceProviderRegisterMethod) GetVariables() []data.Variable { return nil }
func (m *ServiceProviderRegisterMethod) GetReturnType() data.Types     { return nil }
func (m *ServiceProviderRegisterMethod) Call(_ data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}

type ServiceProviderBootMethod struct{}

func (m *ServiceProviderBootMethod) GetName() string               { return "boot" }
func (m *ServiceProviderBootMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ServiceProviderBootMethod) GetIsStatic() bool             { return false }
func (m *ServiceProviderBootMethod) GetParams() []data.GetValue    { return nil }
func (m *ServiceProviderBootMethod) GetVariables() []data.Variable { return nil }
func (m *ServiceProviderBootMethod) GetReturnType() data.Types     { return nil }
func (m *ServiceProviderBootMethod) Call(_ data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}
